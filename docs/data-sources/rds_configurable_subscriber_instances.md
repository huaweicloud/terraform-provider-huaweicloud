---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_configurable_subscriber_instances"
description: |-
  Use this data source to get the list of instances that can be configured as the subscriber.
---

# huaweicloud_rds_configurable_subscriber_instances

Use this data source to get the list of instances that can be configured as the subscriber.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_configurable_subscriber_instances" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `subscriber_instance_id` - (Optional, String) Specifies the subscriber instance ID.

* `subscriber_instance_name` - (Optional, String) Specifies the subscriber instance name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of instances that can be configured as the subscriber.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `instance_id` - Indicates the ID of the instance.

* `instance_name` - Indicates the name of the instance.

* `data_vip` - Indicates the internal IP of the instance.
