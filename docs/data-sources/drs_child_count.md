---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_child_count"
description: |-
  Use this data source to query the number of child tasks of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_child_count

Use this data source to query the number of child tasks of a DRS job within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_drs_child_count" "test" {
  instance_id = var.instance_id
  db_type     = "gaussdbv5"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the child task count.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID of the DDM or GaussDBv5 database.

* `db_type` - (Required, String) Specifies the type of the database instance.  
  The valid values are as follows:
  + **gaussdbv5**
  + **ddm**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `child_count` - The number of child tasks.
