---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_applications"
description: |-
  Use this data source to get the list of CodeArts deploy applications.
---

# huaweicloud_codearts_deploy_applications

Use this data source to get the list of CodeArts deploy applications.

## Example Usage

```hcl
variable "project_id" {}

data "huaweicloud_codearts_deploy_applications" "test" {
  project_id = var.project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the project ID.

* `group_id` - (Optional, String) Specifies the application group ID.
  Enter **no_grouped** to query ungrouped applications.

* `states` - (Optional, List) Specifies the application status list.
  Values can be as follows:
  + **abort**: Deployment suspended.
  + **failed**: Deployment failed.
  + **not_started**: Execution canceled.
  + **pending**: Queuing.
  + **running**: Deployment in progress.
  + **succeeded**: Deployment succeeded.
  + **timeout**: Deployment times out.
  + **not_executed**: Deployment not executed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - Indicates the application list
  The [applications](#attrblock--applications) structure is documented below.

<a name="attrblock--applications"></a>
The `applications` block supports:

* `arrange_infos` - Indicates the deployment task information
  The [arrange_infos](#attrblock--applications--arrange_infos) structure is documented below.

* `id` - Indicates the application ID.

* `name` - Indicates the application name.

* `can_copy` - Indicates whether the user has permission to clone application.

* `can_create_env` - Indicates whether the user has permission to create environment in application.

* `can_delete` - Indicates whether the user has permission to delete application.

* `can_disable` - Indicates whether the user has permission to disable application.

* `can_execute` - Indicates whether the user has permission to deploy.

* `can_manage` - Indicates whether the user has permission to modify application permission.

* `can_modify` - Indicates whether the user has permission to modify application.

* `can_view` - Indicates whether the user has permission to view application.

* `create_tenant_id` - Indicates the created tenant ID.

* `create_user_id` - Indicates the creator user ID.

* `created_at` - Indicates the created time.

* `deploy_system` - Indicates the deployment type.

* `duration` - Indicates the deployment duration.

* `end_time` - Indicates the deployment end time.

* `execution_state` - Indicates the execution status.

* `execution_time` - Indicates the latest execution time.

* `executor_id` - Indicates the executor user ID.

* `executor_nick_name` - Indicates the executor user name.

* `is_care` - Indicates whether application is saved to favorites.

* `is_disable` - Indicates whether the application is disabled.

* `project_name` - Indicates the project name.

* `release_id` - Indicates the release ID.

* `updated_at` - Indicates the updated time.

<a name="attrblock--applications--arrange_infos"></a>
The `arrange_infos` block supports:

* `deploy_system` - Indicates the deployment task type.

* `id` - Indicates the task ID.

* `state` - Indicates the deployment task status.
