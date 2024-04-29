---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_attachments"
description: ""
---

# huaweicloud_er_attachments

Use this data source to filter ER attachments within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_er_attachments" "test" {
  instance_id = var.instance_id

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the ER attachments are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ER instance ID to which the attachment belongs.

* `attachment_id` - (Optional, String) Specifies the specified attachment ID used to query.

* `type` - (Optional, String) Specifies the resource type to be filtered.  
  The valid values are as follows:
  + **vpc**: Virtual private cloud.
  + **vpn**: VPN gateway.
  + **vgw**: Virtual gateway of cloud private line.
  + **peering**: Peering connection, through the cloud connection (CC) to load ERs in different regions to create a
    peering connection.

* `name` - (Optional, String) Specifies the name used to filter the attachments.

* `resource_id` - (Optional, String) Specifies the associated resource ID used to filter the attachments.

* `status` - (Optional, String) Specifies the status used to filter the attachments.
  The valid values are as follows:
  + **available**
  + **failed**
  + **pending_acceptance**
  + **rejected**

* `tags` - (Optional, Map) The key/value pairs used to filter the attachments.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `attachments` - All attachments that match the filter parameters.  
  The [object](#er_data_attachments) structure is documented below.

<a name="er_data_attachments"></a>
The `attachments` block supports:

* `id` - The attachment ID.

* `name` - The attachment name.

* `description` - The description of the attachment.

* `status` - The current status of the attachment.

* `associated` - Whether this attachment has been associated.

* `resource_id` - The associated resource ID.

* `created_at` - The creation time of the attachment.

* `updated_at` - The latest update time of the attachment.

* `tags` - The key/value pairs to associate with the attachment.

* `type` - The attachment type.

* `route_table_id` - The associated route table ID.
