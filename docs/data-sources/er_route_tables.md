---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_route_tables"
description: ""
---

# huaweicloud_er_route_tables

Use this data source to query the route tables under the ER instance within HuaweiCloud.

## Example Usage

### Querying specified route tables under ER instance using name

```hcl
variable "instance_id" {}
variable "route_table_name" {}

data "huaweicloud_er_route_tables" "test" {
  instance_id = var.instance_id
  name        = var.route_table_name
}
```

### Querying specified route tables under ER instance using tags

```hcl
variable "instance_id" {}

data "huaweicloud_er_route_tables" "test" {
  instance_id = var.instance_id

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the ER instance and route table are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the ER instance to which the route tables belongs.

* `route_table_id` - (Optional, String) Specifies the route table ID used to query specified route table.

* `name` - (Optional, String) Specifies the name used to filter the route tables.  
  The name can contain `1` to `64` characters, only English letters, Chinese characters, digits, underscore (_),
  hyphens (-) and dots (.) allowed.

* `tags` - (Optional, Map) Specifies the key/value pairs used to filter the route tables.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `route_tables` - All route tables that match the filter parameters.  
  The [object](#route_tables) structure is documented below.

<a name="route_tables"></a>
The `route_tables` block supports:

* `id` - The route table ID.

* `name` - The name of the route table.

* `description` - The description of the route table.

* `associations` - The association configurations of the route table.  
  The [object](#route_table_relationship) structure is documented below.

* `propagations` - The propagation configurations of the route table.  
  The [object](#route_table_relationship) structure is documented below.

* `routes` - The route details of the route table.  
  The [object](#route_table_routes) structure is documented below.

* `is_default_association` - Whether this route table is the default association route table.

* `is_default_propagation` - Whether this route table is the default propagation route table.

* `status` - The current status of the route table.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

<a name="route_table_relationship"></a>
The `associations` or `propagations` block supports:

* `id` - The ID of the association/propagation.

* `attachment_id` - The attachment ID corresponding to the routing association/propagation.

* `attachment_type` - The attachment type corresponding to the routing association/propagation.

<a name="route_table_routes"></a>
The `routes` block supports:

* `id` - The route ID.

* `destination` - The destination address (CIDR) of the route.

* `is_blackhole` - Whether route is the black hole route.

* `attachments` - The details of the attachment corresponding to the route.  
  The [object](#route_table_route_attachments) structure is documented below.

* `status` - The current status of the route.

<a name="route_table_route_attachments"></a>
The `attachments` block supports:

* `attachment_id` - The ID of the nexthop attachment.

* `attachment_type` - The type of the nexthop attachment.

* `resource_id` - The ID of the resource associated with the attachment.
