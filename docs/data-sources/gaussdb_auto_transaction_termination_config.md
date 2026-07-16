---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_auto_transaction_termination_config"
description: |-
  Use this data source to query the automatic transaction termination configuration of a specified GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_auto_transaction_termination_config

Use this data source to query the automatic transaction termination configuration of a specified GaussDB instance within
HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_auto_transaction_termination_config" "test" {
  instance_id = var.instance_id
  type        = "exec_time"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the automatic transaction termination
  configuration. If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `type` - (Required, String) Specifies the configuration type. The valid values are as follows:
  + **exec_time**: Represents long transactions.
  + **xlog_quantity**: Represents large transactions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `usernames` - The list of filtered usernames.

* `threshold` - The threshold for terminating transactions.

* `auto_stop` - Whether the configuration is enabled.

* `database_names` - The list of all databases in the current instance.

* `database_names_select` - The list of filtered database names.
