---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_instance_upgrade"
description: |-
  Use this resource to upgrade the kernel version of a TaurusDB HTAP StarRocks instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_instance_upgrade

Use this resource to upgrade the kernel version of a TaurusDB HTAP StarRocks instance within HuaweiCloud.

-> This resource is a one-time action resource to upgrade the kernel version of a StarRocks instance. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "starrocks_instance_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_instance_upgrade" "test" {
  instance_id           = var.instance_id
  starrocks_instance_id = var.starrocks_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the TaurusDB instance ID to which the StarRocks
  instance belongs.

* `starrocks_instance_id` - (Required, String, NonUpdatable) Specifies the StarRocks instance ID.

* `delay` - (Optional, String, NonUpdatable) Specifies whether to delay the upgrade.
  The valid values are as follows:
  + **true**
  + **false**

  Defaults to **false**.

* `is_skip_validate` - (Optional, String, NonUpdatable) Specifies whether to skip upgrade verification.
  The valid values are as follows:
  + **true**
  + **false**

  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
