---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_primary_dim_monitored_objects"
description: |-
  Use this data source to query the list of monitored objects for DCS within HuaweiCloud.
---

# huaweicloud_dcs_primary_dim_monitored_objects

Use this data source to query the list of monitored objects for DCS within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_dcs_primary_dim_monitored_objects" "test" {
  dim_name = "dcs_instance_id"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the monitored objects. If omitted, the
  provider-level region will be used.

* `dim_name` - (Required, String) Specifies the primary dimension ID.
  The valid values are **dcs_instance_id** and **dcs_memcached_instance_id**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `router` - The current query dimension route. If it is a primary dimension, the array contains its own ID.

* `children` - The list of sub-dimension objects for the current query dimension.
  The [children](#children_struct) structure is documented below.

* `instances` - The list of monitoring objects for the current query dimension.
  The [instances](#instances_struct) structure is documented below.

<a name="children_struct"></a>
The `children` block supports:

* `dim_name` - The dimension name.

* `dim_route` - The route of the dimension.

<a name="instances_struct"></a>
The `instances` block supports:

* `dcs_instance_id` - The measurement object ID, which is the instance ID.

* `name` - The measurement object name, which is the instance name.

* `status` - The measurement object status, which is the instance status.
