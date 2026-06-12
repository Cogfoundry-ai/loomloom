# TemplateSpec Specification

## 1. Context

TemplateSpec 是 LoomLoom 模板系统的规范格式，用于声明式地定义一个可重复运行的工作流模板。

三个规范级事实：

**TemplateSpec 是唯一事实源。** 模板版本保存的是 TemplateSpec 快照，所有运行时行为从这里派生。

**Workbook 是派生文件。** Workbook 由 TemplateSpec 编译生成，不是主数据。模板版本变化后，Workbook 必须重新派生，不能复用旧版本。

**OpenAPI 不定义 TemplateSpec 语义。** OpenAPI 负责描述接口路径和请求/响应结构；TemplateSpec 的规则、语义和能力边界以本文档为准。

---

## 2. TemplateSpec Pipeline

```text
TemplateSpec
  ↓ compile
WorkflowDefinition + Workbook
  ↓ user fills Workbook
Workbook Submission
  ↓
Run → Executions → Artifacts
```

---

## 3. Terminology

| 术语 | 定义 |
|---|---|
| **Step** | 模板中的一个处理单元，使用一种 ExecutionUnit 能力运行。 |
| **ExecutionUnit** | Step 使用的能力类型，当前支持 `text-generate`、`image-generate`、`video-generate`。 |
| **InputSchema** | 定义用户需要填写的输入字段集合。 |
| **Field** | InputSchema 中的一个输入字段，具有类型、可见性、是否多值等属性。 |
| **FieldBinding** | 将一个 Field 直接映射到某个 step 的某个参数。 |
| **ParamBinding** | 将多个来源（Field、字面量）拼合后映射到某个 step 的某个参数。 |
| **UpstreamBinding** | 将上游 step 的输出端口连接到下游 step 的输入端口。 |
| **Workbook** | TemplateSpec 编译后生成的用户填写界面。 |
| **Execution** | 一次 step 的实际运行实例。fan-out 时一个 step 可产生多个 executions。 |
| **Artifact** | 一次 execution 产生的输出。 |

---

## 4. Top-Level Structure

TemplateSpec 是一个 JSON 对象，包含以下顶层字段：

| 字段 | 类型 | 是否必填 | 说明 |
|---|---|---|---|
| `meta` | object | 必填 | 模板元信息 |
| `steps` | array | 必填，非空 | 执行单元链路 |
| `inputSchema` | object | 必填 | 用户输入字段定义 |
| `fieldBindings` | array | 可选 | Field 到 step 参数的直接绑定；不要用于把 `text_reference` 直接绑定到 `prompt` |
| `paramBindings` | array | 可选 | 多来源组合到 step 参数的绑定 |

---

## 5. Meta

`meta` 描述模板的元信息。

| 字段 | 类型 | 是否必填 | 说明 |
|---|---|---|---|
| `name` | string | 必填 | 模板名称，不能为空 |
| `description` | string | 可选 | 模板说明 |
| `scenario` | string | 可选 | 适用场景分类 |
| `inputSummary` | string | 可选 | 输入内容说明，用于前端展示 |
| `displayOutputType` | string | 可选 | 输出类型的展示名称 |
| `primaryOutputType` | string | 可选 | 如果填写，必须与系统根据终端 step 推导出的输出类型一致 |
| `tags` | array\<string\> | 可选 | 标签列表，用于分类和检索 |

---

## 6. Steps

`steps` 是一个有序数组，每个元素是一个 StepSpec 对象。

### 6.1 StepSpec

| 字段 | 类型 | 是否必填 | 说明 |
|---|---|---|---|
| `stepId` | string | 必填 | 格式：`stp_` 加 6-10 位字母数字（例如 `stp_ab12cd`）。同一 spec 内不能重复。 |
| `displayName` | string | 必填 | 显示名称 |
| `executionUnit` | string | 必填 | 使用的能力类型，必须是 Capability Registry 中已注册的值 |
| `defaultModelRef` | object | 条件必填 | 如果 executionUnit 需要模型，则必须提供 `modelKey` |
| `instruction` | string | 可选 | 模板作者写定的系统级说明，不由用户覆盖 |
| `dependsOn` | array\<string\> | 可选 | 依赖的上游 stepId 列表。可声明多个上游来表达 step-level fan-in。 |
| `upstreamBindings` | array | 可选 | 上游 step 输出到本 step 输入端口的连接，见 6.2 |
| `triggerPolicy` | string | 可选 | 汇聚节点触发策略，见 6.3。空值按 `require_all` 处理 |
| `allowModelOverride` | boolean | 可选 | 是否允许用户覆盖本 step 的模型，默认 `false` |
| `staticParams` | object | 可选 | 写死的运行参数，键必须是该 executionUnit 允许的参数名 |

