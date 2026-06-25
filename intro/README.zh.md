# LoomLoom

> **批量内容生成平台** —— 用自然语言编排 AI 工作流,生成文案、图片和视频。
> 由 CogFoundry 构建 —— [github.com/Cogfoundry-ai/loomloom](https://github.com/Cogfoundry-ai/loomloom)

> 说明:本文件是 `README.md` 的中文译本,仅供参考。实际生效的是英文原版;若两者出现差异,以英文原版为准。

---

## 它是什么

LoomLoom 是一个 CLI 和 agent skill 包,用于 AI 驱动的批量内容工作流。不用手写工作流代码,只需用自然语言描述任务,让 agent 下载模板、准备数据、提交运行、监控进度并下载结果。

常见用例:

- **批量文案** —— 商品描述、改写、摘要、问答,以及文件级文本修改。
- **批量图片生成** —— 电商图、社交素材、概念图,以及逐行视觉生成。
- **批量视频生成** —— 分镜、广告素材,以及文生视频工作流。

---

## 支持的 Agent

安装 LoomLoom 会为所选 agent 添加配套的 skill 包:

| Agent | 状态 |
| --- | --- |
| **Codex**(OpenAI) | 支持 |
| **Claude Code**(Anthropic) | 支持 |
| **OpenClaw** | 支持 |

---

## 快速开始

### Agent 辅助安装

向 Codex、Claude Code 或 OpenClaw 发送类似下面的消息。把 `your-token` 替换为你在 CogFoundry API Keys 页面创建的 token。

```text
Install LoomLoom from this GitHub repository: https://github.com/Cogfoundry-ai/loomloom
My server URL is https://loomloom.cogfoundry.ai/loom/v1, and my token is your-token.
After installation, run doctor to check whether the setup is healthy.
```

### 手动安装

macOS / Linux:

```bash
# 默认安装,带 Codex skill 包
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash

# Claude Code
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --agent claude

# OpenClaw
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --agent openclaw

# 指定版本
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --version v0.2.7

# 最新 beta 或 internal 通道
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta --no-brew

# 指定预发布 tag
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --version v0.2.6-beta.9 --no-brew
```

Windows PowerShell:

```powershell
# 默认安装
irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1 | iex

# Claude Code
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Agent claude

# OpenClaw
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Agent openclaw

# 最新 beta 或 internal 通道
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Channel beta

# 指定预发布 tag
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Version v0.2.6-beta.9
```

Homebrew 分发计划启用,但 tap 仓库和发布所需 Token 尚未完成配置。目前请使用上面的安装脚本。

是否继续保留 Gitee 分发仍待确认。在 owner、repo 和安装地址确定前,不建议使用或传播 Gitee 安装命令。

---

## 配置凭据

```bash
export LOOMLOOM_SERVER="https://loomloom.cogfoundry.ai/loom/v1"
export LOOMLOOM_TOKEN="<your CogFoundry token>"
```

如果不想每个会话都重新设置,把这些值加到 `~/.zshrc`、`~/.bashrc` 或你的 shell profile。CLI 仍兼容旧的 `BATCHJOB_SERVER` 和 `BATCHJOB_TOKEN` 变量,但新配置应使用 `LOOMLOOM_*`。

在 [CogFoundry API Keys](https://console.cogfoundry.ai/api-keys) 页面创建或获取 token。不要在文档、截图、日志或公开对话中暴露真实 token。

> 安全要求:CogFoundry 正式 Token 只能发送到 `https://loomloom.cogfoundry.ai`。正式 API 基础地址是 `https://loomloom.cogfoundry.ai/loom/v1`。不得把正式 Token 发送到其他域名、测试环境、本地地址或第三方代理;使用其他环境时必须使用该环境单独签发的专用 Token。

---

## 验证安装

```bash
loomloom doctor
```

如果环境健康,你就可以开始提交模板运行了。

---

## 模板、私有模板与 SkillBot 的关系

“模板”是 LoomLoom 中可重复执行的 AI 工作流总称。当前产品中需要区分两类模板和一条 Market 发布链路:

```text
模板
├─ 官方模板
│  ├─ 由平台维护并发布
│  ├─ 所有有权限的用户可发现和执行
│  └─ 使用 loomloom template ... 命令
│
└─ 私有模板
   ├─ 用户通过 TemplateSpec 创建和维护
   ├─ 一个模板可以有多个不可变版本
   ├─ 创建者可以直接执行自己的模板版本
   └─ 某个版本可提交到 Market 审核
          ↓
       Market Listing / Listing Version
          ↓ 审核通过并公开
       SkillBot
```

这些名称的含义:

- **官方模板**:平台维护的公共执行入口,例如 `text-v1`。它不是硬编码在 CLI 中的本地模板;CLI 从当前 LoomLoom 服务读取可用列表。
- **私有模板**:用户自己的工作资产。使用 TemplateSpec 创建后会保存为私有模板,后续修改通过新增不可变版本完成。
- **自定义模板**:描述创建方式,不是第三种模板类型。用户完成自定义创作后,得到的是一个私有模板。
- **SkillBot**:私有模板某个版本通过 Market 审核后形成的公开、可付费执行形态。
- **Listing**:SkillBot 在 Market 中的货架对象;一个 Listing 可以连续发布多个版本。
- **Listing Version**:发布时从指定私有模板版本复制出的不可变执行快照。后续修改私有模板不会自动改变已上线 SkillBot。

当前没有独立的“公共模板”资源或命令。需要表达公共可执行对象时,应明确使用“官方模板”或“已发布的 Market SkillBot”,避免混淆。

`loomloom asset list` 是可执行资产聚合视图,当前合并“我的私有模板”和“Market SkillBot”;它不是一种新的模板类型,也不替代 `loomloom template list` 的官方模板列表。

---

## 当前官方模板

以下内容是当前官方模板示例。实际可用模板由目标环境决定,请以 `loomloom template list` 的实时结果为准。

| 模板 ID | 用例 | 输出 | 步骤 |
| --- | --- | --- | --- |
| `text-v1` | 文案、改写、摘要、问答、代码评审 | 文本 / 文件 | 文本生成 |
| `text-image-v1` | 插画、概念图、社交图片 | 图片 | 提示词准备 -> 图片生成 |
| `text-image-video-v1` | 分镜、广告、短视频素材 | 图片 + 视频 | 描述 -> 图片 -> 视频 |

---

## 官方模板标准 Excel 工作流

```bash
# 1. 下载工作簿模板
loomloom template download text-image-v1 --output-file ./task.xlsx

# 2. 填写工作簿并校验
loomloom template validate-file text-image-v1 ./task.xlsx

# 3. 提交工作簿
loomloom template submit-file text-image-v1 ./task.xlsx

# 4. 监控进度
loomloom run watch <run-id>

# 5. 下载服务端生成的结果工作簿
loomloom run result-workbook <run-id> --output-file ./task.result.xlsx

# 6. 下载生成的产物
loomloom artifact download <run-id> --output-dir ./downloads
```

`template backfill-results` 仍可用于旧的本地工作流。新工作流优先用 `run result-workbook`;服务端用提交的输入快照对齐原始行与结果。

---

## 当前官方模板字段

### 文本模板:`text-v1`

| 字段 | 是否必填 | 说明 |
| --- | --- | --- |
| 文本提示词 | 必填 | 主任务提示词,例如“把这段介绍改写为 80-120 字”。 |
| 写作要求 | 可选 | 风格、格式或输出约束。 |
| 参考文本 | 可选 | 直接填短文本,或用 `input-asset upload` 上传大文件并使用返回的 `input_asset_id`。 |

### 图片模板:`text-image-v1`

| 字段 | 是否必填 | 说明 |
| --- | --- | --- |
| 图片提示词 | 必填 | 要生成图片的描述。 |
| 风格要求 | 可选 | 例如水彩、写实或棚拍风格。 |
| 图片宽高比 | 必填 | `1:1`、`4:5`、`16:9` 或 `9:16`。 |

### 视频模板:`text-image-video-v1`

| 字段 | 是否必填 | 说明 |
| --- | --- | --- |
| 场景描述 | 必填 | 视频场景的描述。 |
| 视觉风格要求 | 可选 | 例如电影质感或动漫风格。 |
| 参考图片 URL | 可选 | 一个公开的 HTTP/HTTPS 图片 URL。 |
| 图片宽高比 | 必填 | `1:1`、`4:5`、`16:9` 或 `9:16`。 |
| 视频宽高比 | 必填 | `16:9` 或 `9:16`。 |
| 视频时长 | 必填 | `4`、`6` 或 `8` 秒。 |
| 生成音频 | 必填 | `false` 或 `true`。 |

---

## 输入素材

对于大的参考文件,先上传文件,仅当模板字段的 schema 明确接受素材引用时,才把返回的 `input_asset_id` 填入该字段。

```bash
loomloom input-asset upload ./brief.txt --content-type text/plain
loomloom input-asset upload ./diagram.png --content-type image/png
```

输入素材(input asset)和编排输入(orchestration input)是两回事:`input-asset upload` 返回 `input_asset_id`,用于放进模板字段的参考材料;而 `orchestration-input upload` 返回 `input_file_id`,为 `template-spec run` 提供行数据。

编排输入文件使用 JSONL 格式。对于常见的单根工作流,每个非空行可以是一个扁平 JSON 对象,且所有值都是字符串:

```jsonl
{"prompt":"第一个请求"}
{"prompt":"第二个请求"}
```

当工作流需要明确指定每个步骤的执行输入时,后端也支持 `steps.<step-id>.executions[]` 形式的 unified 行。两种格式中的执行参数值都必须是字符串,并符合私有模板版本允许的输入参数。不得猜测 step ID;只有掌握准确工作流步骤映射时才使用 unified 输入。

---

## 运行状态

用你的 CogFoundry 工作区控制台 URL 在线查看运行进度。

| 状态 | 含义 |
| --- | --- |
| `pending` / `queued` | 运行已受理,正在等待执行。 |
| `running` | 运行进行中。 |
| `completed` | 所有任务成功完成,结果可用。 |
| `partially_failed` | 部分任务失败,但成功结果仍可下载。 |
| `failed` | 运行失败。 |
| `cancelled` | 运行已取消。 |

---

## 命令参考

诸如 `taskFixedFeeT`、`amountT` 这类金额值以 API 单位计,其中 10,000,000 单位等于 1 个货币单位。

### 诊断

| 命令 | 说明 |
| --- | --- |
| `loomloom doctor` | 检查服务可达性、token 接线和版本信息。 |

### 输入

| 命令 | 说明 |
| --- | --- |
| `loomloom input-asset upload <file>` | 上传可复用的原始输入素材(文本/图片),返回 `input_asset_id`。 |
| `loomloom orchestration-input upload <file.jsonl>` | 上传私有模板执行所需的 JSONL 输入,返回 `template-spec run` 所需的 `input_file_id`。 |

### 官方模板

| 命令 | 说明 |
| --- | --- |
| `loomloom template list` | 列出当前环境已发布的官方模板。 |
| `loomloom template schema <id>` | 展示模板字段。 |
| `loomloom template download <id>` | 下载 Excel 工作簿模板。 |
| `loomloom template validate-file <id> <xlsx>` | 校验已填写的工作簿。 |
| `loomloom template precheck-file <id> <xlsx>` | 估算工作簿成本但不提交。 |
| `loomloom template submit-file <id> <xlsx>` | 把已填写的工作簿作为运行提交。 |
| `loomloom template backfill-results <run-id> <xlsx>` | 旧的本地结果回填。 |

### 私有模板(通过 TemplateSpec 创建)

| 命令 | 说明 |
| --- | --- |
| `loomloom template-spec check <spec.json>` | 校验用于创建私有模板的 TemplateSpec。 |
| `loomloom template-spec docs [topic]` | 展示内置 TemplateSpec 文档。 |
| `loomloom template-spec models <step-type>` | 列出某步骤类型的模型。 |
| `loomloom template-spec create <spec.json>` | 创建私有模板。 |
| `loomloom template-spec create-version <template-id> <spec.json>` | 为已有私有模板新增版本。 |
| `loomloom template-spec list` | 列出我的私有模板。 |
| `loomloom template-spec get <template-id>` | 展示单个私有模板及其版本。 |
| `loomloom template-spec versions <template-id>` | 列出私有模板的版本。 |
| `loomloom template-spec download-workbook <template-id> <version-id>` | 下载用户模板工作簿。 |
| `loomloom template-spec validate-workbook <template-id> <version-id> <xlsx>` | 校验用户模板工作簿。 |
| `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx>` | 提交用户模板工作簿。 |
| `loomloom template-spec run <template-id> --version-id <id> --input-file-id <id>` | 用上传的 JSONL 输入运行某私有模板版本。 |

### 运行

| 命令 | 说明 |
| --- | --- |
| `loomloom run submit <id> -f rows.json` | 从 JSON 数组或 JSONL 文件提交输入。 |
| `loomloom run list` | 列出运行(可带 Market 上下文)。 |
| `loomloom run get <run-id>` | 展示单个运行详情。 |
| `loomloom run watch <run-id>` | 监控运行直到终态。 |
| `loomloom run result-rows <run-id>` | 展示对齐的输入行与结果。 |
| `loomloom run result-workbook <run-id>` | 下载服务端生成的结果工作簿。 |

### 产物

| 命令 | 说明 |
| --- | --- |
| `loomloom artifact list <run-id>` | 列出生成的产物。 |
| `loomloom artifact download <run-id>` | 下载生成的产物。 |

### 目录

| 命令 | 说明 |
| --- | --- |
| `loomloom model list --step-type <type>` | 列出某步骤类型的可执行模型。 |
| `loomloom asset list` | 聚合列出我的私有模板与可用 Market SkillBot;不包含官方模板。 |

### Market —— 买家

| 命令 | 说明 |
| --- | --- |
| `loomloom market list` | 浏览已发布的 Market SkillBot。 |
| `loomloom market show <listing-id>` | 展示单个 SkillBot,含其输入 schema。 |
| `loomloom market quote <listing-id> --input-file <json>` | 估算执行成本。 |
| `loomloom market run <listing-id> --input-file <json> --confirm` | 执行 SkillBot(付费)。 |
| `loomloom usage list` | 列出我的 Market SkillBot 使用记录。 |
| `loomloom usage get <run-transaction-id>` | 展示单条使用记录。 |

### Market —— 创作者

| 命令 | 说明 |
| --- | --- |
| `loomloom listing publish <template-id> --template-version-id <id> --display-name <name> --task-fixed-fee-t <fee>` | 提交模板版本进入 Market 审核。 |
| `loomloom listing publish <template-id> --listing-id <listing-id> --template-version-id <new-id> ...` | 为已有 listing 提交新版本。 |
| `loomloom listing list` | 列出我的 Market 上架。 |
| `loomloom listing show <listing-id>` | 展示我的某个上架。 |
| `loomloom listing versions <listing-id>` | 列出我某个上架的版本。 |
| `loomloom listing update <listing-id> --display-name <name>` | 提交公开资料更新进入审核;传展示名、描述或两者。 |
| `loomloom listing unlist <listing-id>` | 停止某上架的新执行。 |
| `loomloom listing relist <listing-id>` | 恢复此前已下架的上架。 |
| `loomloom listing withdraw <listing-id>` | 撤回某上架的待审核请求。 |
| `loomloom creator earnings` | 列出 Market 收益。 |
| `loomloom creator transactions` | 列出 Market 交易。 |
| `loomloom creator review list` | 列出我的审核请求。 |
| `loomloom creator review get <review-request-id>` | 展示单个审核请求。 |
| `loomloom creator review withdraw <review-request-id>` | 撤回一个待审核请求。 |

---

## 私有模板创作

用户使用 TemplateSpec JSON 描述自定义工作流的步骤、输入字段和字段绑定,创建后保存为私有模板。典型的 agent 辅助创作流程是:

```bash
# 1. 查看某执行单元的可用模型
loomloom template-spec models text-generate

# 2. 本地校验 spec
loomloom template-spec check ./my-template.spec.json

# 3. 创建私有模板
loomloom template-spec create ./my-template.spec.json --version-note "initial version"

# 4. 模板变化时新增版本
loomloom template-spec create-version <template-id> ./my-template.spec.json

# 5. 下载、填写、校验并提交工作簿
loomloom template-spec download-workbook <template-id> <version-id> --output-file ./input.xlsx
loomloom template-spec validate-workbook <template-id> <version-id> ./input.xlsx
loomloom template-spec submit-workbook <template-id> <version-id> ./input.xlsx
```

注意:

- TemplateSpec JSON 是事实来源;工作簿是生成的产物。
- 编写用于创建私有模板的 spec 前,用 `loomloom template-spec docs spec` 查看内置规范。
- 用 `loomloom template-spec docs examples` 查看模式示例。
- 用 `loomloom template-spec docs conversation` 进行 agent 辅助的对话式创作。
- 模板变化需要重新下载工作簿。
- `submit-workbook` 会创建真实的托管运行;agent 应在提交前请求明确确认。
- `template-spec run` 同样会创建真实的托管运行,需要相同的确认。

你也可以直接用扁平 JSONL 行运行私有模板版本,无需填写工作簿:

```bash
# 1. 上传行数据并记下返回的 input_file_id
loomloom orchestration-input upload ./rows.jsonl

# 2. 用该输入运行版本
loomloom template-spec run <template-id> --version-id <version-id> --input-file-id <input_file_id>
```

---

## Market / SkillBot

LoomLoom Market 让创作者把自己的私有模板版本发布为付费 SkillBot,并让买家运行它。同一个 CLI 服务两种角色。

发布不是把私有模板本身改成公开状态。系统会从所选私有模板版本生成一个不可变的 Listing Version 执行快照。创作者之后继续修改私有模板或创建新版本,不会自动改变当前线上 SkillBot;需要再次提交新版本审核。

### 买家:发现并运行 SkillBot

```bash
# 1. 浏览并查看 SkillBot
loomloom market list --keyword "tweet"
loomloom market show <listing-id>

# 2. 估算成本(输入文件携带一个对应 listing schema 的 taskInputs 数组)
loomloom market quote <listing-id> --input-file ./request.json

# 3. 执行(付费;必须带 --confirm)
loomloom market run <listing-id> --input-file ./request.json --confirm

# 4. 查看自己的使用记录
loomloom usage list
loomloom usage get <run-transaction-id>
```

`request.json` 示例:

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

先用 listing 详情了解公开字段并取得 `listingVersionId`,但不要从 `inputSchemaSnapshot` 推断内部 step ID:公开快照可能不包含它。当前 CLI 要求精确的 Product API `taskInputs` 请求体;如果用户或集成方没有提供兼容映射,应停止并索要可用的请求 JSON,不得猜测。

### 创作者:发布并管理 SkillBot

```bash
# 1. 提交模板版本进入审核(该版本必须已有一次成功运行)
loomloom listing publish <template-id> \
  --template-version-id <version-id> \
  --display-name "My SkillBot" \
  --task-fixed-fee-t 1000000

# 为已有 listing 提交新版本
loomloom listing publish <template-id> \
  --listing-id <listing-id> \
  --template-version-id <new-version-id> \
  --display-name "My SkillBot" \
  --task-fixed-fee-t 1000000

# 2. 跟踪上架与审核
loomloom listing list
loomloom listing versions <listing-id>
loomloom creator review list

# 3. 管理售卖状态和公开资料
loomloom listing unlist <listing-id>
loomloom listing relist <listing-id>
loomloom listing update <listing-id> --description "Updated description"

# 4. 查看收入
loomloom creator earnings
loomloom creator transactions
```

注意:

- `market run` 会创建真实的付费运行;agent 应先运行 `market quote`,并在执行前请求明确确认。
- `listing publish` 和 `listing update` 是提交变更进入审核;未获批准前不生效。
- `listing publish --listing-id <listing-id>` 会为已有 listing 提交新版本;审核通过前,当前已发布版本继续生效。
- `listing update`、`listing unlist`、`listing relist` 和撤回审核都会改变远程状态;agent 调用前应先总结动作并取得明确确认。
- 金额 `*FeeT` / `*AmountT` / `*PayableT` 值以 API 单位计,其中 10,000,000 单位等于 1 个货币单位。

---

## Agent 命令串联

当一条命令的结果要交给下一条命令时,agent 应优先使用 `--output json`,并原样保存 ID,不得推测或转换:

```text
orchestration-input upload → inputFileId → template-spec run
template-spec run / run submit → runId → run watch / 结果命令
listing publish → reviewRequestId → creator review get/withdraw
market run → runTransactionId 和 runId → usage get / run watch
```

文本输出使用 `input_file_id` 等标签;JSON 输出使用 Product API 字段名,例如 `inputFileId`。

对于 `template submit-file`、`template-spec submit-workbook`、`run submit`、`template-spec run` 和 `market run`,显式传入 `--client-request-id`,并将它和请求一起保存。只有重试完全相同的 payload 时才能复用;payload 变化时必须使用新的 ID。遇到结果不明确的失败后,不要盲目重试付费命令或远程状态变更命令,应先查询原请求是否已经成功。

---

## 卸载

macOS / Linux:

```bash
# 移除 CLI 和 skill 包
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash

# 仅移除 CLI
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash -s -- --cli-only

# 仅移除 skill 包
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash -s -- --skill-only
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.ps1 | iex
```

---

## FAQ

**在哪里获取 token?**

前往 [CogFoundry API Keys](https://console.cogfoundry.ai/api-keys) 创建或获取 token,然后替换示例中的占位符。

CogFoundry 正式 Token 只能用于 `https://loomloom.cogfoundry.ai`。如果目标服务不是这个主机,不要使用正式 Token。

**在哪里查看运行状态?**

使用 [CogFoundry Console](https://console.cogfoundry.ai)。当前暂不提供 Workflow Run 详情页 URL 模板,因此不要自行拼接运行详情链接。

**为什么 `template list` 返回不了模板?**

该账号或环境可能没有可见模板。请联系你的 CogFoundry 工作区管理员确认模板发布情况和权限。

**不用 agent 能手动使用 CLI 吗?**

可以。上面所有工作流都能用 CLI 命令运行。

---

## 链接

- GitHub:[github.com/Cogfoundry-ai/loomloom](https://github.com/Cogfoundry-ai/loomloom)
- LoomLoom API:`https://loomloom.cogfoundry.ai/loom/v1`
- API Keys:[console.cogfoundry.ai/api-keys](https://console.cogfoundry.ai/api-keys)
- CogFoundry Console:[console.cogfoundry.ai](https://console.cogfoundry.ai)
- CogFoundry 官网:[cogfoundry.ai](https://cogfoundry.ai)
- Workflow Run 详情页:当前暂不提供固定 URL 模板。
