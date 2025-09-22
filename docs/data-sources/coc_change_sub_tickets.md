---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_change_sub_tickets"
description: |-
  Use this data source to get the list of COC change work order sub-orders.
---

# huaweicloud_coc_change_sub_tickets

Use this data source to get the list of COC change work order sub-orders.

## Example Usage

```hcl
variable "ticket_type" {}
variable "ticket_id" {}
variable "type" {}

data "huaweicloud_coc_change_sub_tickets" "test" {
  ticket_type = var.ticket_type
  ticket_id   = var.ticket_id
  type        = var.type
}
```

## Argument Reference

The following arguments are supported:

* `ticket_type` - (Required, String) Specifies the work order type. The value can be **change**.

* `ticket_id` - (Required, String) Specifies the change order ticket id.

* `type` - (Required, String) Specifies the resource type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tickets` - Indicates the work order details.

  The [tickets](#data_tickets_struct) structure is documented below.

<a name="data_tickets_struct"></a>
The `tickets` block supports:

* `ep_id` - Indicates the enterprise project ID.

* `resource_id` - Indicates the resource ID.

* `type` - Indicates the resource type.

* `name` - Indicates the resource name.

* `cloud_service_name` - Indicates the cloud service name.

* `domain_id` - Indicates the tenant ID.

* `region_id` - Indicates the region ID.

* `hosting_id` - Indicates the host ID.

* `properties_json` - Indicates the resource attribute information.

* `tags_json` - Indicates the resource tag information.

* `is_deleted` - Indicates whether the change order has been deleted.

* `id` - Indicates the change ticket ID.

* `main_ticket_id` - Indicates the primary key ID of the change work order.

* `parent_ticket_id` - Indicates the parent work order ID.

* `ticket_id` - Indicates the change order work ID.

* `real_ticket_id` - Indicates the change order number.

* `ticket_path` - Indicates the change order work order path.

* `target_value` - Indicates the region information of the change order sub-order.

* `target_type` - Indicates the change order sub-document type.

* `create_time` - Indicates the change order creation time.

* `update_time` - Indicates the change order update time.

* `creator` - Indicates the change order creator.

* `operator` - Indicates the change order operator.
