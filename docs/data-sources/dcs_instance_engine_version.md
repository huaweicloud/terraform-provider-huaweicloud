---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_engine_version"
description: |-
  Use this data source to query the engine version information of a DCS instance within HuaweiCloud.
---

# huaweicloud_dcs_instance_engine_version

Use this data source to query the engine version information of a DCS instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_instance_engine_version" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the engine version. If omitted, the
  provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `engine_minor_version` - The current engine minor version of the DCS instance.

* `latest_engine_minor_version` - The latest engine minor version available for the DCS instance.

* `proxy_minor_version` - The current proxy minor version of the DCS instance.

* `latest_proxy_minor_version` - The latest proxy minor version available for the DCS instance.

* `engine_minor_version_upgradable` - Whether the engine minor version of the DCS instance can be upgraded.

* `proxy_minor_version_upgradable` - Whether the proxy minor version of the DCS instance can be upgraded.
