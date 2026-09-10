# Wails v2 项目开发指南

## 项目定位

- 当前项目：`antdv-next-admin-thin-wails-v2`，Go 模块同名。
- Wails 当前依赖为 `github.com/wailsapp/wails/v2 v2.15.0`；版本以 `go.mod` 为准，CLI 使用 `wails`，不是 `wails3`。
- 前端是 Vue 3 + Antdv Next 无认证后台，子模块分支约定为 `no-auth-wails-v2`（以 `.gitmodules` 为准）。
- 不套用 v3 的 `application.NewService`、`frontend/bindings/` 或 GTK4 后端规则。

## v2 结构与服务注册

- `main.go` 使用 `wails.Run(&options.App{...})`，通过 `Bind` 注册 Go 服务。
- `app.go` 保存应用 context，沿用现有生命周期处理。
- `frontend/wailsjs/` 是 v2 生成绑定；`@wails/*` 指向该目录。
- 系统服务入口对应：

```ts
export { CurrentUsername } from '@wails/go/system/Service';
```

- v2 runtime 使用 `@wails/runtime/runtime`，不要替换为 v3 的 `@wailsio/runtime`。
- `window_restore_windows.go` 监听当前进程的 Windows 最小化结束事件，通过异步客户区重算修复恢复比例；监听线程有消息循环，退出时必须解除监听。
- `window_restore_other.go` 提供非 Windows 空实现；不要破坏 `main.go` 中的清理调用或用页面 zoom 补偿替代原生处理。

## v2 构建与打包

| 命令 | 用途 |
| --- | --- |
| `task dev` | Wails 开发模式 |
| `task build` | 当前系统单文件程序 |
| `task build:package` | Windows NSIS / Linux deb |
| `task generate` | 生成 Wails 前端绑定 |
| `task doctor` | 检查 Wails 开发环境 |

- 配置入口为根 `Taskfile.yml` 与 `wails.json`；任务定义是命令存在与否的依据，不沿用旧文档中已移除的别名。
- 产物位于 `build/bin/`，不是 v3 的 `bin/`。`wails.json` 当前输出名为 `antdv-next-admin-thin-wails`，不要依据仓库名自动加 `-v2`。
- Windows 单文件使用 `wails build`，不能加 `-nopackage` 跳过图标等平台资源。
- Windows 安装包使用 `wails build -nsis`，需 NSIS / makensis；安装器配置在 `build/windows/installer/`。
- EXE 和安装包使用 `build/windows/icon.ico`；Linux 窗口图标由 `main.go` 内嵌 `build/appicon.png` 并传入 Linux Options。
- 产品版本来自 `wails.json` 的 `info.productVersion`；调整输出名、版本或图标时检查安装包脚本及桌面入口是否同步。

### Ubuntu 打包

- v2 使用 GTK3。Taskfile 优先检测 WebKit2GTK 4.1 并加 `webkit2_41` 标签，否则使用 4.0；不要套用 v3 在 24.04 选择 GTK4 的逻辑。
- Ubuntu 24.04 开发依赖为 `libgtk-3-dev libwebkit2gtk-4.1-dev`；22.04 可使用 `libwebkit2gtk-4.0-dev`。还需 `build-essential pkg-config`。
- deb 打包调用 `python3 build/linux/package.py`，由脚本整理程序、PNG 图标、desktop 和 control，最终调用 `dpkg-deb`。需要 `python3 dpkg-dev binutils`，不需第三方 Python 库。
- 脚本读取实际 ELF 架构和 WebKit ABI，声明相应包依赖；目前支持 amd64 / arm64。不要通过改文件扩展名伪装 deb，也不要使用 Windows EXE 充当 Linux 输入。
- 裸 ELF 的文件管理器图标由桌面环境决定；deb 安装的菜单入口关联 hicolor 图标。“单文件”不表示无需 GTK/WebKit 运行库。
- 在目标 Ubuntu 系统上原生构建并验证；未实际构建安装时，不能声称跨版本兼容或安装包测试通过。

## 工作范围与协作

- 本文件适用于当前项目根目录及子目录；编辑前端还须阅读 `frontend/AGENTS.md`。
- 以用户本次明确指定的目录为修改范围，不因相邻存在另一个 Wails 版本而顺带修改。
- 修改前检查根仓库和 `frontend/` 的 Git 状态，保留用户已有改动，不覆盖、不撤销无关修改。
- `frontend/` 是独立 Git 子模块。根仓库中的 `M frontend` 可能表示子模块提交变化或内部未提交修改；须进入子模块检查，不直接重置或更新。
- 未获授权不得访问 Git 远程；代码修改不自动授权提交、推送、创建 PR、合并或发布。
- 获准提交时使用中文标题和正文，可保留 Conventional Commits 类型与技术标识符。涉及子模块时明确其提交与根仓库指针的关系，不能仅提交指针而遗漏子模块源码。
- 在已授权范围内主动完成实现和必要检查，普通可逆选择自行处理；只有影响正确性、范围或授权的缺失信息才澄清。
- 项目内保存自建脚本与产物。允许所需工具正常使用标准缓存及临时目录，不主动清空共享缓存或修改全局配置；仍须遵守运行环境审批要求。

## Go 服务与前端入口

- Go 业务代码按领域放在 `internal/<领域>/`，沿用现有 `internal/system/service.go` 的方式。
- 业务前端统一从手写入口调用服务：

```ts
import { CurrentUsername } from '@/services/system';
```

- `frontend/src/services/system.ts` 使用显式 re-export，保留生成函数的类型、返回值和错误行为，不添加无用途的 Promise 包装或静默降级。
- 新服务按领域增加 `frontend/src/services/<领域>.ts`，只显式导出所需方法；不要建立容易同名冲突的全服务汇总入口。
- 生成绑定的真实路径只出现在服务入口中。不要让组件、store 或业务初始化代码直接依赖 Go 模块生成路径。
- 不手改自动生成绑定来实现适配；Go 导出方法或注册方式变更后，用对应版本的工具重新生成并检查 diff。
- `frontend/src/platform/window.ts` 封装窗口操作。Wails runtime 与业务 Go 服务不是同一层，不为缩短业务导入而随意混合它们。
- 使用 pnpm，沿用现有 Vue Composition API、TypeScript 严格检查和别名；不引入无需求支持的依赖、抽象或额外锁文件。

## 实现与验证

- 新文本文件使用 UTF-8，既有文件保留项目编码约定。Go 使用 gofmt；前端格式遵循 `frontend/AGENTS.md`。
- 保留错误传播、必要日志及资源清理，不能用吞错或伪造成功掩盖问题。
- 未明确要求，不新增或修改测试文件；可以运行相关现有测试。
- Go 改动运行相关构建及 `go vet ./...`；项目通过 embed 引用 `frontend/dist`，缺少前端产物时先走构建任务。
- 前端服务入口或导入改动至少在 `frontend/` 运行 `pnpm run type-check`；按影响补充现有 lint、测试或构建。若要提交前端，遵循前端指南列出的提交前检查。
- 打包改动检查任务解析、实际命令、产物路径、图标和运行依赖；跨系统未经实测不能声称构建或安装已通过。
- 修改窗口恢复逻辑时验证“最大化 → 最小化 → 从任务栏恢复”，检查界面比例、窗口状态与退出清理。
- 验证通过且没有新改动或未解决问题时，不重复检查或扩大任务。交付说明修改内容、已执行的验证和剩余限制。
