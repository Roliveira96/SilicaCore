"""Effective index method (EIM) with an analytic three-layer slab waveguide solver.

Reduces the real cross-section of a channel waveguide (rectangular core) to a 2D profile so the FDTD
simulations run in 2D on a modest machine. Bends use the lateral EIM; the Si3N4 -> TFLN transition
applies the EIM laterally for each taper width.
"""
import math

N_SIO2 = 1.444   # thermal silica at 1550 nm
N_SIN = 1.996    # LPCVD Si3N4 at 1550 nm
N_LN_E = 2.138   # lithium niobate, extraordinary index at 1550 nm (TE in X-cut)
WL = 1.55        # um


def slab_neff(t, n_core, n_sub, n_clad, pol="TE", wl=WL, order=0):
    """Effective index of mode `order` of an asymmetric slab waveguide (core thickness t, in um).

    pol="TE": electric field parallel to the layers; pol="TM": perpendicular.
    Solves the dispersion equation by bisection. Returns None if the mode is cut off.
    """
    k0 = 2 * math.pi / wl
    lo, hi = max(n_sub, n_clad) + 1e-9, n_core - 1e-9

    def f(neff):
        kx = k0 * math.sqrt(n_core ** 2 - neff ** 2)
        gs = k0 * math.sqrt(neff ** 2 - n_sub ** 2)
        gc = k0 * math.sqrt(neff ** 2 - n_clad ** 2)
        if pol == "TM":
            gs *= (n_core / n_sub) ** 2
            gc *= (n_core / n_clad) ** 2
        return kx * t - math.atan(gs / kx) - math.atan(gc / kx) - order * math.pi

    if f(lo) * f(hi) > 0:
        return None
    for _ in range(200):
        mid = (lo + hi) / 2
        if f(lo) * f(mid) <= 0:
            hi = mid
        else:
            lo = mid
    return (lo + hi) / 2


def channel_neff(t, w, n_core=N_SIN, n_clad=N_SIO2, wl=WL):
    """EIM for the quasi-TE mode of a t x w (um) channel waveguide fully surrounded by n_clad.

    Step 1: vertical slab (TE) -> n_slab. Step 2: lateral slab (TM) with core n_slab.
    Returns (n_slab, n_eff_channel).
    """
    n_slab = slab_neff(t, n_core, n_clad, n_clad, "TE", wl)
    n_eff = slab_neff(w, n_slab, n_clad, n_clad, "TM", wl)
    return n_slab, n_eff


def group_index(t, w, dwl=0.005, **kw):
    """Group index by finite difference: n_g = n - lambda dn/dlambda (material dispersion ignored)."""
    n1 = channel_neff(t, w, wl=WL - dwl, **kw)[1]
    n2 = channel_neff(t, w, wl=WL + dwl, **kw)[1]
    n0 = channel_neff(t, w, wl=WL, **kw)[1]
    return n0 - WL * (n2 - n1) / (2 * dwl)


if __name__ == "__main__":
    for name, t, w in [("Thick Si3N4", 0.8, 1.2), ("Thin Si3N4", 0.2, 1.2)]:
        n_slab, n_eff = channel_neff(t, w)
        second = slab_neff(w, n_slab, N_SIO2, N_SIO2, "TM", order=1)
        print(f"{name} {t*1000:.0f} nm x {w:.1f} um: n_slab={n_slab:.4f} n_eff={n_eff:.4f} "
              f"n_g~{group_index(t, w):.3f} 2nd lateral mode: {'cut off' if second is None else f'{second:.4f}'}")
