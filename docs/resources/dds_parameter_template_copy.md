---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_parameter_template_copy"
description: |-
  Manages a DDS parameter template copy resource within HuaweiCloud.
---

# huaweicloud_dds_parameter_template_copy

Manages a DDS parameter template copy resource within HuaweiCloud.

## Example Usage

```hcl
variable "configuration_id" {}
variable "name" {}

resource "huaweicloud_dds_parameter_template_copy" "test" {
  configuration_id = var.configuration_id
  name             = var.name
  description      = "test copy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `configuration_id` - (Required, String, ForceNew) Specifies the parameter template ID.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of replicated parameter template.
  The parameter template name can contain **1** to **64** characters. It can contain only letters, digits, hyphens (-),
  underscores (_), and periods (.).
  Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of replicated parameter template.
  The value is left blank by default. The description must consist of a maximum of **256** characters and cannot contain
  the carriage return character or the following special characters: >!<"&'=
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
