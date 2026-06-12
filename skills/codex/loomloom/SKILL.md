---
name: loomloom
description: 当用户提到 LoomLoom、batchjob、batchflow、批处理、批量处理、模板提交、Excel 模板批量执行、run submit、产物下载或结果回填时使用。
---

# loomloom

当用户想通过 LoomLoom 执行批量内容生成、模板提交、Excel 批量执行、结果查看或产物下载时，使用本 skill。历史名称 `batchjob`、`batchflow`、批处理、批量任务、批量生成、批量跑模板等，都按 LoomLoom 工作流理解。

## 适用场景

- 用户要批量生成文案、图片、视频，或执行文本到图片、文本到图片再到视频的批量流程。
- 用户要查看模板、下载官方 Excel 模板、校验 Excel、提交 Excel、查看 run 状态、下载结果工作簿或下载产物。
- 用户要创建自定义模板，尤其是通过自然语言描述一个可复用工作流。
- 任务可以表达为多行结构化输入，而不是一次性的聊天回答。
- 用户愿意使用开发者工具或由 Agent 代为调用 CLI。

## 不适用场景

- 用户只需要一次即时生成，不需要批处理。
- 需求还停留在纯探索，无法整理成行级输入或模板。
- `LOOMLOOM_TOKEN` 没有配置，并且用户不希望处理环境设置。

## 命令流程

0. 如果用户要求安装内测 CLI，必须显式安装 beta 渠道，不要默认安装 stable：
   `curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta --no-brew`
1. 检查环境：
   `loomloom doctor`
   不要依赖本机默认 server。TemplateSpec 编写、模型发现或模板提交前，先确认 `LOOMLOOM_SERVER` 是用户明确提供的目标环境；如果未配置，应要求用户提供服务地址。
2. 大文件不要贴进上下文，先上传原始输入资产：
   `loomloom input-asset upload <file>`
3. 发现可用模板：
   `loomloom template list`
4. 查看模板结构：
   `loomloom template schema <template-id>`
5. Agent 编写自定义模板时使用 TemplateSpec JSON：
   `loomloom template-spec docs spec`
   `loomloom template-spec docs examples`
   `loomloom template-spec docs conversation`
   `loomloom template-spec models <text-generate|image-generate|video-generate>`
   `loomloom template-spec check <spec.json>`
   `loomloom template-spec create <spec.json>`
   `loomloom template-spec create-version <template-id> <spec.json>`
   `loomloom template-spec download-workbook <template-id> <version-id>`
   `loomloom template-spec validate-workbook <template-id> <version-id> <xlsx-path>`
   `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx-path>`
6. 用户没有明确要求自定义模板时，默认走官方 Excel 工作流：
   `loomloom template download <template-id>`
   `loomloom template validate-file <template-id> <xlsx-path>`
   `loomloom template submit-file <template-id> <xlsx-path>`
   `loomloom run result-workbook <run-id>`
   只有用户明确需要旧的本地回填流程时，才使用：
   `loomloom template backfill-results <run-id> <xlsx-path>`
7. 只有用户明确要求程序化输入时，才走 JSON/JSONL：
   `loomloom run submit <template-id> -f rows.jsonl`
8. 查看执行进度：
   `loomloom run watch <run-id>`
9. 查看或下载对齐后的 run 结果：
   `loomloom run result-rows <run-id>`
   `loomloom run result-workbook <run-id>`
10. 查看或下载产物：
   `loomloom artifact list <run-id>`
   `loomloom artifact download <run-id>`

## 提交确认规则

任何会真正提交任务到托管 BatchJob 服务的命令，都必须在当前对话中取得用户第二次明确确认。

交互按三种状态处理：

1. `default-prep`
   用户还在探索或泛泛表达，例如“帮我跑个批处理”“帮我跑一下”“按你来”。
   此时只准备、下载、校验，不提交。

