---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_associated_quota"
description: |-
  Use this data source to query the application associated quota within HuaweiCloud.
---

# huaweicloud_apig_application_associated_quota

Use this data source to query the application associated quota within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "app_id" {}

data "huaweicloud_apig_application_associated_quota" "test" {
  instance_id = var.instance_id
  app_id      = var.app_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the application is located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the application belongs.

* `app_id` - (Required, String) Specifies the ID of the application to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of data source.

* `app_quota_id` - The ID of the application quota.

* `name` - The name of the application quota.

* `call_limits` - The maximum number of times the application quota can be called.

* `time_unit` - The time unit of the quota limit.

* `time_interval` - The time interval of the quota limit.

* `remark` - The description of the application quota.

* `reset_time` - The first quota reset time point, in RFC3339 format.

* `create_time` - The creation time of the application quota, in RFC3339 format.

* `bound_app_num` - The number of applications bound to the quota policy.
