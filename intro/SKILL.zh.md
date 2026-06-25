---
name: loomloom
description: 当用户提到 LoomLoom、batchjob、batchflow、批处理、批量处理、模板提交、Excel 模板批量执行、run submit、产物下载、结果回填、Market SkillBot、把模板发布/上架到 Market、运行 Market SkillBot、创作者收益或 Market 使用记录时,使用此 skill。
---

> 说明:本文件是 `skills/claude/loomloom/SKILL.md`(详细版,与 openclaw 一致)的中文译本,仅供参考。实际生效的是英文原版;若两者出现差异,以英文原版为准。

# loomloom

当用户想用 LoomLoom 做批量内容生成、模板提交、Excel 批量执行、结果查看或产物下载时,使用此 skill。`batchjob`、`batchflow`、批处理、批量任务、批量生成、批量跑模板等旧称都指 LoomLoom 工作流。

## 何时使用

- 用户想批量生成文案、图片或视频,包括文生图、文生图生视频的批量工作流。
- 用户想列出模板、下载官方 Excel 模板、校验 Excel、提交 Excel、查看运行状态、下载结果工作簿或下载产物。
- 用户想创作并创建自己的私有模板,尤其是从自然语言描述一个可复用工作流。
- 用户想发现、报价或运行一个已发布的 Market SkillBot。
- 用户想把模板发布到 Market、管理自己的上架,或查看创作者收益和审核请求。
- 任务可以表示为结构化的行,而不是一次性的对话回答。
- 用户愿意使用开发者工具,或让 agent 调用 CLI。

## 何时不使用

- 用户只需要一次即时生成,而不是批处理。
- 需求仍处于探索阶段,无法整理成行级输入或模板。
- 未配置 `LOOMLOOM_TOKEN`,且用户不想处理环境设置。

## 核心对象与命令选择

不要把官方模板、私有模板和 SkillBot 当成三套互不相关的工作流系统。它们的关系是:

```text
官方模板 ── 平台维护、直接执行

私有模板 ── 用户通过 TemplateSpec 创建和维护
   └─ 私有模板版本
        └─ 提交 Market 审核
             └─ Listing Version 不可变发布快照
                  └─ 审核通过后作为 SkillBot 供买家执行
```

术语和选择规则:

- **模板**是总称。
- **官方模板**由平台维护。发现和执行官方模板时使用 `template` 命令组。
- **自定义模板**描述创作方式;创建结果是用户的**私有模板**。创作、查看版本和直接执行私有模板时使用 `template-spec` 命令组。
- **SkillBot**是私有模板某个版本通过 Market 审核后的公开、可付费执行形态。买家使用 `market` 命令组。
- **Listing**是 SkillBot 的 Market 货架对象。创作者发布和管理 SkillBot 时使用 `listing` 与 `creator review` 命令组。
- **Listing Version**是发布时从私有模板版本复制出的不可变执行快照。私有模板后续变化不会自动更新线上 SkillBot。
- 当前没有独立的“公共模板”资源。不要用这个词代指官方模板或 Market SkillBot。
- `asset list` 只是“我的私有模板 + Market SkillBot”的可执行资产聚合视图,不是新的模板类型,也不包含官方模板。

根据用户意图选择入口:

- “有哪些平台模板”“用官方 Excel 模板执行” → `template list/schema/download/...`
- “创建我的工作流”“修改我的模板”“运行我的模板版本” → `template-spec ...`
- “把我的模板发布成 SkillBot”“更新或下架我的 SkillBot” → `listing ...`
- “找一个 SkillBot 并付费运行” → `market ...`
- 用户只说“模板”且上下文无法判断官方还是私有时,先澄清,不要猜测。

## 命令流程

0. 默认通过 GitHub 安装脚本安装。Homebrew 计划启用,但 tap 仓库和发布 Token 尚未配置完成;在配置完成前不要引导用户使用 Homebrew。Gitee 是否继续保留仍待确认;在 owner、repo 和安装地址明确前不要主动推荐 Gitee 分发。
   如果用户要求内测/beta 版 CLI,显式安装预发布通道,而不是默认 stable:
   `curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta --no-brew`
