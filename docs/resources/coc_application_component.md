---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_application_component"
description: |-
  Manages a COC application component resource within HuaweiCloud.
---

# huaweicloud_coc_application_component

Manages a COC application component resource within HuaweiCloud.

## Example Usage

```hcl
variable "application_id" {}
variable "name" {}

resource "huaweicloud_coc_application_component" "test" {
  application_id = var.application_id
  name           = var.name
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String, NonUpdatable) Specifies the application ID.

* `name` - (Required, String) Specifies the component name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `code` - Indicates the component code.

* `ep_id` - Indicates the enterprise project ID.

## Import

The COC application component can be imported using `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_coc_application_component.test <id>
```