2. `auto-run-candidate`
   用户明确要求 Agent 执行，例如“你直接帮我自动跑”“直接帮我提交”“你替我跑”“你替我执行”。
   此时仍不要提交，先给出执行摘要并等待确认。

3. `confirmed-to-run`
   用户在看到执行摘要后明确回复“确认提交”“提交吧”“开始跑”“继续执行”等。
   只有此时才可以提交。

该规则适用于：

- `loomloom template submit-file <template-id> <xlsx-path>`
- `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx-path>`
- `loomloom run submit <template-id> -f rows.jsonl`

校验、下载、schema 查看、模型查询、doctor、资产上传、产物列表、结果回填不启动新的付费 run，不需要第二次确认。

执行摘要必须包含：

- template ID
- 输入文件路径或输入来源
- 行数或任务规模
- 将要执行的动作
- 预计费用；如果只能提交后得知费用，要明确说明
- 明确提示：`回复“确认提交”后我才会开始执行`

如果用户说“先别跑”“等一下”等暂停信号，继续停留在准备状态。

## 错误处理

命令失败后，先运行：

`loomloom doctor`

用 doctor 结果先区分：

- 本地环境配置问题
- CLI 版本过旧
- 服务端行为或模型目录问题

不要直接猜测模板、模型或 run 本身有问题。

## 当前 MVP 能力

公开 CLI 当前支持：

- `loomloom doctor`
- `loomloom input-asset upload <file>`
- `loomloom template list`
- `loomloom template schema <template-id>`
- `loomloom template download <template-id>`
- `loomloom template validate-file <template-id> <xlsx-path>`
- `loomloom template submit-file <template-id> <xlsx-path>`
- `loomloom template backfill-results <run-id> <xlsx-path>`
- `loomloom template-spec check <spec.json>`
- `loomloom template-spec docs [spec|authoring|examples|conversation|all]`
- `loomloom template-spec models <text-generate|image-generate|video-generate>`
- `loomloom template-spec create <spec.json>`
- `loomloom template-spec create-version <template-id> <spec.json>`
- `loomloom template-spec download-workbook <template-id> <version-id>`
- `loomloom template-spec validate-workbook <template-id> <version-id> <xlsx-path>`
- `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx-path>`
- `loomloom run submit <template-id> -f rows.jsonl`
- `loomloom run watch <run-id>`
- `loomloom run result-rows <run-id>`
- `loomloom run result-workbook <run-id>`
- `loomloom artifact list <run-id>`
- `loomloom artifact download <run-id>`

## 大文件处理

当用户要批量处理本地代码文件、大文本、本地图片或其他大文件时，避免把文件全文贴进 Agent 上下文。优先：

1. `loomloom input-asset upload <file>`
2. 保存返回的 `input_asset_id`
3. 再分步准备结构化 JSONL 或 Excel 输入

当前阶段只覆盖上传；结构化输入中引用 `input_asset_id` 是后续能力。

## 默认行为

除非用户明确要求 JSON/JSONL，默认使用官方 Excel 模板工作流。

当用户要求创建或定制 workflow/template 时，优先使用 TemplateSpec JSON。TemplateSpec JSON 是源数据，下载的 workbook 是派生产物。模板版本变化后，不要承诺旧 workbook 仍兼容，应重新下载新 workbook。

## 对话式模板创建

当用户用自然语言描述新模板时，不要直接写 TemplateSpec JSON。先运行或参考：

`loomloom template-spec docs conversation`

流程：

1. 问业务问题，不问 TemplateSpec 技术字段。
2. 信息缺失时一次只问一个问题。
3. 面向用户时避免使用 `fieldBindings`、`upstreamBindings`、`fan-in`、`execution`、`outputSchema`、`provider`、`mode` 等技术词。
4. 复杂流程先用业务语言复述，再继续追问。
5. 识别出 workflow steps、角色、视角或生成 agent 后，判断是否存在“每个角色/步骤各自的处理要求”。如果存在，必须先询问模板未来的使用方式，再起草 TemplatePlan。
6. 先起草 TemplatePlan。
7. 展示 TemplatePlan，等待用户确认。
8. 用户确认后再生成 TemplateSpec JSON。
9. 创建前必须运行 `loomloom template-spec check <spec.json>`。
10. check 通过后，再单独请求创建确认。
11. 只有用户明确确认创建模板后，才运行 `loomloom template-spec create <spec.json>`。