### 6.2 UpstreamBinding

| 字段 | 类型 | 是否必填 | 说明 |
|---|---|---|---|
| `inputPort` | string | 必填 | 本 step 的输入端口名 |
| `sourceType` | string | 可选 | `step_output`（默认）或 `initial_input` |
| `sourceStepId` | string | 条件必填 | `sourceType=step_output` 时必填，必须出现在 `dependsOn` 里 |
| `sourcePort` | string | 条件必填 | `sourceType=step_output` 时必填，必须是上游 step 实际有的输出端口名 |
| `sourceInputKey` | string | 条件必填 | `sourceType=initial_input` 时必填，指定从初始输入读取哪个字段 |

### 6.3 triggerPolicy

`triggerPolicy` 控制有上游依赖的 step 何时触发，主要用于多上游汇聚节点。

| 值 | 语义 | 适用场景 |
|---|---|---|
| `require_all` | 默认策略。全部上游 step 成功后才触发；任一上游最终失败、取消或被阻断时，本 step 不执行 | 完整性要求高、强依赖流程 |
| `allow_partial` | 等所有上游进入终态后，只要至少一个上游成功就触发；下游只接收成功上游的产物 | 可以接受部分结果继续汇总的流程 |
| `fail_fast` | 任一上游最终失败、取消或被阻断时，本 step 不执行，并阻断后续依赖路径 | 任一关键输入失败就没有继续价值的流程 |

约束：

- 未声明 `triggerPolicy` 时等同于 `require_all`。
- `allow_partial` 不能用于 root step，必须至少有一个上游依赖。
- 多上游 fan-in 使用 `allow_partial` 时，必须声明显式 `upstreamBindings`。
- `allow_partial` 不会把失败上游的系统错误正文传入业务输入；失败原因仍在 step 状态、错误列或结果视图中查看。
- `allow_partial` 下如果所有上游都失败、取消或被阻断，下游 step 不执行。

---

## 7. Input Schema

`inputSchema` 定义用户需要填写的输入字段。

| 字段 | 类型 | 说明 |
|---|---|---|
| `fields` | array | 输入字段列表，必填，非空 |
| `instructions` | array\<string\> | 填写引导文字，显示在表单上方 |
| `sampleRows` | array\<object\> | 示例数据，键必须使用 field 的 `key`，不能使用 `label` |

### 7.1 TemplateInputField

| 字段 | 类型 | 是否必填 | 说明 |
|---|---|---|---|
| `key` | string | 必填 | 字段唯一标识。保留键（`model`、`provider`、`mode`）不能使用。同一 spec 内不能重复。 |
| `label` | string | 必填 | 显示名称。同一 spec 内不能重复。 |
| `valueType` | string | 必填 | 值类型，见 7.2 |
| `required` | boolean | 可选 | 是否必填，默认 `false` |
| `defaultValue` | string | 条件必填 | 当 `sourceKind` 不是 `user_input` 时，必须非空 |
| `sourceKind` | string | 可选 | 控制字段在 Workbook 里的可见性，见 7.3 |
| `multiValue` | boolean | 可选 | 是否允许多个值，默认 `false` |
| `maxValues` | integer | 条件必填 | `multiValue=true` 时必须填写且 `> 0` |
| `enumValues` | array\<string\> | 条件必填 | `valueType=enum` 时必须填写且非空 |
| `acceptedMimeTypes` | array\<string\> | 条件必填 | `valueType=asset_ref` 或 `text_reference` 时必须填写且非空 |
| `order` | integer | 可选 | 在 Workbook 里的显示顺序 |
| `description` | string | 可选 | 字段说明 |

如果 `multiValue=true`，Workbook 仍然只占用一个单元格。多个值需要在同一单元格中输入，并使用英文分号 `;` 或换行分隔。

