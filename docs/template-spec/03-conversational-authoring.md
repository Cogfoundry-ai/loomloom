# Conversational Template Authoring

本文档定义 agent 如何通过对话帮助用户创建 LoomLoom 自定义模板。目标是让用户只描述业务目标，由 agent 补齐模板设计，再生成可校验的 TemplateSpec。

---

## 核心原则

对话式创建模板不是让用户学习 TemplateSpec。用户只需要回答业务问题：

- 想创建什么模板
- 每次使用模板时要输入什么
- AI 要分几步处理
- 哪些步骤可以并行
- 是否需要汇总
- 中间结果是否展示
- 失败时是否继续

agent 负责把这些回答整理为模板设计稿，再转换为 TemplateSpec。

---

## 固定流程

```text
用户自然语言需求
  -> agent 少量追问
  -> TemplatePlan 设计稿
  -> 用户确认
  -> TemplateSpec JSON
  -> loomloom template-spec check
  -> 创建确认
  -> loomloom template-spec create
```

用户确认设计稿之前，不要创建模板。

注意：用户说“帮我创建一个模板”“创建一个 PRD review 模板”只表示启动创建流程，不等于确认设计稿，也不等于授权创建远端模板。

---

## 三段确认 Gate

对话式创建必须经过三个 gate：

| Gate | 允许做什么 | 禁止做什么 |
|---|---|---|
| 需求收集 | 追问、整理业务流程、输出 TemplatePlan | 生成最终 TemplateSpec、调用 `template-spec create` |
| 设计稿确认 | 用户明确确认 TemplatePlan 后，生成 TemplateSpec 并执行 `template-spec check` | 直接创建模板 |
| 创建确认 | `check` 通过后，用户明确回复“确认创建模板”或等价表达 | 提交批量运行 |

这些表达不算创建确认：

- “帮我创建一个模板”
- “按默认来”
- “继续”
- 只提供 `LOOMLOOM_SERVER` / `LOOMLOOM_TOKEN`
- 只说“生成 spec”

`template-spec create` 前，agent 必须展示：

- 模板名称
- TemplateSpec 文件路径
- `template-spec check` 结果
- 将要执行的创建命令
- 明确提示用户回复“确认创建模板”

---

## 对话规则

agent 必须像模板顾问一样提问，而不是像配置表一样提问。

必须遵守：

1. 一次只问一个问题。
2. 不问技术术语。
3. 用户不确定时，提供默认方案。
4. 用户描述复杂流程时，先复述流程，再继续追问。
5. 不要求用户写 prompt，只让用户描述每一步目标。
6. 先输出 TemplatePlan 设计稿让用户确认。
7. 用户确认后，再生成 TemplateSpec JSON。
8. TemplateSpec check 通过后，再单独请求创建确认。

提问时尽量给用户选项，并提供“按默认”选择。

不要这样问：

```text
是否配置 fieldBindings？
是否声明 upstreamBindings？
是否使用 fan-in？
是否需要 outputSchema？
错误列如何绑定？
```

应该这样问：

```text
这些步骤是依次执行，还是可以同时执行？
这些结果是否需要最后汇总？
中间结果用户要不要看到？
如果某一步失败，要不要继续处理后续步骤？
```

---

## 推荐追问顺序

信息不足时，按下面顺序追问。每次只问一个问题。

1. 你想创建一个什么模板？
2. 用户每次使用模板时，需要提供哪些信息？
3. 你希望 AI 怎么处理这些内容？
4. 过程中有哪些步骤、角色或角度？
5. 中间结果是否需要展示？
6. 最终希望输出哪些结果？
7. 如果某个步骤失败，流程怎么处理？
8. 有没有特别要求？

如果用户已经给出足够信息，不要机械追问全部问题，直接整理设计稿。

示例：

```text
你希望 AI 怎么处理这些内容？
1. 单步处理
2. 连续多步
3. 多个角度分别处理后汇总
4. 生成多个版本
5. 按默认：我根据你的描述判断
```

---

## TemplatePlan 设计稿

TemplatePlan 是 agent 给用户确认的中间设计稿，不是后端 API 主数据。

建议格式：

```json
{
  "templateName": "",
  "goal": "",
  "rowMeaning": "",
  "inputFields": [
    {
      "key": "",
      "label": "",
      "type": "string",
      "required": true,
      "description": ""
    }
  ],
  "steps": [
    {
      "id": "",
      "name": "",
      "executionUnit": "text-generate",
      "goal": "",
      "input": "",
      "output": "",
      "showIntermediateResult": true
    }
  ],
  "relations": [
    {
      "from": ["step_a", "step_b"],
      "to": "step_summary",
      "meaning": "汇总多个上游结果"
    }
  ],
  "outputs": [
    {
      "label": "",
      "sourceStep": "",
      "includeErrorColumn": true
    }
  ],
  "failurePolicy": "默认：单个步骤失败时记录错误原因；能够继续汇总的步骤继续执行，最终结果说明缺失项。",
  "specialRequirements": []
}
```

展示给用户时可以用自然语言表格，不必直接展示 JSON。生成 TemplateSpec 前，agent 可以在内部使用这个结构检查完整性。

---

## 自动补齐规则

用户没有明确说明时，agent 应按以下默认补齐：

