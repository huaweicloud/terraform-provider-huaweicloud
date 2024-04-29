---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eips"
description: ""
---

# huaweicloud_vpc_eips

Use this data source to get a list of EIPs.

## Example Usage

An example filter by name and tag

```hcl
variable "public_ip" {}

data "huaweicloud_vpc_eips" "eip" {
  public_ips = [var.public_ip]

  tags = {
    foo = "bar"
  }
}

output "eip_ids" {
  value = data.huaweicloud_vpc_eips.eip.eips[*].id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available EIPs in the current region.
 All EIPs that meet the filter criteria will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the EIP. If omitted, the provider-level region
  will be used.

* `ids` - (Optional, List) Specifies an array of one or more IDs of the desired EIP.

* `public_ips` - (Optional, List) Specifies an array of one or more public ip addresses of the desired EIP.

* `private_ips` - (Optional, List) Specifies an array of one or more private ip addresses of the desired EIP.

* `port_ids` - (Optional, List) Specifies an array of one or more port ids which bound to the desired EIP.

* `ip_version` - (Optional, Int) Specifies ip version of the desired EIP. The options are:
    + `4`: IPv4.
    + `6`: IPv6.

  The default value is `4`.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID which the desired EIP belongs to.

* `tags` - (Optional, Map) Specifies the included key/value pairs which associated with the desired EIP.

 -> A maximum of 10 tag keys are allowed for each query operation. Each tag key can have up to 10 tag values.
  The tag key cannot be left blank or set to an empty string. Each tag key must be unique, and each tag value in a
  tag must be unique, use commas(,) to separate the multiple values. An empty for values indicates any value.
  The values are in the OR relationship.

## Attribute Reference

The following attributes are exported:

* `id` - Indicates a data source ID.

* `eips` - Indicates a list of all EIPs found. Structure is documented below.

The `eips` block supports:

* `id` - The ID of the EIP.
* `name` - The name of the EIP.
* `public_ip` - The public ip address of the EIP.
* `private_ip` - The private ip address of the EIP.
* `public_ipv6` - The public ipv6 address of the EIP.
* `port_id` - The port id bound to the EIP.
* `ip_version` - The ip version of the EIP.
* `status` - The status of the EIP.
* `type` - The type of the EIP.
* `enterprise_project_id` - The the enterprise project ID of the EIP.
* `bandwidth_id` - The bandwidth id of the EIP.
* `bandwidth_name` - The bandwidth name of the EIP.
* `bandwidth_size` - The bandwidth size of the EIP.
* `bandwidth_share_type` - The bandwidth share type of the EIP.
* `tags` - The key/value pairs which associated with the EIP.
* `created_at` - The create time of the EIP.
