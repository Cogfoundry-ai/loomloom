[English](README.md) | 简体中文

# LoomLoom

> **批量内容生成平台** - 用自然语言编排文案、图片和视频生成等 AI 工作流。
> 由 CogFoundry 构建 - [github.com/Cogfoundry-ai/loomloom](https://github.com/Cogfoundry-ai/loomloom)

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Cogfoundry-ai/loomloom/1-overview)

---

## 它能做什么

LoomLoom 是一个 CLI 和 Agent Skill 包，用于 AI 驱动的批量内容工作流。你不需要手写工作流代码，只需要用自然语言描述任务，让 agent 下载模板、准备数据、提交运行、跟踪进度并下载结果。

常见用途：

- **批量文案生成** - 商品描述、改写、摘要、问答和文件级文本修改。
- **批量图片生成** - 电商图、社媒素材、概念图，以及按行生成图片。
- **批量视频生成** - 分镜、广告素材和文本到视频工作流。

---

## 支持的 Agent

安装 LoomLoom 会为所选 agent 添加对应的 skill 包：

| Agent | 状态 |
| --- | --- |
| **Codex** (OpenAI) | 支持 |
| **Claude Code** (Anthropic) | 支持 |
| **OpenClaw** | 支持 |

---

## 快速开始

### 通过 Agent 辅助安装

把下面这段消息发送给 Codex、Claude Code 或 OpenClaw。把 `your-token` 替换为你在 CogFoundry API Keys 页面创建的 token。

```text
Install LoomLoom from this GitHub repository: https://github.com/Cogfoundry-ai/loomloom
My server URL is https://loomloom.cogfoundry.ai/loom/v1, and my token is your-token.
After installation, run doctor to check whether the setup is healthy.
```

### 手动安装

macOS / Linux：

```bash
# 默认安装 Codex skill 包
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash

# Claude Code
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --agent claude

# OpenClaw
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --agent openclaw

# 指定版本
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --version v0.1.0-beta.1

# 最新 beta 或内部频道
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta --no-brew

# 指定预发布标签
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --version v0.1.0-beta.1 --no-brew
```

Windows PowerShell：

```powershell
# 默认安装
irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1 | iex

# Claude Code
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Agent claude

# OpenClaw
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Agent openclaw

# 最新 beta 或内部频道
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Channel beta

# 指定预发布标签
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Version v0.1.0-beta.1
```

Homebrew 分发计划中，但 tap 仓库和发布所需 token 还没有配置好。目前请使用上面的安装脚本。

---

## 配置凭证

```bash
export LOOMLOOM_SERVER="https://loomloom.cogfoundry.ai/loom/v1"
export LOOMLOOM_TOKEN="<your CogFoundry token>"
```

如果不想每次会话都重新设置，可以把这些值加入 `~/.zshrc`、`~/.bashrc` 或其他 shell profile。CLI 仍兼容旧变量 `BATCHJOB_SERVER` 和 `BATCHJOB_TOKEN`，但新配置应优先使用 `LOOMLOOM_*`。

