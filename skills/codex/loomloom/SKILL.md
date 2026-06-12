---
name: loomloom
description: Use this skill when the user mentions LoomLoom, batchjob, batchflow, batch processing, template submission, Excel template execution, run submit, artifact downloads, or result backfill.
---

# loomloom

Use this skill when the user wants to run batch content generation, submit templates, execute Excel workflows, check run status, download result workbooks, or download generated artifacts through LoomLoom. Legacy names such as `batchjob`, `batchflow`, and "batch processing" refer to the same LoomLoom workflow.

## When To Use

- The user wants to generate copy, images, or videos in batches.
- The user wants to list templates, inspect schemas, download official Excel templates, validate workbooks, submit workbooks, watch runs, download result workbooks, or download artifacts.
- The user wants to create a custom reusable workflow from natural language.
- The task can be represented as structured rows rather than a one-off chat answer.
- The user is willing to use developer tooling or have the agent call the CLI.

## When Not To Use

- The user only needs a single immediate generation.
- The request is still exploratory and cannot be shaped into rows or a template.
- `LOOMLOOM_TOKEN` is not configured and the user does not want to handle environment setup.

## Command Flow

0. If the user asks for a beta or internal CLI, install an explicit pre-release channel instead of stable:
   `curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta --no-brew`
1. Check the environment:
   `loomloom doctor`
   Do not rely on a local default server. Before TemplateSpec authoring, model discovery, or template submission, confirm that `LOOMLOOM_SERVER` is the target environment provided by the user. If it is missing, ask for the service URL.
2. For large inputs, upload the raw input asset instead of pasting full content into the agent context:
   `loomloom input-asset upload <file>`
3. Discover templates:
   `loomloom template list`
4. Inspect template fields:
   `loomloom template schema <template-id>`
5. For custom template authoring, use TemplateSpec JSON:
   `loomloom template-spec docs spec`
   `loomloom template-spec docs examples`
   `loomloom template-spec docs conversation`
   `loomloom template-spec models <execution-unit>`
   `loomloom template-spec check <spec.json>`
   `loomloom template-spec create <spec.json>`
   `loomloom template-spec create-version <template-id> <spec.json>`
   `loomloom template-spec download-workbook <template-id> <version-id>`
   `loomloom template-spec validate-workbook <template-id> <version-id> <xlsx>`
   `loomloom template-spec submit-workbook <template-id> <version-id> <xlsx>`
6. If the user did not explicitly ask for a custom template, use the official Excel workflow by default:
   `loomloom template download <id> --output-file <xlsx>`
   `loomloom template validate-file <id> <xlsx>`
   `loomloom template submit-file <id> <xlsx>`
   Use `loomloom run result-workbook <run-id>` after completion. Use `loomloom template backfill-results` only when the user explicitly needs the older local backfill flow.
7. Use JSON/JSONL only when the user explicitly asks for programmatic input:
   `loomloom run submit <id> -f rows.jsonl`
8. Watch progress:
   `loomloom run watch <run-id>`
9. Inspect or download aligned results:
   `loomloom run result-rows <run-id>`
   `loomloom run result-workbook <run-id> --output-file <xlsx>`
10. Inspect or download artifacts:
   `loomloom artifact list <run-id>`
   `loomloom artifact download <run-id> --output-dir <dir>`

## Submission Confirmation Rule

Any command that submits real work to the hosted LoomLoom service requires a second explicit confirmation from the user in the current conversation.

Handle interactions in three states:

1. Preparation only.
   The user is still exploring or speaking generally, such as "help me run a batch" or "you decide." Prepare, download, and validate only. Do not submit.
2. Execution requested.
   The user says something like "run it for me" or "submit it directly." Still do not submit yet. First show an execution summary and wait for confirmation.
3. Confirmed submission.
   The user replies after seeing the summary with a clear confirmation such as "confirm submit", "submit it", "start", or "continue." Only then may you submit.

This rule applies to:

- `loomloom template submit-file`
- `loomloom run submit`
- `loomloom template-spec submit-workbook`
- `loomloom template-spec create`
- `loomloom template-spec create-version`

Validation, downloads, schema inspection, model lookup, `doctor`, asset upload, artifact listing, and result backfill do not create a new paid hosted run and do not need the second confirmation.

The execution summary must include:

- template ID or template/version ID
- input file path or input source
- row count or task size
- action to be performed
- estimated cost, or a clear note if cost is only known after submission
- the exact confirmation phrase: `Reply "confirm submit" before I start.`

If the user says "do not run yet", "wait", or anything similar, stay in preparation mode.

## Error Handling

After a command fails, run:

```bash
loomloom doctor
```

Use the doctor output to distinguish local configuration issues, old CLI versions, service behavior, or model catalog issues. Do not guess that the template, model, or run is broken without evidence.

## Current MVP Capabilities

The public CLI supports environment checks, input asset upload, template discovery, official template workbook download/validation/submission, TemplateSpec JSON custom template creation, version creation for existing user templates, user-template workbook download/validation/submission, model discovery, JSONL run submission, run watch, result row/workbook download, artifact listing, and artifact download.

## Large File Handling

When the user wants to process local code files, large text, local images, or other large files in batches, avoid pasting full files into the agent context. Prefer:

1. `loomloom input-asset upload <file>`
2. Save the returned `input_asset_id`
3. Prepare structured JSONL or Excel input in later steps

## Default Behavior

Use official Excel template workflows by default unless the user explicitly asks for JSON/JSONL.

