---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_cmdb_component"
description: ""
---

# huaweicloud_aom_cmdb_component

Manages an AOM component resource within HuaweiCloud.

## Example Usage

```hcl
variable "app_id" {}

resource "huaweicloud_aom_cmdb_component" "test" {
  name        = "com_demo"
  model_id    = var.app_id
  model_type  = "APPLICATION"
  description = "component description"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the component. The value can contain 2 to 64 characters.
  Only letters, digits, underscores (_), hyphens (-), and periods (.) are allowed.

* `model_id` - (Required, String, ForceNew) Specifies the application or sub-application ID.
  Changing this parameter will create a new resource.

* `model_type` - (Required, String, ForceNew) Specifies the application type. The valid values are **APPLICATION** and **SUB_APPLICATION**.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description about the component.
  The description can contain a maximum of 255 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
* `register_type` - The register type of the component.
* `app_id` - The application id.
* `sub_app_id` - The sub-application id.
* `created_at` - The creation time.

## Import

The AOM component can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_cmdb_component.test 6bbcad9cbddf4a60abaf5358a9339c98
```
