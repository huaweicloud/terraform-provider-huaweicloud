---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_datasource_connection_privilege"
description: ""
---

# huaweicloud_dli_datasource_connection_privilege

Using this resource to manage the privileges for the DLI datasource **enhanced** connection within HuaweiCloud.

## Example Usage

### Grant queue binding permission to the datasource enhanced connection

```hcl
variable "connection_id" {}
variable "granted_project_id" {}

resource "huaweicloud_dli_datasource_connection_privilege" "test" {
  connection_id = var.connection_id
  project_id    = var.granted_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `connection_id` - (Required, String, ForceNew) Specifies the ID of the connection to be granted.

  Changing this parameter will create a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the ID of the granted project.

  Changing this parameter will create a new resource.

* `privileges` - (Optional, List) The list of permissions granted to the connection.  
  Currently, only **BIND_QUEUE** is supported and is the default value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consisting of `connection_id` and `project_id`, the format is `<connection_id>/<project_id>`.

## Import

The datasource connection privilege detail can be imported using the `connection_id` and `project_id`, separated by a
slash, e.g.

```bash
$ terraform import huaweicloud_dli_datasource_connection_privilege.test <connection_id>/<project_id>
```
