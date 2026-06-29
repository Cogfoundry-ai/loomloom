---
name: loomloom
description: Use this skill when the user mentions LoomLoom, batchjob, batchflow, batch processing, bulk processing, template submission, Excel template batch execution, run submit, artifact downloads, result backfill, Market SkillBots, publishing or listing a template on the Market, running a Market SkillBot, creator earnings, or Market usage records.
---

# loomloom

Use this skill when the user wants to use LoomLoom for batch content generation, template submission, Excel batch execution, result inspection, or artifact download. Legacy names such as `batchjob`, `batchflow`, batch processing, batch tasks, batch generation, and running templates in bulk all refer to LoomLoom workflows.

## When To Use

- The user wants to generate copy, images, or videos in batches, including text-to-image or text-to-image-to-video batch workflows.
- The user wants to list templates, download official Excel templates, validate Excel files, submit Excel files, check run status, download result workbooks, or download artifacts.
- The user wants to create a custom template, especially from a natural-language reusable workflow description.
- The user wants to discover, quote, or run a published Market SkillBot.
- The user wants to publish a template to the Market, manage their listings, or check creator earnings and review requests.
- The task can be represented as structured rows rather than a one-off chat answer.
- The user is willing to use developer tooling or have an agent call the CLI.

## When Not To Use

- The user only needs one immediate generation, not batch processing.
- The request is still exploratory and cannot be shaped into row-level input or a template.
- `LOOMLOOM_TOKEN` is not configured and the user does not want to handle environment setup.

## Core Objects and Command Selection

Do not treat official templates, private templates, and SkillBots as three unrelated workflow systems. Their relationship is:

```text
Official template ── platform-maintained, executed directly

Private template ── created and maintained by a user through TemplateSpec
   └─ Private template version
        └─ Submitted to the Market for review
             └─ Listing Version (immutable publish snapshot)
                  └─ After approval, executable by buyers as a SkillBot
```

Terminology and selection rules:

- **Template** is the umbrella term.
- **Official template** is maintained by the platform. Use the `template` command group to discover and execute official templates.
- **Custom template** describes how something is created; the result is the user's **private template**. Use the `template-spec` command group to author, inspect versions, and directly execute private templates.
- **SkillBot** is the public, paid, executable form of a private template version after it passes Market review. Buyers use the `market` command group.
- **Listing** is the Market shelf object for a SkillBot. Creators use the `listing` and `creator review` command groups to publish and manage SkillBots.
- **Listing Version** is an immutable execution snapshot copied from a private template version at publish time. Later changes to the private template do not automatically update a live SkillBot.
- There is currently no separate "public template" resource. Do not use that term to refer to official templates or Market SkillBots.
- `asset list` is only an aggregated view of executable assets ("my private templates + Market SkillBots"); it is not a new kind of template and does not include official templates.

Choose the entry point by user intent:

- "what platform templates are there", "execute with an official Excel template" → `template list/schema/download/...`
- "create my workflow", "modify my template", "run my template version" → `template-spec ...`
- "publish my template as a SkillBot", "update or unlist my SkillBot" → `listing ...`
- "find a SkillBot and run it for a fee" → `market ...`
- When the user only says "template" and the context cannot tell official from private, clarify first; do not guess.

## Command Flow

0. Install via the GitHub install script by default. On macOS/Linux the default install uses Homebrew. Unless the user explicitly asks not to use Homebrew, do not add `--no-brew`.
   If the user asks for an internal/beta CLI, explicitly install a prerelease channel instead of defaulting to stable:
   `curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta`
   If resolving the latest version fails because the GitHub release API is rate-limited (HTTP 403), retry after a short wait, install a specific version with `--version <tag>` (which skips the rate-limited lookup), or ask the user.
1. Check the environment:
   `loomloom doctor`
   The production default base URL is `https://loomloom.cogfoundry.ai/loom/v1`, but the active server is whatever the user sets in `LOOMLOOM_SERVER` / `--server`. Send the token only to that explicitly configured host, and only over HTTPS.
   Before each request, check the scheme and host of the final target URL. Only send the token to the host the user configured; do not send it to a host the user did not specify, and do not auto-follow redirects to a different domain.
   Use a token issued for the environment you are targeting; do not reuse a token across environments it was not issued for.
   If the token is missing, guide the user to `https://console-dev.cogfoundry.ai/api-keys` to create or get one. Do not echo a real token in replies, logs, or generated files.
