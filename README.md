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

Send a message like this to Codex, Claude Code, or OpenClaw. Replace `your-token` and `<your LoomLoom server URL>` with the values from your CogFoundry workspace.

```text
Install LoomLoom from this GitHub repository: https://github.com/Cogfoundry-ai/loomloom
My server URL is <your LoomLoom server URL>, and my token is your-token.
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
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --version v0.2.7

# Latest beta or internal channel
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta --no-brew

# Specific pre-release tag
curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --version v0.2.6-beta.9 --no-brew
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
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.ps1))) -Version v0.2.6-beta.9
```

Homebrew distribution will be enabled after the CogFoundry tap is published. For now, use the install script above.

---

## Configure Credentials

```bash
export LOOMLOOM_SERVER="<your LoomLoom server URL>"
export LOOMLOOM_TOKEN="your-token"
```

Add these values to `~/.zshrc`, `~/.bashrc`, or your shell profile if you do not want to set them again for every session. The CLI still accepts the legacy `BATCHJOB_SERVER` and `BATCHJOB_TOKEN` variables for compatibility, but new setups should use `LOOMLOOM_*`.

Get the token from your CogFoundry workspace.

---

## Verify Installation

```bash
loomloom doctor
```

If the environment is healthy, you can start submitting template runs.

---

## Built-In Templates

| Template ID | Use case | Output | Steps |
| --- | --- | --- | --- |
| `text-v1` | Copywriting, rewriting, summaries, Q&A, and code review | Text / files | Text generation |
| `text-image-v1` | Illustrations, concept art, and social images | Image | Prompt preparation -> image generation |
| `text-image-video-v1` | Storyboards, ads, and short video assets | Image + video | Description -> image -> video |

---

## Standard Excel Workflow

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

For large reference files, upload the file first and place the returned `input_asset_id` in the template field.

```bash
loomloom input-asset upload ./brief.txt --content-type text/plain
loomloom input-asset upload ./diagram.png --content-type image/png
```

---

## Run Status

Use your CogFoundry workspace console URL to inspect run progress online.

| Status | Meaning |
| --- | --- |
| Queued | The run was submitted and is waiting for scheduling. |
| Running | The run is in progress. |
| Finished | All rows reached a terminal state and results are available. |
| Partially failed | Some rows failed, but successful rows can still be downloaded. |
| Failed | The run failed as a whole. |
| Expired | Results expired and the run should be submitted again if needed. |

---

## Command Reference

| Command | Description |
| --- | --- |
| `loomloom doctor` | Check local configuration. |
| `loomloom template list` | List available templates. |
| `loomloom template schema <id>` | Show template fields. |
| `loomloom template download <id>` | Download an Excel workbook template. |
| `loomloom template validate-file <id> <xlsx>` | Validate a filled workbook. |
| `loomloom template submit-file <id> <xlsx>` | Submit a filled workbook. |
| `loomloom run result-rows <run-id>` | Show aligned input rows and results. |
| `loomloom run result-workbook <run-id>` | Download a server-generated result workbook. |
| `loomloom template backfill-results <run-id> <xlsx>` | Legacy local result backfill. |
| `loomloom run submit <id> -f rows.jsonl` | Submit JSONL input. |
| `loomloom run watch <run-id>` | Watch run progress. |
| `loomloom artifact list <run-id>` | List generated artifacts. |
| `loomloom artifact download <run-id>` | Download generated artifacts. |
| `loomloom input-asset upload <file>` | Upload a large input asset. |
| `loomloom template-spec check <spec.json>` | Validate a custom TemplateSpec. |
| `loomloom template-spec docs [topic]` | Show bundled TemplateSpec documentation. |
| `loomloom template-spec models <step-type>` | List models for a step type. |
| `loomloom template-spec create <spec.json>` | Create a private template. |
| `loomloom template-spec create-version <template-id> <spec.json>` | Add a new version to an existing private template. |
| `loomloom template-spec download-workbook <template-id> <version-id>` | Download a user-template workbook. |
| `loomloom template-spec validate-workbook <template-id> <version-id> <xlsx>` | Validate a user-template workbook. |
| `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx>` | Submit a user-template workbook. |

---

## Custom Templates

Custom templates use TemplateSpec JSON to describe workflow steps, input fields, and field bindings. A typical agent-assisted authoring flow is:

```bash
# 1. Check available models for an execution unit
loomloom template-spec models text-generate

# 2. Validate the spec locally
loomloom template-spec check ./my-template.spec.json

# 3. Create a private template
loomloom template-spec create ./my-template.spec.json --version-note "initial version"

# 4. Add a new version when the template changes
loomloom template-spec create-version <template-id> ./my-template.spec.json

# 5. Download, fill, validate, and submit the workbook
loomloom template-spec download-workbook <template-id> <version-id> --output-file ./input.xlsx
loomloom template-spec validate-workbook <template-id> <version-id> ./input.xlsx
loomloom template-spec submit-workbook <template-id> <version-id> ./input.xlsx
```

Notes:

- TemplateSpec JSON is the source of truth; workbooks are generated artifacts.
- Review the bundled spec with `loomloom template-spec docs spec` before writing custom specs.
- Use `loomloom template-spec docs examples` for patterns.
- Use `loomloom template-spec docs conversation` for agent-assisted conversational authoring.
- Template changes require downloading a new workbook.
- `submit-workbook` creates a real hosted run; agents should ask for explicit confirmation before submitting.

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

Get a token from your CogFoundry workspace and replace `your-token` in the examples.

**Why does `template list` return no templates?**

The account or environment may not have visible templates. Ask your CogFoundry workspace administrator to confirm template publication and permissions.

**Can I use the CLI manually without an agent?**

Yes. Every workflow can be run with the CLI commands above.

---

## Links

- GitHub: [github.com/Cogfoundry-ai/loomloom](https://github.com/Cogfoundry-ai/loomloom)
- Console: use your CogFoundry workspace console URL.
