---
subcategory: "Application Operations Management (AOM)"
---

# huaweicloud_aom_component

Manages an AOM component resource within HuaweiCloud.

## Example Usage

```hcl
variable "topic_urn" {}

resource "huaweicloud_aom_component" "test" {
  description        = "component description"
  model_id           = "%s"
  model_type         = "APPLICATION"
  name               = "component_demo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the component name.
  Changing this creates a new resource.

* `model_id` - (Required, String) Specifies the model id. Application ID
  sub-application ID, ID length cannot exceed 36 characters, consisting of uppercase and lowercase letters and numbers.
  Changing this creates a new resource.

* `model_type` - (Required, String) Specifies the model type. Application, sub-application,
  not case sensitive, Enumeration values: APPLICATION, SUB_APPLICATION
  Changing this creates a new resource.

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `description` - (Optional, String) Specifies the component description.
  The value can be a string of 0 to 1024 characters.

* `aom_id` - (Optional, String) Specifies the aom id.

* `app_id` - (Optional, String) Specifies the application id.

* `register_type` - (Optional,String) Specifies the register type,way to register
  Enumeration values: API, CONSOLE ,SERVICE_DISCOVERY.

* `sub_app_id` - (Optional,String) Specifies the Sub application id.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The creation time. the time when the component was created.

* `creator` - The creator name.

* `modifier` - The modifier name.

* `modified_time` - The modified time. the time when the component was changed.

## Import

The component operations management can be imported using the `id` (name), e.g.

```bash
$ terraform import huaweicloud_aom_component.test component_demo
```
