#!/usr/bin/env python3
"""Gera a figura de erro medido vs previsto por número de saltos (race logic, mapa 16x16).

Sem dependências externas: lê race_logic_hops_16x16.csv (saída de `go run ./cmd/racestats -csv`)
e escreve um SVG estático para o artigo.

Uso: python3 plot_race_logic_hops.py [csv] [svg]
"""
import csv
import math
import os
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
CSV = sys.argv[1] if len(sys.argv) > 1 else os.path.join(HERE, "race_logic_hops_16x16.csv")
SVG = sys.argv[2] if len(sys.argv) > 2 else os.path.join(
    HERE, "..", "..", "docs", "papers", "figuras", "race_logic_erro_por_salto.svg")

# Paleta de referência (dataviz): slots categóricos 1-3, validados par a par no modo claro.
SERIES = [(35, "#2a78d6"), (50, "#eb6834"), (75, "#1baf7a")]
SURFACE, INK, INK2, MUTED, GRID, AXIS = "#fcfcfb", "#0b0b0b", "#52514e", "#898781", "#e1e0d9", "#c3c2b7"
FONT = 'system-ui, -apple-system, "Segoe UI", sans-serif'

W, H = 760, 470
LEFT, RIGHT, TOP, BOTTOM = 84, 96, 92, 62
X_MIN, X_MAX = 0, 33
Y_MIN, Y_MAX = -8, -1  # log10 da taxa de erro


def x_px(h):
    return LEFT + (h - X_MIN) / (X_MAX - X_MIN) * (W - LEFT - RIGHT)


def y_px(rate):
    v = min(max(math.log10(rate), Y_MIN), Y_MAX)
    return TOP + (Y_MAX - v) / (Y_MAX - Y_MIN) * (H - TOP - BOTTOM)


def load(path):
    data = {}
    with open(path) as f:
        for row in csv.DictReader(f):
            unit = int(float(row["unit_ps"]))
            data.setdefault(unit, []).append((int(row["hops"]), int(row["samples"]), int(row["errors"]),
                                              float(row["measured_rate"]), float(row["predicted_rate"])))
    return data


def main():
    data = load(CSV)
    out = []
    add = out.append
    add(f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {W} {H}" width="{W}" height="{H}" '
        f'font-family=\'{FONT}\' role="img" aria-labelledby="t d">')
    add('<title id="t">Erro de decodificação por número de saltos</title>')
    add('<desc id="d">Race logic fotônica, mapa 16x16, 10 chips x 10^4 consultas por unidade de atraso. '
        'Linhas: previsão gaussiana (ruído acumulado sigma vezes raiz de h). Pontos: taxa medida. '
        'Com 100 ps não houve erros em 2,55e7 distâncias.</desc>')
    add(f'<rect width="{W}" height="{H}" fill="{SURFACE}"/>')

    # Título e legenda (linha única acima do gráfico)
    add(f'<text x="{LEFT}" y="30" font-size="16" font-weight="600" fill="{INK}">'
        'Erro de decodificação por número de saltos</text>')
    add(f'<text x="{LEFT}" y="50" font-size="12" fill="{INK2}">'
        'Race logic fotônica, mapa 16×16, 10 chips × 10⁴ consultas por unidade de atraso</text>')
    lx = LEFT
    for unit, color in SERIES:
        add(f'<rect x="{lx}" y="64" width="12" height="12" rx="3" fill="{color}"/>')
        add(f'<text x="{lx + 18}" y="74" font-size="12" fill="{INK2}">{unit} ps</text>')
        lx += 70
    lx += 16
    add(f'<line x1="{lx}" y1="70" x2="{lx + 22}" y2="70" stroke="{MUTED}" stroke-width="2" stroke-linecap="round"/>')
    add(f'<text x="{lx + 28}" y="74" font-size="12" fill="{INK2}">previsto (σ√h)</text>')
    lx += 128
    add(f'<circle cx="{lx + 6}" cy="70" r="4" fill="{MUTED}" stroke="{SURFACE}" stroke-width="2"/>')
    add(f'<text x="{lx + 16}" y="74" font-size="12" fill="{INK2}">medido</text>')

    # Grade e eixos (hairline sólida, recessiva)
    for e in range(Y_MIN, Y_MAX + 1):
        y = y_px(10 ** e)
        add(f'<line x1="{LEFT}" y1="{y:.1f}" x2="{W - RIGHT}" y2="{y:.1f}" stroke="{GRID}" stroke-width="1"/>')
        add(f'<text x="{LEFT - 10}" y="{y + 4:.1f}" font-size="11" fill="{MUTED}" text-anchor="end" '
            f'style="font-variant-numeric: tabular-nums">10<tspan dy="-5" font-size="8">{e}</tspan></text>')
    base = H - BOTTOM
    add(f'<line x1="{LEFT}" y1="{base}" x2="{W - RIGHT}" y2="{base}" stroke="{AXIS}" stroke-width="1"/>')
    for h in range(0, 33, 4):
        x = x_px(h)
        add(f'<text x="{x:.1f}" y="{base + 18}" font-size="11" fill="{MUTED}" text-anchor="middle" '
            f'style="font-variant-numeric: tabular-nums">{h}</text>')
    add(f'<text x="{(LEFT + W - RIGHT) / 2:.1f}" y="{H - 16}" font-size="12" fill="{INK2}" '
        'text-anchor="middle">Saltos no menor caminho</text>')
    add(f'<text transform="translate(22 {(TOP + base) / 2:.1f}) rotate(-90)" font-size="12" fill="{INK2}" '
        'text-anchor="middle">Taxa de erro por distância</text>')

    # Séries: linha prevista (2 px) + pontos medidos (r=4 com anel de superfície de 2 px)
    for unit, color in SERIES:
        rows = sorted(data.get(unit, []))
        pts = [(x_px(h), y_px(p)) for h, _, _, _, p in rows if p >= 10 ** Y_MIN]
        if pts:
            path = " ".join(f"{x:.1f},{y:.1f}" for x, y in pts)
            add(f'<polyline points="{path}" fill="none" stroke="{color}" stroke-width="2" '
                'stroke-linejoin="round" stroke-linecap="round"/>')
        for h, samples, errors, measured, predicted in rows:
            if errors > 0:
                add(f'<circle cx="{x_px(h):.1f}" cy="{y_px(measured):.1f}" r="4" fill="{color}" '
                    f'stroke="{SURFACE}" stroke-width="2"><title>{unit} ps, {h} saltos: {errors} erros em '
                    f'{samples} distâncias ({measured:.2e}); previsto {predicted:.2e}'
                    '</title></circle>')
        # Rótulo direto no fim da linha prevista (texto em tinta secundária, marca colorida ao lado)
        if pts:
            ex, ey = pts[-1]
            add(f'<line x1="{ex + 6:.1f}" y1="{ey:.1f}" x2="{ex + 16:.1f}" y2="{ey:.1f}" stroke="{color}" '
                'stroke-width="2" stroke-linecap="round"/>')
            add(f'<text x="{ex + 20:.1f}" y="{ey + 4:.1f}" font-size="12" fill="{INK2}">{unit} ps</text>')

    add(f'<text x="{LEFT + 10}" y="{TOP + 22}" font-size="11" fill="{MUTED}">'
        '100 ps: 0 erros em 2,55×10⁷ distâncias (fora da escala)</text>')
    add('</svg>')

    os.makedirs(os.path.dirname(SVG), exist_ok=True)
    with open(SVG, "w") as f:
        f.write("\n".join(out) + "\n")
    print(f"figura escrita em {os.path.relpath(SVG)}")


if __name__ == "__main__":
    main()
