---
subcategory: "CodeArts Build"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_build_tasks"
description: |-
  Use this data source to get a list of CodeArts build tasks.
---

# huaweicloud_codearts_build_tasks

Use this data source to get a list of CodeArts build tasks.

## Example Usage

```hcl
variable "codearts_project_id" {}

data "huaweicloud_codearts_build_tasks" "test" {
  project_id = var.codearts_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `build_status` - (Optional, String) Specifies the build status filter condition.

* `by_group` - (Optional, Bool) Specifies whether to group.

* `creator_id` - (Optional, String) Specifies the creator ID.

* `group_path_id` - (Optional, String) Specifies the group ID.

* `search` - (Optional, String) Specifies the search condition.

* `sort_field` - (Optional, String) Specifies the sorting field.

* `sort_order` - (Optional, String) Specifies the sorting order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the task list.
  The [tasks](#attrblock--tasks) structure is documented below.

<a name="attrblock--tasks"></a>
The `tasks` block supports:

* `build_project_id` - Indicates the build project ID.

* `build_time` - Indicates the build time.

* `creator` - Indicates the task creator.

* `disabled` - Indicates whether it is disabled.

* `favorite` - Indicates whether it is favorited.

* `health_score` - Indicates the health score.

* `id` - Indicates the task ID.

* `is_copy` - Indicates whether there is permission to copy the task.

* `is_delete` - Indicates whether there is permission to delete the task.

* `is_execute` - Indicates whether there is permission to execute the task.

* `is_finished` - Indicates whether it has ended.

* `is_forbidden` - Indicates whether there is permission to disable the task.

* `is_modify` - Indicates whether there is permission to modify the task.

* `is_view` - Indicates whether there is permission to view the task.

* `last_build_status` - Indicates the latest build status.

* `last_build_time` - Indicates the latest execution time.

* `last_build_user` - Indicates the last build user.

* `last_build_user_id` - Indicates the last build user ID.

* `last_job_running_status` - Indicates the last build time.

* `name` - Indicates the task name.

* `repo_id` - Indicates the code repository ID.

* `scm_type` - Indicates the code repository type.

* `scm_web_url` - Indicates the code repository web address.

* `source_code` - Indicates the code source.

* `trigger_type` - Indicates the trigger type.

* `user_name` - Indicates the user name.