已有用户模板需要修正时，不要承诺原地更新历史版本。使用 `loomloom template-spec create-version <template-id> <spec.json>` 追加新版本，并让后续 workbook / run 使用新的 `version_id`。

提问时尽量给出选项，并提供“按默认”选择。

### 模板使用方式选择

当模板中存在多个角色、步骤、审核视角、生成 agent，且每个角色/步骤都可能有自己的处理要求时，Agent 必须在生成 TemplatePlan 之前询问：

```text
这个模板未来给别人使用时，你希望审核/生成要求怎么处理？

1. 预设在模板里：别人只填核心材料，系统按你设定好的处理要求自动执行。
2. 让使用者填写：别人可以自己填写或修改每个步骤/角色的处理要求。
3. 两个版本都生成：一个简单版，一个可自定义版。

如果不确定，建议选 1，先做一个简单可用的标准模板。
```

不要在这个问题里暴露 `prompt`、`binding`、`reference`、`field`、`hidden`、`paramBindings`、`fieldBindings` 等技术概念。用户面对的是“核心材料”“处理要求”“标准版”“自定义版”。

必须触发该问题的典型场景：

- PRD 多视角审核
- 合同多角色审核
- 新品发布会多 agent 策划
- 文章按不同风格改写
- 简历从多个面试官视角评审
- 营销内容生成多个渠道版本

三种模式的业务含义：

1. 简单使用模式
   使用者只填写核心材料，系统按模板创建者预设好的处理要求自动执行。适合标准化流程、团队统一审核标准、低门槛使用。

2. 自定义使用模式
   使用者不仅填写核心材料，也可以自己填写或修改每个步骤/角色的处理要求。适合模板分享、不同团队复用、不同场景下调整审核标准。

3. 同时生成两个版本
   生成两个 TemplatePlan，并分别创建两个模板。命名上区分 `标准版` / `自定义版`，例如 `PRD 四视角审核 - 标准版` 和 `PRD 四视角审核 - 自定义版`。

以 PRD 四视角审核为例：

- 简单使用模式：使用者只看到 `PRD 内容`；运营、产品、研发、市场四个审核标准由模板作者预设，使用者不可见。
- 自定义使用模式：使用者看到 `PRD 内容`、`运营审核提示词`、`产品审核提示词`、`研发审核提示词`、`市场审核提示词`，可以按自己的团队标准调整。
- 两个版本都生成：标准版面向普通用户，自定义版面向高级用户或团队负责人。

自动生成规则：

- 用户选简单使用模式时：核心材料是用户输入列；角色/步骤处理要求写入模板预设说明，不作为使用者填写列。
- 用户选自定义使用模式时：核心材料和角色/步骤处理要求都作为用户输入列；系统稳定合成每个步骤需要的完整处理输入。
- 用户选两个版本都生成时：必须先展示两个 TemplatePlan，并让用户确认后再分别生成两个 TemplateSpec；创建远程模板时仍然逐个走“确认创建模板”门禁。

验收标准：

- 多角色、多步骤、多视角模板创建时，Agent 会询问“使用者是否可以修改每个角色/步骤的处理要求”。
- 用户可以选择只填核心材料的简单模式。
- 用户可以选择核心材料加自定义处理要求的自定义模式。
- 用户可以选择同时生成两个版本。
- 简单模式下，使用者不会看到各角色/步骤处理要求输入列。
- 自定义模式下，使用者能看到并修改各角色/步骤处理要求输入列。
- 面向用户的对话不暴露 `prompt`、`binding`、`reference`、`field`、`hidden` 等技术概念。
- 生成的 TemplateSpec 能正确表达简单模式和自定义模式两种模板形态。