1. 检查环境:
   `loomloom doctor`
   CogFoundry 正式环境的 LoomLoom API 基础地址是 `https://loomloom.cogfoundry.ai/loom/v1`。CogFoundry 正式 Token 只能发送到主机 `loomloom.cogfoundry.ai`,并且必须使用 HTTPS。
   在携带正式 Token 发起任何请求前,检查最终目标 URL 的 scheme 和 host。若 host 不是 `loomloom.cogfoundry.ai`,立即停止,不得发送 Token,也不得自动跟随到其他域名。
   测试环境、本地地址、私有部署或第三方代理必须使用对应环境单独签发的专用 Token。不得把 CogFoundry 正式 Token 复用于这些地址。即使用户明确提供其他服务地址,也不能把正式 Token 发送过去。
   若缺少 token,引导用户前往 `https://console.cogfoundry.ai/api-keys` 创建或获取。不得在回复、日志或生成文件中回显真实 token。
2. 不要把大文件粘贴进上下文。先上传原始输入素材:
   `loomloom input-asset upload <file>`
3. 发现当前环境的官方模板:
   `loomloom template list`
4. 查看官方模板字段结构、列出可用模型,或列出可执行资产:
   `loomloom template schema <template-id>`
   `loomloom model list --step-type <text-generate|image-generate|video-generate>`
   `loomloom asset list`
   注意:`asset list` 聚合我的私有模板与 Market SkillBot,不包含官方模板。
5. 创作并保存私有模板时,使用 TemplateSpec JSON:
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
   查看已有的私有模板:
   `loomloom template-spec list`
   `loomloom template-spec get <template-id>`
   `loomloom template-spec versions <template-id>`
   要直接用扁平 JSONL 行运行一个私有模板版本,先上传行数据,再把返回的 input_file_id 传入:
   `loomloom orchestration-input upload <file.jsonl>`
   `loomloom template-spec run <template-id> --version-id <version-id> --input-file-id <input_file_id>`
   对于常见的单根工作流,每个非空行可以是一个所有值均为字符串的扁平 JSON 对象。掌握准确工作流步骤映射时,也可以使用 `steps.<step-id>.executions[]` 形式的 unified 行。两种格式中的执行参数值都必须是字符串,并符合私有模板版本允许的输入参数。不得猜测 step ID。
6. 如果用户没有明确要求创建或使用自己的私有模板,默认走官方 Excel 工作流:
   `loomloom template download <template-id>`
   `loomloom template validate-file <template-id> <xlsx-path>`
   `loomloom template submit-file <template-id> <xlsx-path>`
   `loomloom run result-workbook <run-id>`
   只有在用户明确需要时才使用旧的本地回填流程:
   `loomloom template backfill-results <run-id> <xlsx-path>`
7. 只有在用户明确要求编程式输入时才用 JSON/JSONL:
   `loomloom run submit <template-id> -f rows.jsonl`
8. 监控进度:
   `loomloom run watch <run-id>`
9. 列出、查看或下载运行及其结果:
   `loomloom run list`
   `loomloom run get <run-id>`
   `loomloom run result-rows <run-id>`
   `loomloom run result-workbook <run-id>`
10. 查看或下载产物:
   `loomloom artifact list <run-id>`
   `loomloom artifact download <run-id>`

## Market 生态

LoomLoom Market 让创作者把自己的私有模板版本发布为付费 SkillBot,也让买家运行这些 SkillBot。同一个 CLI 服务两种角色;选择命令前先判断用户处于哪种角色。

发布时,Market 从指定私有模板版本复制不可变的 Listing Version 执行快照。创作者之后修改私有模板不会自动改变线上 SkillBot;要更新线上执行版本,必须对同一 Listing 提交新的模板版本并重新审核。

### 买家角色(使用并为 SkillBot 付费)

- `loomloom market list` —— 浏览已发布的 SkillBot。支持 `--keyword`、`--page-size`、`--page-token`、`--order-by`。
- `loomloom market show <listing-id>` —— 展示单个 SkillBot,包含其 `inputSchemaSnapshot`。构造输入前先读取该 schema。
- `loomloom market quote <listing-id> --input-file <request.json>` —— 估算成本。返回 `estimatedBuyerPayableT`、`taskCount`、`taskFixedFeeT`。
- `loomloom market run <listing-id> --input-file <request.json> --confirm` —— 执行 SkillBot。这是付费操作;参见“提交确认规则”。
- `loomloom usage list` 和 `loomloom usage get <run-transaction-id>` —— 查看买家自己的 SkillBot 调用与结算状态。使用返回的 `runTransactionId`。

quote 和 run 的 `--input-file` JSON 携带一个 `taskInputs` 数组,其形状对应该 listing 的输入 schema,例如:

```json
{
  "listingVersionId": "<listing-version-id>",
  "taskInputs": [
    {
      "steps": {
        "<step-id>": {
          "executions": [
            {
              "prompt": "write a launch tweet"
            }
          ]
        }
      }
    }
  ]
}
```

