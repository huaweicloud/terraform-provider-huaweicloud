---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project_migrate_record"
description: |-
  Use this data source to get the resource move records of EPS within HuaweiCloud.
---

# huaweicloud_enterprise_project_migrate_record

Use this data source to get the resource move records of EPS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_enterprise_project_migrate_record" "test" {}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Optional, String) Specifies the resource ID.
  If omitted, the migration record of all resources will be queried.

* `start_time` - (Optional, String) Specifies the start time.
  This parameter needs to be used in conjunction with `end_time`.
  If omitted, the migration record for the most recent week will be queried.
  The time format is **YYYY-MM-DD hh:mm:ss**.

* `end_time` - (Optional, String) Specifies the end time.
  This parameter needs to be used in conjunction with `start_time`.
  If omitted, the migration record for the most recent week will be queried.
  The time format is **YYYY-MM-DD hh:mm:ss**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of the resource move records.
  The [migrate-record](#migrate_record) structure is documented below.

<a name="migrate_record"></a>
The `migrate-record` block supports:

* `associated` - Whether associated resources are moved.

* `code` - The response code.

* `message` - The response information.

* `project_id` - The project ID.

* `project_name` - The project name.

* `region_id` - The region ID.

* `event_time` - The event time.

* `user_name` - The user name.

* `operate_type` - The move type.

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_type` - The resource type.

* `initiate_ep_id` - The ID of the enterprise project that initiates a resource move.

* `initiate_ep_name` - The name of the enterprise project that initiates a resource move.

* `origin_ep_id` - The ID of the source enterprise project.

* `origin_ep_name` - The name of the source enterprise project.

* `target_ep_id` - The ID of the destination enterprise project.

* `target_ep_name` - The name of the destination enterprise project.

* `exist_failed` - Whether there are failed tasks.
