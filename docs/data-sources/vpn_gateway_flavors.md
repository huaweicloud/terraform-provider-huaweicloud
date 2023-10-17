---
subcategory: "Virtual Private Network (VPN)"
---

# huaweicloud_vpn_gateway_flavors

Use this data source to get the list of VPN gateway flavors.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_vpn_gateway_flavors" "test"{
  availability_zone = "cn-north-4"
  name              = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Required, String) Specifies the availability zone to get the flavors.

* `name` - (Optional, String) Specifies the flavor name.

* `attachment_type` - (Optional, String) Specifies the attachment type. The value can be: **vpc** and **er**.
  Defaults to **vpc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `flavors` - The list of flavors.
