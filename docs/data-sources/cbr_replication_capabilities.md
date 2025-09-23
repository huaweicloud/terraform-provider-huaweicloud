---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_replication_capabilities"
description: |-
  Use this data source to query the CBR replication capabilities within HuaweiCloud.
---

# huaweicloud_cbr_replication_capabilities

Use this data source to query the CBR replication capabilities within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbr_replication_capabilities" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `regions` - List of regions that support replication

  The [regions](#regions_struct) structure is documented below.

<a name="regions_struct"></a>
The `regions` block supports:

* `name` - Region where the cloud service resides

* `replication_destinations` - List of supported destination regions
