---
subcategory: "CodeArts Inspector"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_inspector_host_group"
description: |-
  Manages a CodeArts inspector host group resource within HuaweiCloud.
---

# huaweicloud_codearts_inspector_host_group

Manages a CodeArts inspector host group resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_codearts_inspector_host_group" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the host group name.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The host group can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_inspector_host_group.test <id>
```
