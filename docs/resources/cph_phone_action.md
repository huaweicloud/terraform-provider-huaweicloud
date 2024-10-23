---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_phone_action"
description: |-
  Manages a CPH phone action resource within HuaweiCloud.
---

# huaweicloud_cph_phone_action

Manages a CPH phone action resource within HuaweiCloud.

## Example Usage

```hcl
variable "phone_id" {}

resource "huaweicloud_cph_phone_action" "test" {
  action = "reset"

  phones {
    phone_id = var.phone_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `action` - (Required, String) Specifies the CPH phone action. The valid value can be **reset**, **restart**, **stop**.

* `phones` - (Required, List) Specifies the CPH phones.
  The [phones](#phones_structure) structure is documented below.

* `image_id` - (Optional, String) Specifies the image ID of the CPH phone.

<a name="phones_structure"></a>
The `phones` block supports:

* `phone_id` - (Required, String) Specifies the ID of the CPH phone.

* `property` - (Optional, String)  Specifies the property of the CPH phone.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
