---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_disk_types"
description: |-
  Use this data source to query the disk types supported by each availability zone in CSS within HuaweiCloud.
---

# huaweicloud_css_disk_types

Use this data source to query the disk types supported by each availability zone in CSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_css_disk_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the disk types.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `disk_types` - The disk type list supported by each availability zone.
  The [disk_types](#disk_types_struct) structure is documented below.

<a name="disk_types_struct"></a>
The `disk_types` block supports:

* `availability_zone` - The availability zone name.

* `volume_names` - The storage types supported in the availability zone.
  The valid values are as follows:
  + **SATA**: Common I/O.
  + **SAS**: High I/O.
  + **SSD**: Ultra-high I/O.
  + **ESSD**: Extreme SSD.
