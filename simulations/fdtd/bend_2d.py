"""2D FDTD (Meep) of a 90-degree bend in a Si3N4 waveguide, via the effective index method.

Measures the transmission from the fundamental mode of a straight input waveguide, through a bend
of radius R, into the fundamental mode of the straight output waveguide. Includes bend radiation
loss and the mode mismatch at both straight-bend junctions.

Usage: python bend_2d.py --t 0.8 --w 1.2 --R 20 --res 20 [--csv file]
"""
import argparse
import csv
import math
import os
import time

import meep as mp

from eim import N_SIO2, WL, channel_neff

p = argparse.ArgumentParser()
p.add_argument("--t", type=float, default=0.8, help="Si3N4 core thickness (um)")
p.add_argument("--w", type=float, default=1.2, help="waveguide width (um)")
p.add_argument("--R", type=float, default=20.0, help="bend radius (um)")
p.add_argument("--res", type=int, default=20, help="pixels per um")
p.add_argument("--pol", choices=["ez", "inplane"], default="ez",
               help="ez: E out of the chip plane (original runs); inplane: E in the plane, the quasi-TE mode consistent with the EIM")
p.add_argument("--until", type=float, default=0, help="fixed run time after the source (0 = automatic: 2.5x the transit time)")
p.add_argument("--csv", default=os.path.join(os.path.dirname(__file__), "..", "results", "fdtd_bends.csv"))
a = p.parse_args()

n_slab, n_eff = channel_neff(a.t, a.w)
PARITY = mp.ODD_Z if a.pol == "ez" else mp.EVEN_Z
core, clad = mp.Medium(index=n_slab), mp.Medium(index=N_SIO2)
R, w, L = a.R, a.w, 4.0
dpml, margin = 1.5, 2.0

# Cell: horizontal input at y=-R coming from x<0; quarter circle centred at the origin; vertical output at x=R.
x0, x1 = -(L + dpml), R + w / 2 + margin + dpml
y0, y1 = -(R + w / 2 + margin + dpml), L + dpml
cell = mp.Vector3(x1 - x0, y1 - y0)
center = mp.Vector3((x0 + x1) / 2, (y0 + y1) / 2)


# Geometry built from objects (not a material function) to enable Meep's subpixel averaging: without it
# the curved boundary becomes a pixel staircase that scatters light (R=20 um, w=1.2 um: 0.23 dB at 20 px/um
# vs 0.14 dB at 30 px/um). Later objects override earlier ones.
big = 4 * (R + L + margin + dpml)
geometry = [
    mp.Cylinder(radius=R + w / 2, material=core, center=mp.Vector3() - center),                         # outer ring
    mp.Cylinder(radius=R - w / 2, material=clad, center=mp.Vector3() - center),                         # inner hole
    mp.Block(size=mp.Vector3(big, big), center=mp.Vector3(-big / 2, 0) - center, material=clad),  # x < 0
    mp.Block(size=mp.Vector3(big, big), center=mp.Vector3(big / 2, big / 2) - center, material=clad),  # y > 0
    mp.Block(size=mp.Vector3(big, w), center=mp.Vector3(-big / 2, -R) - center, material=core),   # input
    mp.Block(size=mp.Vector3(w, big), center=mp.Vector3(R, big / 2) - center, material=core),     # output
]

fcen = 1 / WL
df = 0.2 * fcen
src_x = -L - dpml + 1.0
sources = [mp.EigenModeSource(mp.GaussianSource(fcen, fwidth=df),
                              center=mp.Vector3(src_x, -R) - center, size=mp.Vector3(0, 3 * w),
                              eig_band=1, eig_parity=PARITY, direction=mp.X)]
sim = mp.Simulation(cell_size=cell, geometry_center=mp.Vector3(), boundary_layers=[mp.PML(dpml)],
                    geometry=geometry, default_material=clad, sources=sources, resolution=a.res)

mon_in = sim.add_mode_monitor(fcen, 0, 1, mp.ModeRegion(center=mp.Vector3(src_x + 1.0, -R) - center,
                                                         size=mp.Vector3(0, 3 * w)))
mon_out = sim.add_mode_monitor(fcen, 0, 1, mp.ModeRegion(center=mp.Vector3(R, L - 1.0) - center,
                                                          size=mp.Vector3(3 * w, 0)))

start = time.time()
# Fixed run time after the source: 2.5x the pulse transit to the output monitor. The stop_when_fields_decayed
# criterion ended the run before the DFT had accumulated the whole pulse (R=20 um gave 17.8 dB instead of
# 0.23 dB); with a fixed time, 150 and 400 time units give the same T (converged).
transit = (math.pi * R / 2 + 2 * L) * n_eff
sim.run(until_after_sources=a.until if a.until > 0 else 2.5 * transit + 50)
c_in = sim.get_eigenmode_coefficients(mon_in, [1], eig_parity=PARITY, direction=mp.X).alpha[0, 0, 0]
c_out = sim.get_eigenmode_coefficients(mon_out, [1], eig_parity=PARITY, direction=mp.Y).alpha[0, 0, 0]
T = abs(c_out) ** 2 / abs(c_in) ** 2
loss_db = -10 * math.log10(T)
elapsed = time.time() - start

if mp.am_master():
    print(f"RESULT t={a.t} w={a.w} R={R} res={a.res} n_slab={n_slab:.4f} n_eff={n_eff:.4f} "
          f"T={T:.6f} loss_dB_per_90deg={loss_db:.4f} time_s={elapsed:.0f}")
    new = not os.path.exists(a.csv)
    with open(a.csv, "a", newline="") as f:
        wr = csv.writer(f)
        if new:
            wr.writerow(["t_um", "w_um", "R_um", "res_px_per_um", "n_slab", "n_eff", "T", "loss_db_per_90deg", "time_s"])
        wr.writerow([a.t, a.w, R, a.res, f"{n_slab:.4f}", f"{n_eff:.4f}", f"{T:.6f}", f"{loss_db:.4f}", f"{elapsed:.0f}"])
