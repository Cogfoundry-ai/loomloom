English | [简体中文](README.zh-CN.md)

# LoomLoom

> **Batch content generation platform** - Use natural language to orchestrate AI workflows for copy, image, and video generation.
> Built by CogFoundry - [github.com/Cogfoundry-ai/loomloom](https://github.com/Cogfoundry-ai/loomloom)

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Cogfoundry-ai/loomloom/1-overview)

---

## What It Does

LoomLoom is a CLI and agent skill package for AI-driven batch content workflows. Instead of writing workflow code by hand, describe the task in natural language and let an agent download templates, prepare data, submit runs, watch progress, and download results.

Common use cases:

- **Batch copywriting** - product descriptions, rewrites, summaries, Q&A, and file-level text changes.
- **Batch image generation** - e-commerce imagery, social assets, concept art, and row-by-row visual generation.
- **Batch video generation** - storyboards, ad assets, and text-to-video workflows.

---

## Supported Agents

Installing LoomLoom adds a matching skill package for the selected agent:

| Agent | Status |
| --- | --- |
| **Codex** (OpenAI) | Supported |
| **Claude Code** (Anthropic) | Supported |
| **OpenClaw** | Supported |

---

## Quick Start

### Agent-assisted setup

Send a message like this to Codex, Claude Code, or OpenClaw. Replace `your-token` with the token you create on the CogFoundry API Keys page.

```text
Install LoomLoom from this GitHub repository: https://github.com/Cogfoundry-ai/loomloom
My server URL is https://loomloom.cogfoundry.ai/loom/v1, and my token is your-token.
After installation, run doctor to check whether the setup is healthy.
```

### Manual installation

macOS / Linux:

```bash
# Default install with the Codex skill package
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash

# Claude Code
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --agent claude

# OpenClaw
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --agent openclaw

# Specific version
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --version v0.1.0-beta.1

# Latest beta or internal channel
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta

# Specific pre-release tag
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --version v0.1.0-beta.1
```

Windows PowerShell:

```powershell
# Default install
irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1 | iex

# Claude Code
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Agent claude

# OpenClaw
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Agent openclaw

# Latest beta or internal channel
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Channel beta

# Specific pre-release tag
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Version v0.1.0-beta.1
```

Homebrew distribution is planned, but the tap repository and the token required for publishing are not yet configured. For now, use the install script above.

---

## Configure Credentials

```bash
export LOOMLOOM_SERVER="https://loomloom.cogfoundry.ai/loom/v1"
export LOOMLOOM_TOKEN="<your CogFoundry token>"
```

Add these values to `~/.zshrc`, `~/.bashrc`, or your shell profile if you do not want to set them again for every session. The CLI still accepts the legacy `BATCHJOB_SERVER` and `BATCHJOB_TOKEN` variables for compatibility, but new setups should use `LOOMLOOM_*`.

