# Antdv Next Admin Thin · Wails v3

基于 Wails v3、Vue 3 和 Antdv Next 的桌面管理后台脚手架。前端资源内嵌于可执行文件，使用系统 WebView 渲染界面。

当前 Go 依赖为 **Wails v3.0.0-beta.19**。本 README 对应 v3 项目，构建配置入口为 `Taskfile.yml` 和 `build/config.yml`。

## 复用项目：集中配置名称与包信息

产品信息统一维护在 `build/config.yml`。Taskfile 读取配置，Linux 打包任务通过环境变量
传给 nFPM；主程序内嵌同一配置，启动时读取窗口标题和程序标识，因此修改后必须重新构建。

Wails beta.19 内置的 Task 不支持 `mustFromYaml`，因此构建、打包和开发统一使用独立
`task`。`build/config.yml` 的三个开发子进程也调用独立 `task`，不使用 Wails 内置 Task。
根 Taskfile 声明最低版本 `3.51.1`，版本不足会在执行任务前报错并停止，无需额外检查脚本。
Ubuntu 22.04 / 24.04 使用同一规则，GTK 后端仍按下文规则选择。

| 配置字段 | 用途 |
| --- | --- |
| `packaging.appName` | 各平台构建输出基名、Linux 包名、程序标识、desktop 文件名及图标关联、Windows NSIS 安装器文件名 |
| `info.productName` | 窗口标题、Linux 应用菜单显示名称、Windows 产品名称 |
| `info.version` | Linux 包版本、Windows EXE 和 NSIS 产品版本；建议使用 `1.2.3` 三段数字 |
| `info.description` | 应用描述、Linux 包和桌面入口描述、Windows 文件描述 |
| `info.companyName` | Linux vendor、Windows 公司名称 |
| `info.copyright` / `info.comments` | Windows 版权和备注 |
| `packaging.maintainer` | Linux 维护者，格式为 `姓名 <邮箱>`，不再依赖 Git 提交者环境变量 |
| `packaging.homepage` | Linux 包主页 |
| `packaging.release` | Linux 打包修订号，例如同一应用版本的第 `1` 次打包 |
| `packaging.categories` | Linux 应用菜单分类，例如 `Development;`，以分号结尾 |
| `info.productIdentifier` | Wails 平台资源生成使用的产品标识，复用时应更换为自己的反向域名 |

复制项目后，修改上述字段并替换 `build/appicon.png`，再运行 `task build` 或对应打包命令。
`packaging.appName` 使用小写字母、数字和连字符，不使用空格或路径分隔符；中文和空格用于
`info.productName`。字段应使用单行文本，发布前替换默认公司、主页及维护者占位信息。

Windows 构建会自动生成 `bin/windows-info.json` 并用于 `.syso`，NSIS 通过命令参数读取同一份
配置；旧 `build/windows/info.json` 不再用于常规 EXE 构建。Linux 两份 nFPM 配置都读取相同的
环境变量，GTK3/GTK4 的依赖仍分别维护。请通过 Task 打包，直接运行 nFPM/Wails 打包子命令时
需要自行提供环境变量。`task --silent config:show` 输出名称和产物目录的 JSON，CI 用它定位 v3 产物。

仍需单独检查的内容：

- 两份 nFPM 配置中的 `license`、`section`、`priority`：当前 nFPM 不对这些字段展开环境变量，
  不要直接写成 `${APP_LICENSE}`。若变更许可证，还应同步项目许可证文件。
- macOS、iOS、Android、MSIX 等既有平台资源未在本次接入自动同步。修改产品标识后按对应平台
  更新资源并验证；`task common:update:build-assets` 会覆盖构建资源，不能不经检查就用它覆盖当前定制任务和包配置。
- Go 模块名、Go import、前端生成绑定路径及前端自己的标题/版本不由上述打包名称自动重命名。
  前端是独立子模块，应按其指南修改并重新生成绑定。
- GitHub Release 日期标签属于发布工作流，不会替代 `info.version`。

这里只集中产品与打包元数据，仍需在目标 Ubuntu 和 Windows 上验证安装及运行。

## 环境准备

