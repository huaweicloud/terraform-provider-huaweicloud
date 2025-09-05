---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_plugin_version"
description: |-
  Manages a CodeArts pipeline plugin version resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_plugin_version

Manages a CodeArts pipeline plugin version resource within HuaweiCloud.

## Example Usage

```hcl
variable "plugin_name" {}
variable "display_name" {}

resource "huaweicloud_codearts_pipeline_plugin_version" "test" {
  plugin_name         = var.plugin_name
  display_name        = var.display_name
  business_type       = "Build"
  version             = "0.0.1"
  runtime_attribution = "agent"

  execution_info {
    inner_execution_info = jsonencode({
      "execution_type": "COMPOSITE",
      "steps": [{
        "identifier": "1756954499618a927e16a-7a8e-48b8-b732-76751313ea53"
        "task": "official_shell_plugin",
        "name": "Shell",
        "variables": {
          "OFFICIAL_SHELL_SCRIPT_INPUT": "(decode)ZWNobyAiaGVsbG8i"
        }
      }]
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `plugin_name` - (Required, String, NonUpdatable) Specifies the plugin name.

* `display_name` - (Required, String) Specifies the display name.

* `version` - (Required, String, NonUpdatable) Specifies the version.

* `business_type` - (Required, String, NonUpdatable) Specifies the service type.
  Valid values are **Build**, **Gate**, **Deploy**, **Test** and **Normal**.

* `execution_info` - (Required, List) Specifies the execution information.
  The [execution_info](#block--execution_info) structure is documented below.

* `runtime_attribution` - (Required, String, NonUpdatable) Specifies the runtime attributes.
  Valid values are **agent** and **agentless**.

* `business_type_display_name` - (Optional, String, NonUpdatable) Specifies the display name of service type.

* `description` - (Optional, String) Specifies the basic plugin description.

* `icon_url` - (Optional, String) Specifies the icon URL.

* `input_info` - (Optional, List) Specifies the input information.
  The [input_info](#block--input_info) structure is documented below.

* `is_private` - (Optional, Int, NonUpdatable) Specifies whether the plugin is private.
  Valid values are:
  + **1**: private plugin
  + **0**: public plugin

* `version_description` - (Optional, String) Specifies the version description.

* `is_formal` - (Optional, Bool) Specifies whether the plugin is formal. Defaults to **false**.

  -> Only support updating from **false** to **true**, which means to publish the plugin version.

<a name="block--execution_info"></a>
The `execution_info` block supports:

* `inner_execution_info` - (Required, String) Specifies the inner execution information in json format.

<a name="block--input_info"></a>
The `input_info` block supports:

* `name` - (Optional, String) Specifies the name.

* `type` - (Optional, String) Specifies the type.

* `default_value` - (Optional, String) Specifies the default value.

* `layout_content` - (Optional, String) Specifies the style information..

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `unique_id` - Indicates the unique ID.

* `plugin_attribution` - Indicates the plugin attribution.

* `plugin_composition_type` - Indicates the combination extension type.

* `maintainers` - Indicates the maintenance engineer.

* `op_time` - Indicates the operation time.

* `op_user` - Indicates the operator.

## Import

The plugin version can be imported using `plugin_name` and `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_plugin_version.test <plugin_name>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `execution_info`, `is_private`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the plugin version, or the resource definition should be updated to
align with the plugin version. Also you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_pipeline_plugin_version" "test" {
    ...

  lifecycle {
    ignore_changes = [
      execution_info, is_private,
    ]
  }
}
```
