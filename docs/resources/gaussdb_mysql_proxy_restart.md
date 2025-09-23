---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_proxy_restart"
description: |-
  Manages a GaussDB MySQL proxy restart resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_proxy_restart

Manages a GaussDB MySQL proxy restart resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "proxy_id" {}

resource "huaweicloud_gaussdb_mysql_proxy_restart" "test" {
  instance_id = var.instance_id
  proxy_id    = var.proxy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance. Changing this parameter
  will create a new resource.

* `proxy_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL proxy. Changing this parameter will
  create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the `proxy_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
