---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_subscription_regenerate"
description: |-
  Manages an RDS subscription regenerate resource within HuaweiCloud.
---

# huaweicloud_rds_subscription_regenerate

Manages an RDS subscription regenerate resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "subscription_id" {}

resource "huaweicloud_rds_subscription_regenerate" "test" {
  instance_id     = var.instance_id
  subscription_id = var.subscription_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

* `subscription_id` - (Required, String, NonUpdatable) Specifies the ID of the subscription.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
