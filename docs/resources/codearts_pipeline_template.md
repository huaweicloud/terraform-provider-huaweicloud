---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_template"
description: |-
  Manages a CodeArts pipeline template resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_template

Manages a CodeArts pipeline template resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_codearts_pipeline_template" "test" {
  name           = var.name
  language       = "none"
  is_show_source = true
  description    = "test description"
  is_favorite    = true
  definition     = jsonencode({
    "stages": [
      {
        "name": "Stage_1",
        "sequence": "0",
        "jobs": [
          {
            "stage_id": xxx,
            "identifier": "xxx",
            "name": "CodeCheck",
            "depends_on": [],
            "timeout": "",
            "timeout_unit": "",
            "steps": [
              {
                "name": "CodeCheck",
                "task": "official_devcloud_codeCheck_template",
                "sequence": 0,
                "inputs": [
                  {
                    "key": "language",
                    "value": "Java"
                  },
                  {
                    "key": "module_or_template_id",
                    "value": "xxx"
                  }
                ],
                "business_type": "Gate",
                "runtime_attribution": "agent",
                "identifier": "xxx",
                "multi_step_editable": 0,
                "official_task_version": "0.0.1"
              }
            ],
            "resource": "{\"type\":\"system\",\"arch\":\"x86\"}",
            "condition": "$${{ completed() }}",
            "exec_type": "OCTOPUS_JOB",
            "sequence": 0
          }
        ],
        "identifier": "xxx",
        "pre": [
          {
            "task": "official_devcloud_autoTrigger",
            "sequence": 0
          }
        ],
        "post": null,
        "depends_on": [],
        "run_always": false
      }
    ]
  })
  
  variables {
    sequence    = 1
    name        = "test_var"
    type        = "string"
    value       = "test_value"
    description = "test variable"
    is_runtime  = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the template name.

* `language` - (Required, String) Specifies the template language.

* `definition` - (Required, String) Specifies the template definition JSON.

* `is_show_source` - (Required, Bool) Specifies whether to display the pipeline source.

* `description` - (Optional, String) Specifies the template description.

* `variables` - (Optional, List) Specifies the custom variables.
  The [variables](#variables_struct) structure is documented below.

* `is_favorite` - (Optional, Bool) Specifies whether it is a favorite template. Defaults to **false**.

<a name="variables_struct"></a>
The `variables` block supports:

* `name` - (Optional, String) Specifies the custom variable name.

* `sequence` - (Optional, Int) Specifies the parameter sequence, starting from **1**.

* `type` - (Optional, String) Specifies the custom parameter type.
  Valid values are:
  + **autoIncrement**: Auto-increment parameter.
  + **enum**: Enumeration parameter.
  + **string**: String parameter.

* `value` - (Optional, String) Specifies the custom parameter default value.

* `is_secret` - (Optional, Bool) Specifies whether it is a private parameter. Defaults to `false`.

* `description` - (Optional, String) Specifies the parameter description.

* `is_runtime` - (Optional, Bool) Specifies whether the parameters can be set during runtime. Defaults to `false`.

* `is_reset` - (Optional, Bool) Specifies whether to reset.
  + **true**: Uses the edited parameter value.
  + **false**: Uses the auto-increment parameter.

  Defaults to `false`.

* `latest_value` - (Optional, String) Specifies the last parameter value.

* `runtime_value` - (Optional, String) Specifies the value passed in at runtime.

* `limits` - (Optional, List) Specifies the list of enumerated values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - Indicates the creation time.

* `update_time` - Indicates the last update time.

* `creator_id` - Indicates the creator.

* `updater_id` - Indicates the last updater.

* `icon` - Indicates the template icon.

* `manifest_version` - Indicates the manifest version.

* `is_system` - Indicates whether it is a system template.

## Import

The pipeline template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_template.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `is_favorite`.
It is generally recommended running `terraform plan` after importing the template.
You can then decide if changes should be applied to the template, or the resource definition should be updated to
align with the template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_pipeline_template" "test" {
  ...

  lifecycle {
    ignore_changes = [
      is_favorite,
    ]
  }
}
```