先运行 `market show` 了解公开字段并取得 `listingVersionId`,但不要从 `inputSchemaSnapshot` 推断内部 step ID:公开快照可能不包含它。当前 CLI 要求精确的 Product API `taskInputs` 请求体;若没有兼容映射或请求 JSON,应停止并向用户索要,不得猜测。

### 创作者角色(发布并管理 SkillBot)

- `loomloom listing publish <template-id> --template-version-id <id> --display-name <name> --task-fixed-fee-t <fee>` —— 提交一个模板版本进入 Market 审核。该版本必须已至少有一次成功运行。返回带 `reviewStatus: pending_review` 的 `reviewRequestId`。
- 要为已有 listing 提交新版本,使用同一命令并增加 `--listing-id <listing-id>`,同时传入新的 `--template-version-id`。审核通过前,当前已发布版本继续生效。
- `loomloom listing list`、`loomloom listing show <listing-id>`、`loomloom listing versions <listing-id>` —— 查看创作者自己的上架,包括待审、被拒、已下架状态。
- `loomloom listing update <listing-id> --display-name <name> --description <text>` —— 提交公开资料变更进入审核;不改变定价或执行版本。
- `loomloom listing unlist <listing-id>` 和 `loomloom listing relist <listing-id>` —— 停止或恢复已发布上架的新执行。
- `loomloom listing withdraw <listing-id>` —— 撤回该 listing 唯一的待审核请求。若没有待审请求,停止操作;若返回多个,先用 `creator review list` 找到目标请求,再按 ID 明确撤回。
- `loomloom creator review list`、`loomloom creator review get <review-request-id>`、`loomloom creator review withdraw <review-request-id>` —— 直接跟踪和撤回审核请求。
- `loomloom creator earnings` 和 `loomloom creator transactions` —— 查看创作者收入与每次调用的结算。

所有 `*FeeT`、`*AmountT`、`*PayableT` 值都以 API 单位计,其中 10,000,000 单位等于 1 个货币单位。

## 提交确认规则

任何真正向托管 LoomLoom 服务提交工作的命令,都必须在当前对话中获得用户的第二次明确确认。

把交互视为三种状态之一:

1. `default-prep`:用户仍在探索或泛泛而谈。只做准备、下载、校验。不要提交。
2. `auto-run-candidate`:用户明确要求 agent 执行。仍不要提交。先给出执行摘要并等待确认。
3. `confirmed-to-run`:用户看到执行摘要后明确确认。只有此时才可以提交。

此规则适用于 `loomloom template submit-file`、`loomloom template-spec submit-workbook`、`loomloom run submit`、`loomloom template-spec run`、`loomloom market run`。对于 `loomloom market run`,先运行 `loomloom market quote` 并把返回的报价纳入执行摘要。校验、下载、schema 查看、模型查询、报价、`doctor`、素材上传、orchestration-input 上传、产物列举、listing/usage/earnings 读取以及结果回填都不会启动新的付费运行,不需要第二次确认。

执行摘要必须包含模板或 listing ID、输入来源、行数或任务规模、动作、预估成本或明确的成本说明,以及提示语 `Reply "confirm submit" before I start.` 对于 `template submit-file`、`template-spec submit-workbook`、`run submit`、`template-spec run` 和 `market run`,应显式传入稳定的 `--client-request-id`,并保存它,以便安全重试完全相同的 payload。

## 远程状态变更确认规则

创建或修改持久化远程资源前,必须展示确切动作并取得明确确认。此规则适用于:

- `template-spec create` 和 `template-spec create-version`
- `listing publish`(包括带 `--listing-id` 的换版本)
- `listing update`
- `listing unlist` 和 `listing relist`
- `listing withdraw` 和 `creator review withdraw`

只读命令、本地检查、上传、下载和报价不需要此确认。付费运行仍遵循上面的更严格提交确认规则。

## Agent 命令串联

当一条命令的结果要交给下一条命令时,优先使用 `--output json`,并原样保存字段:

- `orchestration-input upload` → `inputFileId` → `template-spec run --input-file-id`
- 提交运行 → `runId` → `run watch` 或结果命令
- `listing publish` → `reviewRequestId` → creator review 命令
- `market run` → `runTransactionId` 和 `runId` → usage/run 命令

