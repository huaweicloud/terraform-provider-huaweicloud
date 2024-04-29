---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_associations"
description: ""
---

# huaweicloud_er_associations

Use this data source to get the list of associations.

## Example Usage

```hcl
variable instance_id {}
variable route_table_id {}

data "huaweicloud_er_associations" "test" {
  instance_id     = var.instance_id
  route_table_id  = var.route_table_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ER instance ID to which the association belongs.

* `route_table_id` - (Required, String) Specifies the route table ID to which the association belongs.

* `attachment_id` - (Optional, String) Specifies the attachment ID corresponding to the association.  

* `attachment_type` - (Optional, String) Specifies the attachment type corresponding to the association.  

* `status` - (Optional, String) Specifies the status of the association. Default value is `available`.
  The valid values are as follows:
  + **available**
  + **failed**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `associations` - All associations that match the filter parameters.
  The [associations](#route_associations) structure is documented below.

<a name="route_associations"></a>
The `associations` block supports:

* `id` - The association ID.

* `route_table_id` -The route table ID corresponding to the association.

* `attachment_id` -The attachment ID corresponding to the association.

* `attachment_type` -The type of the attachment corresponding to the association.

* `resource_id` - The resource ID of the attachment corresponding to the association.

* `route_policy_id` - The route policy ID of the egress IPv4 protocol.

* `status` - The current status of the association.

* `created_at` - The creation time.

* `updated_at` - The latest update time.
