---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_device_monitor_infos"
description: |-
  Use this data source to get the device monitor information list within HuaweiCloud.
---

# huaweicloud_dsc_device_monitor_infos

Use this data source to get the device monitor information list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_device_monitor_infos" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `monitor_infos` - The device monitor information list.
  The [monitor_infos](#monitor_infos_struct) structure is documented below.

<a name="monitor_infos_struct"></a>
The `monitor_infos` block supports:

* `id` - The device ID.

* `name` - The device name.

* `ip` - The device IP address.

* `type` - The device type.

* `status` - The device status.

* `description` - The device description.

* `license_info` - The license information.
  The [license_info](#license_info_struct) structure is documented below.

* `os_info` - The VM resource usage information.
  The [os_info](#os_info_struct) structure is documented below.

<a name="license_info_struct"></a>
The `license_info` block supports:

* `license` - The license information.

* `license_start` - The license effective time.

* `license_end` - The license expiration time.

<a name="os_info_struct"></a>
The `os_info` block supports:

* `cpu` - The number of CPU cores.

* `disk` - The disk size.

* `mem` - The memory size.
