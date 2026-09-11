# Antdv Next Admin Thin · Wails v2

## About

基于 Wails v2、Vue 3 和 Antdv Next 的桌面管理后台，当前 Wails 依赖为 v2.15.0。

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## 复用项目：集中配置名称与包信息

通用产品信息维护在 `wails.json`，Windows EXE / NSIS 沿用 Wails v2 的配置和资源模板。
主程序内嵌该文件，窗口标题和 Linux 程序标识随配置变化，修改后必须重新构建。

| `wails.json` 字段 | 用途 |
| --- | --- |
| `name` | 项目标识、Linux 包名、desktop 文件和图标标识、NSIS 安装器文件名 |
| `outputfilename` | 构建输出文件名、Linux 安装到 `/usr/bin/` 的文件名及 desktop Exec；建议与 `name` 相同，不带 `.exe` |
| `info.productName` | 窗口标题、Linux 菜单名称和包摘要、Windows 产品名称 |
| `info.productVersion` | Linux 包上游版本、Windows EXE / NSIS 版本；建议使用 `1.2.3` 三段数字 |
| `info.companyName` | Windows 公司名称 |
| `info.copyright` | Windows 版权信息 |
| `info.comments` | Windows 备注、Linux 包详细描述和桌面入口描述 |
| `author.name` / `author.email` | Linux deb 维护者，不再固定在 Python 脚本中 |

Linux 专用参数位于 `build/linux/package.json`，不重复保存产品名称或版本：

| 字段 | 用途 |
| --- | --- |
| `homepage` | deb 主页 |
| `release` | 打包修订号，完整版本为 `info.productVersion-release` |
| `section` / `priority` | Debian 软件分类和优先级 |
| `categories` | freedesktop 应用菜单分类，例如 `Office;` |

Linux 专用字段单独保存，是因为 Wails v2 在保存 `wails.json` 时会丢弃其不识别的自定义字段。
Python 打包脚本直接读取两份配置，不需要手动导出环境变量，也不需要第三方 Python 库。
架构和 GTK/WebKit 运行依赖继续从实际 ELF 检测，不作为可随意填写的配置。

复制项目后，修改上述字段，替换 `build/appicon.png` 和 Windows 使用的 `build/windows/icon.ico`，
然后执行 `task build` 或 `task build:package`。`name` 应为合法 Debian 包名；Linux 的
`outputfilename` 使用小写字母、数字、点、下划线或连字符，不含路径分隔符。
中文和空格可用于显示名称；用于 Linux 包和桌面入口的元数据必须是非空单行文本。

Go 模块名、Go import、前端生成绑定路径以及前端自己的标题和版本仍需单独调整，打包名称不会自动修改源码标识。
前端是独立子模块；修改服务路径后使用 v2 CLI 重新生成绑定。macOS 等平台资源和许可证也需单独复核。
v3 的对应配置说明位于 `master` 分支 README，本文件对应 `wailsv2` 分支。

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

产物位于 `build/bin/`。Windows 安装包为 `<name>-amd64-installer.exe`
（ARM64 对应 `arm64`），deb 为 `<name>_<productVersion>-<release>_amd64.deb`
（ARM64 对应 `arm64`）。版本号在 `wails.json` 的 `info.productVersion` 中维护。

Windows 打包安装包前需安装 NSIS，并将 `makensis` 加入 `PATH`。
EXE 和安装包使用 `build/windows/icon.ico`；不要添加 `-nopackage`，否则会跳过 EXE 资源打包。

Ubuntu 24.04 构建依赖：

```sh
sudo apt install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev python3 dpkg-dev binutils
task build
task build:package
sudo apt install ./build/bin/antdv-next-admin-thin-wails-v2_1.0.0-1_amd64.deb
```

Ubuntu 22.04 可使用 `libwebkit2gtk-4.0-dev`。构建自动优先选择已安装的 WebKit2GTK 4.1，
否则使用 4.0；deb 根据可执行文件实际链接的 ABI 声明 WebKit 运行依赖。
请在目标 Ubuntu 版本和架构上构建，不将单个 deb 视为跨所有 Ubuntu 版本通用。

Linux 图标来自 `build/appicon.png`，内嵌于程序；deb 还会安装 `.desktop` 菜单入口和
hicolor 图标。裸 ELF 在文件管理器中通常显示系统默认可执行文件图标，无法像 Windows
EXE 一样仅靠内嵌资源设置文件图标；安装 deb 后可从应用菜单使用带图标的入口。
“单文件”表示无需外置前端资源，仍需系统的 WebView2（Windows）或 GTK/WebKit（Ubuntu）运行环境。
