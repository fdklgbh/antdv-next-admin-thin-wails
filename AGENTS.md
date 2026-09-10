# Wails v3 项目开发指南

## 项目定位

- 当前项目：`antdv-next-admin-thin-wails`，Go 模块同名。
- Wails 当前依赖为 `github.com/wailsapp/wails/v3 v3.0.0-beta.19`；版本以 `go.mod` 为准，CLI 应与依赖匹配。
- 前端是 Vue 3 + Antdv Next 无认证后台，子模块分支约定为 `no-auth-wails-v3`（以 `.gitmodules` 为准）。
- 不套用 Wails v2 的 `wails.json`、`Bind` 或 `@wails/go/...` 方案。

## v3 结构与服务注册

- `main.go` 使用 `application.New`、`application.Options.Services` 和 `application.NewService` 注册 Go 服务。
- `frontend/bindings/` 是 v3 生成绑定；`@wails/*` 指向该目录。
- 系统服务入口对应：

```ts
export { CurrentUsername } from '@wails/antdv-next-admin-thin-wails/internal/system/service';
```

- v3 runtime 使用 `@wailsio/runtime`。即使目录中残留 `frontend/wailsjs/`，也不要将其当作 v3 业务绑定入口。
- `window_restore_windows.go` 在恢复事件后排队重算客户区；`window_restore_other.go` 提供非 Windows 空实现。不要用页面 zoom 补偿替代原生恢复处理。

## v3 构建与打包

| 命令 | 用途 |
| --- | --- |
| `task dev` | 开发模式，默认前端端口 9245 |
| `task build` | 当前系统生产程序 |
| `task run` | 运行已构建程序 |
| `task package` | Windows 默认 NSIS；Linux 为多种包格式 |
| `task linux:create:deb` | 构建 Linux 程序并仅生成 deb |

- 构建入口为根 `Taskfile.yml`；共用任务在 `build/Taskfile.yml`，平台任务在 `build/<平台>/Taskfile.yml`。
- 常规 Windows/Linux 产物位于 `bin/`，不是 v2 的 `build/bin/`。
- Windows 生成 `.syso` 后编译，包含图标、manifest 和版本信息；安装包需要 NSIS，使用 `build/windows/nsis/`。
- 源图标为 `build/appicon.png`，Windows ICO 由图标任务生成。Linux 包包含 `.desktop` 和 hicolor 图标，裸 ELF 文件图标由桌面环境决定。
- 产品与开发配置在 `build/config.yml`；Linux 版本和包元数据在 nFPM 配置中，当前仍有模板信息，不能把它们当作正式发布信息。

### Ubuntu 后端选择

| 系统 | 后端 | 开发依赖 | nFPM 配置 |
| --- | --- | --- | --- |
| Ubuntu 22.04 | GTK3 + WebKit2GTK 4.1 | `libgtk-3-dev libwebkit2gtk-4.1-dev` | `build/linux/nfpm/nfpm-gtk3.yaml` |
| Ubuntu 24.04 | GTK4 + WebKitGTK 6.0 | `libgtk-4-dev libwebkitgtk-6.0-dev` | `build/linux/nfpm/nfpm.yaml` |

- 还需 `build-essential pkg-config`。Taskfile 按 `/etc/os-release` 自动选择；不要将两者统一为 GTK3。
- `GTK_VERSION=3` 或 `GTK_VERSION=4` 同时控制编译标签及包依赖；不要仅使用 `EXTRA_TAGS=gtk3`。
- deb 由 `wails3 tool package` 和 nFPM 配置生成，不使用 v2 的 Python 脚本。程序、desktop 和图标路径必须一致，Linux 名称不带 `.exe`。
- `task package` 在 Linux 会生成 AppImage、deb、RPM、Arch 包；只需 deb 时不必运行全格式打包。
- Ubuntu 22.04 和 24.04 需分别原生验证；不能保证新系统产物在旧系统运行。当前 Docker 交叉编译镜像基于 Debian 13，GTK_VERSION 不会改变其 glibc 基线。
- 当前 Ubuntu 配置仍待实际系统验证。升级 Wails 前复核 GTK3 支持，当前说明指出 v3.1 将移除兼容模式。

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
