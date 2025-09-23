---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_parameter_template_reset"
description: |-
  Manages a DDS parameter template reset resource within HuaweiCloud.
---

# huaweicloud_dds_parameter_template_reset

Manages a DDS parameter template reset resource within HuaweiCloud.

## Example Usage

```hcl
variable "configuration_id" {}

resource "huaweicloud_dds_parameter_template_reset" "test" {
  configuration_id = var.configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `configuration_id` - (Required, String, ForceNew) Specifies the parameter template ID.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
