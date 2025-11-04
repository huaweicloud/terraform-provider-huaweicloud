---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_custom_line"
description: |-
  Manages a DNS custom line resource within HuaweiCloud.
---

# huaweicloud_dns_custom_line

Manages a DNS custom line resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "description" {}
variable "ip_segments" {
  type = list(string)
}

resource "huaweicloud_dns_custom_line" "test" {
  name        = var.name
  description = var.description
  ip_segments = var.ip_segments
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the custom line name.  
  The value consists of `1` to `80` characters including Chinese characters, English letters, digits, hyphens (-),
  underscores (_) and dots (.). The name of each resolution line set by one account must be unique.

* `ip_segments` - (Required, List) Specifies the IP address range.  
  The start IP address is separated from the end IP address with a hyphen (-). The IP address ranges cannot overlap.
  If the start and end IP addresses are the same, there is only one IP address in the range. Set the value to
  `IP1-IP1`. Currently, only IPv4 addresses are supported. You can specify a maximum of `50` IP address ranges.

* `description` - (Optional, String) Specifies the custom line description.  
  A maximum of `255` characters are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also custom line ID.

* `status` - The resource status.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The DNS custom line can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dns_custom_line.test <id>
```
