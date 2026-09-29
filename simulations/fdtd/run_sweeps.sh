#!/usr/bin/env bash
# FDTD sweeps used in the docs. Usage: bash run_sweeps.sh bends|bends_thin|transition
set -euo pipefail
cd "$(dirname "$0")"
export MAMBA_ROOT_PREFIX="${MAMBA_ROOT_PREFIX:-$HOME/micromamba}"
run() { "$HOME/.local/bin/micromamba" run -n meep python "$@"; }
case "${1:-}" in
  bends)          # thick Si3N4 (800 nm): single-mode 0.7 um + multimode 1.2 um for comparison
    rm -f ../results/fdtd_bends.csv
    for R in 10 20 30 50; do run bend_2d.py --t 0.8 --w 0.7 --R "$R" --res 20; done
    run bend_2d.py --t 0.8 --w 0.7 --R 20 --res 30          # convergence
    run bend_2d.py --t 0.8 --w 1.2 --R 20 --res 20          # multimode
    run bend_2d.py --t 0.8 --w 1.2 --R 20 --res 30          # multimode, convergence
    ;;
  bends_thin)     # thin Si3N4 (200 nm, multilayer platform), single-mode at 1.2 um
    for R in 20 30 50 80 100; do run bend_2d.py --t 0.2 --w 1.2 --R "$R" --res 20 --csv ../results/fdtd_bends_thin.csv; done
    ;;
  transition)
    rm -f ../results/fdtd_transition.csv
    for L in 5 10 25 50 100; do run transition_2d.py --L "$L" --res 50; done
    run transition_2d.py --L 25 --res 70                     # convergence
    ;;
  *) echo "usage: $0 bends|bends_thin|transition"; exit 1 ;;
esac
