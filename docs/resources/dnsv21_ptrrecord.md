---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dnsv21_ptrrecord"
description: |-
  Manages a DNS PTR record resource within HuaweiCloud.
---

# huaweicloud_dnsv21_ptrrecord

Manages a DNS PTR record resource within HuaweiCloud.

## Example Usage

```hcl
variable "ptrrecord_names" {}
variable "eip_id" {}

resource "huaweicloud_dns_ptrrecord" "test" {
  names       = var.ptrrecord_names
  publicip_id = var.eip_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `names` - (Required, List) Specifies the domain names of the PTR record.

* `publicip_id` - (Required, String) Specifies the ID of the EIP.

* `description` - (Optional, String) Specifies the description of the PTR record.

* `ttl` - (Optional, Int) Specifies the time to live (TTL) of the record set (in seconds).

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the PTR record.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the PTR record.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `address` - The address of the EIP.

* `status` - The status of the PTR record.

## Import

The PTR record can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dnsv21_ptrrecord.test <id>
```