- Go：满足 `go.mod` 声明的版本（当前为 `1.26.7`）。
- Node.js：推荐 22.12+，以及 pnpm。
- Task：3.51.1 或更新的兼容 3.x 版本，`task` 必须在 PATH 中。
- Wails v3 CLI：建议与项目依赖保持同版本。

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.19
```

确保 Go 的 bin 目录在 PATH 中，并按下文安装目标系统的原生依赖。
可用 `task --version` 检查版本；安装已验证版本使用
`go install github.com/go-task/task/v3/cmd/task@v3.51.1`。
日常使用 `task dev`、`task build`、`task package`，不要使用 `wails3 build/package/task`。
构建任务会安装前端依赖、构建前端及生成所需资源；前端默认使用 pnpm。

## 开发与常用命令

在项目根目录执行：

| 命令 | 用途 |
| --- | --- |
| `task dev` | 启动开发模式及热更新，前端端口默认 9245 |
| `task build` | 为当前系统构建生产程序 |
| `task run` | 运行已构建的程序 |
| `task package` | 为当前系统构建并生成安装包 |
| `task linux:create:deb` | 构建 Linux 程序并仅生成 deb |

修改开发端口：`task dev WAILS_VITE_PORT=9246`。
需要调试构建时使用 `task build DEV=true`；指定架构使用 `ARCH=amd64` 或 `ARCH=arm64`。
常规 Windows/Linux 产物输出到 **`bin/`**，不是 `build/bin/`。

## Windows

运行需要 Microsoft Edge WebView2 Runtime。构建 NSIS 安装包还需要安装 NSIS，并将 `makensis` 加入 PATH。

```sh
# 带图标、版本信息和 manifest 的单文件 EXE
task build

# NSIS 安装包（默认安装到计算机范围）
task package

# 仅为当前用户安装的 NSIS 包
task package INSTALL_SCOPE=user
```

主要产物：

- `bin/antdv-next-admin-thin-wails.exe`
- `bin/antdv-next-admin-thin-wails-amd64-installer.exe`（ARM64 对应 `arm64`）

安装包包含 WebView2 引导安装流程，缺少运行环境时可能需要联网。
EXE 内嵌前端资源，但“单文件”不表示同时内嵌完整的 WebView2 Runtime。

项目包含 Windows 最小化恢复后的客户区重算处理，解决最大化窗口最小化再打开后的界面比例异常；其他系统使用空实现。

## Ubuntu 22.04 / 24.04

Linux Taskfile 读取 `/etc/os-release`，自动选择后端，并同步选择安装包依赖：

| 构建系统 | 后端 | deb 配置 |
| --- | --- | --- |
| Ubuntu 22.04 | GTK3 + WebKit2GTK 4.1 | `build/linux/nfpm/nfpm-gtk3.yaml` |
| Ubuntu 24.04 | GTK4 + WebKitGTK 6.0（Wails 默认） | `build/linux/nfpm/nfpm.yaml` |

### 安装开发依赖

Ubuntu 22.04：

```sh
sudo apt update
sudo apt install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev
```

Ubuntu 24.04：

```sh
sudo apt update
sudo apt install build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev
```

### 构建与打包

在对应 Ubuntu 系统上执行相同命令：

```sh
# 单文件 ELF 程序
task build

# 构建程序并生成 deb
task linux:create:deb
```

运行单文件程序：

```sh
./bin/antdv-next-admin-thin-wails
```

deb 也输出到 `bin/`，当前 Wails 打包命令生成 `<appName>.deb`，版本和架构记录在包元数据中。
使用 `sudo apt install ./bin/<实际文件名>.deb` 安装，以便同时解析系统运行依赖。
包版本读取 `build/config.yml` 的 `info.version`，当前为 `0.1.0`。

- `task package` 会依次生成 AppImage、deb、RPM 和 Arch 包；只要 deb 时使用 `task linux:create:deb`。
- deb 由 `wails3 tool package` 生成，不需要额外 Python 打包脚本。
- 开发与生产构建使用同一后端选择规则；其他系统默认 GTK4。
- 显式覆盖时使用 `task build GTK_VERSION=3` 或 `task linux:create:deb GTK_VERSION=4`。
  `GTK_VERSION` 同时控制编译标签和包依赖，不要仅通过 `EXTRA_TAGS=gtk3` 切换。

请分别在目标 Ubuntu 系统上原生构建和验证。24.04 构建的程序不保证能在 22.04 上运行，
因为 GTK 后端和 glibc 等系统库版本可能不同。当前 Ubuntu 22.04 / 24.04 的实际构建、安装及运行仍待验证。
按当前 Wails 说明，GTK3 兼容模式将在 v3.1 移除，升级时需要重新评估 22.04 支持。

### 设置抽屉花屏：Linux WebKit 渲染方案

如果打开设置抽屉时出现画面重复、压缩或错位，可对比以下两种方案。此类现象可能与
WebKitGTK、显卡驱动或显示环境的合成渲染有关，仅凭截图不能确认根因；两种方案均需在目标 Ubuntu 上验证。

#### 方案一：仅关闭 DMA-BUF 渲染路径，保留 GPU 加速

先将 `main.go` 中 `application.LinuxWindow` 的 GPU 策略改为：

```go
WebviewGpuPolicy: application.WebviewGpuPolicyAlways,
```

重新构建，再通过环境变量启动：

```sh
task build
WEBKIT_DISABLE_DMABUF_RENDERER=1 ./bin/antdv-next-admin-thin-wails

