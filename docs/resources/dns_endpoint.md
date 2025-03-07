---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_endpoint"
description: |-
  Manages a DNS endpoint resource within HuaweiCloud.
---

# huaweicloud_dns_endpoint

Manages a DNS endpoint resource within HuaweiCloud.

-> For the same subnet, only one of the `huaweicloud_dns_endpoint` and `huaweicloud_dns_endpoint_assignment` resources
   is allowed to manage an endpoint. We recommend you use the `huaweicloud_dns_endpoint_assignment` resource.

## Example Usage

### Create an endpoint in the same subnet and associate two IPs

```hcl
variable "endpoint_name" {}
variable "subnet_id" {}

resource "huaweicloud_dns_endpoint" "test" {
  name      = var.endpoint_name
  direction = "inbound"

  ip_addresses {
    subnet_id = var.subnet_id
  }
  ip_addresses {
    subnet_id = var.subnet_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DNS endpoint. If omitted,
  the `region` argument of the provider will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the DNS endpoint.  
  The name valid length is limited from `1` to `64` characters. Only Chinese and English characters, digits and
  special characters (-._) are allowed.

* `direction` - (Required, String, ForceNew) Specifies the direction of the endpoint.  
  Changing this parameter will create a new resource.  
  The valid values are as follows:
  + **inbound**
  + **outbound**

* `ip_addresses` - (Required, List) Specifies the list of the IP addresses of the endpoint.  
  The valid length of the IP address list ranges form `2` to `6`.  
  The [ip_address](#endpoint_ip_addresses) structure is documented below.

<a name="endpoint_ip_addresses"></a>
The `ip_address` block supports:

* `subnet_id` - (Required, String) Specifies the subnet ID of the IP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `vpc_id` - The ID of the VPC to which the subnet belongs.

* `resolver_rule_count` - The number of resolver rules.

-> The newly added resolver rules needs to wait until the next time `terraform refresh` command is executed before that value
   can be refreshed.

* `status` - The status of endpoint.

* `created_at` - The creation time of the endpoint.

* `updated_at` - The latest update time of the endpoint.

* `ip_addresses` - The list of the IP addresses of the endpoint.
  The [ip_addresses](#endpoint_ip_addresses_attr) structure is documented below.

<a name="endpoint_ip_addresses_attr"></a>
The `ip_address` block supports:

* `ip_address_id` - The ID of the IP address.

* `status` - The status of IP address.

* `created_at` - The creation time of the IP address.

* `updated_at` - The latest update time of the IP address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Endpoint can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dns_endpoint.test <id>
```
