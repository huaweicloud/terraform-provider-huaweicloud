---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_group_tags"
description: |-
  Use this data source to query the tags of the server group within HuaweiCloud Workspace.
---

# huaweicloud_workspace_app_server_group_tags

Use this data source to query the tags of the server group within HuaweiCloud Workspace.

## Example Usage

### Querying tags of all server groups

```hcl
data "huaweicloud_workspace_app_server_group_tags" "test" {}
```

### Querying tags of the specific server group

```hcl
variable "server_group_id" {}

data "huaweicloud_workspace_app_server_group_tags" "test" {
  server_group_id = var.server_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the server group tags are located.  
  If omitted, the provider-level region will be used.

* `server_group_id` - (Optional, String) Specifies the ID of the server group to which the tags belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tag list of the server group.  
  The [tags](#workspace_app_server_group_tags_attr) structure is documented below.

<a name="workspace_app_server_group_tags_attr"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - The value list of the tag.
