---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_application_deployment_records"
description: |-
  Use this data source to get the list of CodeArts deploy application deployment records.
---

# huaweicloud_codearts_deploy_application_deployment_records

Use this data source to get the list of CodeArts deploy application deployment records.

## Example Usage

```hcl
variable "project_id" {}
variable "task_id" {}
variable "start_date" {}
variable "end_date" {}

data "huaweicloud_codearts_deploy_application_deployment_records" "test" {
  project_id = var.project_id
  task_id    = var.task_id
  start_date = var.start_date
  end_date   = var.end_date
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the project ID for CodeArts service.

* `task_id` - (Required, String) Specifies the deployment task ID.

* `start_date` - (Required, String) Specifies the start time. The value format is **yyyy-mm-dd**.

* `end_date` - (Required, String) Specifies the end time. The value format is **yyyy-mm-dd**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the record list.
  The [records](#attrblock--records) structure is documented below.

<a name="attrblock--records"></a>
The `records` block supports:

* `id` - Indicates the record ID.

* `duration` - Indicates the deployment duration.

* `operator` - Indicates the operator user name.

* `release_id` - Indicates the deployment record sequence number.

* `state` - Indicates the application status.

* `start_time` - Indicates the start time of application deployment. The value format is **yyyy-mm-dd hh:mm:ss**.

* `end_time` - Indicates the end time of application deployment. The value format is **yyyy-mm-dd hh:mm:ss**.

* `type` - Indicates the deployment type.
