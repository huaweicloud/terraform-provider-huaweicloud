---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_resource_tags"
description: |-
  Use this data source to query the tag list of a specifies resource within HuaweiCloud.
---

# huaweicloud_er_resource_tags

Use this data source to query the tag list of a specifies resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_er_resource_tags" "test" {
  resource_type = "instance"
  resource_id   = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource tags.  
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type to which the tags belong that to be queried.  
  The valid values are as follows:
  + **instance**: enterprise router instance.
  + **route-table**: route table.
  + **vpc-attachment**: VPC connection.
  + **vgw-attachment**: virtual gateway connection.
  + **peering-attachment**: peering connection.
  + **vpn-attachment**: VPN gateway connection.
  + **ecn-attachment**: enterprise network connection.
  + **cfw-attachment**: cloud firewall connection.

* `resource_id` - (Required, String) Specifies the resource ID to which the tags belong that to be queried.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tags of a specified resource.
