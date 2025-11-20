---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_image_signature_policy_execution_record_tasks"
description: |-
  Use this data source to get the list of SWR enterprise image signature policy execution record tasks.
---

# huaweicloud_swr_enterprise_image_signature_policy_execution_record_tasks

Use this data source to get the list of SWR enterprise image signature policy execution record tasks.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "policy_id" {}
variable "execution_id" {}

data "huaweicloud_swr_enterprise_image_signature_policy_execution_record_tasks" "test" {
  instance_id    = var.instance_id
  namespace_name = var.namespace_name
  policy_id      = var.policy_id
  execution_id   = var.execution_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String) Specifies the namespace name.

* `policy_id` - (Required, String) Specifies the policy ID.

* `execution_id` - (Required, String) Specifies the execution record ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the execution records.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - Indicates the execution task ID.

* `execution_id` - Indicates the execution record ID.

* `job_id` - Indicates the job ID.

* `namespace` - Indicates the namespace.

* `repository` - Indicates the repository name.

* `status` - Indicates the status.

* `status_text` - Indicates the status detail.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.