2. Do not paste large files into context. Upload raw input assets first:
   `loomloom input-asset upload <file>`
3. Discover the official templates in the current environment:
   `loomloom template list`
4. Inspect official template fields, list available models, or list executable assets:
   `loomloom template schema <template-id>`
   `loomloom model list --step-type <text-generate|image-generate|video-generate>`
   `loomloom asset list`
   Note: `asset list` aggregates my private templates and Market SkillBots; it does not include official templates.
5. When authoring and saving a private template, use TemplateSpec JSON:
   `loomloom template-spec docs spec`
   `loomloom template-spec docs examples`
   `loomloom template-spec docs conversation`
   `loomloom template-spec models <text-generate|image-generate|video-generate>`
   `loomloom template-spec check <spec.json>`
   `loomloom template-spec create <spec.json>`
   `loomloom template-spec create-version <template-id> <spec.json>`
   `loomloom template-spec download-workbook <template-id> <version-id>`
   `loomloom template-spec validate-workbook <template-id> <version-id> <xlsx-path>`
   `loomloom template-spec precheck-workbook <template-id> <version-id> <xlsx-path>`
   `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx-path>`
   Inspect existing private templates:
   `loomloom template-spec list`
   `loomloom template-spec get <template-id>`
   `loomloom template-spec versions <template-id>`
   To run a private template version directly from flat JSONL rows, upload the rows first and pass the returned input_file_id:
   `loomloom orchestration-input upload <file.jsonl>`
   `loomloom template-spec precheck <template-id> --version-id <version-id> --input-file-id <input_file_id>`
   `loomloom template-spec run <template-id> --version-id <version-id> --input-file-id <input_file_id>`
   Use `precheck-workbook` before `submit-workbook`, and `precheck` before `template-spec run`, to estimate model/API cost and balance without creating a run.
   For common single-root workflows, each non-empty line may be a flat JSON object with string values. Unified rows using `steps.<step-id>.executions[]` are also supported when exact workflow step mappings are available. In either format, execution parameter values must be strings and allowed by the private template version. Never guess step IDs.
6. If the user did not explicitly ask to create or use their own private template, use the official Excel workflow by default:
   `loomloom template download <template-id>`
   `loomloom template validate-file <template-id> <xlsx-path>`
   `loomloom template submit-file <template-id> <xlsx-path>`
   `loomloom run result-workbook <run-id>`
   Use the older local backfill flow only when the user explicitly needs it:
   `loomloom template backfill-results <run-id> <xlsx-path>`
7. Use JSON/JSONL only when the user explicitly asks for programmatic input:
   `loomloom run submit <template-id> -f rows.jsonl`
8. Watch progress:
   `loomloom run watch <run-id>`
9. List, inspect, or download runs and their results:
   `loomloom run list`
   `loomloom run get <run-id>`
   `loomloom run result-rows <run-id>`
   `loomloom run result-workbook <run-id>`
10. Inspect or download artifacts:
   `loomloom artifact list <run-id>`
   `loomloom artifact download <run-id>`

## Market Ecosystem

LoomLoom Market lets creators publish their private template versions as paid SkillBots, and lets buyers run those SkillBots. The same CLI serves two roles; decide which role the user is in before choosing commands.

At publish time, the Market copies an immutable Listing Version execution snapshot from the chosen private template version. Later changes to the private template do not automatically change a live SkillBot; to update the live execution version, submit a new template version for the same Listing and have it reviewed again.

### Buyer role (use and pay for a SkillBot)

- `loomloom market list` — browse published SkillBots. Supports `--keyword`, `--page-size`, `--page-token`, and `--order-by`.
- `loomloom market show <listing-id>` — show one SkillBot, including its `inputSchemaSnapshot`. Read the schema before building input.
- `loomloom market quote <listing-id> --input-file <request.json>` — estimate cost. Returns `estimatedBuyerPayableT`, `taskCount`, and `taskFixedFeeT`.
- `loomloom market run <listing-id> --input-file <request.json> --confirm` — execute the SkillBot. This is a paid action; see the Submission Confirmation Rule.
- `loomloom usage list` and `loomloom usage get <run-transaction-id>` — review the buyer's own SkillBot calls and settlement status. Use the returned `runTransactionId`.

