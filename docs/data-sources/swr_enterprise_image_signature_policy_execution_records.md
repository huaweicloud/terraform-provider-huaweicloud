---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_image_signature_policy_execution_records"
description: |-
  Use this data source to get the list of SWR enterprise image signature policy execution records.
---

# huaweicloud_swr_enterprise_image_signature_policy_execution_records

Use this data source to get the list of SWR enterprise image signature policy execution records.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "policy_id" {}

data "huaweicloud_swr_enterprise_image_signature_policy_execution_records" "test" {
  instance_id    = var.instance_id
  namespace_name = var.namespace_name
  policy_id      = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String) Specifies the namespace name.

* `policy_id` - (Required, String) Specifies the policy ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `executions` - Indicates the execution records.

  The [executions](#executions_struct) structure is documented below.

<a name="executions_struct"></a>
The `executions` block supports:

* `id` - Indicates the execution record ID.

* `policy_id` - Indicates the policy ID.

* `trigger` - Indicates the trigger type.

* `total` - Indicates the total count.

* `stopped` - Indicates the count that is stopped.

* `succeed` - Indicates the count that is succeed.

* `failed` - Indicates the count that is failed.

* `in_progress` - Indicates the count that is in progress.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.

* `status` - Indicates the status.

* `status_text` - Indicates the status detail.
