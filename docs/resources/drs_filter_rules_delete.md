---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_filter_rules_delete"
description: |-
  Use this resource to delete data filter rules of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_filter_rules_delete

Use this resource to delete data filter rules of a DRS job within HuaweiCloud.

-> This resource is a one-time action resource for deleting data filter rules of a DRS job. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_filter_rules_delete" "test" {
  job_id   = var.job_id
  rule_ids = [1, 2]
  type     = "common"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to delete the filter rules.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS job to delete filter rules from.

* `rule_ids` - (Required, List, NonUpdatable) Specifies the list of filter rule IDs to delete.

* `type` - (Required, String, NonUpdatable) Specifies the type of the data filter rules to delete.  
  The valid values are as follows:
  + **common**
  + **config**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the job ID.
