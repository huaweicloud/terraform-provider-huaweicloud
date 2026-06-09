---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_component_template"
description: |-
  Manages a component template resource within HuaweiCloud.
---

# huaweicloud_secmaster_component_template

Manages a component template resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "component_id" {}
variable "template_name" {}
variable "task_config" {}

resource "huaweicloud_secmaster_component_template" "test" {
  workspace_id  = var.workspace_id
  component_id  = var.component_id
  template_name = var.template_name
  task_config   = var.task_config
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the component template
  belongs.

* `component_id` - (Required, String) Specifies the ID of the component to which the template belongs.

* `template_name` - (Required, String) Specifies the name of the component template.

* `task_config` - (Required, String) Specifies the action configuration content of the component template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the component template.

* `update_time` - The latest update time of the component template.

## Import

The component template can be imported using the `workspace_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_component_template.test <workspace_id>/<id>
```
