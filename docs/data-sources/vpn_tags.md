---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_tags"
description: |-
  Use this data source to get the list of VPN tags.
---

# huaweicloud_vpn_tags

Use this data source to get the list of VPN tags.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_vpn_tags" "test" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  Valid values are **vpn-gateway**, **customer-gateway**, **vpn-connection**, **p2c-vpn-gateways**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the list of resource tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates a tag key.

* `values` - Indicates a tag value.