| 项目 | 默认规则 |
|---|---|
| 每行代表什么 | 一行代表一次独立批处理任务 |
| 输入字段 | 使用用户描述中的业务名，字段 key 用稳定 snake_case |
| 步骤 prompt | 根据步骤目标生成固定 `instruction`，不要让用户写完整 prompt |
| 输出列 | 每个展示给用户的 step 都有一个结果列 |
| 错误列 | 每个展示给用户的 step 都有一个错误原因列 |
| Step 连接关系 | 根据业务步骤自动整理串行、并行和汇总关系 |
| 输入合并方式 | 多上游汇总默认使用目标 port 的 merge policy，例如文本汇总使用 `concat_text` |
| allow_partial 说明 | 如果允许部分结果继续，在汇聚 step 上设置 `triggerPolicy=allow_partial`；不要依赖自然语言 prompt 承担平台触发策略 |
| 系统错误展示 | Excel 中为用户可见 step 自动保留错误原因列 |
| 业务停止条件 | 如果某步骤可判定业务停止，最终输出中展示停止原因和已完成结果 |
| 中间结果 | 多步骤审核、拆解、生成类模板默认展示；纯内部整理步骤默认不展示 |
| 失败策略 | 默认使用 `require_all`；用户明确接受部分结果继续时使用 `allow_partial`；任一关键步骤失败就应停止时使用 `fail_fast` |
| 模型 | 先运行 `loomloom template-spec models <execution-unit>`，选择可用默认模型 |
| 路由参数 | 不暴露 `provider` 和 `mode` |

用户不需要主动说“帮我加错误原因列”。只要某个 step 的结果会展示给用户，agent 就应自动规划对应错误原因列。

---

## 业务模式映射

### 单步处理

用户说：

```text
输入一段商品描述，生成一版营销文案。
```

映射：

```text
1 个 text-generate step
用户输入字段 -> step.prompt
```

### 串行多步

用户说：

```text
先整理需求，再生成正式文档。
```

映射：

```text
stp_extract -> stp_write
下游 step 通过 upstreamBindings 接收上游 output
```

### 多角度并行后汇总

用户说：

```text
产品、研发、风险三个角度分别评审，最后汇总。
```

映射：

```text
stp_product_review  ↘
stp_engineering_review -> stp_summary
stp_risk_review     ↗
```

实现规则：

- 三个评审 step 是 sibling steps。
- 每个评审 step 都有自己的输入绑定。
- 它们可以绑定同一个输入字段，也可以绑定不同输入字段。
- summary step 使用 `dependsOn` 声明三个上游。
- summary step 使用多个 `upstreamBindings` 接收三个上游 `output`。
- 不要把这个场景表达为一个 step 内部 `expanded`。

TemplatePlan 应明确展示为：

```text
Step 1：产品审核
Step 2：运营审核
Step 3：研发审核
Step 4：结果汇总

输出列：
产品审核结果 / 产品审核错误原因
运营审核结果 / 运营审核错误原因
研发审核结果 / 研发审核错误原因
结果汇总 / 结果汇总错误原因
```

### 生成多个版本

用户说：

```text
根据一个主题生成三种不同风格的文案。
```

映射：

- 如果用户明确要固定三个风格，优先建三个 sibling steps，分别写不同 instruction。
- 如果用户要从一个多值字段动态展开，可以使用 `multiValue=true` + `bindMode=expanded`。

多视角结构化工作流优先使用多个 sibling steps；动态数量展开才使用 `expanded`。

---

## 生成 TemplateSpec 前的检查

生成 TemplateSpec 前，agent 至少确认：

1. TemplatePlan 已经给用户看过。
2. 用户已经确认设计稿。
3. 每个 step 都有清晰目标。
4. 每个 root step 都有用户输入或静态 instruction。
5. 每个下游 step 的输入来源清楚。
6. 多上游汇总使用 step-level `dependsOn` + `upstreamBindings`。
7. 没有暴露 `provider` 或 `mode`。
8. 需要模型时已经查询过 `loomloom template-spec models`。

---

## CLI 验证闭环

生成 spec 后必须执行：

```bash
loomloom template-spec check ./spec.json
```

如果失败，先根据错误修正 spec，不要直接创建模板。

通过后不要立刻创建模板。先向用户展示将要创建的模板和 check 结果，并等待明确确认：

```text
回复“确认创建模板”后，我会执行 template-spec create。
```

用户确认后再执行：

```bash
loomloom template-spec create ./spec.json --version-note "initial version"
```

如果是在修正已有用户模板，追加新版本而不是修改历史版本：

```bash
loomloom template-spec create-version <template-id> ./spec.json --version-note "fix version"
```

创建成功后，建议继续下载 workbook 验证：

```bash
loomloom template-spec download-workbook <template-id> <version-id> --output-file ./input.xlsx
```

提交真实批量运行前，仍必须向用户二次确认。

---

## 当前不支持

第一版对话式创建模板只覆盖 TemplateSpec 当前可表达的工作流。

不要承诺：

- 循环
- 条件分支
- 动态路由
- 长期 memory
- 运行中人工介入
- 任意工具调用 agent
- 任意自定义 execution unit

如果用户提出这些需求，agent 应说明当前版本暂不支持，并给出可用的近似方案。