不得把 `inputAssetId`(`ia_xxx`)转换为 `inputFileId`,也不得根据名称猜测 ID。
对于“提交确认规则”中列出的五个支持该参数的提交命令,显式传入 `--client-request-id`,将它和请求一起保存,且仅在重试完全相同的 payload 时复用。payload 变化时必须使用新的 ID。

## 错误处理

先检查命令返回的错误,再选择恢复方式:

- 本地 flag、文件、JSON 或 schema 错误:修正输入后重试,无需运行 `doctor`。
- 鉴权、endpoint、网络、服务版本或意外服务端错误:运行 `loomloom doctor`。
- 不得虚构缺失的 ID、隐藏的 step ID 或服务端状态。
- 结果不明确的失败后,不要盲目重试付费命令或远程状态变更命令。先查询相关 run、listing 或 review 状态。重试完全相同的提交时复用原 `--client-request-id`;payload 改变时使用新的 ID。

## 当前 MVP 能力

当前公开 CLI 支持以下命令组:

- 环境:`doctor`。
- 输入:`input-asset upload`、`orchestration-input upload`。
- 发现:`template list`、`template schema`、`model list`、`asset list`。
- 官方 Excel 工作流:`template download`、`template validate-file`、`template precheck-file`、`template submit-file`、`template backfill-results`。
- 私有模板创作与执行:`template-spec docs`、`template-spec check`、`template-spec models`、`template-spec create`、`template-spec create-version`、`template-spec list`、`template-spec get`、`template-spec versions`、`template-spec download-workbook`、`template-spec validate-workbook`、`template-spec submit-workbook`、`template-spec run`。
- 运行:`run submit`、`run list`、`run get`、`run watch`、`run result-rows`、`run result-workbook`。
- 产物:`artifact list`、`artifact download`。
- Market(买家):`market list`、`market show`、`market quote`、`market run`、`usage list`、`usage get`。
- Market(创作者):`listing publish`、`listing list`、`listing show`、`listing versions`、`listing update`、`listing unlist`、`listing relist`、`listing withdraw`、`creator earnings`、`creator transactions`、`creator review list`、`creator review get`、`creator review withdraw`。

## 大文件处理

当用户想批量处理本地代码文件、大段文本、本地图片或其他大文件时,避免把整个文件粘贴进 agent 上下文。优先:

1. `loomloom input-asset upload <file>`
2. 保存返回的 `input_asset_id`
3. 只把它填入 schema 明确接受素材引用的字段

对于 TemplateSpec,使用兼容的 `asset_ref` 或 `text_reference` 字段,并遵循内置绑定说明。`input_asset_id` 是参考素材 ID,绝不能作为 `template-spec run` 的 `inputFileId`。

## 默认行为

除非用户明确要求 JSON/JSONL,默认使用官方 Excel 模板工作流。

当用户要求创建或定制自己的工作流/模板时,优先 TemplateSpec JSON。创建结果是私有模板。TemplateSpec JSON 是源数据;下载的工作簿是派生产物。当模板版本变化时,不要承诺旧工作簿仍然兼容。下载新工作簿。

## 对话式模板创作

当用户用自然语言描述新模板时,不要立刻写 TemplateSpec JSON。先运行或参考 `loomloom template-spec docs conversation`。

流程:

1. 问业务问题,不要问 TemplateSpec 技术字段问题。
2. 一次只问一个缺失细节。
3. 避免面向用户的技术术语,如 `fieldBindings`、`upstreamBindings`、`fan-in`、`execution`、`outputSchema`、`provider`、`mode`。
4. 继续之前,用业务语言复述复杂工作流。
5. 识别出工作流步骤、角色、视角或生成的 agent 后,判断每个角色/步骤是否有各自的处理要求。如果有,在起草 TemplatePlan 前先询问模板未来的使用方式。
6. 先起草 TemplatePlan。
7. 展示 TemplatePlan 并等待用户确认。
8. 仅在确认后才生成 TemplateSpec JSON。
9. 创建前,运行 `loomloom template-spec check <spec.json>`。
10. check 通过后,单独再要一次创建确认。
11. 仅在明确的创建确认后,才运行 `loomloom template-spec create <spec.json>`。

当已有的用户模板需要修复时,不要承诺原地修改历史版本。用 `loomloom template-spec create-version <template-id> <spec.json>` 追加一个新版本,并让后续工作簿/运行使用新的 `version_id`。

### 模板使用模式选择

当模板有多个角色、步骤、评审视角或生成的 agent,且每个角色/步骤可能有各自的处理要求时,agent 必须在生成 TemplatePlan 前询问:

