---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_installation_scripts"
description: |-
  Use this data source to get the installation scripts of SecMaster nodes.
---

# huaweicloud_secmaster_installation_scripts

Use this data source to get the installation scripts of SecMaster nodes.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_installation_scripts" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of installation scripts.
  The [records](#secmaster_installation_scripts) structure is documented below.

<a name="secmaster_installation_scripts"></a>
The `records` block supports:

* `os_type` - The operating system type.

* `commands` - The installation commands for the specified operating system.
