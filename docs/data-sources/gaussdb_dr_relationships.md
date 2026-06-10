---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_relationships"
description: |-
  Use this data source to get the list of GaussDB disaster recovery relationships within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_relationships

Use this data source to get the list of GaussDB disaster recovery relationships within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_dr_relationships" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_name` - (Optional, String) Specifies the instance name to filter the local instance.

* `instance_id` - (Optional, String) Specifies the instance ID to filter the local instance ID.

* `dr_role` - (Optional, String) Specifies the disaster recovery role.
  The valid values are as follows:
  + **master**: Primary instance.
  + **disaster**: Disaster recovery instance.

* `dr_type` - (Optional, String) Specifies the disaster recovery type.
  The valid values are as follows:
  + **stream**: Stream disaster recovery.

* `dr_status` - (Optional, String) Specifies the disaster recovery status.
  The valid values are as follows:
  + **normal**: The disaster recovery relationship is normal.
  + **failover**: The disaster recovery has been promoted to primary.
  + **pending**: Task in progress.
  + **pre_check_failed**: Disaster recovery pre-check failed.
  + **pre_checking**: Disaster recovery pre-check in progress.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `relations` - The list of disaster recovery relations.
  The [relations](#gaussdb_dr_relations) structure is documented below.

<a name="gaussdb_dr_relations"></a>
The `relations` block supports:

* `disaster_type` - The disaster recovery type.

* `name` - The disaster recovery task name.

* `disaster_role` - The disaster recovery role.

* `created` - The creation time.

* `updated` - The update time.

* `id` - The disaster recovery relation ID.

* `synchronization_id` - The ID of the disaster recovery relationship.

* `status` - The disaster recovery status.

* `precheck_failed_reason` - The pre-check failed reason.

* `instance_id` - The instance ID.

* `instance_name` - The instance name.

* `instance_status` - The instance status.

* `actions` - The list of actions.

* `slave_region_instance_info` - The slave region instance information.
  The [slave_region_instance_info](#gaussdb_dr_region_instance_info) structure is documented below.

* `master_region_instance_info` - The master region instance information.
  The [master_region_instance_info](#gaussdb_dr_region_instance_info) structure is documented below.

<a name="gaussdb_dr_region_instance_info"></a>
The `master_region_instance_info` and `slave_region_instance_info` block supports:

* `region_code` - The region code.

* `instance_id` - The instance ID.

* `project_id` - The project ID.

* `project_name` - The project name.

* `ip_address` - The list of data IP addresses, separated by commas.
