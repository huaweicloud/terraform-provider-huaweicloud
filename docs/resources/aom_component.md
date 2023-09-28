---
subcategory: "Application Operations Management (AOM)"
---

# huaweicloud_aom_component

Manages an AOM component resource within HuaweiCloud.

## Example Usage

```hcl
variable "model_id" {}

resource "huaweicloud_aom_component" "test" {
  description = "component description"
  model_id    = var.model_id
  model_type  = "APPLICATION"
  name        = "component_demo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `name` - (Required, String) Specifies the component name.

* `model_id` - (Required, String) Specifies the application or sub-application ID, which can contain up to 36
  characters. Only letters and digits are allowed.

* `model_type` - (Required, String) Specifies the application type. Value options: **APPLICATION**, **SUB_APPLICATION**.
  The value is case-insensitive.

* `description` - (Optional, String) Specifies the component description. The value can be a string of 0 to 1024
  characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The component id.

* `create_time` - The creation time, the time when the component was created.

* `creator` - The creator name.

* `modifier` - The modifier name.

* `modified_time` - The modified time, the time when the component was changed.

* `aom_id` - The aom id.

* `app_id` - The application id.

* `register_type` - The registration method. Enumeration values: **API**, **CONSOLE**, **SERVICE_DISCOVERY**.

* `sub_app_id` - The sub application id.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The component operations management can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_component.test <id>
```
