---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_networks_by_tags"
description: |-
  Use this data source to get the list of central networks by tag.
---

# huaweicloud_cc_central_networks_by_tags

Use this data source to get the list of central networks by tag.

## Example Usage

```hcl
data "huaweicloud_cc_central_networks_by_tags" "test"{
  tags {
    key    = "key"
    values = ["value"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `tags` - (Required, List) Specifies the included tags.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key. The key can contain a maximum of **128** Unicode characters, including
  letters, digits, hyphens (-), and underscores (_).

* `values` - (Required, List) Specifies the list of values with the same key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `central_networks` - Indicates the list of central network.
  The [central_networks](#central_networks_struct) structure is documented below.

<a name="central_networks_struct"></a>
The `central_networks` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the instance name.

* `description` - Indicates the resource description.

* `state` - Indicates the central network status. The value can be:
  + **AVAILABLE**: The central network is available.
  + **UPDATING**: The central network is being updated.
  + **FAILED**: The operation on the central network failed.
  + **CREATING**: The central network is being created.
  + **DELETING**: The central network is being deleted.
  + **DELETED**: The central network is deleted.
  + **RESTORING**: The central network is being restored.

* `enterprise_project_id` - Indicates the ID of the enterprise project that the resource belongs to.

* `tags` - Indicates the resource tags.
  The [tags](#central_networks_tags_struct) structure is documented below.

* `default_plane_id` - Indicates the ID of the default central network plane.

* `auto_associate_route_enabled` - Indicates whether the auto associate route is enabled.

* `auto_propagate_route_enabled` - Indicates whether the auto propagate route is enabled.

* `planes` - Indicates the list of central network planes.
  The [planes](#planes_struct) structure is documented below.

* `er_instances` - Indicates the list of enterprise routers on a central network.
  The [er_instances](#er_instances_struct) structure is documented below.

* `created_at` - Indicates the time when the resource was created. The UTC time is in the **yyyy-MM-ddTHH:mm:ss** format.

* `updated_at` - Indicates the time when the resource was updated. The UTC time is in the **yyyy-MM-ddTHH:mm:ss** format.

<a name="central_networks_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `value` - Indicates the tag value.

<a name="planes_struct"></a>
The `planes` block supports:

* `id` - Indicates the instance ID.

* `name` - Indicates the instance name.

* `associate_er_tables` - Indicates the List of the enterprise routers on a central network.
  The [associate_er_tables](#associate_er_tables_struct) structure is documented below.

* `exclude_er_connections` - Indicates whether to exclude the connections to enterprise routers on the central network.
  The [exclude_er_connections](#exclude_er_connections_struct) structure is documented below.

* `is_full_mesh` - Indicates whether is full mesh.

<a name="associate_er_tables_struct"></a>
The `associate_er_tables` block supports:

* `project_id` - Indicates the project ID.

* `region_id` - Indicates the region ID.

* `enterprise_router_id` - Indicates the enterprise router ID.

* `enterprise_router_table_id` - Indicates the ID of the enterprise router route table.

<a name="exclude_er_connections_struct"></a>
The `exclude_er_connections` block supports:

* `exclude_er_instances` - Indicates the connections between enterprise routers managed by the central network plane.
  The [exclude_er_instances](#exclude_er_instances_struct) structure is documented below.

<a name="exclude_er_instances_struct"></a>
The `exclude_er_instances` block supports:

* `project_id` - Indicates the project ID.

* `region_id` - Indicates the region ID.

* `enterprise_router_id` - Indicates the enterprise router ID.

<a name="er_instances_struct"></a>
The `er_instances` block supports:

* `enterprise_router_id` - Indicates the enterprise router ID.

* `project_id` - Indicates the project ID.

* `region_id` - Indicates the region ID.

* `asn` - Indicates the ASN of the network instance when BGP is used for routing.

* `site_code` - Indicates the site code.
