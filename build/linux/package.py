"""Package the native Wails Linux build as a deb; no third-party Python modules."""

import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile


def main():
    project = Path(__file__).resolve().parents[2]
    config = json.loads((project / "wails.json").read_text(encoding="utf-8"))
    name = config["name"]
    info = config["info"]
    version = info["productVersion"]
    output = project / "build" / "bin"
    binary = output / config["outputfilename"]
    if not binary.is_file():
        raise FileNotFoundError(f"Build the Linux executable first: {binary}")

    # Inspect the actual binary so the package requires the WebKit ABI it uses.
    environment = dict(os.environ, LC_ALL="C")
    headers = subprocess.check_output(
        ["readelf", "-h", str(binary)], text=True, env=environment
    )
    if "Advanced Micro Devices X86-64" in headers:
        arch = "amd64"
    elif "AArch64" in headers:
        arch = "arm64"
    else:
        raise ValueError("Only Linux amd64 and arm64 packages are supported")
    dynamic = subprocess.check_output(
        ["readelf", "-d", str(binary)], text=True, env=environment
    )
    if "libwebkit2gtk-4.1.so" in dynamic:
        webkit = "libwebkit2gtk-4.1-0"
    elif "libwebkit2gtk-4.0.so" in dynamic:
        webkit = "libwebkit2gtk-4.0-37"
    else:
        raise ValueError("The executable does not link a supported WebKit2GTK ABI")

    with tempfile.TemporaryDirectory(prefix="deb-", dir=output) as temporary:
        root = Path(temporary)
        executable = root / "usr" / "bin" / name
        executable.parent.mkdir(parents=True)
        shutil.copyfile(binary, executable)
        executable.chmod(0o755)

        # Use the actual PNG dimensions for the hicolor icon theme directory.
        source_icon = project / "build/appicon.png"
        png = source_icon.read_bytes()
        if png[:8] != b"\x89PNG\r\n\x1a\n":
            raise ValueError("build/appicon.png must be a PNG image")
        width = int.from_bytes(png[16:20], "big")
        height = int.from_bytes(png[20:24], "big")
        icon = root / f"usr/share/icons/hicolor/{width}x{height}/apps/{name}.png"
        icon.parent.mkdir(parents=True)
        icon.write_bytes(png)
        icon.chmod(0o644)

        desktop = root / "usr/share/applications" / f"{name}.desktop"
        desktop.parent.mkdir(parents=True)
        desktop.write_text(
            "[Desktop Entry]\nType=Application\n"
            f"Name={info['productName']}\nExec=/usr/bin/{name}\nIcon={name}\n"
            f"StartupWMClass={name}\nTerminal=false\nCategories=Office;\n",
            encoding="utf-8",
        )
        desktop.chmod(0o644)

        control = root / "DEBIAN/control"
        control.parent.mkdir()
        control.write_text(
            f"Package: {name}\nVersion: {version}\nArchitecture: {arch}\n"
            "Section: utils\nPriority: optional\n"
            "Maintainer: fdklgbh <fdklgbh@users.noreply.github.com>\n"
            f"Depends: libc6, libgtk-3-0t64 | libgtk-3-0, {webkit}\n"
            f"Description: {info['productName']}\n"
            " Desktop administration application built with Wails.\n",
            encoding="utf-8",
        )
        control.chmod(0o644)
        # TemporaryDirectory defaults to 0700; package directories must be traversable.
        for directory, _, _ in os.walk(root):
            Path(directory).chmod(0o755)
        package = output / f"{name}_{version}_{arch}.deb"
        subprocess.run(
            ["dpkg-deb", "--root-owner-group", "--build", str(root), str(package)],
            check=True,
        )
        print(f"Created {package}")


if __name__ == "__main__":
    main()