# 若测试 deb 安装的程序，需先重新打包并安装上述配置的版本
WEBKIT_DISABLE_DMABUF_RENDERER=1 /usr/bin/antdv-next-admin-thin-wails
```

该变量只作用于本次启动，不会自动应用到应用菜单启动。它关闭 DMA-BUF 渲染路径，
不等于关闭全部硬件加速，适合优先排查部分驱动、虚拟机或 Wayland 环境的兼容问题。
实际支持及效果取决于系统 WebKitGTK 版本。当前默认策略为 `Never`，只设置环境变量不会重新启用 GPU，
因此对比前必须修改上述策略并重新构建。

#### 方案二：完全禁用 WebView GPU 加速（当前默认）

当前 `main.go` 在 `application.LinuxWindow` 中显式设置：

```go
WebviewGpuPolicy: application.WebviewGpuPolicyNever,
```

此方案使用软件渲染，规避 WebView GPU 合成路径的兼容问题，无需额外环境变量；
可能增加 CPU 占用，并影响复杂动画或图形内容的性能。该配置仅作用于 Linux 窗口。

从方案一切回时，将策略恢复为 `Never`，重新执行 `task build`；需要 deb 时执行
`task linux:create:deb` 并重新安装。两种方案都应验证设置抽屉的打开、关闭和滚动，
以及窗口缩放后的显示效果；当前尚未完成目标 Ubuntu 环境的实际验证。

## 图标与包信息

- `build/appicon.png`：原生 AN 图标，来自 `frontend/public/logo.png`，当前为 256×256。
- `build/windows/icon.ico`：由图标任务生成，用于 Windows EXE 和 NSIS 安装包。
- 更换品牌图标时同步页面 Logo 和原生 PNG，再生成 ICO；两份 nFPM 的 hicolor 目录尺寸应与 PNG 一致。
- `build/linux/antdv-next-admin-thin-wails.desktop`：Linux 构建生成的应用菜单入口。
- 两份 nFPM 配置将程序安装到 `/usr/bin/`，将桌面入口和图标安装到系统应用目录及 hicolor 图标主题目录。

Linux 裸 ELF 在文件管理器中的图标由桌面环境决定，不能像 Windows EXE 一样保证显示自定义文件图标；
deb 安装后通过应用菜单入口关联图标。

### 替换图标并重新打包（Windows / Ubuntu）

1. 准备同一套品牌图片：将页面 Logo 放到 `frontend/public/logo.png`，将方形 PNG 放到
   `build/appicon.png`。原生 PNG 推荐 256×256；只替换前端 Logo 不会改变系统图标。
2. 若原生 PNG 尺寸变化，同步修改 `build/linux/nfpm/nfpm.yaml` 和 `nfpm-gtk3.yaml` 中
   `hicolor/<宽>x<高>/apps/` 的目录，保持与图片实际尺寸一致。保留 pixmaps 安装项和安装后缓存刷新。
3. 生成派生图标，并在目标系统重新构建。不要仅重命名 PNG 的扩展名来制作 ICO。

```sh
# 强制由 build/appicon.png 重新生成 ICO / ICNS
task --force common:generate:icons

# Windows：生成带新图标的 EXE；需要安装包再执行第二条
task --force build
task --force package
```

在 Ubuntu 项目根目录重新生成并安装 deb：

```sh
task --force linux:create:deb
sudo apt install --reinstall ./bin/antdv-next-admin-thin-wails.deb
```

以上文件名按默认 `packaging.appName` 编写，修改名称后使用自己的实际文件名。
跨机器打包时，必须同步 PNG、两份 nFPM 配置及相关源码后，在 Ubuntu 重新执行打包命令。
`sync-to-ubuntu.bat` 同步的是源码，不会自动重新构建远端 `bin/` 中的旧 deb。

若安装后仍显示旧图标，先核对源文件、包内文件和已安装文件，不要直接认定是缓存：

```sh
sha256sum build/appicon.png
dpkg-deb -c bin/antdv-next-admin-thin-wails.deb | grep -E 'png|desktop'
dpkg-deb --fsys-tarfile bin/antdv-next-admin-thin-wails.deb \
  | tar -xOf - ./usr/share/pixmaps/antdv-next-admin-thin-wails.png | sha256sum
