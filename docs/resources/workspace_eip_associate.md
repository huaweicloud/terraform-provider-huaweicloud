---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_eip_associate"
description: ""
---

# huaweicloud_workspace_eip_associate

Using this resource to associate an EIP to a specified desktop within HuaweiCloud.

## Example Usage

### Associate with an EIP ID

```hcl
variable "desktop_id" {}
variable "eip_id" {}

resource "huaweicloud_workspace_eip_associate" "test" {
  desktop_id = var.desktop_id
  eip_id     = var.eip_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to associate the EIP. If omitted, the provider-level
  region will be used. Changing this will create a new resource.

* `desktop_id` - (Required, String, ForceNew) Specifies the desktop ID. Changing this will create a new resource.

* `eip_id` - (Required, String, ForceNew) Specifies the EIP ID to associate. Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is `<desktop_id>/<eip_id>`.

* `enterprise_project_id` - The enterprise project ID to which the EIP associated.

* `public_ip` - The IP address of the EIP.

## Import

EIP association can be imported using the `desktop_id` and associated `eip_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_workspace_eip_associate.test <desktop_id>/<eip_id>
```