```text
当别人以后使用此模板时,评审/生成要求应如何处理?

1. 预置在模板里:用户只填核心材料,系统自动遵循你预定义的要求。
2. 让用户填写:用户可以为每个步骤/角色填写或修改要求。
3. 生成两个版本:一个简单版和一个可定制版。

如果不确定,先选 1,做一个简单可用的标准模板。
```

不要暴露技术概念,如 `prompt`、`binding`、`reference`、`field`、`hidden`、`paramBindings`、`fieldBindings`。典型触发场景包括多视角 PRD 评审、多角色合同评审、多 agent 发布活动策划、把一篇文章改写为多种风格、从多个面试官视角评审简历,以及为多个渠道生成营销内容。

三种模式:

1. 简单使用模式:用户只填核心材料,系统用预置的处理要求执行。
2. 自定义使用模式:用户既填核心材料,也填每个步骤/角色的处理要求。
3. 生成两个版本:生成两个 TemplatePlan 并创建两个模板,命名为 `Standard Version` / `Custom Version`。

自动生成规则:简单模式下,角色/步骤的处理要求是模板预置指令,不是用户可填列;自定义模式下,核心材料和角色/步骤的处理要求都是用户输入列;生成两个版本时,先展示两个 TemplatePlan,获得用户确认后再生成两个 TemplateSpec。

验收标准:多角色/多步骤/多视角的模板创建必须询问使用模式;简单模式隐藏角色/步骤要求列;自定义模式暴露并允许编辑这些列;面向用户的对话不暴露技术概念;生成的 TemplateSpec 能正确表达简单和自定义两种形态。

创建确认门槛:

- “创建一个 PRD 评审模板”只是开始流程;它不构成对远程创建的确认。
- 环境变量、token、服务 URL 是配置,不是创建确认。
- “生成 spec”只意味着生成并本地 check spec;它不构成对运行 `template-spec create` 的确认。
- 运行 `template-spec create` 之前,展示模板名、spec 路径、check 结果、确切的创建命令,并请用户回复 `confirm create template`。

TemplatePlan 应覆盖:模板名称与目标、行的含义、输入字段、工作流步骤、串行/并行/汇总关系、模板使用模式、用户可见输出、失败策略、错误列、默认模型,以及特殊要求。

默认建模规则:

- “产品、工程、风险分别评审,然后汇总”应建模为多个并行的 `text-generate` 步骤加一个下游汇总步骤,使用步骤级 `dependsOn` 和 `upstreamBindings`。
- 不要把多角色评审建模为一个 `expanded` 步骤。`expanded` 只用于动态多值输入。
- 默认为每个用户可见步骤添加结果列和错误列。
- 如果允许部分完成,汇总步骤应说明缺失的上游并暴露失败步骤的错误列。
- 保持 `provider` 和 `mode` 为内部字段。

TemplateSpec 创作约束:

- 编写用于创建私有模板的 TemplateSpec 前,运行 `loomloom template-spec docs spec`,以 CLI 内置文档作为当前契约。
- 用 `loomloom template-spec docs examples` 查看示例和模式。
- 用 `loomloom template-spec docs conversation` 查看自然语言创建流程。
- 默认只用 `text-generate`、`image-generate`、`video-generate`。
- 使用 OpenAPI lowerCamel 字段,如 `meta.name`、`steps[].stepId`、`defaultModelRef.modelKey`。
- 用步骤级 `dependsOn` 和 `upstreamBindings` 连接步骤;源输出端口通常是 `output`。
- 选择 `defaultModelRef.modelKey` 前,运行 `loomloom template-spec models <execution-unit>` 并使用返回的 `model_id`。
- 仅当步骤设置 `allowModelOverride=true` 且某字段绑定到 `paramKey=model` 时,才暴露模型列。
- 不要绑定 `provider` 或 `mode`。

## 可信结果来源

提交的工作簿和服务端运行输入快照是事实来源。运行完成后,优先用 `loomloom run result-workbook <run-id>`,因为服务端会对齐原始输入快照和产物。只有用户明确需要旧的本地 Excel 回填流程时,才用 `template backfill-results`。

## 控制台访问

CogFoundry Console 入口是 `https://console.cogfoundry.ai`。用户需要查看运行状态时,可以提供该 Console 首页。

当前暂不提供 Workflow Run 详情页 URL 模板。不要根据 `runId` 猜测或拼接详情页链接;如果 CLI 返回了服务端明确提供的 URL,可以原样使用,否则只提供 Console 首页和 CLI 查询命令。

CogFoundry 官网是 `https://cogfoundry.ai`。
