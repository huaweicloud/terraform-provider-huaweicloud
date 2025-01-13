---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_parameter_template_reset"
description: |-
  Manages a GaussDB OpenGauss parameter template reset resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_parameter_template_reset

Manages a GaussDB OpenGauss parameter template reset resource within HuaweiCloud.

## Example Usage

```hcl
variable "config_id" {}

resource "huaweicloud_gaussdb_opengauss_parameter_template_reset" "test" {
  config_id = var.config_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `config_id` - (Required, String, ForceNew) Specifies the ID of the source parameter template to be reset. Changing
  this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals `config_id`.
