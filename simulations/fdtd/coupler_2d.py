"""2D FDTD check of the supermode coupling length used in spiral_crosstalk.py.

Two parallel Si3N4 waveguides (effective index method, same reduction as the
bend simulations) with a given gap. Light is launched in guide 1 before guide 2
starts; a the steady-state complex field gives and the power fraction in
guide 2 along the propagation. Its first maximum sits at the coupling length L_c,
which is compared with the analytic transfer-matrix value.

Usage: micromamba run -n meep python coupler_2d.py --gap 0.4 --res 20
"""
import argparse
import csv
import os
import time

import meep as mp
import numpy as np

from eim import N_SIO2, WL, channel_neff
from spiral_crosstalk import THICKNESS, WIDTH, coupling_length_mm

p = argparse.ArgumentParser()
p.add_argument("--gap", type=float, default=0.4, help="edge-to-edge gap (um)")
p.add_argument("--res", type=int, default=20, help="pixels per um")
p.add_argument("--csv", default=os.path.join(os.path.dirname(__file__), "..", "results", "fdtd_coupler.csv"))
a = p.parse_args()

n_slab, _ = channel_neff(THICKNESS, WIDTH)
lc_analytic_um = coupling_length_mm(a.gap)[0] * 1000
core, clad = mp.Medium(index=n_slab), mp.Medium(index=N_SIO2)

dpml, lead = 1.5, 6.0
length = 2.6 * lc_analytic_um
sx = 2 * dpml + lead + length
sy = 2 * dpml + 2 * WIDTH + a.gap + 4.0
y1 = (WIDTH + a.gap) / 2  # guide 1 centre
x_start = -sx / 2 + dpml + lead  # guide 2 starts here

geometry = [
    mp.Block(size=mp.Vector3(mp.inf, WIDTH), center=mp.Vector3(0, y1), material=core),
    mp.Block(size=mp.Vector3(sx, WIDTH), center=mp.Vector3(x_start + sx / 2, -y1), material=core),
]
fcen = 1 / WL
src_x = -sx / 2 + dpml + 1.0
sources = [mp.EigenModeSource(mp.ContinuousSource(fcen, width=20), center=mp.Vector3(src_x, y1),
                              size=mp.Vector3(0, WIDTH + a.gap), eig_band=1, eig_parity=mp.EVEN_Z, direction=mp.X)]
sim = mp.Simulation(cell_size=mp.Vector3(sx, sy), boundary_layers=[mp.PML(dpml)], geometry=geometry,
                    default_material=clad, sources=sources, resolution=a.res, force_complex_fields=True)

t0 = time.time()
transit = sx * n_slab
sim.run(until=3 * transit)
# Complex CW fields at steady state: |E|^2 is the time-averaged intensity.
region_center = mp.Vector3(x_start + length / 2, 0)
region_size = mp.Vector3(length, sy - 2 * dpml)
ex = sim.get_array(center=region_center, size=region_size, component=mp.Ex, cmplx=True)
ey = sim.get_array(center=region_center, size=region_size, component=mp.Ey, cmplx=True)
intensity = np.abs(ex) ** 2 + np.abs(ey) ** 2  # shape (nx, ny)
ny = intensity.shape[1]
ys = np.linspace(-(sy - 2 * dpml) / 2, (sy - 2 * dpml) / 2, ny)
upper = intensity[:, ys > 0].sum(axis=1)
lower = intensity[:, ys < 0].sum(axis=1)
frac = lower / (upper + lower)
xs = np.linspace(0, length, len(frac))
first_peak = int(np.argmax(frac[: int(len(frac) * 0.6)]))
lc_fdtd_um = xs[first_peak]
elapsed = time.time() - t0

if mp.am_master():
    print(f"RESULT gap={a.gap} res={a.res} Lc_analytic={lc_analytic_um:.2f} um Lc_fdtd={lc_fdtd_um:.2f} um "
          f"max_transfer={frac[first_peak]:.3f} time={elapsed:.0f}s")
    new = not os.path.exists(a.csv)
    with open(a.csv, "a", newline="") as f:
        w = csv.writer(f)
        if new:
            w.writerow(["gap_um", "res_px_per_um", "Lc_analytic_um", "Lc_fdtd_um", "max_transfer", "time_s"])
        w.writerow([a.gap, a.res, f"{lc_analytic_um:.2f}", f"{lc_fdtd_um:.2f}", f"{frac[first_peak]:.3f}", f"{elapsed:.0f}"])
