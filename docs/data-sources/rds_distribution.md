---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_distribution"
description: |-
  Use this data source to get the distribution of the RDS instance.
---

# huaweicloud_rds_distribution

Use this data source to get the distribution of the RDS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_distribution" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - Indicates the status of the distributions.
  The value can be: **normal**, **abnormal**, **creating**, **createfail**.

* `distributor_instance_id` - Indicates the ID of the distribution instance.

* `distributor_instance_name` - Indicates the name of the distribution instance.
