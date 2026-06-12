# TemplateSpec Specification

TemplateSpec describes a reusable LoomLoom workflow. It defines template metadata, workflow steps, input fields, field bindings, step dependencies, and output behavior. Agents should use this document as the source of truth when generating custom templates.

## 1. Top-Level Shape

```json
{
  "meta": {},
  "steps": [],
  "inputSchema": {},
  "fieldBindings": [],
  "paramBindings": []
}
```

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `meta` | object | Yes | Template metadata. |
| `steps` | array | Yes | Ordered workflow steps. |
| `inputSchema` | object | Yes | Workbook input fields. |
| `fieldBindings` | array | Optional | Direct field-to-step parameter bindings. |
| `paramBindings` | array | Optional | Composite parameter bindings from multiple sources. |

## 2. Metadata

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | string | Yes | Human-readable template name. |
| `description` | string | Optional | Template summary. |
| `scenario` | string | Optional | Product scenario key. |
| `inputSummary` | string | Optional | Short input summary. |
| `displayOutputType` | string | Optional | Human-readable output type. |
| `primaryOutputType` | string | Optional | Must match the inferred terminal output type if present. |
| `tags` | array<string> | Optional | Search and grouping tags. |

## 3. Steps

Each item in `steps` is a StepSpec.

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `stepId` | string | Yes | Stable ID, formatted as `stp_` plus 6-10 letters or digits. Must be unique. |
| `displayName` | string | Yes | User-facing step name. |
| `executionUnit` | string | Yes | Capability type registered in the capability registry. |
| `defaultModelRef` | object | Conditional | Required when the execution unit needs a model. |
| `instruction` | string | Optional | Template-author instruction that users cannot override. |
| `dependsOn` | array<string> | Optional | Upstream step IDs for step-level dependencies. |
| `upstreamBindings` | array | Optional | Connections from upstream outputs into this step. |
| `triggerPolicy` | string | Optional | Fan-in trigger policy. Defaults to `require_all`. |
| `allowModelOverride` | boolean | Optional | Whether workbook users may override this step's model. |
| `staticParams` | object | Optional | Fixed runtime parameters allowed by the execution unit. |

## 4. Upstream Bindings

`upstreamBindings` connect upstream step outputs or initial inputs into a target step input port.

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `inputPort` | string | Yes | Target input port on the current step. |
| `sourceType` | string | Optional | `step_output` by default, or `initial_input`. |
| `sourceStepId` | string | Conditional | Required for `step_output`; must be listed in `dependsOn`. |
| `sourcePort` | string | Conditional | Required for `step_output`; usually `output`. |
| `sourceInputKey` | string | Conditional | Required for `initial_input`. |

## 5. Trigger Policies

| Value | Meaning | Use case |
| --- | --- | --- |
| `require_all` | Wait for all upstream steps to succeed. If any upstream fails, the step does not run. | Strict workflows. |
| `allow_partial` | Wait for all upstream steps to finish, then run if at least one upstream succeeded. | Summaries that can tolerate missing inputs. |
| `fail_fast` | Stop the downstream path as soon as a required upstream fails. | Workflows where any critical failure makes later work useless. |

Constraints:

- Missing `triggerPolicy` means `require_all`.
- `allow_partial` cannot be used on a root step.
- Multi-upstream fan-in with `allow_partial` must declare explicit `upstreamBindings`.
- Failed upstream system errors are not passed as business input text.

## 6. Input Schema

`inputSchema.fields` defines workbook columns.

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `key` | string | Yes | Stable field key. Reserved keys such as `model`, `provider`, and `mode` cannot be used. |
| `label` | string | Yes | User-facing label. |
| `valueType` | string | Yes | Field value type. |
| `required` | boolean | Optional | Defaults to `false`. |
| `defaultValue` | string | Conditional | Required when `sourceKind` is not `user_input`. |
| `sourceKind` | string | Optional | `user_input`, `default_value`, or `hidden`. |
| `multiValue` | boolean | Optional | Whether one cell may contain multiple values. |
| `maxValues` | integer | Conditional | Required and positive when `multiValue=true`. |
| `enumValues` | array<string> | Conditional | Required for `valueType=enum`. |
| `acceptedMimeTypes` | array<string> | Conditional | Required for `asset_ref` or `text_reference`. |
| `order` | integer | Optional | Workbook display order. |
| `description` | string | Optional | Field help text. |

Supported `valueType` values:

| Value | Meaning |
| --- | --- |
| `string` | Plain text. |
| `enum` | Value from `enumValues`. |
| `image_url` | HTTP/HTTPS image URL. |
| `asset_ref` | Uploaded asset ID such as `ia_...`. |
| `text_reference` | Inline text or uploaded text asset reference. |

## 7. Field Bindings

`fieldBindings` map one input field directly to one step parameter.

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `fieldKey` | string | Yes | Field key from `inputSchema.fields`. |
| `stepId` | string | Yes | Target step. |
| `paramKey` | string | Yes | Target parameter. |
| `bindMode` | string | Yes | `shared` or `expanded`. |

`multiValue=false` fields must use `shared`. `multiValue=true` fields must use `expanded`.

Routing parameters have special rules:

- `paramKey=model` requires `allowModelOverride=true` on the target step.
- `provider` and `mode` must not be exposed through bindings.

## 8. Param Bindings

`paramBindings` combine multiple sources into one target parameter.

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `stepId` | string | Yes | Target step. |
| `paramKey` | string | Yes | Target parameter; cannot be `model`, `provider`, or `mode`. |
| `bindMode` | string | Yes | `shared` or `expanded`. |
| `separator` | string | Optional | Separator used when joining sources. |
| `sources` | array | Yes | Non-empty list of `field_ref` or `literal` sources. |

A ParamBinding may contain at most three regular visible field sources. At most one source may refer to a `multiValue=true` field; when present, `bindMode` must be `expanded`.

## 9. Binding Target Rules

Within the same step, a target parameter or input port must have a single source. Do not bind both a user field and an upstream output to the same target.

Invalid example:

```text
FieldBinding -> stp_summary.prompt
UpstreamBinding -> stp_summary.prompt
```

## 10. Fan-Out And Fan-In

`bindMode=expanded` fans out a multi-value field into multiple executions of the same step.

Step-Level Fan-In connects multiple upstream outputs into one downstream input port. It is valid only when the target port allows multiple inputs, source and target types match, and the target port defines merge behavior.

## 11. Built-In Execution Units

| Unit | Input ports | Output port | Model override |
| --- | --- | --- | --- |
| `text-generate` | `prompt` text, optional `reference` text | `output` text/plain | Supported |
| `image-generate` | `prompt` text | `output` image/png | Supported |
| `video-generate` | `prompt` text, optional `image` | `output` video/mp4 | Not supported |

## 12. Version Boundary

This document describes TemplateSpec v1 (`template-spec/v1`).

v1 supports step-level fan-in/fan-out, fan-in trigger policies, single-target binding uniqueness, explicit multi-source upstream binding rules, and merge behavior defined by the target input port.
