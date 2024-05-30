---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_quotas"
description: |-
  Use this data source to query the application quotas within HuaweiCloud.
---

# huaweicloud_apig_application_quotas

Use this data source to query the application quotas within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "quota_id" {}

data "huaweicloud_apig_application_quotas" "test" {
  instance_id = var.instance_id
  quota_id    = var.quota_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the application quotas belong.

* `name` - (Optional, String) Specifies the name of the application quota to be queried.

* `quota_id` - (Optional, String) Specifies the ID of the application quota.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - All application quotas that match the filter parameters.

  The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `id` - The ID of the application quota.

* `name` - The name of the application quota.

* `description` - The description of the application quota.

* `call_limits` - The maximum number of times a application quota can be called.

* `time_unit` - The time unit.

* `time_interval` - The time limit of a quota.

* `bound_app_num` - The number of applications bound to the quota policy.

* `created_at` - The creation time of the application quota, in RFC3339 format.
