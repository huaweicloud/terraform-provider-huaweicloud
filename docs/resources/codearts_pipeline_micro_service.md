---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_micro_service"
description: |-
  Manages a CodeArts pipeline mirco service resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_micro_service

Manages a CodeArts pipeline mirco service resource within HuaweiCloud.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "name" {}

resource "huaweicloud_codearts_pipeline_micro_service" "test" {
  project_id = var.codearts_project_id
  name       = var.name
  type       = "microservice"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.

* `name` - (Required, String, NonUpdatable) Specifies the micro service type.

* `type` - (Required, String, NonUpdatable) Specifies the micro service name.

* `parent_id` - (Optional, String, NonUpdatable) Specifies the micro service parent ID.

* `repos` - (Optional, List) Specifies the repository information.
  The [repos](#block--repos) structure is documented below.

* `description` - (Optional, String) Specifies the micro service description.

* `is_followed` - (Optional, Bool) Specifies whether the micro service is followed.

<a name="block--repos"></a>
The `repos` block supports:

* `branch` - (Required, String) Specifies the branch.

* `git_url` - (Required, String) Specifies the Git address of the Git repository.

* `http_url` - (Required, String) Specifies the HTTP address of the Git repository.

* `language` - (Required, String) Specifies the language.

* `repo_id` - (Required, String) Specifies the repository ID.

* `type` - (Required, String) Specifies the repository type.

* `endpoint_id` - (Optional, String) Specifies the endpoint ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `creator_id` - Indicates the creator ID.

* `updater_id` - Indicates the updater ID.

* `status` - Indicates the micro service status.

* `create_time` - Indicates the create time.

* `update_time` - Indicates the update time.

* `creator_name` - Indicates the creator name.

* `updater_name` - Indicates the updater name.

## Import

The micro service can be imported using `project_id` and `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_micro_service.test <project_id>/<id>
```
