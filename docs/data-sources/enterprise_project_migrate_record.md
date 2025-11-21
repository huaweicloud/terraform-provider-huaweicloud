---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project_migrate_record"
description: |-
  Use this data source to get the resource move record of EPS within HuaweiCloud.
---

# huaweicloud_enterprise_project_migrate_record

Use this data source to get the resource move record of EPS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_enterprise_project_migrate_record" "test" {
  start_time  = "2020-10-10 01:00:00"
  end_time    = "2025-11-19 15:30:00"
  resource_id = "3dde353d-0117-4e1a-a09c-4750f61a3c5d"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Optional, String) Specifies the resource id.  
  If omitted, the migration record of all resources will be queried.

* `start_time` - (Optional, String) Specifies the start time.  
  This parameter needs to be used in conjunction with end_time.  
  If omitted, the migration record for the most recent week will be queried.  
  The time format is **YYYY-MM-DD hh:mm:ss**.

* `end_time` - (Optional, String) Specifies the end time.  
  This parameter needs to be used in conjunction with start_time.  
  If omitted, the migration record for the most recent week will be queried.  
  The time format is **YYYY-MM-DD hh:mm:ss**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of the resource move record.
  The [migrate-record](#migrate_record) structure is documented below.

<a name="migrate_record"></a>
The `migrate-record` block supports:

* `associated` - Whether associated resources are moved.

* `code` - Response code.

* `message` - Response information.

* `project_id` - Project ID.

* `project_name` - Project name.

* `region_id` - Region ID.

* `event_time` - Event time.

* `user_name` - User name.

* `operate_type` - Move type.

* `resource_id` - Resource ID.

* `resource_name` - Resource name.

* `resource_type` - Resource type.

* `initiate_ep_id` - ID of the enterprise project that initiates a resource move.

* `initiate_ep_name` - Name of the enterprise project that initiates a resource move.

* `origin_ep_id` - ID of the source enterprise project.

* `origin_ep_name` - Name of the source enterprise project.

* `target_ep_id` - ID of the destination enterprise project.

* `target_ep_name` - Name of the destination enterprise project.

* `exist_failed` - Whether there are failed tasks.
