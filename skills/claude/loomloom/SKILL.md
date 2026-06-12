---
name: loomloom
description: Use this skill when the user mentions LoomLoom, batchjob, batchflow, batch processing, template submission, Excel template execution, run submit, artifact downloads, or result backfill.
---

# loomloom

Use this skill when the user wants to run batch content generation, submit templates, execute Excel workflows, check run status, download result workbooks, or download generated artifacts through LoomLoom. Legacy names such as `batchjob`, `batchflow`, and "batch processing" refer to the same LoomLoom workflow.

## Command Flow

0. If the user asks for a beta or internal CLI, install an explicit pre-release channel instead of stable:
   `curl -fsSL https://raw.githubusercontent.com/Cogfoundry-ai/loomloom/main/install.sh | bash -s -- --channel beta --no-brew`
1. Check the environment:
   `loomloom doctor`
   Confirm that `LOOMLOOM_SERVER` is the target environment provided by the user. If it is missing, ask for the service URL.
2. For large inputs, upload assets instead of pasting full files into the agent context:
   `loomloom input-asset upload <file>`
3. Use official Excel templates by default:
   `loomloom template list`
   `loomloom template schema <template-id>`
   `loomloom template download <id> --output-file <xlsx>`
   `loomloom template validate-file <id> <xlsx>`
   `loomloom template submit-file <id> <xlsx>`
4. For custom templates, use TemplateSpec JSON:
   `loomloom template-spec docs spec`
   `loomloom template-spec docs examples`
   `loomloom template-spec docs conversation`
   `loomloom template-spec check <spec.json>`
   `loomloom template-spec create <spec.json>`
   `loomloom template-spec create-version <template-id> <spec.json>`
5. Watch and download results:
   `loomloom run watch <run-id>`
   `loomloom run result-workbook <run-id> --output-file <xlsx>`
   `loomloom artifact download <run-id> --output-dir <dir>`

## Confirmation Rule

Any command that submits real work to the hosted LoomLoom service requires a second explicit confirmation from the user in the current conversation. This applies to:

- `loomloom template submit-file`
- `loomloom run submit`
- `loomloom template-spec submit-workbook`
- `loomloom template-spec create`
- `loomloom template-spec create-version`

Prepare, download, validate, inspect schemas, upload assets, list artifacts, and backfill results without a second confirmation because those operations do not create a new hosted paid run.

Before submission or remote template creation, show a concise execution summary with the template, input source, task size, action, cost note, and exact command. Ask the user to reply `confirm submit` or `confirm create template`.

## Custom Template Authoring

When the user describes a template in natural language:

1. Ask business questions, not TemplateSpec technical-field questions.
2. Ask one missing-detail question at a time.
3. Restate complex workflows in business language.
4. Draft a TemplatePlan and wait for user confirmation.
5. Generate TemplateSpec JSON only after the TemplatePlan is confirmed.
6. Run `loomloom template-spec check <spec.json>`.
7. Ask for a separate creation confirmation before `template-spec create`.

For multi-role, multi-step, or multi-view templates, ask whether future users should use preset requirements, editable per-role requirements, or both versions. Do not expose technical terms such as `fieldBindings`, `upstreamBindings`, `fan-in`, `provider`, or `mode` in user-facing questions.

## Error Handling

After a command fails, run:

```bash
loomloom doctor
```

Use the doctor output to distinguish local configuration issues, old CLI versions, service behavior, or model catalog issues.

## Console Access

Console links shown to the user must come from user-provided CogFoundry workspace information. Do not default to local or historical service URLs.