在 [CogFoundry API Keys](https://console-dev.cogfoundry.ai/api-keys) 页面创建或获取 token。不要在文档、截图、日志或公开对话中暴露真实 token。

> 安全要求：只把 token 发送到你明确配置的 `LOOMLOOM_SERVER` URL（或传入的 `--server` 值），并使用 HTTPS。生产默认 base URL 是 `https://loomloom.cogfoundry.ai/loom/v1`。不要把 token 发送到未指定的主机，也不要跟随会把 token 转发到其他域名的重定向。请使用目标环境签发的 token。

---

## 验证安装

```bash
loomloom doctor
```

如果环境健康，就可以开始提交模板运行。

---

## 模板、私有模板和 SkillBot

“Template” 是 LoomLoom 中可重复 AI 工作流的统称。当前产品区分两类模板，以及一条 Market 发布路径：

```text
Template
├─ Official template
│  ├─ 由平台维护和发布
│  ├─ 任何有权限的用户都可以发现和执行
│  └─ 使用 loomloom template ... 命令
│
└─ Private template
   ├─ 用户通过 TemplateSpec 创建和维护
   ├─ 一个模板可以有多个不可变版本
   ├─ 创建者可以直接执行自己的模板版本
   └─ 某个版本可以提交到 Market 审核
          ↓
       Market Listing / Listing Version
          ↓ 审核通过并公开
       SkillBot
```

这些名称的含义：

- **Official template**：平台维护的公开执行入口，例如 `text-v1`。它不是 CLI 中硬编码的本地模板；CLI 会从当前 LoomLoom 服务读取可用列表。
- **Private template**：用户自己的工作资产。通过 TemplateSpec 创建并保存为私有模板，之后通过新增不可变版本进行变更。
- **Custom template**：描述创建方式，不是第三类模板。用户完成自定义创作后，结果就是 private template。
- **SkillBot**：私有模板版本通过 Market 审核后的公开、付费、可执行形式。
- **Listing**：SkillBot 在 Market 上的货架对象；一个 Listing 可以随着时间发布多个版本。
- **Listing Version**：发布时从指定私有模板版本复制出来的不可变执行快照。之后修改私有模板不会自动改变线上 SkillBot。

当前没有单独的 “public template” 资源或命令。提到公开可执行对象时，请明确说 “official template” 或 “published Market SkillBot”，避免混淆。

`loomloom asset list` 是可执行资产的聚合视图，目前合并了 “my private templates” 和 “Market SkillBots”。它不是新的模板类型，也不会替代 `loomloom template list` 返回的官方模板列表。

---

## 当前官方模板

下面是当前官方模板示例。实际可用模板取决于目标环境，请以 `loomloom template list` 的实时结果为准。

| Template ID | 用途 | 输出 | 步骤 |
| --- | --- | --- | --- |
| `text-v1` | 文案、改写、摘要、问答和代码审查 | 文本 / 文件 | 文本生成 |
| `text-image-v1` | 插画、概念图和社媒图片 | 图片 | 提示词准备 -> 图片生成 |
| `text-image-video-v1` | 分镜、广告和短视频素材 | 图片 + 视频 | 描述 -> 图片 -> 视频 |

---

## 标准 Excel 工作流（官方模板）

```bash
# 1. 下载工作簿模板
loomloom template download text-image-v1 --output-file ./task.xlsx

# 2. 填写工作簿并验证
loomloom template validate-file text-image-v1 ./task.xlsx

# 3. 提交工作簿
loomloom template submit-file text-image-v1 ./task.xlsx

# 4. 观察进度
loomloom run watch <run-id>

# 5. 下载服务端生成的结果工作簿
loomloom run result-workbook <run-id> --output-file ./task.result.xlsx

# 6. 下载生成的 artifact
loomloom artifact download <run-id> --output-dir ./downloads
```

`template backfill-results` 仍可用于较旧的本地工作流。新工作流建议优先使用 `run result-workbook`；服务端会使用提交时的输入快照来对齐原始行和结果。

---

## 模板字段

### 文本模板：`text-v1`

| 字段 | 是否必填 | 说明 |
| --- | --- | --- |
| Text prompt | 必填 | 主要任务提示，例如 “Rewrite this introduction in 80-120 words.” |
| Writing requirements | 可选 | 风格、格式或输出约束。 |
| Reference text | 可选 | 短文本，或先用 `input-asset upload` 上传大文件，再使用返回的 `input_asset_id`。 |

### 图片模板：`text-image-v1`

| 字段 | 是否必填 | 说明 |
| --- | --- | --- |
| Image prompt | 必填 | 要生成图片的描述。 |
| Style requirements | 可选 | 例如 watercolor、photorealistic 或 studio style。 |
| Image aspect ratio | 必填 | `1:1`、`4:5`、`16:9` 或 `9:16`。 |

### 视频模板：`text-image-video-v1`

| 字段 | 是否必填 | 说明 |
| --- | --- | --- |
| Scene description | 必填 | 视频场景描述。 |
| Visual style requirements | 可选 | 例如 cinematic tone 或 anime style。 |
| Reference image URL | 可选 | 一个公开 HTTP/HTTPS 图片 URL。 |
| Image aspect ratio | 必填 | `1:1`、`4:5`、`16:9` 或 `9:16`。 |
| Video aspect ratio | 必填 | `16:9` 或 `9:16`。 |
| Video duration | 必填 | `4`、`6` 或 `8` 秒。 |
| Generate audio | 必填 | `false` 或 `true`。 |

---

## 输入资产

对于大型参考文件，请先上传文件，并且只把返回的 `input_asset_id` 放入 schema 接受 asset reference 的模板字段中。

```bash
loomloom input-asset upload ./brief.txt --content-type text/plain
loomloom input-asset upload ./diagram.png --content-type image/png
```

Input assets 和 orchestration inputs 是不同概念：`input-asset upload` 返回用于模板字段中参考材料的 `input_asset_id`，而 `orchestration-input upload` 返回用于 `template-spec run` 行数据的 `input_file_id`。

Orchestration input 文件是 JSONL。对于常见的单根工作流，每个非空行可以是值为字符串的扁平 JSON 对象：

```jsonl
{"prompt":"first request"}
{"prompt":"second request"}
```

当工作流需要显式的逐步骤执行输入时，后端也支持形如 `steps.<step-id>.executions[]` 的统一行格式。在两种格式中，执行参数值都必须是字符串，并且必须符合私有模板版本允许的输入参数。不要编造 step ID；只有在确切知道工作流步骤映射时才使用统一输入。

---

## 运行状态

使用 [CogFoundry Console](https://console-dev.cogfoundry.ai/quickstart) 在线查看运行进度。

| 状态 | 含义 |
| --- | --- |
| `pending` / `queued` | 运行已被接受，正在等待执行。 |
| `running` | 正在执行。 |
| `completed` | 所有任务成功完成，结果可用。 |
| `partially_failed` | 部分任务失败，但成功结果仍可下载。 |
| `failed` | 运行失败。 |
| `cancelled` | 运行已取消。 |

---

## 命令参考

`taskFixedFeeT`、`amountT` 等金额值使用 API units，其中 10,000,000 units 等于 1 个货币单位。

### 诊断

| 命令 | 说明 |
| --- | --- |
| `loomloom doctor` | 检查服务可达性、token 配置和版本信息。 |

### 输入

| 命令 | 说明 |
| --- | --- |
| `loomloom input-asset upload <file>` | 上传可复用的原始输入资产（文本/图片），并获得 `input_asset_id`。 |
| `loomloom orchestration-input upload <file.jsonl>` | 上传扁平 JSONL 行，并获得 `template-spec precheck` 和 `template-spec run` 所需的 `input_file_id`。 |

### 官方模板

| 命令 | 说明 |
| --- | --- |
| `loomloom template list` | 列出当前环境发布的官方模板。 |
| `loomloom template schema <id>` | 显示模板字段。 |
| `loomloom template download <id>` | 下载 Excel 工作簿模板。 |
| `loomloom template validate-file <id> <xlsx>` | 验证填好的工作簿。 |
| `loomloom template precheck-file <id> <xlsx>` | 在不提交的情况下估算工作簿成本。 |
| `loomloom template submit-file <id> <xlsx>` | 把填好的工作簿提交为一次运行。 |
| `loomloom template backfill-results <run-id> <xlsx>` | 旧版本地结果回填。 |

### 私有模板（通过 TemplateSpec 创建）

| 命令 | 说明 |
| --- | --- |
| `loomloom template-spec check <spec.json>` | 验证用于创建私有模板的 TemplateSpec。 |
| `loomloom template-spec docs [topic]` | 显示内置 TemplateSpec 文档。 |
| `loomloom template-spec models <step-type>` | 列出某种 step type 可用的模型。 |
| `loomloom template-spec create <spec.json>` | 创建私有模板。 |
| `loomloom template-spec create-version <template-id> <spec.json>` | 给现有私有模板添加新版本。 |
| `loomloom template-spec list` | 列出我的私有模板。 |
| `loomloom template-spec get <template-id>` | 显示一个私有模板及其版本。 |
| `loomloom template-spec versions <template-id>` | 列出私有模板的版本。 |
| `loomloom template-spec download-workbook <template-id> <version-id>` | 下载用户模板工作簿。 |
| `loomloom template-spec validate-workbook <template-id> <version-id> <xlsx>` | 验证用户模板工作簿。 |
| `loomloom template-spec precheck-workbook <template-id> <version-id> <xlsx>` | 在不提交的情况下预估用户模板工作簿费用和余额。 |
| `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx>` | 提交用户模板工作簿。 |
| `loomloom template-spec precheck <template-id> --version-id <id> --input-file-id <id>` | 在不提交的情况下预估已上传 JSONL 输入的费用和余额。 |
| `loomloom template-spec run <template-id> --version-id <id> --input-file-id <id>` | 使用已上传的 JSONL 输入运行私有模板版本。 |

### 运行

| 命令 | 说明 |
| --- | --- |
| `loomloom run submit <id> -f rows.json` | 从 JSON 数组或 JSONL 文件提交输入。 |
| `loomloom run list` | 列出运行，可包含 Market 上下文。 |
| `loomloom run get <run-id>` | 显示一次运行的详情。 |
| `loomloom run watch <run-id>` | 观察运行进度，直到进入终态。 |
| `loomloom run result-rows <run-id>` | 显示对齐后的输入行和结果。 |
| `loomloom run result-workbook <run-id>` | 下载服务端生成的结果工作簿。 |

### Artifacts

| 命令 | 说明 |
| --- | --- |
| `loomloom artifact list <run-id>` | 列出生成的 artifacts。 |
| `loomloom artifact download <run-id>` | 下载生成的 artifacts。 |

### 目录

| 命令 | 说明 |
| --- | --- |
| `loomloom model list --step-type <type>` | 列出某种 step type 可执行的模型。 |
| `loomloom asset list` | 聚合列出我的私有模板和可用 Market SkillBots；不包含官方模板。 |

### Market - 买家

| 命令 | 说明 |
| --- | --- |
| `loomloom market list` | 浏览已发布的 Market SkillBots。 |
| `loomloom market show <listing-id>` | 显示某个 SkillBot，包括输入 schema。 |
| `loomloom market quote <listing-id> --input-file <json>` | 估算执行成本。 |
| `loomloom market run <listing-id> --input-file <json> --confirm --client-request-id <id>` | 使用 JSON 输入行执行 SkillBot（付费）。 |
| `loomloom market workbook download <listing-id> --output-file <xlsx>` | 下载 Market 工作簿模板。 |
| `loomloom market workbook validate <listing-id> --file <xlsx>` | 验证已填写的 Market 工作簿。 |
| `loomloom market workbook quote <listing-id> --file <xlsx>` | 估算工作簿执行成本。 |
| `loomloom market workbook run <listing-id> --file <xlsx> --confirm --client-request-id <id>` | 使用工作簿执行 SkillBot（付费）。 |
| `loomloom usage list` | 列出我的 Market SkillBot 使用记录。 |
| `loomloom usage get <run-transaction-id>` | 显示一条使用记录。 |

### Market - 创建者

| 命令 | 说明 |
| --- | --- |
| `loomloom listing publish <template-id> --template-version-id <id> --display-name <name> --task-fixed-fee-t <fee>` | 提交模板版本到 Market 审核。 |
| `loomloom listing publish <template-id> --listing-id <listing-id> --template-version-id <new-id> ...` | 给现有 listing 提交新版本。 |
| `loomloom listing list` | 列出我的 Market listings。 |
| `loomloom listing show <listing-id>` | 显示我的某个 listing。 |
| `loomloom listing versions <listing-id>` | 列出我的某个 listing 的版本。 |
| `loomloom listing update <listing-id> --display-name <name>` | 提交公开资料更新审核；可传 display name、description 或两者都传。 |
| `loomloom listing unlist <listing-id>` | 停止 listing 的新执行。 |
| `loomloom listing relist <listing-id>` | 恢复之前下架的 listing。 |
| `loomloom listing withdraw <listing-id>` | 撤回某个 listing 的待审核请求。 |
| `loomloom creator earnings` | 列出 Market 收益。 |
| `loomloom creator transactions` | 列出 Market 交易。 |
| `loomloom creator review list` | 列出我的审核请求。 |
| `loomloom creator review get <review-request-id>` | 显示一个审核请求。 |
| `loomloom creator review withdraw <review-request-id>` | 撤回一个待审核请求。 |

---

## 私有模板创作

用户通过 TemplateSpec JSON 描述自定义工作流的步骤、输入字段和字段绑定；创建后会保存为私有模板。典型的 agent 辅助创作流程如下：

```bash
# 1. 查看某个 execution unit 可用的模型
loomloom template-spec models text-generate

# 2. 本地验证 spec
loomloom template-spec check ./my-template.spec.json

# 3. 创建私有模板
loomloom template-spec create ./my-template.spec.json --version-note "initial version"

# 4. 模板变更时添加新版本
loomloom template-spec create-version <template-id> ./my-template.spec.json

# 5. 下载、填写、验证、预估并提交工作簿
loomloom template-spec download-workbook <template-id> <version-id> --output-file ./input.xlsx
loomloom template-spec validate-workbook <template-id> <version-id> ./input.xlsx
loomloom template-spec precheck-workbook <template-id> <version-id> ./input.xlsx
loomloom template-spec submit-workbook <template-id> <version-id> ./input.xlsx
```

注意：

- TemplateSpec JSON 是事实来源；工作簿是生成产物。
- 编写自定义 spec 前，请用 `loomloom template-spec docs spec` 查看内置规范。
- 用 `loomloom template-spec docs examples` 查看模式示例。
- 用 `loomloom template-spec docs conversation` 查看 agent 辅助的对话式创作说明。
- 模板变更后需要下载新的工作簿。
- 如果想在提交前预估模型/API 费用和余额，可以先运行 `precheck-workbook`；预估不会创建 run，也不会执行任务。
- 预估文本输出包含 `estimated_cost`、`available_balance` 和 `sufficient`；JSON 输出使用 `estimatedTotalCostT`。
- `submit-workbook` 会创建真实的托管运行；agent 提交前应请求明确确认。
- `template-spec run` 也会创建真实的托管运行，并需要同样的确认。

也可以不填写工作簿，直接用扁平 JSONL 行运行私有模板版本：

```bash
# 1. 上传行数据并记录返回的 input_file_id
loomloom orchestration-input upload ./rows.jsonl

# 2. 预估费用和余额
loomloom template-spec precheck <template-id> --version-id <version-id> --input-file-id <input_file_id>

# 3. 使用该输入运行版本
loomloom template-spec run <template-id> --version-id <version-id> --input-file-id <input_file_id>
```

---

## Market / SkillBot

LoomLoom Market 允许创建者把私有模板版本发布为付费 SkillBot，也允许买家运行它。同一个 CLI 服务于两种角色。

发布不会把私有模板本身变成公开对象。系统会从选定的私有模板版本生成一个不可变的 Listing Version 执行快照。之后编辑私有模板或新增版本，都不会自动改变线上 SkillBot；你需要再次提交新版本审核。

### 买家：发现并运行 SkillBot

```bash
# 1. 浏览并查看 SkillBots
loomloom market list --keyword "tweet"
loomloom market show <listing-id>

# 2A. JSON 输入：用公开 inputRows 估算成本
loomloom market quote <listing-id> --input-file ./request.json

# 3A. JSON 输入：确认后执行（付费）
loomloom market run <listing-id> --input-file ./request.json --confirm --client-request-id <stable-id>

# 2B. 工作簿输入：下载、填写、验证并估价
loomloom market workbook download <listing-id> --output-file ./market-input.xlsx
loomloom market workbook validate <listing-id> --file ./market-input.xlsx
loomloom market workbook quote <listing-id> --file ./market-input.xlsx

# 3B. 工作簿输入：确认后执行（付费）
loomloom market workbook run <listing-id> --file ./market-input.xlsx --confirm --client-request-id <stable-id>

# 4. 查看自己的使用记录并下载结果
loomloom usage list
loomloom usage get <run-transaction-id>
loomloom run result-workbook <run-id> --output-file ./market-result.xlsx
```

示例 `request.json`：

```json
{
  "inputRows": [
    {
      "prompt": "write a launch tweet"
    }
  ]
}
```

构造 JSON 输入前，先用 `market show` 查看公开字段和示例。展示给用户看 `inputSchemaSnapshot.fields[].label`，提交时用 `inputSchemaSnapshot.fields[].key` 作为 `inputRows` 的 key，并按 `fields[].required` 做必填。不要向 Market 买方执行接口传 `taskInputs`、`workflowDefinition`、`templateSpec` 或隐藏的 Core / TemplateSpec 内部结构。

### 创建者：发布和管理 SkillBot

```bash
# 1. 发布模板版本进行审核（必须已有一次成功运行）
loomloom listing publish <template-id> \
  --template-version-id <version-id> \
  --display-name "My SkillBot" \
  --task-fixed-fee-t 1000000

# 给现有 listing 提交新版本
loomloom listing publish <template-id> \
  --listing-id <listing-id> \
  --template-version-id <new-version-id> \
  --display-name "My SkillBot" \
  --task-fixed-fee-t 1000000

# 2. 跟踪 listing 和审核
loomloom listing list
loomloom listing versions <listing-id>
loomloom creator review list

# 3. 管理销售状态和公开资料
loomloom listing unlist <listing-id>
loomloom listing relist <listing-id>
loomloom listing update <listing-id> --description "Updated description"

# 4. 查看收入
loomloom creator earnings
loomloom creator transactions
```

注意：

- `market run` 会创建真实的付费运行；agent 应先运行 `market quote`，并在执行前请求明确确认。
- Market 工作簿执行前应先验证和估价，执行后用 `run result-workbook` 下载结果。
- `listing publish` 和 `listing update` 会提交审核；审核通过前不会生效。
- `listing publish --listing-id <listing-id>` 会给现有 listing 提交新版本。当前已发布版本会保持可用，直到新审核通过。
- `listing update`、`listing unlist`、`listing relist` 和撤回审核都会改变远程状态。agent 调用前应总结操作并请求明确确认。
- 对于 Market SkillBot，`market quote` 可在执行前预估买家应付金额。平台从每笔调用费中抽取 10%，创作者实收为调用费的 90%。
- 只有在重试完全相同的付费 Market payload 时，才复用同一个 `--client-request-id`。只要输入变化，就使用新的 ID。
- 工作簿 `content` 会作为 JSON 里的 base64 发送；不要打印完整 base64。结果行中的 `accessUrl` 是临时签名 URL，不要写入长期日志或文档。
- 金额类 `*FeeT` / `*AmountT` / `*PayableT` 值使用 API units，其中 10,000,000 units 等于 1 个货币单位。

---

## Agent 命令串联

命令结果要喂给后续命令时，agent 应优先使用 `--output json`。请完整保留 ID，不要推断或转换。

```text
orchestration-input upload -> inputFileId -> template-spec precheck -> template-spec run
template-spec run / run submit -> runId -> run watch / result commands
listing publish -> reviewRequestId -> creator review get/withdraw
market run -> runTransactionId and runId -> usage get / run watch
market workbook run -> runTransactionId and runId -> usage get / run watch / result-workbook
```

文本输出使用 `input_file_id` 这类标签；JSON 输出使用 `inputFileId` 这类 Product API 字段名。

对于 `template submit-file`、`template-spec submit-workbook`、`run submit`、`template-spec run`、`market run` 和 `market workbook run`，请传入明确的 `--client-request-id`，随请求保存它，并且只在重试完全相同 payload 时复用。如果 payload 变化，请使用新 ID。对于付费或会改变远程状态的命令，不要在结果不明确时盲目重试；应先检查原请求是否已成功。

---

## 卸载

macOS / Linux：

```bash
# 移除 CLI 和 skill 包
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash

# 只移除 CLI
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash -s -- --cli-only

# 只移除 skill 包
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash -s -- --skill-only
```

Windows PowerShell：

```powershell
irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.ps1 | iex
```

---

## FAQ

**我从哪里获取 token？**

在 [CogFoundry API Keys](https://console-dev.cogfoundry.ai/api-keys) 页面创建或获取 token，然后替换示例中的占位符。

只把 token 用于你配置的 `LOOMLOOM_SERVER`。不要把它发送到未明确设置的主机。

**在哪里查看运行状态？**

使用 [CogFoundry Console](https://console-dev.cogfoundry.ai/quickstart)。目前没有 Workflow Run 详情页的固定 URL 模板，所以不要自行拼接 run detail 链接。

**为什么 `template list` 没有返回模板？**

账号或环境可能没有可见模板。请联系 CogFoundry workspace 管理员确认模板发布和权限。

**可以不通过 agent，手动使用 CLI 吗？**

可以。上面所有工作流都可以直接通过 CLI 命令运行。

---

## 链接

- GitHub: [github.com/Cogfoundry-ai/loomloom](https://github.com/Cogfoundry-ai/loomloom)
- LoomLoom API: `https://loomloom.cogfoundry.ai/loom/v1`
- API Keys: [console-dev.cogfoundry.ai/api-keys](https://console-dev.cogfoundry.ai/api-keys)
- CogFoundry Console: [console-dev.cogfoundry.ai/quickstart](https://console-dev.cogfoundry.ai/quickstart)
- CogFoundry website: [cogfoundry.ai](https://cogfoundry.ai)
- Workflow Run detail page: 暂无固定 URL 模板。
