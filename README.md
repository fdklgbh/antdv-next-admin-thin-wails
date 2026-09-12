# Antdv Next Admin Thin · Wails v2

## About

基于 Wails v2、Vue 3 和 Antdv Next 的桌面管理后台，当前 Wails 依赖为 v2.15.0。

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## 拉取仓库与前端子模块

v2 与 v3 使用同一个后端仓库，v2 位于 `wailsv2` 分支，`frontend/` 是独立子模块，
对应前端分支 `no-auth-wails-v2`。首次克隆需同时拉取子模块；以下 SSH 方式要求已配置 GitHub SSH 密钥：

```sh
git clone --branch wailsv2 --recurse-submodules git@github.com:fdklgbh/antdv-next-admin-thin-wails.git antdv-next-admin-thin-wails-v2
cd antdv-next-admin-thin-wails-v2
```

没有配置 SSH 时可用 HTTPS。由于 `.gitmodules` 保存的是 SSH 地址，子模块命令需临时转换地址；
以下配置仅对该次命令生效，不修改全局 Git 设置：

```sh
git clone --branch wailsv2 https://github.com/fdklgbh/antdv-next-admin-thin-wails.git antdv-next-admin-thin-wails-v2
cd antdv-next-admin-thin-wails-v2
git -c 'url.https://github.com/.insteadOf=git@github.com:' submodule update --init --recursive
```

如果已经克隆但 `frontend/` 为空，在项目根目录补拉：

```sh
git submodule sync --recursive
git submodule update --init --recursive
```

后续更新前先检查并保存主仓库和前端子模块中的本地修改，再执行：

```sh
git switch wailsv2
git pull --ff-only origin wailsv2
git submodule sync --recursive
git submodule update --init --recursive
git submodule status --recursive
```

HTTPS 用户将上述 `submodule update` 命令替换为带 `-c` 的版本即可。
普通子模块更新检出的是后端提交锁定的前端提交，不是前端分支的最新提交；
子模块处于 detached HEAD 是正常现象。日常获取项目不要使用 `git submodule update --remote`，
以免前后端版本不匹配。v3 应另外克隆 `master` 分支，放到不同目录。

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
当前原生 PNG 和 Windows ICO 均使用与 `frontend/public/logo.png` 相同的 AN 标志。
更换品牌图标时应同步这三处；只替换页面 Logo 不会更新任务栏、窗口或安装包图标。

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
hicolor 图标和 `/usr/share/pixmaps/` 图标，并在安装后刷新图标缓存和桌面数据库。
更新此打包逻辑后需重新生成并安装 deb；已有安装不会自动获得新增图标文件。
裸 ELF 在文件管理器中通常显示系统默认可执行文件图标，无法像 Windows
EXE 一样仅靠内嵌资源设置文件图标；安装 deb 后可从应用菜单使用带图标的入口。
“单文件”表示无需外置前端资源，仍需系统的 WebView2（Windows）或 GTK/WebKit（Ubuntu）运行环境。

## 替换图标并重新打包（Windows / Ubuntu）

1. 将页面 Logo 替换为 `frontend/public/logo.png`，将相同品牌的方形 PNG 放到
   `build/appicon.png`，推荐 256×256。只替换前端 Logo 不会改变任务栏、窗口或安装包图标。
2. Windows 还需更新 `build/windows/icon.ico`。使用图标转换工具生成包含
   16、32、48、64、128、256 像素尺寸的 ICO；不能只把 PNG 扩展名改为 `.ico`。
   也可以先备份并移走旧 ICO，再在 Windows 执行 `task build`，让 Wails v2 从 PNG 自动生成。
   已存在的 ICO 不会因为 PNG 改变而自动更新。
3. Linux 打包脚本会读取 PNG 实际尺寸，自动设置 hicolor 目录，并将同一图片安装到 pixmaps；
   无需手动填写图片尺寸。重新生成 deb 后才能更新已安装图标。

Windows 项目根目录执行：

```sh
task build
# 需要 NSIS 安装包时执行
task build:package
```

Ubuntu 项目根目录执行：

```sh
task build:package
sudo apt install --reinstall ./build/bin/antdv-next-admin-thin-wails-v2_1.0.0-1_amd64.deb
```

包名由 `wails.json` 的 `name`、`info.productVersion` 和 `build/linux/package.json` 的
`release` 以及实际架构决定；上面是当前默认示例，应以本次生成文件为准。
跨机器构建时先同步源码和图标，再在 Ubuntu 打包，不能重复安装同步前生成的旧 deb。

若图标仍旧，可在 Ubuntu 对照检查（将 deb 文件名替换为实际名称）：

```sh
sha256sum build/appicon.png
dpkg-deb -c build/bin/antdv-next-admin-thin-wails-v2_1.0.0-1_amd64.deb | grep -E 'png|desktop'
dpkg-deb --fsys-tarfile build/bin/antdv-next-admin-thin-wails-v2_1.0.0-1_amd64.deb \
  | tar -xOf - ./usr/share/pixmaps/antdv-next-admin-thin-wails-v2.png | sha256sum
sha256sum /usr/share/pixmaps/antdv-next-admin-thin-wails-v2.png
```

三处 PNG 的 SHA256 应一致。源文件与包内不一致时重新打包；包内与安装路径不一致时重新安装。
更改 `name` 后也要同步替换上述安装路径。安装完成后完全退出旧程序，再从应用菜单启动。
文件一致但固定图标仍旧时，取消固定后重新固定；Ubuntu 必要时注销后重新登录。
Windows 需运行新 EXE 或重新安装新安装包，旧进程不会自动更新图标。
