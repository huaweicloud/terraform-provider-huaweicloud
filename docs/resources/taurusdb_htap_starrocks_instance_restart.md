---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_instance_restart"
description: |-
  Manages a TaurusDB HTAP StarRocks instance restart resource within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_instance_restart

Manages a TaurusDB HTAP StarRocks instance restart resource within HuaweiCloud.

-> This resource is a one-time action resource to restart a StarRocks instance. Deleting this resource will not clear
the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "taurusdb_instance_id"
variable "starrocks_instance_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_instance_restart" "test" {
  taurusdb_instance_id  = var.taurusdb_instance_id
  starrocks_instance_id = var.starrocks_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `taurusdb_instance_id` - (Required, String, NoneUpdatable) Specifies the TaurusDB instance ID.

* `starrocks_instance_id` - (Required, String, NoneUpdatable) Specifies the StarRocks instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `starrocks_instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
