---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_group"
description: |-
  Manages a Workspace APP group resource within HuaweiCloud.
---

# huaweicloud_workspace_app_group

Manages a Workspace APP group resource within HuaweiCloud.

## Example Usage

```hcl
variable "app_group_name" {}

resource "huaweicloud_workspace_app_group" "test" {
  name = var.app_group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `name` - (Required, String) Specifies the name of the application group.  
  The valid length is limited from `1` to `64`, only Chinese and English characters, digits, underscores (_) and
  hyphens (-) are allowed.

* `type` - (Optional, String, ForceNew) Specifies the type of the application group.
  If omitted, the defult value is **COMMON_APP**.  
  The valid values are as follows:
  + **COMMON_APP**
  + **SESSION_DESKTOP_APP**

* `server_group_id` - (Optional, String) Specifies the server group ID associated with the application group.

* `description` - (Optional, String) Specifies the description of the application group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the application group, in RFC3339 format.

## Import

The Workspace APP group resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_app_group.test <id>
```
