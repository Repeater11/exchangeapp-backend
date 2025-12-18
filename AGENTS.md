# Repository Guidelines

## 项目结构与模块组织
- `cmd/server/main.go` 入口 entrypoint，加载配置并现阶段打印端口；在此扩展服务启动流程。
- `internal/config/config.go` 定义配置结构体与 Viper 加载器；新增配置先在此声明，再在业务侧使用。
- `configs/config.yml` 存放运行默认值，视为模板，勿提交真实机密或生产凭据。
- `go.mod` 指定 Go 1.24.x 与 `exchangeapp` 模块依赖，保持版本一致性。

## 构建、测试与开发命令
- `go run ./cmd/server` 本地运行服务，读取 `configs/config.yml`。
- `go build -o bin/server ./cmd/server` 构建二进制（生成 `bin/`）。
- `go fmt ./...` 按 Go 标准格式化；提交前必跑。
- `go test ./...` 运行单测；故障排查可加 `-v`。
- `go vet ./...` 静态检查常见问题；开 PR 前执行。

## 代码风格与命名约定
- 遵循 `go fmt`（tab 缩进、规范导入），导出标识符用 `CamelCase`，接收者命名简洁（`c`、`cfg`）。
- 包按关注点拆分，目录名小写；避免名称重复或口吃（如 `config.Config`）。
- `config.yml` 配置键保持小写+下划线，便于与结构体标签对齐。
- 倾向提前返回和小函数；致命日志仅在进程边界使用（如配置加载失败）。

## 测试规范
- 测试与源码同目录，命名 `_test.go`，使用 Go `testing` 包和表驱动用例。
- 测试命名 `Test函数名场景`，对不同配置用子测试。
- 需要样例时放在 `testdata/`，保持精简。
- 合并前覆盖配置解析和核心业务逻辑的关键路径。

## 提交与 Pull Request 指南
- 当前仓库无可见历史，建议用 Conventional Commits（如 `feat: add config validation`）或简洁祈使句描述意图。
- PR 说明需包含范围、风险、验证步骤（如 `go test ./...`、手动配置加载）。
- 关联相关 issue，必要时附 CLI 输出或截图，并注明配置变更以便审阅。
- 聚焦差异，尽量将重构与行为变更拆分提交。

## 安全与配置提示
- 开发和测试使用本地 `configs/config.yml` 样例，生产请改用环境变量或独立配置文件并更新 `.gitignore`，避免误传敏感数据。
- 数据库等机密不要写入 git；优先用 env vars、CI secrets、或 vault，在代码侧通过 Viper `AutomaticEnv` 等方式注入。
- 变更配置字段时同步更新 `config.yml` 默认值、结构体注释、以及相关文档，并运行 `go vet` 确认未使用的键被清理。
- 日志中避免输出密码或密钥；如需排查，使用 mock 配置或局部 debug 日志，并在提交前移除。
