"""2D FDTD (Meep) of the adiabatic Si3N4 -> thin-film lithium niobate (TFLN) transition, in a vertical cut.

Simulated plane: x = propagation, Meep's y = height. The lateral Si3N4 width tapers linearly from w_in
to w_tip along L_taper; the effective index method turns each width into a layer index (lateral TM
slab). The LN film (extraordinary index, TE) sits above, separated by a SiO2 gap. The electric field is
lateral (out of the simulated plane, Ez polarization in Meep).

Measures the transmission from the fundamental supermode at the input (confined in Si3N4) into the
fundamental supermode at the output (confined in LN). Representative geometry (assumption), not a
specific foundry process.

Usage: python transition_2d.py --L 50 --res 50 [--csv file]
"""
import argparse
import bisect
import csv
import math
import os
import time

import meep as mp

from eim import N_LN_E, N_SIN, N_SIO2, WL, slab_neff

p = argparse.ArgumentParser()
p.add_argument("--L", type=float, default=50.0, help="taper length (um)")
p.add_argument("--res", type=int, default=50, help="pixels per um")
p.add_argument("--t_sin", type=float, default=0.8, help="Si3N4 thickness (um)")
p.add_argument("--w_in", type=float, default=1.2, help="Si3N4 width at the input (um)")
p.add_argument("--w_tip", type=float, default=0.15, help="taper tip width (um)")
p.add_argument("--gap", type=float, default=0.1, help="SiO2 gap between Si3N4 and LN (um)")
p.add_argument("--t_ln", type=float, default=0.3, help="LN film thickness (um)")
p.add_argument("--csv", default=os.path.join(os.path.dirname(__file__), "..", "results", "fdtd_transition.csv"))
a = p.parse_args()

# Si3N4 layer index as a function of width (lateral EIM, TM), tabulated for fast interpolation.
W_TABLE = [a.w_tip + (a.w_in - a.w_tip) * i / 400 for i in range(401)]
N_TABLE = [slab_neff(w, N_SIN, N_SIO2, N_SIO2, "TM") for w in W_TABLE]


def n_sin_layer(w):
    i = min(max(bisect.bisect_left(W_TABLE, w), 1), len(W_TABLE) - 1)
    w0, w1, n0, n1 = W_TABLE[i - 1], W_TABLE[i], N_TABLE[i - 1], N_TABLE[i]
    return n0 + (n1 - n0) * (w - w0) / (w1 - w0)


dpml, straight, clad = 1.0, 4.0, 1.5
L = a.L
length = L + 2 * straight + 2 * dpml
z_sin0 = -a.t_sin / 2                     # Si3N4 centred at y=0
z_sin1 = a.t_sin / 2
z_ln0 = z_sin1 + a.gap
z_ln1 = z_ln0 + a.t_ln
height = (z_ln1 + clad + dpml) - (z_sin0 - clad - dpml)
yc = ((z_ln1 + clad + dpml) + (z_sin0 - clad - dpml)) / 2
cell = mp.Vector3(length, height)
x_start = -length / 2 + dpml + straight   # taper start


def width_at(x):
    if x <= x_start:
        return a.w_in
    if x >= x_start + L:
        return a.w_tip
    return a.w_in + (a.w_tip - a.w_in) * (x - x_start) / L


media = {}


def medium(n):
    key = round(n, 4)
    if key not in media:
        media[key] = mp.Medium(index=key)
    return media[key]


def eps_at(pt):
    z = pt.y + yc
    if z_sin0 <= z <= z_sin1:
        return medium(n_sin_layer(width_at(pt.x)))
    if z_ln0 <= z <= z_ln1:
        return medium(N_LN_E)
    return medium(N_SIO2)


fcen = 1 / WL
src_x = -length / 2 + dpml + 1.0
mon_in_x = src_x + 1.5
mon_out_x = length / 2 - dpml - 1.5
yspan = mp.Vector3(0, height - 2 * dpml)
sources = [mp.EigenModeSource(mp.GaussianSource(fcen, fwidth=0.2 * fcen), center=mp.Vector3(src_x, 0),
                              size=yspan, eig_band=1, eig_parity=mp.EVEN_Z, direction=mp.X)]
sim = mp.Simulation(cell_size=cell, boundary_layers=[mp.PML(dpml)], material_function=eps_at,
                    sources=sources, resolution=a.res)
mon_in = sim.add_mode_monitor(fcen, 0, 1, mp.ModeRegion(center=mp.Vector3(mon_in_x, 0), size=yspan))
mon_out = sim.add_mode_monitor(fcen, 0, 1, mp.ModeRegion(center=mp.Vector3(mon_out_x, 0), size=yspan))

start = time.time()
# Fixed run time after the source: 2.5x the transit to the output monitor (see bend_2d.py on the decay criterion).
transit = (mon_out_x - src_x) * N_SIN
sim.run(until_after_sources=2.5 * transit + 50)
c_in = sim.get_eigenmode_coefficients(mon_in, [1], eig_parity=mp.EVEN_Z, direction=mp.X).alpha[0, 0, 0]
c_out = sim.get_eigenmode_coefficients(mon_out, [1], eig_parity=mp.EVEN_Z, direction=mp.X).alpha[0, 0, 0]
T = abs(c_out) ** 2 / abs(c_in) ** 2
loss_db = -10 * math.log10(T)
elapsed = time.time() - start

if mp.am_master():
    n_ln_slab = slab_neff(a.t_ln, N_LN_E, N_SIO2, N_SIO2, "TE")
    print(f"RESULT L={L} res={a.res} n_sin_layer(in)={n_sin_layer(a.w_in):.4f} n_sin_layer(tip)={n_sin_layer(a.w_tip):.4f} "
          f"T={T:.6f} loss_dB={loss_db:.4f} time_s={elapsed:.0f} (LN slab n_eff={n_ln_slab:.4f})")
    new = not os.path.exists(a.csv)
    with open(a.csv, "a", newline="") as f:
        wr = csv.writer(f)
        if new:
            wr.writerow(["L_taper_um", "res_px_per_um", "t_sin_um", "w_in_um", "w_tip_um", "gap_um", "t_ln_um",
                         "T", "loss_db", "time_s"])
        wr.writerow([L, a.res, a.t_sin, a.w_in, a.w_tip, a.gap, a.t_ln, f"{T:.6f}", f"{loss_db:.4f}", f"{elapsed:.0f}"])
