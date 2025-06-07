---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_upgrading_minor_version"
description: |-
  Use this data source to get the list of log configuration (LTS) configs.
---

# huaweicloud_rds_upgrading_minor_version

Use this data source to get the list of log configuration (LTS) configs.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_upgrading_minor_version" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `is_delayed` - (Optional, String) Specifies the instance name.
  + **true**: The upgrade is delayed and performed within the maintenance window.
  + **false**: The upgrade is performed immediately. This is the default value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `job_id` - Indicates the task ID for the upgrade operation, which can be used to track the job status.