### 7.2 valueType

| 值 | 含义 | 额外约束 |
|---|---|---|
| `string` | 普通文本 | 无 |
| `enum` | 枚举值 | `enumValues` 必须非空 |
| `image_url` | 图片 URL | 值必须是合法的 HTTP/HTTPS URL |
| `asset_ref` | 已上传文件的 asset ID | 值必须是 `ia_` 开头的 UUID；`acceptedMimeTypes` 必须非空 |
| `text_reference` | 内联文本或 asset ref | 支持内联文本和 asset ref，不支持 URL；`acceptedMimeTypes` 必须非空；如果可能填写 `ia_xxx`，必须通过 `upstreamBindings` 的 `sourceType=initial_input` 绑定到文本 input port |

`text_reference` 不能只通过 `fieldBindings` 直接绑定到 `prompt`。否则用户填写 `ia_xxx`
时，模型只会收到 asset ID 字符串，而不是文件内容。需要让模型读取上传文本时，应使用：

```json
{
  "inputPort": "reference",
  "sourceType": "initial_input",
  "sourceInputKey": "reference_text"
}
```

### 7.3 sourceKind

| 值 | Workbook 行为 | 额外约束 |
|---|---|---|
| `user_input`（默认） | 用户可见、可编辑 | 无 |
| `default_value` | 可见，不可编辑，由默认值驱动 | `defaultValue` 必须非空 |
| `hidden` | 不展示给用户 | `defaultValue` 必须非空 |

`hidden=true` 与 `sourceKind=hidden` 等效。

---

## 8. Bindings

### 8.1 FieldBinding

将一个 Field 直接映射到一个 step 的一个参数。

| 字段 | 类型 | 是否必填 | 说明 |
|---|---|---|---|
| `fieldKey` | string | 必填 | 引用 `inputSchema.fields` 中的 `key` |
| `stepId` | string | 必填 | 目标 step 的 `stepId` |
| `paramKey` | string | 必填 | 目标参数名 |
| `bindMode` | string | 必填 | `shared` 或 `expanded`，必须与 field 的 `multiValue` 一致 |

`multiValue=false` 的字段必须使用 `bindMode=shared`；`multiValue=true` 的字段必须使用 `bindMode=expanded`。

路由参数特殊规则：`paramKey=model` 要求对应 step 必须声明 `allowModelOverride=true`；`paramKey=provider` 和 `paramKey=mode` 不允许通过 binding 暴露。

### 8.2 ParamBinding

将多个来源拼合后映射到一个 step 的一个参数。

| 字段 | 类型 | 是否必填 | 说明 |
|---|---|---|---|
| `stepId` | string | 必填 | 目标 step 的 `stepId` |
| `paramKey` | string | 必填 | 目标参数名。不能是路由参数（`model`、`provider`、`mode`）。 |
| `bindMode` | string | 必填 | `shared` 或 `expanded` |
| `separator` | string | 可选 | 多个 source 拼合时的分隔符 |
| `sources` | array | 必填，非空 | 来源列表，见下方 ParamSource |

**ParamSource：**

| 字段 | 类型 | 说明 |
|---|---|---|
| `kind` | string | `field_ref` 或 `literal` |
| `fieldKey` | string | `kind=field_ref` 时，引用的 field key |
| `literal` | string | `kind=literal` 时，写死的字符串，不能为空 |

一个 ParamBinding 最多允许一个 `multiValue=true` 的 field_ref 来源。如果存在 `multiValue=true` 的 field_ref，则 `bindMode` 必须是 `expanded`。

一个 ParamBinding 可以包含最多 3 个普通可见字段来源。例如可以把“正文内容”“风格要求”“输出格式”组合成同一个 `prompt` 参数。多字段组合应通过 `ParamBinding.sources` 表达，不应创建多条 `FieldBinding` 重复写入同一个 `stepId + paramKey`。

当前推荐的产品化范围是组合文本字段生成模型类执行器的 `prompt` 参数。图片、视频或文件素材输入仍应通过 `UpstreamBinding` / input port 进入，并受模型能力限制；不要把多图片输入数量理解为 ParamBinding 能力。

### 8.3 Binding Target 与多上游输入

