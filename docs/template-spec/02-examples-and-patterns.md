# Examples and Patterns

---

## Pattern 1：单步生成

**适用**：单步文案生成、改写、翻译、分类。

```json
{
  "meta": {
    "name": "Single Text Generate",
    "description": "Single-step text generation"
  },
  "steps": [
    {
      "stepId": "stp_text",
      "displayName": "Generate",
      "executionUnit": "text-generate",
      "defaultModelRef": { "modelKey": "google/gemini-2.5-flash" }
    }
  ],
  "inputSchema": {
    "fields": [
      {
        "key": "prompt",
        "label": "Prompt",
        "required": true,
        "valueType": "string",
        "order": 1
      }
    ]
  },
  "fieldBindings": [
    {
      "fieldKey": "prompt",
      "stepId": "stp_text",
      "paramKey": "prompt",
      "bindMode": "shared"
    }
  ]
}
```

---

## Pattern 2：线性链路

**适用**：两步串联，第二步以第一步的输出为输入。例如：起草 → 润色。

```json
{
  "meta": {
    "name": "Draft and Polish",
    "description": "Draft text, then polish the draft"
  },
  "steps": [
    {
      "stepId": "stp_draft",
      "displayName": "Draft",
      "executionUnit": "text-generate",
      "defaultModelRef": { "modelKey": "google/gemini-2.5-flash" },
      "instruction": "Write a draft based on the prompt."
    },
    {
      "stepId": "stp_polish",
      "displayName": "Polish",
      "executionUnit": "text-generate",
      "defaultModelRef": { "modelKey": "google/gemini-2.5-flash" },
      "dependsOn": ["stp_draft"],
      "upstreamBindings": [
        {
          "inputPort": "prompt",
          "sourceStepId": "stp_draft",
          "sourcePort": "output"
        }
      ],
      "instruction": "Polish and improve the draft above."
    }
  ],
  "inputSchema": {
    "fields": [
      {
        "key": "topic",
        "label": "Topic",
        "required": true,
        "valueType": "string",
        "order": 1
      }
    ]
  },
  "fieldBindings": [
    {
      "fieldKey": "topic",
      "stepId": "stp_draft",
      "paramKey": "prompt",
      "bindMode": "shared"
    }
  ]
}
```

**要点**：

- `stp_polish` 的 `prompt` 来自 UpstreamBinding，不再有 FieldBinding 写入 `prompt`，遵守 binding target 唯一性。
- `instruction` 是模板作者写定的系统级说明，用户不可见、不可修改。

---

## Pattern 3：Step-Level Fan-In Review Summary

**适用**：对同一份材料从多个视角并行评审，然后汇总结果。典型场景：PRD 多角色评审。

**结构**：

```text
多个并行 step（分别从产品、研发、风险视角评审 PRD）
  ↓
一个汇总 step（upstreamBindings 接收多个上游 outputs）
```

