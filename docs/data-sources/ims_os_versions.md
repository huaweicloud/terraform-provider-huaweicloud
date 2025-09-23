---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_os_versions"
description: |-
  Use this data source to get the list of OS versions supported by IMS image within HuaweiCloud.
---

# huaweicloud_ims_os_versions

Use this data source to get the list of OS versions supported by IMS image within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_ims_os_versions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `tag` - (Optional, String) Specifies the label value for the OS versions.
  Multiple tags separated by commas, e.g. **bms,uefi**.  
  The valid values are as follows:
  + **bms**: OS versions that supports BMS image type.
  + **uefi**: OS versions that supports UEFI boot mode.
  + **arm**: OS versions based on ARM architecture.
  + **x86**: OS versions based on x86 architecture.

  If omitted, query all supported OS versions in the current region.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `os_versions` - The OS version details.

  The [os_versions](#os_versions_struct) structure is documented below.

<a name="os_versions_struct"></a>
The `os_versions` block supports:

* `platform` - The operating system platform.

* `versions` - The operating system details.

  The [versions](#os_versions_versions_struct) structure is documented below.

<a name="os_versions_versions_struct"></a>
The `versions` block supports:

* `platform` - The operating system platform.

* `os_version_key` - The operating system key value.
  By default, the value of `os_version` is taken as the `os_version_key` value.

* `os_version` - The complete information of the operating system.

* `os_bit` - The number of bits for the operating system.

* `os_type` - The type of operating system.
