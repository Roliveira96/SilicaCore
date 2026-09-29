#!/usr/bin/env python3
"""Bottom-up budget of the photonic tensor unit (PTU) floor.

Broadcast-and-weight architecture: N inputs on N wavelengths, TFLN intensity
modulators, 1:M broadcast, one Sb2Se3-trimmed microring per weight and one
balanced photodetector + TIA + ADC per output column.

For each target precision (effective bits) the script finds the full-scale
photocurrent that keeps the analog noise (TIA thermal, shot and laser RIN)
at the quantization-noise level, then derives optical power, laser power,
converter power, total power, TOPS/W, energy per MAC and area.

Usage: python3 ptu_budget.py [--rin -150] [--size 64] [--rate 10e9]
"""

import argparse
import math

ELECTRON_CHARGE = 1.602e-19

# Optical loss from laser to photodetector, dB (assumptions, see doc 15).
LOSS_DB = {
    "laser-to-chip coupling": 1.5,
    "TFLN modulator": 2.0,
    "wavelength multiplexer": 1.5,
    "1:64 splitter excess": 1.5,
    "propagation": 1.0,
    "weight ring": 1.0,
    "photodetector coupling": 0.5,
}


def required_full_scale(bits, sigma_thermal, bandwidth, rin_linear):
    """Smallest full-scale current whose total noise equals the quantization noise."""
    target = 2.0 / (2**bits * math.sqrt(12.0))
    if math.sqrt(rin_linear * bandwidth) >= target:
        return None  # RIN alone already exceeds the budget
    current = 1e-7
    while current < 10.0:
        noise = math.sqrt(
            sigma_thermal**2
            + 2 * ELECTRON_CHARGE * current * bandwidth
            + rin_linear * bandwidth * current**2
        )
        if noise <= target * current:
            return current
        current *= 1.005
    return None


def main():
    parser = argparse.ArgumentParser(description="PTU power, noise and area budget")
    parser.add_argument("--size", type=int, default=64, help="array size N = M")
    parser.add_argument("--rate", type=float, default=10e9, help="vector rate in samples/s")
    parser.add_argument("--rin", type=float, default=-150.0, help="laser RIN in dB/Hz")
    parser.add_argument("--wpe", type=float, default=0.20, help="laser wall-plug efficiency")
    parser.add_argument("--tia-noise", type=float, default=20e-12, help="TIA input noise, A/sqrt(Hz)")
    parser.add_argument("--dac-energy", type=float, default=1.0e-12, help="DAC + driver energy per sample, J")
    parser.add_argument("--adc-fom", type=float, default=30e-15, help="ADC Walden figure of merit, J/step")
    parser.add_argument("--tia-power", type=float, default=15e-3, help="power per TIA, W")
    parser.add_argument("--digital-power", type=float, default=0.3, help="control and buffers, W")
    args = parser.parse_args()

    n = m = args.size
    bandwidth = args.rate / 2
    responsivity = 1.0
    loss_db = sum(LOSS_DB.values())
    eta = 10 ** (-loss_db / 10)
    sigma_thermal = args.tia_noise * math.sqrt(bandwidth)
    rin_linear = 10 ** (args.rin / 10)
    macs_per_s = n * m * args.rate
    tops = 2 * macs_per_s / 1e12
    ring_control = (0.1e-3 * n * m, 1.0e-3 * n * m)

    print(f"PTU {n}x{m} at {args.rate / 1e9:.0f} GS/s: {tops:.1f} TOPS")
    print(f"optical loss {loss_db:.1f} dB (eta {eta:.3f}), TIA noise {sigma_thermal * 1e6:.2f} uA rms, RIN {args.rin:.0f} dB/Hz")
    print()
    header = "bits,full_scale_uA,power_per_wavelength_mW,laser_W,dac_W,adc_W,tia_W,ring_control_W,total_W,TOPS_per_W,fJ_per_MAC"
    print(header)
    for bits in (4, 5, 6, 7, 8):
        full_scale = required_full_scale(bits, sigma_thermal, bandwidth, rin_linear)
        if full_scale is None:
            print(f"{bits},RIN limited,,,,,,,,,")
            continue
        per_wavelength = full_scale / responsivity / eta
        laser = n * per_wavelength / args.wpe
        dac = n * args.rate * args.dac_energy
        adc = m * args.rate * args.adc_fom * 2**bits
        tia = m * args.tia_power
        fixed = laser + dac + adc + tia + args.digital_power
        low, high = fixed + ring_control[0], fixed + ring_control[1]
        print(
            f"{bits},{full_scale * 1e6:.0f},{per_wavelength * 1e3:.2f},{laser:.2f},{dac:.2f},{adc:.2f},{tia:.2f},"
            f"{ring_control[0]:.2f}-{ring_control[1]:.2f},{low:.2f}-{high:.2f},"
            f"{tops / high:.1f}-{tops / low:.1f},{low / macs_per_s * 1e15:.0f}-{high / macs_per_s * 1e15:.0f}"
        )

    ring_cell_mm2 = (80e-3) ** 2
    modulator_mm2 = 0.5
    misc_mm2 = 10.0
    area = n * m * ring_cell_mm2 + n * modulator_mm2 + misc_mm2
    print()
    print(f"area per PTU {area:.1f} mm2 (rings {n * m * ring_cell_mm2:.1f}, modulators {n * modulator_mm2:.0f}, other {misc_mm2:.0f}), {tops / area:.2f} TOPS/mm2")

    group_index = 2.06
    ring_radius = 30e-6
    fsr_nm = (1.55e-6) ** 2 / (group_index * 2 * math.pi * ring_radius) * 1e9
    half_width = 0.16 / 2
    span = 7 * 0.8
    next_order = 10 * math.log10(half_width**2 / ((fsr_nm - span) ** 2 + half_width**2))
    neighbor = 10 * math.log10(half_width**2 / (0.8**2 + half_width**2))
    print(f"ring FSR (R = 30 um) {fsr_nm:.2f} nm for an 8-channel span of {span:.1f} nm")
    print(f"crosstalk: adjacent channel {neighbor:.1f} dB, next ring order {next_order:.1f} dB")

    io_bits = 6
    io_tbps = 2 * n * args.rate * io_bits / 1e12
    print(f"electrical I/O per PTU at {io_bits} bits: {io_tbps:.2f} Tb/s (in + out)")


if __name__ == "__main__":
    main()
