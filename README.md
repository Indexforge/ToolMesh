# ToolMesh

## ToolMesh 项目实施计划

### 0. 工程技术基线（强约束）

#### 0.1 Web 与基础设施约束
- **API Server**
  - 使用 `net/http`
  - 明确区分 `http.Handler` / `http.HandlerFunc`
  - 使用标准 `http.Server` 控制生命周期
- **路由策略**
  - 初期使用 `http.ServeMux`
  - 明确禁止引入第三方路由库
- **日志系统**
  - 使用 Go 1.21+ 的 `log/slog`
  - 全链路结构化日志
  - 禁止 `log.Printf`、禁止非结构化日志
- **工程目标**
  - 可控依赖
  - 可预测性能
  - 可调试性优先于“开发爽感”

---

### 1. 总体架构规划（修订）

#### 1.1 架构定位

ToolMesh 是一个 AI 能力编排平台，不是 Web 框架驱动的应用。

HTTP 层的职责被严格限制为：
- 请求接收
- 参数反序列化
- Context 传递
- 响应序列化

任何业务决策均不允许出现在 Handler 中。

#### 1.2 核心模块划分（保持不变，职责更清晰）
1. **HTTP API Layer（`net/http`）**
   - 路由注册
   - 请求校验
   - Context 注入（request id / trace id）
   - 错误映射（domain error → HTTP）
2. **Orchestrator（核心编排层）**
   - Prompt 选择
   - RAG / MCP 调度
   - 多轮推理控制
   - Tool 调用合法性校验
3. **RAG Service**
   - 文档检索
   - 向量库访问
   - 只返回文本，不做推理
4. **MCP Tool Server**
   - 独立进程
   - 强类型输入 / 输出
   - 内部访问 MySQL
5. **LLM Gateway**
   - OpenAI-compatible API Client
   - Tool / Function calling 支持

---

### 2. HTTP API 设计原则（新增重点）

#### 2.1 Handler 设计铁律

每个 HTTP Handler 只允许做 5 件事：
1. 解析请求（JSON → struct）
2. 校验参数（基本合法性）
3. 调用 Orchestrator
4. 映射错误
5. 写回响应

禁止：
- Prompt 拼装
- Tool 判断
- 业务分支逻辑
- SQL / RAG 直接调用

#### 2.2 示例 Handler 结构（概念）

```go
func ChatHandler(orch *Orchestrator, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        var req ChatRequest
        if err := decodeJSON(r.Body, &req); err != nil {
            writeError(w, err)
            return
        }

        resp, err := orch.HandleChat(ctx, req)
        if err != nil {
            writeError(w, err)
            return
        }

        writeJSON(w, resp)
    }
}
```

Handler 像“插座”，不允许像“变压器”。

---

### 3. 日志系统设计（slog）

#### 3.1 slog 使用原则
- 全局统一 logger
- 禁止随意 `slog.Default()`
- 所有日志必须带上下文信息

#### 3.2 推荐字段规范

| 字段 | 说明 |
| --- | --- |
| `request_id` | 每个 HTTP 请求唯一 |
| `module` | http / orchestrator / rag / mcp |
| `tool` | Tool 名称（如有） |
| `duration_ms` | 调用耗时 |
| `error` | error 对象 |

#### 3.3 Logger 初始化（概念）

```go
logger := slog.New(
    slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }),
)
```

#### 3.4 日志分层策略
- HTTP 层：请求开始 / 结束
- Orchestrator：决策日志（选择了什么 Prompt / Tool）
- Tool 层：输入 / 输出摘要（不记录敏感数据）
- Error：必须结构化输出 error

---

### 4. 核心能力拆解（不变，但强调边界）

#### 4.1 RAG
- HTTP 层永远不直接调用
- 只允许 Orchestrator 调用
- slog 记录：
  - query
  - top-k
  - 命中文档 id

#### 4.2 MCP Tool
- Tool Server 自身也是 `net/http` 或 stdio（MCP）
- Tool 调用：
  - 参数校验在 Go 侧完成
- slog 记录：
  - tool name
  - input schema version
  - execution time

---

### 5. Orchestrator 设计（进一步强调）

#### 5.1 Orchestrator 的唯一职责
- 状态机 + 调度器
- 不关心 HTTP
- 不关心日志输出格式
- 不关心存储实现

它只回答一个问题：

> “当前这个问题，下一步该调用谁？”

---

### 6. 错误处理与 HTTP 映射（新增）

#### 6.1 错误分层
- `ErrBadRequest`
- `ErrToolFailed`
- `ErrModelUnavailable`
- `ErrNoContextFound`

#### 6.2 HTTP 映射示例

| Domain Error | HTTP |
| --- | --- |
| BadRequest | 400 |
| ToolFailed | 502 |
| ModelError | 503 |
| NoContext | 200 + 明确说明 |

---

### 7. 里程碑（微调）

**MVP 阶段（原生可控）**
- `net/http` API
- `slog` 全链路日志
- 单 Prompt 模板
- 单 MCP Tool
- 单模型

目标不是“快”，而是“不会重构”。

---

### 8. 工程结论（架构判断）

你这套约束意味着三件事：
1. ToolMesh 是平台，不是 Demo
2. 你在为 2–3 年后的维护负责
3. 你明确拒绝“框架绑架”