在同一个 step 内，FieldBinding / ParamBinding 写入的 `paramKey` 仍然必须唯一。

FieldBinding / ParamBinding 使用 `paramKey` 表示目标参数；UpstreamBinding 使用 `inputPort` 表示目标输入端口。在 v1 中，两者属于同一个目标命名空间。

因此，下面的结构仍然是非法的：

```text
FieldBinding  → stp_summary.prompt
UpstreamBinding → stp_summary.prompt（同时绑定）
```

也就是说，同一个 step 的同一个目标不能同时来自用户字段和上游产物。

UpstreamBinding 支持多个上游写入同一个 `inputPort`，但必须满足该目标输入端口声明 `allowMultiple=true`，并且每个上游输出类型都能被目标输入端口的 `accepts` 接受。

例如多个 `text-generate.output`（`text/plain`）可以同时绑定到下游 `text-generate.prompt`（`accepts=text/*`，`allowMultiple=true`）。

---

## 9. Execution Semantics

### 9.1 shared

字段的值被该 step 的所有 executions 共享。适用于所有并行 executions 使用同一个值的场景。

### 9.2 expanded 与 fan-out

`bindMode=expanded` 将 `multiValue` 字段的多个值展开为同一 step 的多个独立 executions，每个值对应一次运行。这称为 **fan-out**。

```text
输入：["值1", "值2", "值3"]（bindMode=expanded）
  → execution 1：值1
  → execution 2：值2
  → execution 3：值3
```

### 9.3 Step-Level Fan-In

当下游 step 通过多个 UpstreamBinding 连接多个上游 step 时，系统将这些上游 outputs 汇入同一个下游输入端口。

```text
stp_product.output     ↘
stp_engineering.output  → stp_summary.prompt
stp_risk.output        ↗
```

条件：

- `dependsOn` 必须声明所有上游 step。
- 每个 UpstreamBinding 的 `sourceStepId` 必须属于 `dependsOn`。
- `sourcePort` 必须存在。
- 目标 `inputPort` 必须存在。
- `sourcePort.type` 必须匹配目标 `inputPort.accepts`。
- 多个 binding 写入同一个 `inputPort` 时，该输入端口必须 `allowMultiple=true`。

---

## 10. Capability Registry

### `text-generate`

| 属性 | 值 |
|---|---|
| 输入端口 | `prompt`（text/\*，必填，AllowMultiple=true）；`reference`（text/\*，可选） |
| 输出端口 | `output`（text/plain） |
| 模型覆盖 | 支持，需要 `allowModelOverride=true` |

### `image-generate`

| 属性 | 值 |
|---|---|
| 输入端口 | `prompt`（text/\*，必填，AllowMultiple=true） |
| 输出端口 | `output`（image/png） |
| 模型覆盖 | 支持，需要 `allowModelOverride=true` |

### `video-generate`

| 属性 | 值 |
|---|---|
| 输入端口 | `prompt`（text/\*，必填，AllowMultiple=false）；`image`（image/\*，可选，AllowMultiple=false） |
| 输出端口 | `output`（video/mp4） |
| 模型覆盖 | 不支持 |

---

## 11. Unsupported Patterns

### 类型不兼容的 Step 依赖

step 可以声明多个上游，但依赖必须能建立有效的数据连接。下面的结构是非法的：

```text
image-generate.output（image/png） → text-generate.prompt（只接受 text/*）
```

原因：source output type 不匹配 target input accepts。

### 单输入端口的多来源绑定

下面的结构也是非法的：

```text
stp_a.output ↘
stp_b.output  → video-generate.prompt
```

原因：`video-generate.prompt` 只接受一个输入来源，`allowMultiple=false`。

---

## 12. Versioning

本文档描述的规范版本为 **TemplateSpec v1**（内部标识 `template-spec/v1`）。

v1 的主要边界：

- 支持 step-level fan-in / fan-out。
- 支持 fan-in 汇聚节点触发策略：`require_all`、`allow_partial`、`fail_fast`。
- FieldBinding / ParamBinding 仍保持单目标唯一。
- UpstreamBinding 多来源写入同一 inputPort 时，必须由目标 port 的 `allowMultiple` 和 `accepts` 明确允许。
- merge 语义由目标 input port 的 `mergePolicy` 决定。
