---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_endpoints"
description: |-
  Use this data source to get the list of DNS endpoints.
---

# huaweicloud_dns_endpoints

Use this data source to get the list of DNS endpoints.

## Example Usage

```hcl
variable "direction" {}

data "huaweicloud_dns_endpoints" "test" {
  direction = var.direction
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `direction` - (Required, String) Specifies the direction of the endpoint.
  The valid values can be **inbound** or **outbound**.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to which the endpoint belongs.

* `name` - (Optional, String) Specifies the name of the endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `endpoints` - Indicates the returned endpoints.

  The [endpoints](#endpoints_struct) structure is documented below.

<a name="endpoints_struct"></a>
The `endpoints` block supports:

* `id` - Indicates the endpoint ID, which is a UUID used to identify the endpoint.

* `name` - Indicates the endpoint name.

* `vpc_id` - Indicates the ID of the VPC to which the endpoint belongs.

* `direction` - Indicates the direction of the endpoint.

* `ipaddress_count` - Indicates the number of IP addresses of the endpoint.

* `resolver_rule_count` - Indicates the number of endpoint rules in the endpoint.

* `status` - Indicates the resource status.
  The value can be **PENDING_CREATE**, **ACTIVE**, **PENDING_DELETE**, or **ERROR**.

* `create_time` - Indicates the creation time. Format is **yyyy-MM-dd'T'HH:mm:ss.SSS**.

* `update_time` - Indicates the update time. Format is **yyyy-MM-dd'T'HH:mm:ss.SSS**.
