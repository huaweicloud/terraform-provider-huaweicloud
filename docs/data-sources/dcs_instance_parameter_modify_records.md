---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_parameter_modify_records"
description: |-
  Use this data source to query the list of DCS instance configuration parameter modification records within HuaweiCloud.
---

# huaweicloud_dcs_instance_parameter_modify_records

Use this data source to query the list of DCS instance configuration parameter modification records within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_instance_parameter_modify_records" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DCS instance is located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histories` - The list of configuration parameter modification records.  
  The [histories](#dcs_config_histories_history) structure is documented below.

<a name="dcs_config_histories_history"></a>
The `histories` block supports:

* `history_id` - The modification record ID.

* `type` - The modification type.  
  Valid value is **config_param** (instance configuration parameter).

* `created_at` - The modification time.

* `status` - The modification status.  
  Valid values are:
  + **SUCCESS**: Instance parameter modification succeeded.
  + **FAILURE**: Instance parameter modification failed.
  + **EXECUTING**: Instance parameter modification is in progress.
