---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_endpoint_assignment"
description: |-
  Manages a DNS endpoint assignment resource within HuaweiCloud.
---

# huaweicloud_dns_endpoint_assignment

Manages a DNS endpoint assignment resource within HuaweiCloud.

-> For the same subnet, only one of the `huaweicloud_dns_endpoint` and `huaweicloud_dns_endpoint_assignment` resources
   is allowed to manage an endpoint. We recommend using this resource to replace the `huaweicloud_dns_endpoint` resource.

## Example Usage

### Create an endpoint in the same subnet

```hcl
variable "endpoint_name" {}
variable "subnet_id" {}
variable "ip_addresses" {
  type = list(string)
}

resource "huaweicloud_dns_endpoint_assignment" "test" {
  name      = var.endpoint_name
  direction = "inbound"

  dynamic "assignments" {
    for_each = range(length(var.ip_addresses))
    content {
      subnet_id  = var.subnet_id
      ip_address = var.ip_addresses[assignments.key]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the endpoint.  
  The name valid length is limited from `1` to `64` characters. Only Chinese and English characters, digits and
  special characters (-._) are allowed.

* `direction` - (Required, String, ForceNew) Specifies the direction of the endpoint.  
  Changing this parameter will create a new resource.  
  The valid values are as follows:
  + **inbound**
  + **outbound**
  
* `assignments` - (Required, List) Specifies the list of the IP addresses of the endpoint.  
  The valid length of the `assignments` ranges from `2` to `6`.  
  The [assignments](#endpoint_assignments) structure is documented below.

<a name="endpoint_assignments"></a>
The `assignments` block supports:

* `subnet_id` - (Required, String) Specifies the subnet ID to which the IP address belongs.

* `ip_address` - (Required, String) Specifies the IP address associated with the endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also endpoint ID.

* `assignments` - The list of the IP addresses of the endpoint.
  The [assignments](#endpoint_attr_assignments) structure is documented below.

* `vpc_id` - The VPC ID associated with the endpoint.

* `status` - The current status of the endpoint.

* `created_at` - The creation time of the endpoint, in RFC3339 format.

<a name="endpoint_attr_assignments"></a>
The `assignments` block supports:

* `ip_address_id` - The ID of the IP address associated with the endpoint.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The DNS endpoint resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dns_endpoint_assignment.test <id>
```
