---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_baseline_check_results"
description: |-
  Use this data source to get the list of SecMaster baseline check results.
---

# huaweicloud_secmaster_baseline_check_results

Use this data source to get the list of SecMaster baseline check results.

## Example Usage

```hcl
variable "workspace_id" {}
variable "from_date" {}
variable "to_date" {}

data "huaweicloud_secmaster_baseline_check_results" "test" {
  workspace_id = var.workspace_id
  from_date    = var.from_date
  to_date      = var.to_date
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `from_date` - (Optional, String) Specifies the start time of the baseline check.
  For example: **2024-08-12T14:00:00.000+08:00**.

* `to_date` - (Optional, String) Specifies the end time of the baseline check.
  For example: **2024-08-12T14:00:00.000+08:00**.

* `condition` - (Optional, Map) Specifies the condition expression.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `baseline_check_results` - The list of baseline check result.
