"""Crosstalk between neighbouring turns of a Si3N4 delay spiral.

Two parallel 800 nm x 0.7 um Si3N4 waveguides separated by `gap` exchange power
with coupling length L_c = lambda / (2 (n_even - n_odd)). The even and odd
supermodes are found with the effective index method (eim.py): the vertical
slab gives n_slab, then a lateral five-layer TM slab (clad | core | gap | core |
clad) is solved exactly with transfer matrices, using the mirror symmetry.

In a spiral, light in one turn can couple into the next outer turn, skip it and
reach the node one turn early. Early echoes above the detector threshold would
fire a race node too soon. This script estimates, per edge:
  - per-turn transfer: sin^2(pi L_turn / (2 L_c)) for each pair of neighbouring turns;
  - total early-echo power: sum over the turns (incoherent, turn lengths differ).

Usage: python3 spiral_crosstalk.py [--length-mm 222] [--csv ../results/spiral_crosstalk.csv]
"""
import argparse
import csv
import math

from eim import N_SIN, N_SIO2, WL, channel_neff

THICKNESS = 0.8  # um
WIDTH = 0.7  # um


def _layer(h, g, n, neff, d, k0):
    """Propagate (H, G = H'/n^2) of a TM field through a uniform layer of thickness d."""
    if n > neff:
        k = k0 * math.sqrt(n * n - neff * neff)
        c, s = math.cos(k * d), math.sin(k * d)
        return c * h + (n * n / k) * s * g, -(k / (n * n)) * s * h + c * g
    q = k0 * math.sqrt(neff * neff - n * n)
    c, s = math.cosh(q * d), math.sinh(q * d)
    return c * h + (n * n / q) * s * g, (q / (n * n)) * s * h + c * g


def supermode(n_core, n_clad, width, gap, parity, wl=WL):
    """Effective index of the even or odd lateral TM supermode of two identical cores."""
    k0 = 2 * math.pi / wl

    def mismatch(neff):
        h, g = (1.0, 0.0) if parity == "even" else (0.0, 1.0)
        h, g = _layer(h, g, n_clad, neff, gap / 2, k0)
        h, g = _layer(h, g, n_core, neff, width, k0)
        q = k0 * math.sqrt(neff * neff - n_clad * n_clad)
        norm = math.hypot(h, g) or 1.0
        return (g + (q / (n_clad * n_clad)) * h) / norm  # zero for a field decaying outside

    lo, hi = n_clad + 1e-12, n_core - 1e-12
    grid = [hi - (hi - lo) * i / 4000 for i in range(4001)]
    prev_n, prev_f = grid[0], mismatch(grid[0])
    for n in grid[1:]:
        f = mismatch(n)
        if prev_f * f <= 0:  # first root from the top = fundamental supermode
            a, b, fa = n, prev_n, f
            for _ in range(200):
                mid = (a + b) / 2
                fm = mismatch(mid)
                if fa * fm <= 0:
                    b = mid
                else:
                    a, fa = mid, fm
            return (a + b) / 2
        prev_n, prev_f = n, f
    return None


def coupling_length_mm(gap, wl=WL):
    n_slab, _ = channel_neff(THICKNESS, WIDTH, wl=wl)
    n_even = supermode(n_slab, N_SIO2, WIDTH, gap, "even", wl)
    n_odd = supermode(n_slab, N_SIO2, WIDTH, gap, "odd", wl)
    return wl / (2 * (n_even - n_odd)) / 1000, n_even, n_odd


def spiral_echo(length_mm, pitch_um, lc_mm):
    """Per-turn and total early-echo power for a square spiral of the given length and pitch."""
    area_mm2 = length_mm * pitch_um / 1000
    side_mm = math.sqrt(area_mm2)
    turns = max(1, int(side_mm / (2 * pitch_um / 1000)))
    worst_turn = 0.0
    total = 0.0
    for i in range(turns):
        turn_len = 4 * (side_mm - 2 * i * pitch_um / 1000)
        if turn_len <= 0:
            break
        p = math.sin(math.pi * turn_len / (2 * lc_mm)) ** 2
        worst_turn = max(worst_turn, p)
        total += p
    return turns, side_mm, worst_turn, min(total, 1.0)


def db(x):
    return 10 * math.log10(x) if x > 0 else float("-inf")


def main():
    parser = argparse.ArgumentParser(description="Crosstalk between neighbouring turns of a Si3N4 delay spiral")
    parser.add_argument("--length-mm", type=float, default=222.0, help="spiral length of the longest edge")
    parser.add_argument("--csv", default="", help="optional CSV output")
    args = parser.parse_args()

    n_slab, n_single = channel_neff(THICKNESS, WIDTH)
    print(f"Si3N4 {THICKNESS * 1000:.0f} nm x {WIDTH} um, lambda {WL} um: n_slab {n_slab:.4f}, single-guide n_eff {n_single:.4f}")
    print(f"Longest edge: {args.length_mm:.0f} mm of spiral\n")
    header = ["pitch_um", "gap_um", "Lc_mm", "turns", "side_mm", "worst_turn_dB", "total_early_echo_dB"]
    print(f"{'pitch':>7} {'gap':>6} {'L_c':>12} {'turns':>6} {'side':>8} {'worst turn':>11} {'early echo':>11}")
    rows = []
    for pitch in (2.0, 2.5, 3.0, 3.5, 4.0, 5.0, 6.0, 8.0):
        gap = pitch - WIDTH
        lc, _, _ = coupling_length_mm(gap)
        turns, side, worst, total = spiral_echo(args.length_mm, pitch, lc)
        rows.append([pitch, round(gap, 2), round(lc, 3), turns, round(side, 3), round(db(worst), 1), round(db(total), 1)])
        print(f"{pitch:5.1f}um {gap:4.1f}um {lc:9.2f} mm {turns:6d} {side:6.2f}mm {db(worst):8.1f} dB {db(total):8.1f} dB")
    if args.csv:
        with open(args.csv, "w", newline="") as f:
            w = csv.writer(f)
            w.writerow(header)
            w.writerows(rows)
    print("\nEarly echo = power that skips turns and reaches the node before the main pulse (incoherent sum).")
    print("It must stay well below the comparator threshold (about -3 dB of the received pulse).")


if __name__ == "__main__":
    main()
