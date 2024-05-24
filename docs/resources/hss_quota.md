---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_quota"
description: |-
  Manages an HSS quota resource within HuaweiCloud.
---

# huaweicloud_hss_quota

Manages an HSS quota resource within HuaweiCloud.

## Example Usage

```hcl
variable "quota_version" {}

resource "huaweicloud_hss_quota" "test" {
  version     = var.quota_version
  period_unit = "month"
  period      = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the HSS quota resource belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies protection quota version. Changing this parameter will create a
  new resource.  
  The valid values are as follows:
  + **hss.version.basic**: Basic version.
  + **hss.version.advanced**: Professional version.
  + **hss.version.enterprise**: Enterprise version.
  + **hss.version.premium**: Ultimate version.
  + **hss.version.wtp**: Web page tamper prevention version.
  + **hss.version.container.enterprise**: Container version.

* `period_unit` - (Required, String, ForceNew) Specifies the charging period unit of the quota.
  Valid values are **month** and **year**. Changing this parameter will create a new resource.

* `period` - (Required, Int, ForceNew) Specifies the charging period of the quota. Changing this parameter will
  create a new resource.  
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.  
  If `period_unit` is set to **year**, the valid values are `1`, `2`, `3`, or `5`.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the HSS quota.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the HSS quota belongs.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of quota. The value can be **normal**, **expired**, or **freeze**.

* `used_status` - The usage status of quota. The value can be **idle** or **used**.

* `host_id` - The host ID for quota binding.

* `host_name` - The host name for quota binding.

* `charging_mode` - The charging mode of quota.  
  The valid values are as follows:
  + **prePaid**: The yearly/monthly billing mode.

* `expire_time` - The expiration time of quota, `-1` indicates no expiration date.

* `shared_quota` - Is it a shared quota. The value can be **shared** or **unshared**.

* `is_trial_quota` - Is it a trial quota. The value can be **true** or **false**.

* `enterprise_project_name` - The enterprise project name to which the quota belongs.

## Import

The quota can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_hss_quota.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `period_unit`, `period`, `auto_renew`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_quota" "test" { 
  ...

  lifecycle {
    ignore_changes = [
      period_unit, period, auto_renew,
    ]
  }
}
```
