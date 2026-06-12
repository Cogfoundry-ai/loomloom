# Conversational Template Authoring

This document defines how an agent should help a user create a LoomLoom custom template through conversation. The user should describe business intent; the agent turns that intent into a TemplatePlan and then a valid TemplateSpec.

## Core Principle

Conversational authoring is not about teaching the user TemplateSpec. The user only needs to answer business questions:

- What template they want to create
- What each workbook row represents
- What users will provide each time
- How the AI should process the inputs
- Which steps can run in parallel
- Whether results should be summarized
- Which intermediate results should be visible
- Whether the workflow should continue after partial failure

## Required Flow

```text
natural-language request
  -> agent asks a few business questions
  -> TemplatePlan draft
  -> user confirms
  -> TemplateSpec JSON
  -> template-spec check
  -> creation confirmation
  -> template-spec create
```

Do not create the hosted template before the user confirms the design and then confirms creation.

## Three Confirmation Gates

| Gate | Allowed | Not allowed |
| --- | --- | --- |
| Requirements | Ask questions, summarize workflow, draft TemplatePlan | Generate final TemplateSpec or call `template-spec create` |
| Plan confirmation | Generate TemplateSpec and run `template-spec check` | Create the hosted template |
| Creation confirmation | Run `template-spec create` after explicit confirmation | Submit a batch run |

These are not creation confirmations:

- "Create a template"
- "Use the defaults"
- "Continue"
- Providing only `LOOMLOOM_SERVER` or `LOOMLOOM_TOKEN`
- "Generate the spec"

Before `template-spec create`, show the template name, TemplateSpec path, check result, exact command, and ask the user to reply `confirm create template`.

## Dialogue Rules

Ask like a workflow consultant, not a configuration form.

Rules:

1. Ask one question at a time.
2. Avoid technical terms.
3. Offer defaults when the user is unsure.
4. Restate complex workflows before asking follow-up questions.
5. Do not ask the user to write prompts; ask them to describe each step's goal.
6. Show a TemplatePlan before generating TemplateSpec JSON.
7. Generate TemplateSpec only after the user confirms the TemplatePlan.
8. Ask for separate creation confirmation after `template-spec check` passes.

## Recommended Question Order

Ask only what is missing:

1. What kind of template do you want to create?
2. What information will users provide for each row?
3. How should the AI process that information?
4. Are there multiple steps, roles, or viewpoints?
5. Should intermediate results be visible?
6. What final outputs should users get?
7. What should happen if a step fails?
8. Any special requirements?

## TemplatePlan Draft

TemplatePlan is an intermediate design for user confirmation. It is not backend source data.

It should include:

- template name and goal
- one-row meaning
- user input fields
- workflow steps
- serial, parallel, and summary relationships
- template usage mode
- visible outputs
- failure policy
- error columns
- default model choices
- special requirements

## Defaults

When the user does not specify details, use these defaults:

| Item | Default |
| --- | --- |
| Row meaning | One row is one independent batch task. |
| Field keys | Use stable snake_case keys based on business names. |
| Step prompt | Generate fixed `instruction` text from the step goal. |
| Output columns | Add a result column for each user-visible step. |
| Error columns | Add an error-reason column for each user-visible step. |
| Step relationships | Infer serial, parallel, and summary relationships from the business flow. |
| Partial failures | Use `require_all` by default; use `allow_partial` only when the user accepts partial results. |
| Model | Run `loomloom template-spec models <execution-unit>` and choose an available default. |
| Routing parameters | Do not expose `provider` or `mode`. |

## Business Pattern Mapping

### Single Step

User says: "Generate marketing copy from a product description."

Map to one `text-generate` step with the user field bound to `prompt`.

### Serial Steps

User says: "First organize the requirements, then generate a final document."

Map to two steps and pass the first output into the second step with `upstreamBindings`.

### Parallel Reviews With Summary

User says: "Review from product, engineering, and risk viewpoints, then summarize."

Map to sibling review steps plus one summary step. The summary step lists all review steps in `dependsOn` and uses multiple `upstreamBindings`.

### Multiple Versions

If the user wants fixed versions, use sibling steps. If the user wants a dynamic list of values, use `multiValue=true` with `bindMode=expanded`.

## Pre-Generation Checklist

Before generating TemplateSpec JSON, confirm:

1. The TemplatePlan was shown to the user.
2. The user confirmed the design.
3. Each step has a clear goal.
4. Each root step has user input or static instruction.
5. Each downstream step has a clear input source.
6. Multi-upstream summaries use step-level `dependsOn` and `upstreamBindings`.
7. `provider` and `mode` are not exposed.
8. Model choices were checked with `loomloom template-spec models`.

## CLI Validation Loop

After generating the spec, run:

```bash
loomloom template-spec check <spec.json>
```

If validation fails, fix the spec and check again. Do not create the template until validation passes and the user confirms creation.

## Unsupported In v1

Do not promise:

- loops
- conditional branches
- dynamic routing
- long-term memory
- human intervention during a run
- arbitrary agent tool calls
- arbitrary custom execution units

If the user asks for these, explain the current limitation and offer the closest supported workflow.
