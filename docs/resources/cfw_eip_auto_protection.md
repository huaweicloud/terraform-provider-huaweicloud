---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_eip_auto_protection"
description: |-
  Manages the EIP auto protection resource within HuaweiCloud.
---

# huaweicloud_cfw_eip_auto_protection

Manages the EIP auto protection resource within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "object_id" {}

resource "huaweicloud_cfw_eip_auto_protection" "test" {
  fw_instance_id = var.fw_instance_id
  object_id      = var.object_id
  status         = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `object_id` - (Required, String, NonUpdatable) Specifies the protection object ID.

* `status` - (Required, Int, NonUpdatable) Specifies whether to enable the addition of EIP automatic protection.  
  The valid values are as follows:
  + **1**: Enable auto protection for new EIPs.
  + **0**: Disable auto protection for new EIPs.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `object_id`.

* `available_eip_count` - The number of EIPs that can be protected.

* `beyond_max_count` - Whether the EIP count limit is exceeded.

* `eip_protected_self` - The number of protected EIPs.

* `eip_total` - The total number of EIPs.

* `eip_un_protected` - The number of unprotected EIPs.

## Import

The EIP auto protection resource can be imported using `fw_instance_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_eip_auto_protection.test <fw_instance_id>/<id>
```
