---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_data_map_security_level"
description: |-
  Use this data source to get the security level of the current region within HuaweiCloud.
---

# huaweicloud_dsc_data_map_security_level

Use this data source to get the security level of the current region within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_data_map_security_level" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `level` - The security level of the current region.
