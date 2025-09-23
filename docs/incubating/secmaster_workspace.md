---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workspace"
description: |-
  Manages a SecMaster workspace resource within HuaweiCloud.
---

# huaweicloud_secmaster_workspace

Manages a SecMaster workspace resource within HuaweiCloud.

-> Destroying this resource will not change the status of the workspace resource.

## Example Usage

### Basic Example

```hcl
variable "name" {}
variable "project_name" {}

resource "huaweicloud_secmaster_workspace" "test" {
  name         = var.name
  project_name = var.project_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the workspace.

* `project_name` - (Required, String, NonUpdatable) Specifies name of the project in which to create the workspace.

* `description` - (Optional, String) Specifies the description of the workspace.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies ID of the enterprise project
  in which to create the workspace.

* `enterprise_project_name` - (Optional, String, NonUpdatable) Specifies name of the enterprise project
  in which to create the workspace.

* `view_bind_id` - (Optional, String) Specifies the space ID bound to the view.

* `is_view` - (Optional, String, NonUpdatable) Specifies whether the view is used.

* `tags` - (Optional, Map, NonUpdatable) Specifies the tags of the workspace in key/pair format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time.

* `update_time` - The update time.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

* `domain_id` - The domain ID.

* `domain_name` - The domain name.

* `view_bind_name` - The bound view name.

* `workspace_agency_list` - The workspace agency list.
  The [workspace_agency_list](#workspace_agency_list) structure is documented below.

<a name="workspace_agency_list"></a>
The `workspace_agency_list` block supports:

* `project_id` - The project ID.

* `id` - The ID of the workspace agency.

* `name` - The name of the workspace agency.

* `region_id` - The region ID of the workspace agency.

* `workspace_attribution` - The workspace attribution of the workspace agency.

* `agency_version` - The agency version of the workspace agency.

* `domain_id` - The domain ID of the workspace agency.

* `domain_name` - The domain name of the workspace agency.

* `iam_agency_id` - The IAM agency ID of the workspace agency.

* `iam_agency_name` - The IAM agency name of the workspace agency.

* `resource_spec_code` - The resource spec code of the workspace agency.

* `selected` - Whether to be selected.

## Import

SecMaster workspace can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_secmaster_workspace.test <id>
```
