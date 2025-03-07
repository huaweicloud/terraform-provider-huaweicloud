---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_ptrrecord"
description: |-
  Manages a DNS PTR record resource within HuaweiCloud.
---

# huaweicloud_dns_ptrrecord

Manages a DNS PTR record resource within HuaweiCloud.

## Example Usage

```hcl
variable "ptrrecord_name" {}
variable "eip_id" {}
variable "description" {}

resource "huaweicloud_dns_ptrrecord" "test" {
  name          = var.ptrrecord_name
  floatingip_id = var.eip_id
  description   = var.description
  ttl           = 3000

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the PTR record. If omitted,
  the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the domain name of the PTR record.  
  A domain name is case-insensitive. Uppercase letters will also be converted into lowercase letters.

* `floatingip_id` - (Required, String, ForceNew) Specifies the ID of the FloatingIP/EIP.  
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the PTR record.

* `ttl` - (Optional, Int) Specifies the time to live (TTL) of the record set (in seconds), defaults to `300`.  
  The valid value is range from `1` to `2,147,483,647`.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the PTR record.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the PTR record.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the PTR record ID, the format is `{region}:{floatingip_id}`.

* `address` - The address of the FloatingIP/EIP.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The PTR record resource can be imported using `id` (consists of the region and the EIP ID (`floatingip_id`)), e.g.

```bash
$ terraform import huaweicloud_dns_ptrrecord.test <id>
```

You can also use `region` and `floatingip_id` instead of `id`, separated by a colon (:), e.g.

```bash
$ terraform import huaweicloud_dns_ptrrecord.test <region>:<floatingip_id>
```
