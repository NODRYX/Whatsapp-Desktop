"""Regenerate app icons from logo.avif (single source of truth).

The 740x740 AVIF is downscaled first (never used at full size), then packed:
  - logo.png  : 256x256 PNG for Windows toast notifications
  - icon.ico  : multi-size ICO (16/32/48/64/128/256, PNG-compressed entries)
                referenced by resource.rc and embedded via rsrc.syso

Requires: ffmpeg on PATH (decodes AVIF).

Usage:  python gen_icon.py
Then:   rsrc -manifest app.manifest -ico icon.ico -o rsrc.syso
"""

import struct
import subprocess
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
SOURCE = HERE / "logo.avif"
# Downscaled sizes: 256 for toast + full ICO set. Source is 740x740.
SIZES = [16, 32, 48, 64, 128, 256]


def run(cmd: list[str]) -> None:
    r = subprocess.run(cmd, capture_output=True, text=True)
    if r.returncode != 0:
        print("FAILED:", " ".join(cmd), file=sys.stderr)
        print(r.stderr[-2000:], file=sys.stderr)
        sys.exit(1)


def main() -> None:
    if not SOURCE.exists():
        print(f"missing source: {SOURCE}", file=sys.stderr)
        sys.exit(1)

    pngs: dict[int, bytes] = {}
    for size in SIZES:
        out = HERE / f".icon-{size}.png"
        # Downscale first with high-quality lanczos; AVIF input is 740x740.
        run([
            "ffmpeg", "-hide_banner", "-loglevel", "error", "-y",
            "-i", str(SOURCE),
            "-vf", f"scale={size}:{size}:flags=lanczos",
            "-frames:v", "1",
            str(out),
        ])
        pngs[size] = out.read_bytes()

    # logo.png for toast notifications (256px downscaled).
    (HERE / "logo.png").write_bytes(pngs[256])

    # icon.ico: ICONDIR + entries + PNG payloads (Vista+ supports PNG in ICO).
    count = len(SIZES)
    header = struct.pack("<HHH", 0, 1, count)
    entries = b""
    offset = 6 + 16 * count
    for size in SIZES:
        data = pngs[size]
        w = h = 0 if size == 256 else size
        entries += struct.pack("<BBBBHHII", w, h, 0, 0, 1, 32, len(data), offset)
        offset += len(data)
    with open(HERE / "icon.ico", "wb") as f:
        f.write(header + entries + b"".join(pngs[s] for s in SIZES))

    for size in SIZES:
        (HERE / f".icon-{size}.png").unlink(missing_ok=True)

    print(f"wrote logo.png + icon.ico ({count} sizes) from {SOURCE.name}")


if __name__ == "__main__":
    main()
