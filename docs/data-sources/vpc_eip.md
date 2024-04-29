---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip"
description: ""
---

# huaweicloud_vpc_eip

Use this data source to get the details of an available EIP.

## Example Usage

```hcl
data "huaweicloud_vpc_eip" "by_address" {
  public_ip = "123.60.208.163"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the EIP.
  If omitted, the provider-level region will be used.

* `public_ip` - (Optional, String) Specifies the public **IPv4** address of the EIP.

* `port_id` - (Optional, String) Specifies the port id of the EIP.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the EIP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the EIP.

* `status` - The status of the EIP.

* `type` - The type of the EIP.

* `private_ip` - The private ip of the EIP.

* `ip_version` - The IP version, either 4 or 6.

* `ipv6_address` - The IPv6 address of the EIP.

* `bandwidth_id` - The bandwidth id of the EIP.

* `bandwidth_name` - The bandwidth name of the EIP.

* `bandwidth_size` - The bandwidth size of the EIP.

* `bandwidth_share_type` - The bandwidth share type of the EIP.

* `created_at` - The create time of the EIP.
