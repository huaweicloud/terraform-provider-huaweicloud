---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_endpoint"
description: ""
---

# huaweicloud_dns_endpoint

Manages a DNS Endpoint in the HuaweiCloud DNS Service.

## Example Usage

```hcl
variable subnet_id {}
variable ip {}

resource "huaweicloud_dns_endpoint" "test" {
  name      = "test"
  direction = "inbound"

  ip_addresses {
    subnet_id = var.subnet_id
    ip        = var.ip
  }
  ip_addresses {
    subnet_id = var.subnet_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DNS endpoint. If omitted,
  the `region` argument of the provider will be used. Changing this creates a new DNS endpoint.

* `name` - (Required, String) Specifies the name of the DNS endpoint resource.

* `direction` - (Required, String, ForceNew) Specifies the direction of the endpoint. The value can be *inbound* or *outbound*.
  Changing this creates a new DNS endpoint.

* `ip_addresses` - (Required, List) Specifies the IP address list of the DNS endpoint.
  The List length limit range form 2 to 6.
  The [ip_address](#Address) structure is documented below.

<a name="Address"></a>
The `ip_address` block supports:

* `subnet_id` - (Required, String) Specifies the subnet id of the IP address.

* `ip` - (Optional, String) Specifies the unique IP of the IP address. If `ip` is not declare, it will be distributed by
  cloud service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `vpc_id` - The VPC ID which the subnet belongs.

* `resolver_rule_count` - The number of resolver rules.

* `status` - The status of endpoint.

* `created_at` - The created time.

* `updated_at` - The last updated time.

* `ip_addresses` - The IP address list of the DNS endpoint.

The `ip_address` block supports:

* `ip_address_id` - The ID of the IP address.

* `ip` - The IP which is associated to the endpoint.

* `status` - The status of IP address.

* `created_at` - The created time of the IP address.

* `updated_at` - The last updated time of the IP address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Endpoint can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dns_endpoint.test ff8080828a94313a018bdc88d3f3447d
```
