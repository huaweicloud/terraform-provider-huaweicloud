---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_parser"
description: |-
  Manages a collector parser resource within HuaweiCloud.
---

# huaweicloud_secmaster_collector_parser

Manages a collector parser resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "title" {}

resource "huaweicloud_secmaster_collector_parser" "test" {
  workspace_id = var.workspace_id
  title        = var.title
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the collector parser belongs.

* `title` - (Required, String, NonUpdatable) Specifies the title of the collector parser.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the collector parser.

* `parser_id` - (Optional, String, NonUpdatable) Specifies the ID of the parser.

* `modules` - (Optional, List, NonUpdatable) Specifies the list of module information.

  The [modules](#modules_struct) structure is documented below.

<a name="modules_struct"></a>
The `modules` block supports:

* `name` - (Optional, String) Specifies the module name.

* `template_id` - (Optional, String) Specifies the template UUID.

* `connection_module_id` - (Optional, String) Specifies the associated module UUID.

* `module_id` - (String) The unique UUID of the current module.

* `children` - (Optional, List) Specifies the list of sub-modules.

  The [children](#children_struct) structure is documented below.

* `fields` - (Optional, List) Specifies the list of field configurations.

  The [fields](#fields_struct) structure is documented below.

<a name="children_struct"></a>
The `children` block supports:

* `name` - (Optional, String) Specifies the module name.

* `template_id` - (Optional, String) Specifies the template UUID.

* `connection_module_id` - (Optional, String) Specifies the associated module UUID.

* `fields` - (Optional, List) Specifies the list of field configurations.

  The [fields](#fields_struct) structure is documented below.

<a name="fields_struct"></a>
The `fields` block supports:

* `name` - (Optional, String) Specifies the field name.

* `type` - (Optional, String) Specifies the field type.

* `value` - (Optional, String) Specifies the field value.

* `other` - (Optional, String) Specifies other supplementary information.

* `template_field_id` - (Optional, String) Specifies the template field UUID.

* `connection_module_id` - (Optional, String) Specifies the associated module UUID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the parser UUID returned by the API.

* `channel_refer_count` - The reference count of the parser.

## Import

The collector parser can be imported using the `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_collector_parser.test <workspace_id>/<id>
```
