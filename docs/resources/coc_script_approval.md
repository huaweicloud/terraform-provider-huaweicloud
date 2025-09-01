---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_approval"
description: |-
  Manages a COC script approve resource within HuaweiCloud.
---

# huaweicloud_coc_script_approval

Manages a COC script approve resource within HuaweiCloud.

~> Deleting script approve resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "script_uuid" {}

resource "huaweicloud_coc_script_approval" "test" {
  script_uuid = var.script_uuid
  status      = "APPROVED"
}
```

## Argument Reference

The following arguments are supported:

* `script_uuid` - (Required, String, NonUpdatable) Specifies the approval script ID.

* `status` - (Required, String, NonUpdatable) Specifies the approval status.
  Values can be **APPROVED** or **REJECTED**.

* `comments` - (Optional, String, NonUpdatable) Specifies the approval comments.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The script ID, which equals to `script_uuid`.
