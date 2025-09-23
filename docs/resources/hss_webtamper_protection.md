---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_webtamper_protection"
description: |-
  Manages an HSS web tamper protection resource within HuaweiCloud.
---

# huaweicloud_hss_webtamper_protection

Manages an HSS web tamper protection resource within HuaweiCloud.

-> Currently, for HSS web tamper protection, after enabling protection, you need to manually set the protection
  directory in the console for it to take effect. Additionally, after enabling dynamic protection, you also need to
  restart Tomcat for it to take effect. You can log in to the HSS console, select **Prevention**,
  then choose **Web Tamper Protection**, click on **Configure Protection**, and set the protection directory.

## Example Usage

```hcl
variable "host_id" {}
variable "quota_id" {}

resource "huaweicloud_hss_webtamper_protection" "test" {
  host_id  = var.host_id
  quota_id = var.quota_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the HSS web tamper protection resource belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `host_id` - (Required, String, ForceNew) Specifies the host ID for the web tamper protection.
  Changing this parameter will create a new resource.

  -> Before using HSS web tamper protection, it is necessary to ensure that the agent status of the host is **online**.

* `quota_id` - (Optional, String, ForceNew) Specifies quota ID (yearly/monthly billing mode) for web tamper protection.
  If omitted, an existing quota will be randomly selected. Changing this parameter will create a new resource.

* `is_dynamics_protect` - (Optional, Bool) Specifies whether to enable dynamic web tamper protection.
  Setting to **true** means enabling dynamic protection, while leaving it blank or setting to **false** means disabling
  dynamic protection.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  web tamper protection belongs. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID same as `host_id`.

* `host_name` - The host name of the web tamper protection.

* `public_ip` - The elastic public IP address of the host for web tamper protection.

* `private_ip` - The private IP address of the host for web tamper protection.

* `group_name` - The name of the host group to which the host for web tamper protection belongs.

* `os_bit` - The operating system bits of the host for web tamper protection.

* `os_type` - The operating system type of the host for web tamper protection.

* `protect_status` - The protection status of the host for web tamper protection.
  The value can be **closed** or **opened**.

* `rasp_protect_status` - The dynamic protection status of the host for web tamper protection.
  The value can be **closed** or **opened**.

* `anti_tampering_times` - The number of defended tampering attacks.

* `detect_tampering_times` - The number of detected tampering attacks.

## Import

The web tamper protection can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_hss_webtamper_protection.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `quota_id`, `is_dynamics_protect`,
`enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_webtamper_protection" "test" {
  ...

  lifecycle {
    ignore_changes = [
      quota_id, is_dynamics_protect, enterprise_project_id,
    ]
  }
}
```