When the user asks to create or customize a workflow/template, prefer TemplateSpec JSON. TemplateSpec JSON is the source data; downloaded workbooks are generated artifacts. If a template version changes, do not promise that old workbooks remain compatible. Download a new workbook.

## Conversational Template Authoring

When the user describes a new template in natural language, do not immediately write TemplateSpec JSON. First use or reference:

```bash
loomloom template-spec docs conversation
```

Flow:

1. Ask business questions, not TemplateSpec technical-field questions.
2. Ask one question at a time when information is missing.
3. Avoid exposing technical terms such as `fieldBindings`, `upstreamBindings`, `fan-in`, `execution`, `outputSchema`, `provider`, or `mode` to the user.
4. Restate complex workflows in business language before asking follow-up questions.
5. If the workflow has multiple roles, steps, review viewpoints, or generation agents, ask how future users should provide per-role or per-step requirements before drafting the TemplatePlan.
6. Draft a TemplatePlan first.
7. Show the TemplatePlan and wait for confirmation.
8. Generate TemplateSpec JSON only after the user confirms the TemplatePlan.
9. Run `loomloom template-spec check <spec.json>` before creation.
10. Ask for a separate creation confirmation after the check passes.
11. Run `loomloom template-spec create <spec.json>` only after the user explicitly confirms template creation.

For existing user templates, do not promise in-place historical updates. Use `loomloom template-spec create-version <template-id> <spec.json>` and make later workbooks/runs use the new `version_id`.

### Template Usage Mode

If a template contains multiple roles, steps, review viewpoints, or generation agents, and each role or step may need its own requirements, ask before generating the TemplatePlan:

```text
When other people use this template later, how should the review or generation requirements be handled?

1. Preset in the template: users fill only the core material, and the system runs your predefined requirements.
2. Editable by users: users can fill or modify the requirements for each step or role.
3. Generate both versions: one simple version and one customizable version.
```

If the user is unsure, recommend option 1 as the simplest useful standard template.

Do not expose terms such as `prompt`, `binding`, `reference`, `field`, `hidden`, `paramBindings`, or `fieldBindings` in that question. Common trigger scenarios include multi-view PRD review, multi-role contract review, launch-event planning with multiple agents, multi-style article rewriting, resume review from multiple interviewer viewpoints, and multi-channel marketing generation.

Rules:

- Simple mode: users fill core material only; per-role or per-step requirements live in the template instructions.
- Customizable mode: users fill core material and may also edit requirements for each role or step.
- Both versions: show two TemplatePlans and get confirmation before generating two TemplateSpecs.

Acceptance criteria:

- Multi-role, multi-step, or multi-view templates ask the usage-mode question.
- Simple mode does not expose role/step requirement columns to template users.
- Customizable mode exposes editable input columns for those requirements.
- User-facing dialogue avoids technical concepts.
- The generated TemplateSpec expresses both simple and customizable modes correctly.

Creation gate:

- "Create a PRD review template" only starts the authoring flow; it is not remote creation confirmation.
- Environment variables, tokens, and server URLs are configuration, not creation confirmation.
- "Generate the spec" means generate and check locally; it is not permission to run `template-spec create`.
- Before `template-spec create`, show the template name, spec path, check result, exact creation command, and ask the user to reply `confirm create template`.

TemplatePlan should cover the template name and goal, what each workbook row means, user input fields, workflow steps, serial/parallel/summary relationships, template usage mode, visible outputs, failure policy, error columns, default model, and special requirements.

Default modeling rules:

- "Product, engineering, and risk review separately, then summarize" should be modeled as multiple parallel `text-generate` steps plus a downstream summary step using step-level `dependsOn` and `upstreamBindings`.
- Do not model multi-role review as one `expanded` step. Use `expanded` only for dynamic multi-value input.
- Add result and error columns for every user-visible step by default.
- If partial completion is allowed, the summary step should explain missing upstream results and expose error columns for failed steps.
- Keep `provider` and `mode` internal.

TemplateSpec constraints:

- Before writing a custom TemplateSpec, run `loomloom template-spec docs spec` and use the bundled docs as the current contract.
- Use `loomloom template-spec docs examples` for patterns.
- Use `loomloom template-spec docs conversation` for the natural-language creation flow.
- Use only `text-generate`, `image-generate`, and `video-generate` as default execution units.
- Use lowerCamel OpenAPI fields such as `meta.name`, `steps[].stepId`, and `defaultModelRef.modelKey`.
- Connect steps with step-level `dependsOn` and `upstreamBindings`; the usual source output port is `output`.
- Before choosing `defaultModelRef.modelKey`, run `loomloom template-spec models <execution-unit>` and use a returned `model_id`.
- Expose model columns only when the step has `allowModelOverride=true` and a field binding to `paramKey=model`.
- Do not bind `provider` or `mode`.

## Result Source Of Truth

The submitted workbook and the server-side run input snapshot are the source of truth. After a run completes, prefer `loomloom run result-workbook <run-id>` because the server aligns original input rows with generated artifacts. Use `template backfill-results` only when the user explicitly needs the older local Excel backfill flow.

## Console Access

Console links shown to the user must come from user-provided CogFoundry workspace information. Do not default to local or historical service URLs.

After a successful submit, watch, or backfill operation, you may add:

- If the user provided a CogFoundry console URL, they can check run status in that console.
- If a `run_id` exists, they can search for it or inspect the latest run record.
