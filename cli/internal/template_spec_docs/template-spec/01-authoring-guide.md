# Authoring Guide

本文档面向模板作者，说明如何写出一个合法的 TemplateSpec，以及如何排查常见错误。字段的完整定义和约束以 `00-template-spec.md` 为准。

---

## 写一个 spec 的顺序

1. 定义 steps：你的模板有哪些步骤，每步用什么能力
2. 定义 inputSchema：用户需要填写什么
3. 选择绑定方式：用户输入如何到达步骤
4. 连接步骤（多步骤时）：上游结果如何流入下游
5. 开放模型覆盖（可选）

---

## 第一步：定义 steps

从模板要做什么出发，先确定有几个步骤、每步用什么 executionUnit。

**命名规则**：`stepId` 格式为 `stp_` 加 6-10 位字母数字，例如 `stp_ab12cd`。同一 spec 内不能重复。

**instruction vs defaultValue**：如果需要给 step 一个固定的系统指令（例如"请将以下内容翻译为英文"），用 `instruction`，不要通过 field 传入。`instruction` 由模板作者写定，用户不可见、不可覆盖。

**多步骤时**：在依赖下游的 step 上声明 `dependsOn`。一个 step 可以声明多个上游，用来表达 step-level fan-in。汇聚节点未声明 `triggerPolicy` 时，默认按 `require_all` 处理。

---

## 第二步：定义 inputSchema

确定用户需要填写什么，以及每个字段的类型和可见性。

### valueType 选择

| 场景 | 选用 |
|---|---|
| 普通文本输入 | `string` |
| 固定选项 | `enum`（配合 `enumValues`） |
| 图片 URL | `image_url` |
| 用户上传的文本文件 | `text_reference`（配合 `acceptedMimeTypes`），并通过 `upstreamBindings` 的 `sourceType=initial_input` 绑定到文本 input port |
| 用户上传的非文本文件 | `asset_ref`（配合 `acceptedMimeTypes`） |

### sourceKind 选择

| 场景 | 选用 |
|---|---|
| 需要用户填写 | `user_input`（默认，不需要显式写） |
| 有合理默认值，不强制用户修改 | `default_value`（`defaultValue` 必须非空） |
| 模板内部使用，不展示给用户 | `hidden`（`defaultValue` 必须非空） |

### multiValue

当 step 需要并行处理多个输入时，将对应字段设置 `multiValue=true`，并指定 `maxValues`。这个字段后续必须通过 `bindMode=expanded` 绑定，触发 fan-out。

在 Workbook 中，`multiValue=true` 不会展开成多列。模板使用者需要在同一个单元格里填写多个值，并使用英文分号 `;` 或换行分隔。

---

## 第三步：选择绑定方式

绑定规则：**同一个 step 的同一个输入目标只能有一种来源**（FieldBinding、ParamBinding、UpstreamBinding 互斥）。

### FieldBinding vs ParamBinding

**用 FieldBinding**：一个字段直接对应一个 step 参数，一对一，无需拼接。

**用 ParamBinding**：需要将多个来源拼合成一个参数，例如在用户输入前加固定前缀，或将最多 3 个普通可见文本字段合并为同一个 `prompt`。

如果一个模型步骤需要用户分别填写“正文内容”“风格要求”“输出格式”，应使用一条 `ParamBinding` 写入该 step 的 `prompt`，不要创建多条 `FieldBinding` 重复写入同一个 `stepId + prompt`。

图片、视频或文件素材输入不属于文本参数拼装，仍应通过 input port 连接，并受具体模型能力约束。

### bindMode 选择

- `multiValue=false` 的字段 → `bindMode=shared`
- `multiValue=true` 的字段 → `bindMode=expanded`

两者必须一致，不匹配会在 `create` 时报校验错误。

---

## 第四步：连接 steps（多步骤时）

用 `upstreamBindings` 将上游 step 的输出端口连接到下游 step 的输入端口。

**注意**：UpstreamBinding 的 `inputPort` 和 FieldBinding/ParamBinding 的 `paramKey` 属于同一个目标命名空间。如果 UpstreamBinding 已经将 `prompt` 连接到本 step，就不能再通过 FieldBinding/ParamBinding 向同一个 `prompt` 写入。

多个上游可以绑定到同一个 `inputPort`，前提是目标输入端口允许多来源，并且所有上游输出类型都被目标输入端口接受。例如多个 `text-generate.output` 可以汇入下游 `text-generate.prompt`。

如果下游 step 需要额外的指令，用 `instruction` 写定，不要通过额外 binding 传入。

### triggerPolicy 选择

