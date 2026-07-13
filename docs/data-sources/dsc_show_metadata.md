---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_show_metadata"
description: |-
  Use this data source to get the asset metadata information required by the situation dashboard within HuaweiCloud.
---

# huaweicloud_dsc_show_metadata

Use this data source to get the asset metadata information required by the situation dashboard within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_show_metadata" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `column_num` - The total number of data columns.

* `file_num` - The total number of files.

* `sensitive_column_num` - The number of sensitive columns.

* `sensitive_file_num` - The number of sensitive files.

* `table_num` - The total number of data tables.
