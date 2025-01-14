---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_host_protection"
description: |-
  Manages an HSS host protection resource within HuaweiCloud.
---
# huaweicloud_hss_host_protection

Manages an HSS host protection resource within HuaweiCloud.

## Example Usage

```hcl
variable "host_id" {}
variable "version" {}
variable "charging_mode" {}

resource "huaweicloud_hss_host_protection" "test" {
  host_id       = var.host_id
  version       = var.version
  charging_mode = var.charging_mode
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the HSS host protection resource belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `host_id` - (Required, String, ForceNew) Specifies the host ID for the host protection.
  Changing this parameter will create a new resource.

  -> Before using host protection, it is necessary to ensure that the agent status of the host is **online**.

* `version` - (Required, String) Specifies the protection version enabled by the host.  
  The valid values are as follows:
  + **hss.version.basic**: Basic version.
  + **hss.version.advanced**: Professional version.
  + **hss.version.enterprise**: Enterprise version.
  + **hss.version.premium**: Ultimate version.

* `charging_mode` - (Required, String) Specifies the charging mode for host protection.  
  The valid values are as follows:
  + **prePaid**: The yearly/monthly billing mode.
  + **postPaid**: The pay-per-use billing mode.

* `quota_id` - (Optional, String) Specifies quota ID for host protection.
  If omitted, randomly select quota for the corresponding version.
  This field is valid only when `charging_mode` is set to **prePaid**.

* `is_wait_host_available` - (Optional, Bool) Specifies whether to wait for the host agent status to become **online**.
  The value can be **true** or **false**. Defaults to **false**.

  -> If this field is set to **true**, the program will wait for a maximum of `30` minutes until the host's agent status
  becomes **online**, and then enable host protection.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the host
  protection belongs. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID same as `host_id`.

* `host_name` - The host name.

* `host_status` - The host status. The value can be **ACTIVE**, **SHUTOFF**, **BUILDING**, or **ERROR**.

* `private_ip` - The private IP address of the host.

* `agent_id` - The agent ID installed on the host.

* `agent_status` - The agent status of the host. The value can be **installed**, **not_installed**, **online**,
  **offline**, **install_failed**, or **installing**.

* `os_type` - The operating system type of the host. The value can be **Linux** or **Windows**.

* `status` - The protection status of the host. The value can be **closed** or **opened**.

* `detect_result` - The security detection result of the host. The value can be **undetected**, **clean**, **risk**,
  or **scanning**.

* `asset_value` - The asset importance. The value can be **important**, **common**, or **test**.

* `open_time` - The time to enable host protection.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The host protection can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_hss_host_protection.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `quota_id`, `is_wait_host_available`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_host_protection" "test" { 
  ...
  
  lifecycle {
    ignore_changes = [
      quota_id, is_wait_host_available,
    ]
  }
}
```
