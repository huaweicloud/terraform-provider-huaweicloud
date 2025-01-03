---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_parameter_template_apply"
description: |-
  Manages a GaussDB OpenGauss parameter template apply resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_parameter_template_apply

Manages a GaussDB OpenGauss parameter template apply resource within HuaweiCloud.

## Example Usage

```hcl
variable "config_id" {}
variable "instance_id" {}

resource "huaweicloud_gaussdb_opengauss_parameter_template_apply" "test" {
  config_id   = var.config_id
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `config_id` - (Required, String, ForceNew) Specifies the parameter template ID.

* `instance_id` - (Required, String, ForceNew) Specifies the GaussDB OpenGauss instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<config_id>/<instance_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
