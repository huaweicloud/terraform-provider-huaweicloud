---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_basic_plugin"
description: |-
  Manages a CodeArts pipeline basic plugin resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_basic_plugin

Manages a CodeArts pipeline basic plugin resource within HuaweiCloud.

## Example Usage

```hcl
variable "plugin_name" {}
variable "display_name" {}
variable "business_type" {}

resource "huaweicloud_codearts_pipeline_basic_plugin" "test" {
  plugin_name   = var.plugin_name
  display_name  = var.display_name
  business_type = var.business_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `plugin_name` - (Required, String, NonUpdatable) Specifies the basic plugin name.

* `display_name` - (Required, String) Specifies the display name.

* `business_type` - (Required, String, NonUpdatable) Specifies the service type.

* `description` - (Optional, String) Specifies the basic plugin description.

* `icon_url` - (Optional, String) Specifies the icon URL.

* `business_type_display_name` - (Optional, String, NonUpdatable) Specifies the display name of service type.

* `is_private` - (Optional, Int, NonUpdatable) Specifies whether the basic plugin is private or not.

* `plugin_composition_type` - (Optional, String, NonUpdatable) Specifies the combination extension type.

* `runtime_attribution` - (Optional, String, NonUpdatable) Specifies the runtime attributes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `plugin_name`.

* `maintainers` - Indicates the maintenance engineer.

* `unique_id` - Indicates the unique ID.

## Import

The basic plugin can be imported using `plugin_name`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_basic_plugin.test <plugin_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `is_private`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the basic plugin, or the resource definition should be updated to
align with the basic plugin. Also you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_pipeline_basic_plugin" "test" {
    ...

  lifecycle {
    ignore_changes = [
      is_private,
    ]
  }
}
```
