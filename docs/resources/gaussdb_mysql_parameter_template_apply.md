---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_parameter_template_apply"
description: |-
  Manages a GaussDB MySQL parameter template apply resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_parameter_template_apply

Manages a GaussDB MySQL parameter template apply resource within HuaweiCloud.

## Example Usage

```hcl
variable "configuration_id" {}
variable "instance_id" {}

resource "huaweicloud_gaussdb_mysql_parameter_template_apply" "test" {
  configuration_id = var.configuration_id
  instance_id      = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `configuration_id` - (Required, String, ForceNew) Specifies the parameter template ID.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<configuration_id>/<instance_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
