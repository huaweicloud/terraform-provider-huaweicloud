---
subcategory: CodeArts
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_project"
description: ""
---

# huaweicloud_codearts_project

Manages a Project resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_codearts_project" "test" {
  name = "demo_project"
  type = "scrum"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The project name.  
  The name can contain `1` to `128` characters.

* `type` - (Required, String, ForceNew) The type of project.  
  The valid values are **scrum**, **xboard**, **basic**, **phoenix**.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description about the project.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project ID of the project.  
  Value 0 indicates the default enterprise project.

  Changing this parameter will create a new resource.

* `source` - (Optional, String, ForceNew) The source of project.

  Changing this parameter will create a new resource.

* `template_id` - (Optional, Int, ForceNew) The template id which used to create project.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `archive` - Whether the project is archived.

* `project_code` - The project code.

* `project_num_id` - The number id of project.

## Import

The project can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_project.test 0ce123456a00f2591fabc00385ff1234
```
