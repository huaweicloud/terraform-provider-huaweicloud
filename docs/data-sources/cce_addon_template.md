---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_addon_template"
description: ""
---

# huaweicloud_cce_addon_template

Use this data source to get available HuaweiCloud CCE add-on template.

## Example Usage

```hcl
variable "cluster_id" {}

variable "addon_name" {}

variable "addon_version" {}

data "huaweicloud_cce_addon_template" "test" {
  cluster_id = var.cluster_id
  name       = var.addon_name
  version    = var.addon_version
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE add-ons. If omitted, the provider-level
  region will be used.

* `cluster_id` - (Required, String) Specifies the ID of container cluster.

* `name` - (Required, String) Specifies the add-on name.

* `version` - (Required, String) Specifies the add-on version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID of the add-on template.

* `description` - The description of the add-on.

* `spec` - The detail configuration of the add-on template.

* `stable` - Whether the add-on template is a stable version.

* `support_version` - The cluster information.
  + `virtual_machine` - The cluster (Virtual Machine) version that the add-on template supported.
  + `bare_metal` - The cluster (Bare Metal) version that the add-on template supported.
