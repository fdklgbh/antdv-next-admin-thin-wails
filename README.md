# README

## About

This is the official Wails Vanilla template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

在目标系统上运行以下命令（需要 Go、Node.js、pnpm、Wails v2 CLI 和 Task）：

| 命令 | Windows | Ubuntu |
| --- | --- | --- |
| `task build` | 带图标的单文件 EXE | 内嵌窗口图标的 ELF 可执行文件 |
| `task build:package` | 带图标的 NSIS 安装包，同时生成 EXE | deb 安装包，同时生成 ELF |

产物位于 `build/bin/`。Windows 安装包为 `antdv-next-admin-thin-wails-amd64-installer.exe`
（ARM64 对应 `arm64`），deb 为 `antdv-next-admin-thin-wails_1.0.0_amd64.deb`
（ARM64 对应 `arm64`）。版本号在 `wails.json` 的 `info.productVersion` 中维护。
`task build:nsis` 保留为 Windows 安装包命令的别名。

Windows 打包安装包前需安装 NSIS，并将 `makensis` 加入 `PATH`。
EXE 和安装包使用 `build/windows/icon.ico`；不要添加 `-nopackage`，否则会跳过 EXE 资源打包。

Ubuntu 24.04 构建依赖：

```sh
sudo apt install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev python3 dpkg-dev binutils
task build
task build:package
sudo apt install ./build/bin/antdv-next-admin-thin-wails_1.0.0_amd64.deb
```

Ubuntu 22.04 可使用 `libwebkit2gtk-4.0-dev`。构建自动优先选择已安装的 WebKit2GTK 4.1，
否则使用 4.0；deb 根据可执行文件实际链接的 ABI 声明 WebKit 运行依赖。
请在目标 Ubuntu 版本和架构上构建，不将单个 deb 视为跨所有 Ubuntu 版本通用。

Linux 图标来自 `build/appicon.png`，内嵌于程序；deb 还会安装 `.desktop` 菜单入口和
hicolor 图标。裸 ELF 在文件管理器中通常显示系统默认可执行文件图标，无法像 Windows
EXE 一样仅靠内嵌资源设置文件图标；安装 deb 后可从应用菜单使用带图标的入口。
“单文件”表示无需外置前端资源，仍需系统的 WebView2（Windows）或 GTK/WebKit（Ubuntu）运行环境。