Create or get a token on the [CogFoundry API Keys](https://console-dev.cogfoundry.ai/api-keys) page. Never expose a real token in documentation, screenshots, logs, or public conversations.

> Security requirement: Only send your token to the `LOOMLOOM_SERVER` URL you explicitly configured (or the `--server` value you pass), and use HTTPS. The production default base URL is `https://loomloom.cogfoundry.ai/loom/v1`. Do not send the token to a host you did not specify, and do not follow redirects that point the token at a different domain. Use a token issued for the environment you are targeting.

---

## Verify Installation

```bash
loomloom doctor
```

If the environment is healthy, you can start submitting template runs.

---

## Templates, Private Templates, and SkillBots

"Template" is the umbrella term for a repeatable AI workflow in LoomLoom. The current product distinguishes two kinds of templates and one Market publishing path:

```text
Template
├─ Official template
│  ├─ Maintained and published by the platform
│  ├─ Discoverable and executable by any permitted user
│  └─ Uses the loomloom template ... commands
│
└─ Private template
   ├─ Created and maintained by a user through TemplateSpec
   ├─ One template can have multiple immutable versions
   ├─ The creator can execute their own template versions directly
   └─ A version can be submitted to the Market for review
          ↓
       Market Listing / Listing Version
          ↓ approved and made public
       SkillBot
```

What these names mean:

- **Official template**: A platform-maintained public execution entry, such as `text-v1`. It is not a local template hardcoded into the CLI; the CLI reads the available list from the current LoomLoom service.
- **Private template**: A user's own working asset. Created with TemplateSpec and saved as a private template, then changed by adding new immutable versions.
- **Custom template**: Describes how it is created, not a third kind of template. When a user finishes custom authoring, the result is a private template.
- **SkillBot**: The public, paid, executable form of a private template version after it passes Market review.
- **Listing**: The Market shelf object for a SkillBot; one Listing can publish multiple versions over time.
- **Listing Version**: An immutable execution snapshot copied from a specific private template version at publish time. Later changes to the private template do not automatically change a live SkillBot.

There is currently no separate "public template" resource or command. To refer to a public executable object, explicitly say "official template" or "published Market SkillBot" to avoid confusion.

`loomloom asset list` is an aggregated view of executable assets; it currently merges "my private templates" and "Market SkillBots". It is not a new kind of template, and it does not replace the official template list from `loomloom template list`.

---

## Current Official Templates

The following are current official template examples. The templates actually available depend on the target environment; rely on the live result of `loomloom template list`.

| Template ID | Use case | Output | Steps |
| --- | --- | --- | --- |
| `text-v1` | Copywriting, rewriting, summaries, Q&A, and code review | Text / files | Text generation |
| `text-image-v1` | Illustrations, concept art, and social images | Image | Prompt preparation -> image generation |
| `text-image-video-v1` | Storyboards, ads, and short video assets | Image + video | Description -> image -> video |

---

## Standard Excel Workflow (Official Templates)

```bash
# 1. Download a workbook template
loomloom template download text-image-v1 --output-file ./task.xlsx

# 2. Fill the workbook and validate it
loomloom template validate-file text-image-v1 ./task.xlsx

# 3. Submit the workbook
loomloom template submit-file text-image-v1 ./task.xlsx

# 4. Watch progress
loomloom run watch <run-id>

# 5. Download the server-generated result workbook
loomloom run result-workbook <run-id> --output-file ./task.result.xlsx

# 6. Download generated artifacts
loomloom artifact download <run-id> --output-dir ./downloads
```

`template backfill-results` is still available for older local workflows. For new workflows, prefer `run result-workbook`; the server uses the submitted input snapshot to align original rows with results.

---

## Template Fields

### Text template: `text-v1`

| Field | Required | Description |
| --- | --- | --- |
| Text prompt | Required | Main task prompt, such as "Rewrite this introduction in 80-120 words." |
| Writing requirements | Optional | Style, format, or output constraints. |
| Reference text | Optional | Inline short text, or upload a large file with `input-asset upload` and use the returned `input_asset_id`. |

### Image template: `text-image-v1`

| Field | Required | Description |
| --- | --- | --- |
| Image prompt | Required | Description of the image to generate. |
| Style requirements | Optional | For example, watercolor, photorealistic, or studio style. |
| Image aspect ratio | Required | `1:1`, `4:5`, `16:9`, or `9:16`. |

### Video template: `text-image-video-v1`

| Field | Required | Description |
| --- | --- | --- |
| Scene description | Required | Description of the video scene. |
| Visual style requirements | Optional | For example, cinematic tone or anime style. |
| Reference image URL | Optional | One public HTTP/HTTPS image URL. |
| Image aspect ratio | Required | `1:1`, `4:5`, `16:9`, or `9:16`. |
| Video aspect ratio | Required | `16:9` or `9:16`. |
| Video duration | Required | `4`, `6`, or `8` seconds. |
| Generate audio | Required | `false` or `true`. |

---

## Input Assets

For large reference files, upload the file first and place the returned `input_asset_id` only in a template field whose schema accepts an asset reference.

```bash
loomloom input-asset upload ./brief.txt --content-type text/plain
loomloom input-asset upload ./diagram.png --content-type image/png
```

Input assets and orchestration inputs are different things: `input-asset upload` returns an `input_asset_id` for reference material placed inside a template field, while `orchestration-input upload` returns an `input_file_id` that supplies the row data for a `template-spec run`.

An orchestration input file is JSONL. For the common single-root workflow, each non-empty line can be a flat JSON object whose values are strings:

```jsonl
{"prompt":"first request"}
{"prompt":"second request"}
```

The backend also supports unified rows shaped as `steps.<step-id>.executions[]` when a workflow needs explicit per-step execution input. In both formats, execution parameter values must be strings and must match the private template version's allowed input parameters. Do not invent step IDs; use unified input only when the exact workflow step mapping is available.

---

## Run Status

Use the [CogFoundry Console](https://console-dev.cogfoundry.ai/quickstart) to inspect run progress online.

| Status | Meaning |
| --- | --- |
| `pending` / `queued` | The run was accepted and is waiting for execution. |
| `running` | The run is in progress. |
| `completed` | All tasks completed successfully and results are available. |
| `partially_failed` | Some tasks failed, but successful results can still be downloaded. |
| `failed` | The run failed. |
| `cancelled` | The run was cancelled. |

---

## Command Reference

Monetary values such as `taskFixedFeeT` and `amountT` are in API units, where 10,000,000 units equal 1 currency unit.

### Diagnostics

| Command | Description |
| --- | --- |
| `loomloom doctor` | Check server reachability, token wiring, and version info. |

### Inputs

| Command | Description |
| --- | --- |
| `loomloom input-asset upload <file>` | Upload a reusable raw input asset (text/image) and get an `input_asset_id`. |
| `loomloom orchestration-input upload <file.jsonl>` | Upload flat JSONL rows and get the `input_file_id` required by `template-spec precheck` and `template-spec run`. |

### Official templates

| Command | Description |
| --- | --- |
| `loomloom template list` | List official templates published in the current environment. |
| `loomloom template schema <id>` | Show template fields. |
| `loomloom template download <id>` | Download an Excel workbook template. |
| `loomloom template validate-file <id> <xlsx>` | Validate a filled workbook. |
| `loomloom template precheck-file <id> <xlsx>` | Estimate cost for a workbook without submitting. |
| `loomloom template submit-file <id> <xlsx>` | Submit a filled workbook as a run. |
| `loomloom template backfill-results <run-id> <xlsx>` | Legacy local result backfill. |

### Private templates (created via TemplateSpec)

| Command | Description |
| --- | --- |
| `loomloom template-spec check <spec.json>` | Validate a TemplateSpec used to create a private template. |
| `loomloom template-spec docs [topic]` | Show bundled TemplateSpec documentation. |
| `loomloom template-spec models <step-type>` | List models for a step type. |
| `loomloom template-spec create <spec.json>` | Create a private template. |
| `loomloom template-spec create-version <template-id> <spec.json>` | Add a new version to an existing private template. |
| `loomloom template-spec list` | List my private templates. |
| `loomloom template-spec get <template-id>` | Show one private template and its versions. |
| `loomloom template-spec versions <template-id>` | List versions of a private template. |
| `loomloom template-spec download-workbook <template-id> <version-id>` | Download a user-template workbook. |
| `loomloom template-spec validate-workbook <template-id> <version-id> <xlsx>` | Validate a user-template workbook. |
| `loomloom template-spec precheck-workbook <template-id> <version-id> <xlsx>` | Estimate cost and balance for a user-template workbook without submitting. |
| `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx>` | Submit a user-template workbook. |
| `loomloom template-spec precheck <template-id> --version-id <id> --input-file-id <id>` | Estimate cost and balance for an uploaded JSONL input without submitting. |
| `loomloom template-spec run <template-id> --version-id <id> --input-file-id <id>` | Run a private template version from an uploaded JSONL input. |

### Runs

| Command | Description |
| --- | --- |
| `loomloom run submit <id> -f rows.json` | Submit input from a JSON array or JSONL file. |
| `loomloom run list` | List runs with optional Market context. |
| `loomloom run get <run-id>` | Show one run's detail. |
| `loomloom run watch <run-id>` | Watch run progress until a terminal state. |
| `loomloom run result-rows <run-id>` | Show aligned input rows and results. |
| `loomloom run result-workbook <run-id>` | Download a server-generated result workbook. |

### Artifacts

| Command | Description |
| --- | --- |
| `loomloom artifact list <run-id>` | List generated artifacts. |
| `loomloom artifact download <run-id>` | Download generated artifacts. |

### Catalog

| Command | Description |
| --- | --- |
| `loomloom model list --step-type <type>` | List executable models for a step type. |
| `loomloom asset list` | Aggregated list of my private templates and available Market SkillBots; does not include official templates. |

### Market — buyer

| Command | Description |
| --- | --- |
| `loomloom market list` | Browse published Market SkillBots. |
| `loomloom market show <listing-id>` | Show one SkillBot, including its input schema. |
| `loomloom market quote <listing-id> --input-file <json>` | Estimate execution cost. |
| `loomloom market run <listing-id> --input-file <json> --confirm` | Execute a SkillBot (paid). |
| `loomloom usage list` | List my Market SkillBot usage records. |
| `loomloom usage get <run-transaction-id>` | Show one usage record. |

### Market — creator

| Command | Description |
| --- | --- |
| `loomloom listing publish <template-id> --template-version-id <id> --display-name <name> --task-fixed-fee-t <fee>` | Submit a template version for Market review. |
| `loomloom listing publish <template-id> --listing-id <listing-id> --template-version-id <new-id> ...` | Submit a new version for an existing listing. |
| `loomloom listing list` | List my Market listings. |
| `loomloom listing show <listing-id>` | Show one of my listings. |
| `loomloom listing versions <listing-id>` | List versions of one of my listings. |
| `loomloom listing update <listing-id> --display-name <name>` | Submit a public-profile update for review; pass a display name, description, or both. |
| `loomloom listing unlist <listing-id>` | Stop new executions of a listing. |
| `loomloom listing relist <listing-id>` | Restore a previously unlisted listing. |
| `loomloom listing withdraw <listing-id>` | Withdraw the pending review request for a listing. |
| `loomloom creator earnings` | List Market earnings. |
| `loomloom creator transactions` | List Market transactions. |
| `loomloom creator review list` | List my review requests. |
| `loomloom creator review get <review-request-id>` | Show one review request. |
| `loomloom creator review withdraw <review-request-id>` | Withdraw a pending review request. |

---

## Private Template Authoring

Users describe a custom workflow's steps, input fields, and field bindings with TemplateSpec JSON; once created, it is saved as a private template. A typical agent-assisted authoring flow is:

```bash
# 1. Check available models for an execution unit
loomloom template-spec models text-generate

# 2. Validate the spec locally
loomloom template-spec check ./my-template.spec.json

# 3. Create a private template
loomloom template-spec create ./my-template.spec.json --version-note "initial version"

# 4. Add a new version when the template changes
loomloom template-spec create-version <template-id> ./my-template.spec.json

# 5. Download, fill, validate, precheck, and submit the workbook
loomloom template-spec download-workbook <template-id> <version-id> --output-file ./input.xlsx
loomloom template-spec validate-workbook <template-id> <version-id> ./input.xlsx
loomloom template-spec precheck-workbook <template-id> <version-id> ./input.xlsx
loomloom template-spec submit-workbook <template-id> <version-id> ./input.xlsx
```

Notes:

- TemplateSpec JSON is the source of truth; workbooks are generated artifacts.
- Review the bundled spec with `loomloom template-spec docs spec` before writing custom specs.
- Use `loomloom template-spec docs examples` for patterns.
- Use `loomloom template-spec docs conversation` for agent-assisted conversational authoring.
- Template changes require downloading a new workbook.
- `precheck-workbook` estimates model/API cost and balance; it does not create a run.
- Precheck text output includes `estimated_cost`, `available_balance`, and `sufficient`; JSON output uses `estimatedTotalCostT`.
- `submit-workbook` creates a real hosted run; agents should ask for explicit confirmation before submitting.
- `template-spec run` also creates a real hosted run and requires the same confirmation.

You can also run a private template version directly from flat JSONL rows, without filling a workbook:

```bash
# 1. Upload the rows and capture the returned input_file_id
loomloom orchestration-input upload ./rows.jsonl

# 2. Estimate cost and balance
loomloom template-spec precheck <template-id> --version-id <version-id> --input-file-id <input_file_id>

# 3. Run the version with that input
loomloom template-spec run <template-id> --version-id <version-id> --input-file-id <input_file_id>
```

---

## Market / SkillBot

LoomLoom Market lets a creator publish a private template version as a paid SkillBot, and lets a buyer run it. The same CLI serves both roles.

Publishing does not turn the private template itself into a public object. The system generates an immutable Listing Version execution snapshot from the chosen private template version. Later edits to the private template, or new versions, do not automatically change the live SkillBot; you must submit a new version for review again.

### Buyer: discover and run a SkillBot

```bash
# 1. Browse and inspect SkillBots
loomloom market list --keyword "tweet"
loomloom market show <listing-id>

# 2. Estimate cost (the input file carries a taskInputs array shaped to the listing schema)
loomloom market quote <listing-id> --input-file ./request.json

# 3. Execute (paid; --confirm is required)
loomloom market run <listing-id> --input-file ./request.json --confirm

# 4. Review your own usage
loomloom usage list
loomloom usage get <run-transaction-id>
```

Example `request.json`:

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

Use the listing detail to understand public fields and obtain `listingVersionId`, but do not infer internal step IDs from `inputSchemaSnapshot`: the public snapshot may not expose them. The current CLI expects an exact Product API `taskInputs` payload. If the user or integration does not provide that mapping, stop and ask for a compatible request JSON instead of guessing.

### Creator: publish and manage a SkillBot

```bash
# 1. Publish a template version for review (it must already have one successful run)
loomloom listing publish <template-id> \
  --template-version-id <version-id> \
  --display-name "My SkillBot" \
  --task-fixed-fee-t 1000000

# Submit a new version for an existing listing
loomloom listing publish <template-id> \
  --listing-id <listing-id> \
  --template-version-id <new-version-id> \
  --display-name "My SkillBot" \
  --task-fixed-fee-t 1000000

# 2. Track listings and reviews
loomloom listing list
loomloom listing versions <listing-id>
loomloom creator review list

# 3. Manage sale status and public profile
loomloom listing unlist <listing-id>
loomloom listing relist <listing-id>
loomloom listing update <listing-id> --description "Updated description"

# 4. Review income
loomloom creator earnings
loomloom creator transactions
```

Notes:

- `market run` creates a real paid run; agents should run `market quote` first and ask for explicit confirmation before executing.
- `listing publish` and `listing update` submit changes for review; they do not take effect until approved.
- `listing publish --listing-id <listing-id>` submits a new version for the existing listing. The currently published version stays active until the new review is approved.
- `listing update`, `listing unlist`, `listing relist`, and review withdrawal change remote state. Agents should summarize the action and ask for explicit confirmation before invoking them.
- Monetary `*FeeT` / `*AmountT` / `*PayableT` values are in API units where 10,000,000 units equal 1 currency unit.

---

## Agent Command Chaining

Agents should prefer `--output json` for commands whose results feed later commands. Preserve IDs exactly; never infer or transform them.

```text
orchestration-input upload → inputFileId → template-spec precheck → template-spec run
template-spec run / run submit → runId → run watch / result commands
listing publish → reviewRequestId → creator review get/withdraw
market run → runTransactionId and runId → usage get / run watch
```

Text output uses labels such as `input_file_id`; JSON output uses Product API field names such as `inputFileId`.

For `template submit-file`, `template-spec submit-workbook`, `run submit`, `template-spec run`, and `market run`, pass an explicit `--client-request-id`, retain it with the request, and reuse it only when retrying the identical payload. Use a new ID if the payload changes. Do not blindly retry paid or remote-state-changing commands after an ambiguous failure; first check whether the original request succeeded.

---

## Uninstall

macOS / Linux:

```bash
# Remove CLI and skill package
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash

# Remove CLI only
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash -s -- --cli-only

# Remove skill package only
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.sh | bash -s -- --skill-only
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/uninstall.ps1 | iex
```

---

## FAQ

**Where do I get a token?**

Create or get a token on the [CogFoundry API Keys](https://console-dev.cogfoundry.ai/api-keys) page, then replace the placeholder in the examples.

Use your token only with the `LOOMLOOM_SERVER` you configured. Do not send it to a host you did not explicitly set.

**Where do I check run status?**

Use the [CogFoundry Console](https://console-dev.cogfoundry.ai/quickstart). There is currently no URL template for a Workflow Run detail page, so do not construct run-detail links yourself.

**Why does `template list` return no templates?**

The account or environment may not have visible templates. Ask your CogFoundry workspace administrator to confirm template publication and permissions.

**Can I use the CLI manually without an agent?**

Yes. Every workflow can be run with the CLI commands above.

---

## Links

- GitHub: [github.com/Cogfoundry-ai/loomloom](https://github.com/Cogfoundry-ai/loomloom)
- LoomLoom API: `https://loomloom.cogfoundry.ai/loom/v1`
- API Keys: [console-dev.cogfoundry.ai/api-keys](https://console-dev.cogfoundry.ai/api-keys)
- CogFoundry Console: [console-dev.cogfoundry.ai/quickstart](https://console-dev.cogfoundry.ai/quickstart)
- CogFoundry website: [cogfoundry.ai](https://cogfoundry.ai)
- Workflow Run detail page: no fixed URL template is available yet.
