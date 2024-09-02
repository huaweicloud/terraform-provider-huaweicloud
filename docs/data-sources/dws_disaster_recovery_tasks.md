---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_disaster_recovery_tasks"
description: |-
  Use this data source to get the list of DR tasks.
---

# huaweicloud_dws_disaster_recovery_tasks

Use this data source to get the list of DR tasks.

## Example Usage

```hcl
variable "dr_name" {}

data "huaweicloud_dws_disaster_recovery_tasks" "test" {
  name = var.dr_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `dr_type` - (Optional, String) Specifies the type of the DR task. Only support **az** now.

* `name` - (Optional, String) Specifies the name of the DR task.

* `primary_cluster_name` - (Optional, String) Specifies the name of the primary cluster.

* `primary_cluster_region` - (Optional, String) Specifies the region of the primary cluster.

* `standby_cluster_name` - (Optional, String) Specifies the name of the standby cluster.

* `standby_cluster_region` - (Optional, String) Specifies the region of the standby cluster.

* `status` - (Optional, String) Specifies the status of the DR task. The valid values are:
  + **unstart**
  + **running**
  + **stopped**
  + **start_failed**
  + **stop_failed**
  + **abnormal**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of DR tasks.
  The [tasks](#attrblock_disaster_recovery) structure is documented below.

<a name="attrblock_disaster_recovery"></a>
The `tasks` block supports:

* `id` - The DR task ID.

* `name` - The name of the DR task.

* `dr_type` - The type of the DR task.

* `primary_cluster_id` - The primary cluster ID.

* `primary_cluster_name` - The name of the primary cluster.

* `primary_cluster_project_id` - The project ID of the primary cluster.

* `primary_cluster_region` - The region of the primary cluster.

* `primary_cluster_role` - The role of the primary cluster.

* `primary_cluster_status` - The status of the primary cluster.

* `standby_cluster_id` - The standby cluster ID.

* `standby_cluster_name` - The name of the standby cluster.

* `standby_cluster_project_id` - The project ID of the standby cluster.

* `standby_cluster_region` - The region of the standby cluster.

* `standby_cluster_role` - The role of the standby cluster.

* `standby_cluster_status` - The status of the standby cluster.

* `start_at` - The start time of the DR task, in UTC format.

* `status` - The status of the DR task.

* `last_disaster_time` - The lasted success synchronized time, in UTC format.

* `create_at` - The creation time of the DR task, in UTC format.