创建确认门禁：

- “创建一个 PRD 审核模板”只表示进入创建流程，不等于确认远程创建。
- 环境变量、token、server URL 只是配置信息，不等于确认创建。
- “生成 spec”只表示生成并本地 check，不等于确认运行 `template-spec create`。
- 运行 `template-spec create` 前，必须展示模板名、spec 路径、check 结果、准确创建命令，并要求用户回复“确认创建模板”。

TemplatePlan 应覆盖：

- 模板名和目标
- 一行 workbook 代表什么
- 用户输入字段
- workflow steps 和每个 step 的目标
- 串行、并行、汇总关系
- 模板使用方式：简单使用模式、自定义使用模式，或同时生成两个版本
- 用户可见的中间输出
- 最终输出
- 失败策略
- 业务停止条件
- Excel 中的系统错误展示
- 默认模型选择
- 特殊要求

默认建模规则：

- “产品、研发、风险分别审核，然后汇总”应建模为多个并行 `text-generate` step 加一个下游汇总 step，使用 step 级 `dependsOn` 和 `upstreamBindings`。
- 不要把多角色审核建模成一个 `expanded` step；`expanded` 只用于动态多值输入。
- TemplatePlan 中应把多角色审核列成多个 step，而不是“一个会 fan-out 的审核 step”。
- 默认给每个用户可见 step 增加结果列和错误列。
- 用户不需要主动要求错误列；每个用户可见 step 都默认增加。
- 如果允许部分完成，汇总 step 应解释缺失上游，并暴露失败 step 的错误列。
- `provider` 和 `mode` 保持内部，不暴露给模板用户。

TemplateSpec 编写约束：

- 编写自定义 TemplateSpec 前，先运行 `loomloom template-spec docs spec`，以 CLI 内置文档作为当前契约。
- 示例和模式参考 `loomloom template-spec docs examples`。
- 自然语言创建流程参考 `loomloom template-spec docs conversation`。
- 已安装 skill 中也有 `docs/template-spec/` 备份，但优先使用 CLI docs 命令，因为它反映当前 CLI 携带的文档。
- execution unit 默认只使用 `text-generate`、`image-generate`、`video-generate`，除非用户有明确文档化的自定义 unit。
- 使用 OpenAPI 的 lowerCamel 字段，如 `meta.name`、`steps[].stepId`、`defaultModelRef.modelKey`。
- step 之间使用 step 级 `dependsOn` 和 `upstreamBindings`；source output port 通常是 `output`。
- 只有目标 input port 允许多来源时，才能把多个上游汇入同一个 input port：必须满足 `AllowMultiple=true`、`Accepts` 类型匹配，并且有明确 `MergePolicy`。
- 选择 `defaultModelRef.modelKey` 前，必须运行 `loomloom template-spec models <execution-unit>`，并使用返回的 `model_id`。
- 只有当 step 设置 `allowModelOverride=true` 且字段绑定到 `paramKey=model` 时，才暴露模型列。
- 不要绑定 `provider` 或 `mode`，这些路由控制不通过模板暴露。

## 结果可信来源

使用以下命令时：

- `loomloom template submit-file <template-id> <xlsx-path>`
- `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx-path>`
- `loomloom run result-workbook <run-id>`

提交的 workbook 和服务端 run input snapshot 是事实来源。run 完成后优先使用 `run result-workbook`，因为服务端会把原始输入快照和产物对齐。只有用户明确需要旧的本地 Excel 回填时，才使用 `template backfill-results`。

## 控制台访问

返回给用户的控制台链接应来自用户提供的 CogFoundry workspace 信息。不要默认使用本机或历史服务地址。

提交、watch 或 backfill 成功后，可以补充：

- 如果用户提供了 CogFoundry 控制台地址，提示用户可在控制台查看 run 状态。
- 如果已有 `run_id`，提示用户可在页面搜索或查看最新记录。
