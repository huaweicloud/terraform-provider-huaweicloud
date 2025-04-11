---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_tags"
description: |-
  Use this data source to query the tag list of all resources of the same type within HuaweiCloud.
---

# huaweicloud_er_tags

Use this data source to query the tag list of all resources of the same type within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_er_tags" "test" {}
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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of all tags for resources of the same type.  
  The [tags](#er_project_tags) structure is documented below.

<a name="er_project_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

* `values` - All values corresponding to the key.
