# Antdv Next Admin Thin · Wails v3

基于 Wails v3、Vue 3 和 Antdv Next 的桌面管理后台脚手架。前端资源内嵌于可执行文件，使用系统 WebView 渲染界面。

当前 Go 依赖为 **Wails v3.0.0-beta.19**。本 README 对应 v3 项目，构建配置入口为 `Taskfile.yml` 和 `build/config.yml`。

## 环境准备

- Go：满足 `go.mod` 声明的版本（当前为 `1.26.7`）。
- Node.js：推荐 22.12+，以及 pnpm。
- Task：提供 `task` 命令。
- Wails v3 CLI：建议与项目依赖保持同版本。

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.19
```

确保 Go 的 bin 目录在 PATH 中，并按下文安装目标系统的原生依赖。
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

deb 也输出到 `bin/`，文件名由包名、版本、release 和架构组成。
使用 `sudo apt install ./bin/<实际文件名>.deb` 安装，以便同时解析系统运行依赖。
当前两份 nFPM 配置中的版本为 `0.1.0`。

- `task package` 会依次生成 AppImage、deb、RPM 和 Arch 包；只要 deb 时使用 `task linux:create:deb`。
- deb 由 `wails3 tool package` 生成，不需要额外 Python 打包脚本。
- 开发与生产构建使用同一后端选择规则；其他系统默认 GTK4。
- 显式覆盖时使用 `task build GTK_VERSION=3` 或 `task linux:create:deb GTK_VERSION=4`。
  `GTK_VERSION` 同时控制编译标签和包依赖，不要仅通过 `EXTRA_TAGS=gtk3` 切换。

请分别在目标 Ubuntu 系统上原生构建和验证。24.04 构建的程序不保证能在 22.04 上运行，
因为 GTK 后端和 glibc 等系统库版本可能不同。当前 Ubuntu 22.04 / 24.04 的实际构建、安装及运行仍待验证。
按当前 Wails 说明，GTK3 兼容模式将在 v3.1 移除，升级时需要重新评估 22.04 支持。

## 图标与包信息

- `build/appicon.png`：源图标；当前为 1024×1024。
- `build/windows/icon.ico`：由图标任务生成，用于 Windows EXE 和 NSIS 安装包。
- `build/linux/antdv-next-admin-thin-wails.desktop`：Linux 构建生成的应用菜单入口。
- 两份 nFPM 配置将程序安装到 `/usr/bin/`，将桌面入口和图标安装到系统应用目录及 hicolor 图标主题目录。

Linux 裸 ELF 在文件管理器中的图标由桌面环境决定，不能像 Windows EXE 一样保证显示自定义文件图标；
deb 安装后通过应用菜单入口关联图标。

`build/config.yml` 保存 Wails 产品信息；Windows 资源还涉及 `build/windows/info.json` 和 NSIS 配置。
Linux 包版本、维护者及依赖以两份 nFPM 配置为准。发布前需将模板中的公司、产品和维护者信息改为实际值，
并同步维护两份 nFPM 的版本。维护者字段使用 `GIT_COMMITTER_NAME` 和 `GIT_COMMITTER_EMAIL` 环境变量。

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

- `release_tag`：必填，例如 `v1.0.0`，使用尚未发布的新标签。
- `v3_ref`：默认 `master`，可指定 v3 分支、标签或提交。
- `v2_ref`：默认 `wailsv2`，可指定 v2 分支、标签或提交。

一次运行并行执行 6 个 amd64 任务，单个失败不会取消其他任务：

| 系统 | v2 产物 | v3 产物 |
| --- | --- | --- |
| Windows Server 2022 runner | 单文件 EXE + NSIS 安装包 | 单文件 EXE + NSIS 安装包 |
| Ubuntu 22.04 | GTK3 单文件程序 + deb | GTK3 单文件程序 + deb |
| Ubuntu 24.04 | GTK3 单文件程序 + deb | GTK4 单文件程序 + deb |

在运行详情的 **Artifacts** 下载对应版本和系统的文件，保留 14 天。
全部构建成功后，Release 中也会提供十二个程序／安装包附件和 `SHA256SUMS.txt`，不受 Artifacts 的 14 天保留期限制。
发布任务单独使用 `contents: write`，先创建草稿并上传附件，再正式发布；不会覆盖同名 Release。
上传失败留下草稿时，检查并处理该草稿后再重跑发布任务。新标签默认指向触发工作流的提交，
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
