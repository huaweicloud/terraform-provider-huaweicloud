---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_parameter_template_reset"
description: |-
  Manages a GeminiDB parameter template reset resource within HuaweiCloud.
---

# huaweicloud_geminidb_parameter_template_reset

Manages a GeminiDB parameter template reset resource within HuaweiCloud.

-> This resource is used to reset a custom parameter template to its default values.
This is a one-time operation and cannot be undone.

## Example Usage

```hcl
variable "config_id" {}

resource "huaweicloud_geminidb_parameter_template_reset" "test" {
  config_id = var.config_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `config_id` - (Required, String, NoneUpdatable) Specifies the ID of the parameter template to reset.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the configuration ID.
