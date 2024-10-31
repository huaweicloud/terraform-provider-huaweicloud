---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_diagnosis_statistics"
description: |-
  Use this data source to get the abnormal instances by each metric.
---

# huaweicloud_gaussdb_mysql_diagnosis_statistics

Use this data source to get the abnormal instances by each metric.

## Example Usage

```hcl
data "huaweicloud_gaussdb_mysql_diagnosis_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `diagnosis_info` - Specifies the diagnosis information list.

  The [diagnosis_info](#diagnosis_info_struct) structure is documented below.

<a name="diagnosis_info_struct"></a>
The `diagnosis_info` block supports:

* `metric_name` - Specifies the metric name.

* `count` - Specifies the number of instances.