The `--input-file` JSON for quote and run carries a `taskInputs` array shaped to the listing's input schema, for example:

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

Read `market show` first to understand public fields and obtain `listingVersionId`, but do not infer internal step IDs from `inputSchemaSnapshot`: it may not expose them. The current CLI expects an exact Product API `taskInputs` payload. If no compatible mapping or request JSON is available, stop and ask for it instead of guessing.

### Creator role (publish and manage a SkillBot)

- `loomloom listing publish <template-id> --template-version-id <id> --display-name <name> --task-fixed-fee-t <fee>` — submit a template version for Market review. The version must already have at least one successful run. Returns a `reviewRequestId` with `reviewStatus: pending_review`.
- To submit a new version for an existing listing, run the same command with `--listing-id <listing-id>` and the new `--template-version-id`. The published version stays active until approval.
- `loomloom listing list`, `loomloom listing show <listing-id>`, `loomloom listing versions <listing-id>` — inspect the creator's own listings, including pending, rejected, and unlisted states.
- `loomloom listing update <listing-id> --display-name <name> --description <text>` — submit a public-profile change for review; it does not change pricing or the execution version.
- `loomloom listing unlist <listing-id>` and `loomloom listing relist <listing-id>` — stop or resume new executions of a published listing.
- `loomloom listing withdraw <listing-id>` — withdraw the single pending review request for that listing. If none exists, stop. If multiple are reported, use `creator review list` and withdraw the intended request explicitly.
- `loomloom creator review list`, `loomloom creator review get <review-request-id>`, `loomloom creator review withdraw <review-request-id>` — track and withdraw review requests directly.
- `loomloom creator earnings` and `loomloom creator transactions` — review creator income and per-call settlement.

All `*FeeT`, `*CostT`, `*AmountT`, and `*PayableT` values are in API units where 10,000,000 units equal 1 currency unit.

## Submission Confirmation Rule

Any command that actually submits work to the hosted LoomLoom service must receive a second explicit confirmation from the user in the current conversation.

Treat the interaction as one of three states:

1. `default-prep`: the user is still exploring or speaking generally. Prepare, download, and validate only. Do not submit.
2. `auto-run-candidate`: the user explicitly asks the agent to execute. Still do not submit. First provide an execution summary and wait for confirmation.
3. `confirmed-to-run`: after seeing the execution summary, the user explicitly confirms. Only then may you submit.

This rule applies to `loomloom template submit-file`, `loomloom template-spec submit-workbook`, `loomloom run submit`, `loomloom template-spec run`, and `loomloom market run`. For `loomloom market run`, first run `loomloom market quote` and include the returned estimate in the execution summary. Validation, precheck commands, downloads, schema inspection, model lookup, quoting, `doctor`, asset upload, orchestration-input upload, artifact listing, listing/usage/earnings reads, and result backfill do not start a new paid run and do not require the second confirmation.

The execution summary must include the template or listing ID, input source, row count or task size, action, estimated cost or a clear cost note, and the prompt `Reply "confirm submit" before I start.` For `template submit-file`, `template-spec submit-workbook`, `run submit`, `template-spec run`, and `market run`, pass an explicit stable `--client-request-id` and retain it for safe retry of the identical payload.

Present the confirmation summary in plain business language (what will happen, which template or SkillBot, how many tasks, the cost). Do not show the raw CLI command in the confirmation unless the user explicitly asks to see it.

## Remote State Change Confirmation Rule

Before creating or changing persistent remote resources, show the exact action and ask for explicit confirmation. This applies to `template-spec create`, `template-spec create-version`, `listing publish` (including `--listing-id` version updates), `listing update`, `listing unlist`, `listing relist`, `listing withdraw`, and `creator review withdraw`.

Describe the action in plain business language; "show the exact action" does not mean printing the CLI command. Do not show the raw CLI command unless the user explicitly asks to see it.

Read-only commands, local checks, uploads, downloads, and quotes do not require this confirmation. A paid run still follows the stricter Submission Confirmation Rule above.

## Agent Command Chaining

Prefer `--output json` whenever one command feeds another. Preserve exact fields:

- `orchestration-input upload` → `inputFileId` → `template-spec precheck --input-file-id` → `template-spec run --input-file-id`
- run submission → `runId` → run watch/result commands
- `listing publish` → `reviewRequestId` → creator review commands
- `market run` → `runTransactionId` and `runId` → usage/run commands

