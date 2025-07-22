---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_micro_services"
description: |-
  Use this data source to get a list of CodeArts pipeline groups.
---

# huaweicloud_codearts_pipeline_micro_services

Use this data source to get a list of CodeArts pipeline groups.

## Example Usage

```hcl
variable "codearts_project_id" {}

data "huaweicloud_codearts_pipeline_micro_services" "test" {
  project_id = var.codearts_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `name` - (Optional, String) Specifies the micro service name.

* `sort_dir` - (Optional, String) Specifies the sorting sequence.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `micro_services` - Indicates the micro service list.
  The [micro_services](#attrblock--micro_services) structure is documented below.

<a name="attrblock--micro_services"></a>
The `micro_services` block supports:

* `id` - Indicates the micro service ID.

* `name` - Indicates the micro service name.

* `description` - Indicates the micro service description.

* `is_followed` - Indicates whether the micro service is followed.

* `repos` - Indicates the repository information.
  The [repos](#attrblock--micro_services--repos) structure is documented below.

* `parent_id` - Indicates the micro service parent ID.

* `status` - Indicates the micro service status.

* `type` - Indicates the micro service type.

* `create_time` - Indicates the create time.

* `creator_id` - Indicates the creator ID.

* `creator_name` - Indicates the creator name.

* `update_time` - Indicates the update time.

* `updater_id` - Indicates the updater ID.

* `updater_name` - Indicates the updater name.

<a name="attrblock--micro_services--repos"></a>
The `repos` block supports:

* `branch` - Indicates the branch.

* `endpoint_id` - Indicates the endpoint ID.

* `git_url` - Indicates the Git address of the Git repository.

* `http_url` - Indicates the HTTP address of the Git repository.

* `language` - Indicates the language.

* `repo_id` - Indicates the repository ID.

* `type` - Indicates the repository type.
