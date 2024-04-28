---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_line_group"
description: ""
---

# huaweicloud_dns_line_group

Manages a DNS line group resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "description" {}
variable "lines" {
  type = list(string)
}

resource "huaweicloud_dns_line_group" "test" {
  name        = var.name
  lines       = var.lines
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the line group name.
  The value consists of `1` to `64` characters including chinese and english letters, digits, hyphens (-), underscores (_),
  and periods (.). The name of each resource set by one account must be unique.

* `lines` - (Required, List) Specifies the list of the resolution line IDs. You should specify at least `2` different lines.

* `description` - (Optional, String) Specifies the line group description. A maximum of `255` characters are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The resource status. The value can be **ACTIVE**, **ERROR**, **FREEZE**, **DISABLE**.

* `created_at` - Resource creation time.

* `updated_at` - Resource update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The DNS line group can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_dns_line_group.test <id>
```
