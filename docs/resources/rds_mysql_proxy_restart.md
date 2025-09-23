---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_proxy_restart"
description: |-
  Manages an RDS MySQL proxy restart resource within HuaweiCloud.
---

# huaweicloud_rds_mysql_proxy_restart

Manages an RDS MySQL proxy restart resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "proxy_id" {}

resource "huaweicloud_rds_mysql_proxy_restart" "test" {
  instance_id = var.instance_id
  proxy_id    = var.proxy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS MySQL instance.

* `proxy_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS MySQL proxy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the `proxy_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
