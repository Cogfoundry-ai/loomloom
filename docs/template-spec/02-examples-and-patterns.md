# TemplateSpec Examples And Patterns

## Pattern 1: Single-Step Generation

Use this for one-step copywriting, rewriting, translation, classification, or code review.

```json
{
  "meta": {
    "name": "Simple Copy Generator",
    "description": "Generate copy from one prompt.",
    "scenario": "text",
    "displayOutputType": "Text",
    "primaryOutputType": "text"
  },
  "steps": [
    {
      "stepId": "stp_text01",
      "displayName": "Generate text",
      "executionUnit": "text-generate",
      "instruction": "Generate clear, useful text for the user's request.",
      "defaultModelRef": { "modelKey": "google/gemini-2.5-flash" },
      "allowModelOverride": true
    }
  ],
  "inputSchema": {
    "fields": [
      { "key": "prompt", "label": "Prompt", "required": true, "valueType": "string", "order": 1 }
    ]
  },
  "fieldBindings": [
    { "fieldKey": "prompt", "stepId": "stp_text01", "paramKey": "prompt", "bindMode": "shared" }
  ]
}
```

## Pattern 2: Linear Chain

Use this when a later step consumes an earlier step's output.

```text
stp_draft.output -> stp_polish.prompt
```

Key points:

- The downstream `prompt` comes from `upstreamBindings`.
- Do not also bind a user field to the same downstream `prompt`.
- Use `instruction` for fixed author guidance.

## Pattern 3: Step-Level Fan-In Review Summary

Use this when multiple steps process the same material from different viewpoints and a final step summarizes the results.

```text
stp_prod01.output  \
stp_eng01.output    -> stp_summary.prompt
stp_risk01.output /
```

Example step IDs:

- `stp_prod01`
- `stp_eng01`
- `stp_risk01`
- `stp_summary`

Implementation notes:

- The review steps are sibling steps.
- Each review step may bind the same source field or its own source field.
- The summary step lists all review steps in `dependsOn`.
- The summary step receives all review outputs through multiple `upstreamBindings`.
- This is true Step-Level Fan-In, not one `expanded` step.
- Use `triggerPolicy=require_all` for strict summaries.
- Use `triggerPolicy=allow_partial` only when partial summaries are acceptable.

## Pattern 4: Multiple Versions

If the user wants three fixed styles, prefer three sibling steps with different instructions. If the number of styles is dynamic and user-provided, use a multi-value field with `bindMode=expanded`.