sha256sum /usr/share/pixmaps/antdv-next-admin-thin-wails.png
```

三处 PNG 的 SHA256 应一致。源文件与包内不一致说明 deb 未更新或打包目录不对；
包内与已安装文件不一致说明尚未安装该包。安装后完全退出旧进程，再从应用菜单启动。
文件都一致但固定图标仍旧时，取消固定后重新固定；Ubuntu 必要时注销后重新登录。
Windows 同样需要运行新 EXE 或重新安装新安装包，已运行的旧进程不会自动换图标。

`build/config.yml` 保存 Wails 产品信息；Windows 资源还涉及 `build/windows/info.json` 和 NSIS 配置。
Linux 包版本、维护者等产品信息读取 `build/config.yml`，系统依赖仍由两份 nFPM 配置维护。
发布前需将配置中的公司、产品、主页和维护者信息改为实际值。

## 跨平台构建

优先在目标系统原生构建。Linux 需要 CGO；从非 Linux 系统、缺少 C 编译器或跨架构构建时，现有任务会使用 Docker：

```sh
task setup:docker
task build GOOS=linux ARCH=amd64 GTK_VERSION=4
```

现有交叉编译镜像基于 Debian 13。`GTK_VERSION` 只选择图形后端，不会改变镜像的 glibc 等系统库基线，
因此不能据此保证 Docker 产物兼容 Ubuntu 22.04 或 24.04。

## 项目结构

```text
frontend/                    Vue 前端及测试
internal/system/             系统服务
main.go                      应用入口与窗口配置
window_restore_windows.go    Windows 窗口恢复处理
window_restore_other.go      其他平台的对应空实现
Taskfile.yml                 开发、构建和打包入口
build/config.yml             Wails 产品及开发模式配置
build/Taskfile.yml           前端、绑定、图标等共用任务
build/windows/               Windows 资源及安装包配置
build/linux/                 Linux 构建、桌面入口及 nFPM 配置
bin/                         构建产物
```

## 检查

### GitHub Actions 手动打包 v2 / v3

工作流位于 `.github/workflows/build-desktop.yml`，仅通过 `workflow_dispatch` 手动触发，
不会因 push 或 PR 自动构建；手动运行后，六组构建全部成功才会汇总发布到同一个 GitHub Release。
该文件需要先提交到远程默认分支，才能在 GitHub Actions 页面显示手动运行入口。

在 **Actions → 构建 Wails v2 和 v3 安装包 → Run workflow** 中填写：

- `v3_ref`：默认 `master`，可指定 v3 分支、标签或提交。
- `v2_ref`：默认 `wailsv2`，可指定 v2 分支、标签或提交。

无需填写版本号。发布阶段自动按北京时间生成 `v年.月.日` 标签和标题，例如 `v2026.09.11`。

一次运行并行执行 6 个 amd64 任务，单个失败不会取消其他任务：

| 系统 | v2 产物 | v3 产物 |
| --- | --- | --- |
| Windows Server 2022 runner | 单文件 EXE + NSIS 安装包 | 单文件 EXE + NSIS 安装包 |
| Ubuntu 22.04 | GTK3 单文件程序 + deb | GTK3 单文件程序 + deb |
| Ubuntu 24.04 | GTK3 单文件程序 + deb | GTK4 单文件程序 + deb |

在运行详情的 **Artifacts** 下载对应版本和系统的文件，保留 14 天。
全部构建成功后，Release 中也会提供十二个程序／安装包附件和 `SHA256SUMS.txt`，不受 Artifacts 的 14 天保留期限制。
发布任务单独使用 `contents: write`。六组构建成功且十二个附件校验通过后，删除当天已有的
同名 Release、旧附件和标签，再创建草稿、上传新附件并正式发布；其他日期的 Release 不受影响。
若删除后上传失败，当天旧发布不会自动恢复，可重跑发布任务。新标签指向触发工作流的提交，
v2、v3 的实际源码提交分别记录在构建摘要中。发布标签不会自动改写应用内的产品版本号。
每份 artifact 同时包含单文件程序与安装包；Linux 单文件程序封装为 `*-standalone.tar.gz`，
解压后保留可执行权限。deb、EXE、安装包名称包含 v2/v3 和系统标识，避免相互覆盖。

每个任务检出所选后端提交及其记录的前端子模块提交，不追踪子模块分支最新版本。
公开 GitHub 子模块的 SSH 地址由 checkout 自动转换为 HTTPS，无需额外 SSH 密钥或 PAT。
前端的 `pnpm-workspace.yaml` 必须提交到对应子模块分支，其中 `allowBuilds` 显式允许
`@parcel/watcher` 的构建脚本。只在本地执行 `pnpm approve-builds` 而未提交配置，会导致 CI 安装失败。
修改该配置后，先提交并推送前端子模块，再更新对应后端分支记录的子模块提交。
构建摘要记录实际后端、前端提交及 Wails CLI 版本。工作流安装 NSIS 或对应 GTK/WebKit 开发库，
调用现有 Taskfile 构建并上传产物；通过 CI 构建不等于已经完成安装包及窗口行为实测。

### 本地检查

Go 静态检查：

```sh
go vet ./...
```

前端检查（在 `frontend/` 中运行）：

```sh
pnpm run lint
pnpm run type-check
pnpm run test:unit:run
pnpm run build:check
```

Wails v3 文档：[v3.wails.io](https://v3.wails.io/)。