Never convert `inputAssetId` (`ia_xxx`) into `inputFileId`, and never guess IDs from names.
For the five supported submission commands listed in the Submission Confirmation Rule, pass an explicit `--client-request-id`, retain it with the request, and reuse it only for an identical retry. A changed payload requires a new ID.

## Error Handling

Inspect the returned error before choosing a recovery step:

- For local flag, file, JSON, or schema errors, correct the input and retry without running `doctor`.
- For authentication, endpoint, network, service-version, or unexpected server errors, run `loomloom doctor`.
- Never invent missing IDs, hidden step IDs, or server state.
- Do not blindly retry paid or remote-state-changing commands after an ambiguous failure. First query the relevant run, listing, or review state. If the identical submission must be retried, reuse its original `--client-request-id`; use a new ID only for a changed payload.

## Current MVP Capabilities

The public CLI currently supports these command groups:

- Environment: `doctor`.
- Inputs: `input-asset upload`, `orchestration-input upload`.
- Discovery: `template list`, `template schema`, `model list`, `asset list`.
- Official Excel workflow: `template download`, `template validate-file`, `template precheck-file`, `template submit-file`, `template backfill-results`.
- Custom templates: `template-spec docs`, `template-spec check`, `template-spec models`, `template-spec create`, `template-spec create-version`, `template-spec list`, `template-spec get`, `template-spec versions`, `template-spec download-workbook`, `template-spec validate-workbook`, `template-spec precheck-workbook`, `template-spec submit-workbook`, `template-spec precheck`, `template-spec run`.
- Runs: `run submit`, `run list`, `run get`, `run watch`, `run result-rows`, `run result-workbook`.
- Artifacts: `artifact list`, `artifact download`.
- Market (buyer): `market list`, `market show`, `market quote`, `market run`, `usage list`, `usage get`.
- Market (creator): `listing publish`, `listing list`, `listing show`, `listing versions`, `listing update`, `listing unlist`, `listing relist`, `listing withdraw`, `creator earnings`, `creator transactions`, `creator review list`, `creator review get`, `creator review withdraw`.

## Large File Handling

When the user wants to batch-process local code files, large text, local images, or other large files, avoid pasting full files into the agent context. Prefer:

1. `loomloom input-asset upload <file>`
2. Save the returned `input_asset_id`
3. Put it only in a schema field that accepts an asset reference

For TemplateSpec, use a compatible `asset_ref` or `text_reference` field and follow the bundled binding guidance. An `input_asset_id` is reference material; it is never the `inputFileId` used by `template-spec run`.

## Default Behavior

Unless the user explicitly asks for JSON/JSONL, use the official Excel template workflow by default.

When the user asks to create or customize a workflow/template, prefer TemplateSpec JSON. TemplateSpec JSON is source data; downloaded workbooks are derived artifacts. When a template version changes, do not promise that old workbooks remain compatible. Download a new workbook.

## Conversational Template Authoring

When the user describes a new template in natural language, do not immediately write TemplateSpec JSON. First run or reference `loomloom template-spec docs conversation`.

Flow:

1. Ask business questions, not TemplateSpec technical-field questions.
2. Ask one missing-detail question at a time.
3. Avoid user-facing technical terms such as `fieldBindings`, `upstreamBindings`, `fan-in`, `execution`, `outputSchema`, `provider`, and `mode`.
4. Restate complex workflows in business language before continuing.
5. After identifying workflow steps, roles, perspectives, or generated agents, determine whether each role/step has its own processing requirements. If yes, ask about future template usage before drafting TemplatePlan.
6. Draft TemplatePlan first.
7. Show TemplatePlan and wait for user confirmation.
8. Generate TemplateSpec JSON only after confirmation.
9. Before creation, run `loomloom template-spec check <spec.json>`.
10. After check passes, ask for separate creation confirmation.
11. Run `loomloom template-spec create <spec.json>` only after explicit creation confirmation.

When an existing user template needs a fix, do not promise in-place mutation of historical versions. Use `loomloom template-spec create-version <template-id> <spec.json>` to append a new version, and make later workbooks/runs use the new `version_id`.

### Template Usage Mode Selection

When a template has multiple roles, steps, review perspectives, or generated agents, and each role/step may have its own processing requirements, the agent must ask before generating TemplatePlan:

