# API 文档维护说明

本项目使用 `docs/openapi.yaml` 作为接口契约（OpenAPI 3.0）。不使用 Swagger 注解，接口有变更时手动维护该文件。

## 维护规则
- 新增/修改/删除接口时，同步更新 `docs/openapi.yaml` 中的 `paths`
- DTO 结构调整时，同步更新 `components/schemas`
- 返回状态码变更时，同步更新对应 `responses`
- 鉴权接口需要标注 `security` 与 `bearerAuth`

## 本地查看方式（可选）
如需本地可视化查看，可用 Swagger UI 或 Redoc 打开 `docs/openapi.yaml`。

示例（Swagger UI）：
1. 安装 `swagger-ui` 或使用在线工具
2. 指向本地文件 `docs/openapi.yaml`

说明：此步骤不是项目必须依赖，仅用于阅读文档。
