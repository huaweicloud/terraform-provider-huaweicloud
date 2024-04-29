---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_attachments"
description: |-
  Use this data source to get the list of CC attachments on the central network.
---

# huaweicloud_cc_central_network_attachments

Use this data source to get the list of CC attachments on the central network.

## Example Usage

```hcl
variable "central_network_id" {}

data "huaweicloud_cc_central_network_attachments" "test" { 
  central_network_id = var.central_network_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `central_network_id` - (Required, String) Specifies the central network ID.

* `attachment_id` - (Optional, String) Specifies the attachment ID.

* `name` - (Optional, String) Specifies the attachment name.

* `attachment_instance_type` - (Optional, String) Specifies the type of attachment instance.
  The valid values are **GDGW** and **ER_ROUTE_TABLE**.

* `state` - (Optional, String) Specifies the attachment status.
  The valid values are as follows:
  + **AVAILABLE**: The attachment is available.
  + **CREATING**: The attachment is being created.
  + **UPDATING**: The attachment is being updated.
  + **DELETING**: The attachment is being deleted.
  + **FREEZING**: The attachment is being frozen.
  + **UNFREEZING**: The attachment is being unfrozen.
  + **RECOVERING**: The attachment is being recovered.
  + **FAILED**: The operation on the attachment failed.
  + **DELETED**: The attachment is deleted.
  + **APPROVING**: The attachment is being approved.
  + **APPROVED**: The attachment is approved.
  + **UNAPPROVED**: The approval failed.

* `attachment_instance_id` - (Optional, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `central_network_attachments` - List of attachments on the central network.

  The [central_network_attachments](#central_network_attachments_struct) structure is documented below.

<a name="central_network_attachments_struct"></a>
The `central_network_attachments` block supports:

* `id` - The attachment ID.

* `name` - The attachment name.

* `description` - The attachment description.

* `domain_id` - The domain ID.

* `state` - The attachment status.

* `created_at` - Time when the resource was created.

* `updated_at` - Time when the resource was updated.

* `central_network_id` - The central network ID.

* `central_network_plane_id` - The central network plane ID.

* `global_connection_bandwidth_id` - The global connection bandwidth ID.

* `bandwidth_type` - The bandwidth type.

* `bandwidth_size` - The bandwidth size.

* `is_frozen` - Whether the resource is frozen.

* `enterprise_router_id` - The enterprise router ID.

* `enterprise_router_project_id` - The project ID to which the enterprise router belongs.

* `enterprise_router_region_id` - The region ID to which the enterprise router belongs.

* `enterprise_router_attachment_id` - The enterprise router attachment ID.

* `attachment_instance_type` - The attachment instance type.

* `attachment_instance_id` - The attachment instance ID.

* `attachment_id` - The ID of the enterprise router connection.

* `attachment_instance_project_id` - The project ID of the attachment instance.

* `attachment_instance_region_id` - The Region ID of the attachment instance.

* `attachment_instance_site_code` - The attachment instance site code.

* `enterprise_router_site_code` - The enterprise router site code.

* `specification_value` - Additional information about an attachment.

  The [specification_value](#central_network_attachments_specification_value_struct) structure is documented below.

<a name="central_network_attachments_specification_value_struct"></a>
The `specification_value` block supports:

* `enterprise_router_table_id` - The enterprise router table ID.

* `attached_er_id` - The attached enterprise router ID.

* `approved_state` - Approval status.

* `hosted_cloud` - Huawei Cloud or partner cloud.

* `reason` - Reason for rejecting attachment creation.
