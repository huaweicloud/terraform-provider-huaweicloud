---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_resource"
description: ""
---

# huaweicloud_dataarts_factory_resource

Manages DataArts Factory resource within HuaweiCloud.

## Example Usage

### Create a resource with OBS path

```hcl
variable "workspace_id" {}
variable "name" {}
variable "directory" {}

resource "huaweicloud_dataarts_factory_resource" "test" {
  workspace_id = var.workspace_id
  name         = var.name
  type         = "jar"
  location     = "obs://test/main.jar"
  directory    = var.directory
}
```

### Create a resource with HDFS path

```hcl
variable "workspace_id" {}
variable "name" {}
variable "directory" {}

resource "huaweicloud_dataarts_factory_resource" "test" {
  workspace_id = var.workspace_id
  name         = var.name
  type         = "jar"
  location     = "hdfs://test/main.jar"
  directory    = var.directory
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID which the resource in.
  Changing this creates a new resource.

* `directory` - (Required, String) Specifies the directory where the resource is located.

* `name` - (Required, String) Specifies the resource name.

* `type` - (Required, String) Specifies the resource type. The valid values are **archive**, **file**,
  **jar** and **pyFile**.

* `location` - (Required, String) Specifies the path of the file. Currently, only OBS paths and HDFS paths
  are supported.

* `depend_packages` - (Optional, List) Specifies an array of dependent files.
  The [depend_packages](#DataArts_Factory_Resource_Depend_Packages) structure is documented below.

* `description` - (Optional, String) Specifies the resource description.

<a name="DataArts_Factory_Resource_Depend_Packages"></a>
The `depend_packages` block supports:

* `location` - (Required, String) Specifies the path of the dependent file. Currently, only OBS paths is
  supported.

* `type` - (Required, String) Specifies the type of the dependent file. The valid values are **file**,
  **jar** and **pyFile**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

DataArts Factory resource can be imported using `<workspace_id>/<id>`, e.g.

```bash
$ terraform import huaweicloud_dataarts_Factory_resource.test <workspace_id>/<id>
```