```json
{
  "meta": {
    "name": "PRD Review and Summary",
    "description": "Review one PRD from product, engineering, and risk perspectives, then summarize"
  },
  "steps": [
    {
      "stepId": "stp_prod01",
      "displayName": "Product Review",
      "executionUnit": "text-generate",
      "instruction": "Review the PRD from a product perspective. Focus on goal clarity, scope boundaries, requirement consistency, user value, and ambiguous behaviors.",
      "defaultModelRef": { "modelKey": "google/gemini-2.5-flash-lite" }
    },
    {
      "stepId": "stp_eng001",
      "displayName": "Engineering Review",
      "executionUnit": "text-generate",
      "instruction": "Review the PRD from an engineering perspective. Focus on implementation complexity, technical risks, edge cases, dependencies, and missing constraints.",
      "defaultModelRef": { "modelKey": "google/gemini-2.5-flash-lite" }
    },
    {
      "stepId": "stp_risk01",
      "displayName": "Risk Review",
      "executionUnit": "text-generate",
      "instruction": "Review the PRD from a risk and rollout perspective. Focus on launch risk, support burden, data or compliance concerns, and operational readiness.",
      "defaultModelRef": { "modelKey": "google/gemini-2.5-flash-lite" }
    },
    {
      "stepId": "stp_summary",
      "displayName": "Summary",
      "executionUnit": "text-generate",
      "defaultModelRef": { "modelKey": "google/gemini-2.5-flash-lite" },
      "dependsOn": ["stp_prod01", "stp_eng001", "stp_risk01"],
      "triggerPolicy": "require_all",
      "upstreamBindings": [
        {
          "inputPort": "prompt",
          "sourceStepId": "stp_prod01",
          "sourcePort": "output"
        },
        {
          "inputPort": "prompt",
          "sourceStepId": "stp_eng001",
          "sourcePort": "output"
        },
        {
          "inputPort": "prompt",
          "sourceStepId": "stp_risk01",
          "sourcePort": "output"
        }
      ],
      "instruction": "Summarize the review results above into one final PRD review. Consolidate overlapping findings, call out the most important risks, and end with a concise conclusion."
    }
  ],
  "inputSchema": {
    "fields": [
      {
        "key": "prd_content",
        "label": "PRD Content",
        "required": true,
        "valueType": "string",
        "order": 1
      }
    ],
    "instructions": [
      "Each row represents one PRD review task.",
      "Fill PRD Content with the PRD text. The same field is bound to each review step.",
      "The summary step depends on all review steps and receives all review outputs."
    ]
  },
  "fieldBindings": [
    {
      "fieldKey": "prd_content",
      "stepId": "stp_prod01",
      "paramKey": "prompt",
      "bindMode": "shared"
    },
    {
      "fieldKey": "prd_content",
      "stepId": "stp_eng001",
      "paramKey": "prompt",
      "bindMode": "shared"
    },
    {
      "fieldKey": "prd_content",
      "stepId": "stp_risk01",
      "paramKey": "prompt",
      "bindMode": "shared"
    }
  ]
}
```

**要点**：

- `prd_content` 是三个 review step 各自的输入绑定。三个 step 收到相同字段值，不代表新增了 input fan-out 语法。
- 如果希望三个 step 使用不同材料，可以改成 `prd_product`、`prd_engineering`、`prd_risk` 三个字段，分别绑定到三个 root step。
- `stp_summary` 是真正的 step-level fan-in：它依赖三个上游 step，并把三个 `output` 都绑定到自己的 `prompt`。
- `text-generate.prompt` 允许多来源 text 输入，所以这个 fan-in 合法。
- `triggerPolicy` 控制汇总节点的触发策略。这里显式写 `require_all`，表示三个视角全部成功后才汇总；不写时默认也是 `require_all`。

如果业务允许部分视角失败后继续汇总，可以把汇总 step 改为：

```json
{
  "stepId": "stp_summary",
  "displayName": "Summary",
  "executionUnit": "text-generate",
  "defaultModelRef": {
    "modelKey": "deepseek/deepseek-v4-flash"
  },
  "dependsOn": ["stp_prod01", "stp_eng001", "stp_risk01"],
  "triggerPolicy": "allow_partial",
  "upstreamBindings": [
    {
      "inputPort": "prompt",
      "sourceStepId": "stp_prod01",
      "sourcePort": "output"
    },
    {
      "inputPort": "prompt",
      "sourceStepId": "stp_eng001",
      "sourcePort": "output"
    },
    {
      "inputPort": "prompt",
      "sourceStepId": "stp_risk01",
      "sourcePort": "output"
    }
  ],
  "instruction": "Create a consolidated review from available upstream perspectives."
}
```

`allow_partial` 会等所有上游进入终态后再判断：至少一个上游成功才执行汇总；全部上游失败时汇总不执行。失败上游的错误原因不会作为业务正文传给汇总 step。