```text
When other people use this template later, how should the review/generation requirements be handled?

1. Preset in the template: users only fill the core material, and the system follows your predefined requirements automatically.
2. Let users fill them: users can fill or modify the requirements for each step/role.
3. Generate both versions: one simple version and one customizable version.

If you are unsure, choose 1 first to make a simple, usable standard template.
```

Do not expose technical concepts such as `prompt`, `binding`, `reference`, `field`, `hidden`, `paramBindings`, or `fieldBindings`. Typical trigger scenarios include multi-perspective PRD review, multi-role contract review, multi-agent launch event planning, rewriting an article in multiple styles, resume review from multiple interviewer perspectives, and generating marketing content for multiple channels.

Three modes:

1. Simple usage mode: users fill only core material, and the system executes with preset processing requirements.
2. Custom usage mode: users fill both core material and processing requirements for each step/role.
3. Generate both versions: generate two TemplatePlans and create two templates, naming them as `Standard Version` / `Custom Version`.

Automatic generation rules: in simple mode, role/step processing requirements are template preset instructions, not user-fillable columns; in custom mode, core material and role/step processing requirements are both user input columns; when generating both versions, show both TemplatePlans first and get user confirmation before generating two TemplateSpecs.

Acceptance criteria: multi-role/multi-step/multi-perspective template creation must ask about usage mode; simple mode hides role/step requirement columns; custom mode exposes and allows editing those columns; user-facing conversation does not expose technical concepts; generated TemplateSpec correctly expresses both simple and custom shapes.

Creation confirmation gate:

- "Create a PRD review template" only starts the flow; it does not confirm remote creation.
- Environment variables, token, and server URL are configuration, not creation confirmation.
- "Generate spec" only means generate and locally check the spec; it does not confirm running `template-spec create`.
- Before running `template-spec create`, describe in plain language what will be created (template name, what it does, the local check result) and ask the user to reply `confirm create template`. Do not show the raw CLI command unless the user explicitly asks to see it.

TemplatePlan should cover template name and goal, row meaning, input fields, workflow steps, serial/parallel/summary relationships, template usage mode, user-visible outputs, failure policy, error columns, default model, and special requirements.

Default modeling rules:

- "Product, engineering, and risk review separately, then summarize" should be modeled as multiple parallel `text-generate` steps plus one downstream summary step, using step-level `dependsOn` and `upstreamBindings`.
- Do not model multi-role review as one `expanded` step. `expanded` is only for dynamic multi-value input.
- Add result and error columns for every user-visible step by default.
- If partial completion is allowed, the summary step should explain missing upstreams and expose failed step error columns.
- Keep `provider` and `mode` internal.

TemplateSpec authoring constraints:

- Before writing custom TemplateSpec, run `loomloom template-spec docs spec` and use the CLI-bundled document as the current contract.
- Use `loomloom template-spec docs examples` for examples and patterns.
- Use `loomloom template-spec docs conversation` for the natural-language creation flow.
- Use only `text-generate`, `image-generate`, and `video-generate` by default.
- Use OpenAPI lowerCamel fields such as `meta.name`, `steps[].stepId`, and `defaultModelRef.modelKey`.
- Connect steps with step-level `dependsOn` and `upstreamBindings`; the source output port is usually `output`.
- Before choosing `defaultModelRef.modelKey`, run `loomloom template-spec models <execution-unit>` and use a returned `model_id`.
- Expose a model column only when the step sets `allowModelOverride=true` and a field binds to `paramKey=model`.
- Do not bind `provider` or `mode`.

## Trusted Result Sources

The submitted workbook and server-side run input snapshot are the sources of truth. After the run completes, prefer `loomloom run result-workbook <run-id>` because the server aligns original input snapshots and artifacts. Use `template backfill-results` only when the user explicitly needs the older local Excel backfill flow.

## Console Access

The CogFoundry Console entry point is `https://console-dev.cogfoundry.ai/quickstart`. When a user needs to check run status, you can give them this Console page.

There is currently no URL template for a Workflow Run detail page. Do not guess or construct a detail-page link from a `runId`; if the CLI returns a URL explicitly provided by the server, use it as-is, otherwise only provide the Console home page and the CLI query commands.

The CogFoundry website is `https://cogfoundry.ai`.