当 step 有上游依赖时，可以通过 `triggerPolicy` 控制触发策略：

| 策略 | 何时使用 | 运行时行为 |
|---|---|---|
| `require_all` | 默认选择。汇总必须基于全部上游结果 | 全部上游成功才执行；任一上游失败则不执行 |
| `allow_partial` | 部分角色、部分分支失败时，仍希望基于成功结果继续汇总 | 等所有上游结束后，只要至少一个成功就执行；失败上游不进入业务输入 |
| `fail_fast` | 任一关键上游失败后，后续结果没有意义 | 任一上游失败后，本节点及后续依赖路径不执行 |

`allow_partial` 是 workflow runtime 的触发策略，不是要求模板作者在自然语言里手写“如果部分失败请说明”。系统会给下游 step 附带结构化输入完整性元信息；失败上游的系统错误原因仍应从 step 状态、错误列或结果视图读取，不应混入业务正文。

---

## 第五步：开放模型覆盖（可选）

让用户在 Workbook 里选择模型，需要同时满足两个条件：

**条件一**：step 声明 `allowModelOverride=true`

**条件二**：在 `fieldBindings` 里将某个字段绑定到该 step 的 `paramKey=model`

只满足其中一个条件，模型列不会出现在 Workbook 中。`provider` 和 `mode` 不支持作为模板参数暴露。

---

## 校验层级

提交 TemplateSpec 时，系统依次进行四层校验：

| 层级 | 检查内容 | 发生时机 |
|---|---|---|
| Schema 校验 | JSON 结构是否合法，字段类型是否正确 | `check` 接口 |
| Semantic 校验 | 字段引用是否有效，valueType 约束是否满足 | `check` 接口 |
| Authoring 校验 | Workflow shape 是否在 v1 支持范围内 | `create` 接口 |
| Runtime 校验 | 模型可用性等运行时条件 | 实际执行时 |

**schema valid ≠ authoring valid ≠ runtime guaranteed**。通过前一层不代表后一层也会通过。

---

## `check` 通过，`create` 失败

JSON 结构和字段引用合法，但 workflow shape 违反了 v1 约束。常见原因：

| 原因 | 说明 |
|---|---|
| 多上游 step 缺少显式 binding | 多上游 fan-in 必须通过 `upstreamBindings` 明确每个上游 output 如何进入下游 input port |
| `inputPort` 不存在 | 指定了该 executionUnit 没有的输入端口 |
| `sourcePort` 不存在 | 指定了上游 step 没有的输出端口 |
| `sourceStepId` 不在 `dependsOn` 里 | UpstreamBinding 引用的 step 未声明依赖 |
| Binding target 冲突 | FieldBinding / ParamBinding 与 UpstreamBinding 同时写入同一个目标，或多个上游写入了 `allowMultiple=false` 的 input port |
| 上下游类型不匹配 | 上游 output type 不在下游 input port 的 `accepts` 范围内 |
| `triggerPolicy` 非法 | 只允许 `require_all`、`allow_partial`、`fail_fast` |
| `allow_partial` 用在 root step | `allow_partial` 必须至少有一个上游依赖 |
| 模型绑定未开放 | 对未声明 `allowModelOverride=true` 的 step 绑定 `paramKey=model` |
| `bindMode` 与 `multiValue` 不匹配 | 两者必须一致 |

## `create` 成功，运行时失败

模板定义合法，但运行时条件不满足。常见原因：

- 模型 ID 不存在、未开通或已下线
- Provider 服务异常
- 输入内容触发业务侧限制

---

## 运行结果读取

运行完成后，优先使用服务端结果视图，而不是依赖本地 Excel 回填。

| 结果视图 | 适用场景 | 说明 |
|---|---|---|
| `result-rows` | 所有带输入快照的新 run | 返回输入快照、状态、错误和 artifacts，用于程序化读取 |
| `result-workbook` | workbook/Excel 提交的 run | 下载服务端生成的结果 Excel，保留原始输入列并追加结果列 |

Artifact URL 是结果产物的规范访问方式。对于 `text/*` 产物，服务端会尽量返回
`inlineText` 作为轻量预览；当前只内联不超过 4KB 的文本。超过阈值的文本产物仍是
文本结果，但 `result-rows` / `result-workbook` 会回退展示 artifact URL。

CLI 对应命令：

```bash
loomloom run result-rows <run-id>
loomloom run result-workbook <run-id> --output-file ./result.xlsx
```

`template backfill-results` 是旧的本地回填流程。新模板和内测验证应优先使用 `run result-workbook`，因为它使用服务端保存的输入快照和原始 workbook，能避免本地文件版本与 run 结果不一致。
