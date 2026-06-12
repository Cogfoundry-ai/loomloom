# TemplateSpec Authoring Guide

This guide explains how agents and developers should create LoomLoom custom templates.

## Recommended Workflow

1. Clarify the business goal.
2. Identify what one workbook row represents.
3. List user input fields.
4. Break the workflow into steps.
5. Decide whether steps are sequential, parallel, or summarized through fan-in.
6. Query available models with `loomloom template-spec models <execution-unit>`.
7. Write TemplateSpec JSON.
8. Run `loomloom template-spec check <spec.json>`.
9. Create the template only after explicit user confirmation.

## Modeling Rules

- Use `text-generate`, `image-generate`, and `video-generate` as the default execution units.
- Prefer fixed `instruction` text for template-author requirements.
- Use workbook fields only for values future users should edit.
- Use `dependsOn` and `upstreamBindings` for step-to-step data flow.
- Keep `provider` and `mode` internal.
- Expose model selection only when a step sets `allowModelOverride=true`.

## User Inputs

Use stable lower_snake_case field keys and clear labels. Do not use reserved keys such as `model`, `provider`, or `mode`.

For large reference text or files, prefer `asset_ref` or `text_reference` fields and upload the file with `loomloom input-asset upload`.

## Serial Workflows

For a workflow such as "draft, then polish", model two steps:

```text
stp_draft.output -> stp_polish.prompt
```

The downstream step should receive upstream output through `upstreamBindings`, not through a duplicate field binding to the same target.

## Parallel Review Workflows

For "product, engineering, and risk review, then summarize", model three sibling review steps and one downstream summary step. The summary step declares all three review steps in `dependsOn` and receives their outputs through multiple `upstreamBindings`.

Use `triggerPolicy=allow_partial` only when the business process accepts partial upstream success.

## Creation Safety

`template-spec create` and `template-spec create-version` modify hosted state. Agents must show the template name, spec path, check result, and command, then wait for explicit creation confirmation.
