---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_filter_rules_collect"
description: |-
  Use this resource to collect the data filter rules of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_filter_rules_collect

Use this resource to collect the data filter rules of a DRS job within HuaweiCloud.

-> This resource is a one-time action resource for collecting the data filter rules of a DRS job. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_filter_rules_collect" "test" {
  job_id       = var.job_id
  filter_type  = "common"
  is_all_rules = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to collect the filter rules.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS job to collect filter rules from.

* `filter_type` - (Required, String, NonUpdatable) Specifies the type of the data filter rules to collect.  
  The valid values are as follows:
  + **common**
  + **config**

* `is_all_rules` - (Optional, Bool, NonUpdatable) Specifies whether to query all filter rules (including synchronized
  and unsent ones).

* `limit` - (Optional, Int, NonUpdatable) Specifies the number of records to return.

* `offset` - (Optional, Int, NonUpdatable) Specifies the offset of the records to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the async query task ID returned by the API.
