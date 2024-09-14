---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_ptrrecord"
description: ""
---

# huaweicloud_dns_ptrrecord

Manages a DNS PTR record in the HuaweiCloud DNS Service.

## Example Usage

```hcl
resource "huaweicloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_dns_ptrrecord" "ptr_1" {
  name          = "ptr.example.com."
  description   = "An example PTR record"
  floatingip_id = huaweicloud_vpc_eip.eip_1.id
  ttl           = 3000

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the PTR record. If omitted,
  the `region` argument of the provider will be used. Changing this creates a new PTR record.

* `name` - (Required, String) Specifies the domain name of the PTR record. A domain name is case-insensitive.
  Uppercase letters will also be converted into lowercase letters.

* `floatingip_id` - (Required, String, ForceNew) Specifies the ID of the FloatingIP/EIP.
  Changing this creates a new PTR record.

* `description` - (Optional, String) Specifies the description of the PTR record.

* `ttl` - (Optional, Int) Specifies the time to live (TTL) of the record set (in seconds).
  The valid value is range from `1` to `2,147,483,647`.

* `tags` - (Optional, Map) Tags key/value pairs to associate with the PTR record.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the PTR record. Changing this
  creates a new PTR record.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The PTR record ID, which is in {region}:{floatingip_id} format.

* `address` - The address of the FloatingIP/EIP.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

PTR records can be imported using region and floatingip/eip ID, separated by a colon(:), e.g.

```bash
$ terraform import huaweicloud_dns_ptrrecord.ptr_1 cn-north-1:d90ce693-5ccf-4136-a0ed-152ce412b6b9
```
