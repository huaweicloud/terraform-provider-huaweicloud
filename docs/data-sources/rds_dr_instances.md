---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_dr_instances"
description: |-
  Use this data source to get the list of RDS DR instances.
---

# huaweicloud_rds_dr_instances

Use this data source to get the list of RDS DR instances.

## Example Usage

```hcl
data "huaweicloud_rds_dr_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_dr_relations` - Indicates the DR instance list.

  The [instance_dr_relations](#instance_dr_relations_struct) structure is documented below.

<a name="instance_dr_relations_struct"></a>
The `instance_dr_relations` block supports:

* `master_instance` - Indicates the master instance information.

  The [master_instance](#instance_dr_relations_master_instance_struct) structure is documented below.

* `slave_instances` - Indicates the DR instance information list.

  The [slave_instances](#instance_dr_relations_slave_instances_struct) structure is documented below.

* `instance_id` - Indicates the current region instance ID.

<a name="instance_dr_relations_master_instance_struct"></a>
The `master_instance` block supports:

* `project_name` - Indicates the project name.

* `instance_id` - Indicates the instance ID.

* `region` - Indicates the region.

* `project_id` - Indicates the project ID.

<a name="instance_dr_relations_slave_instances_struct"></a>
The `slave_instances` block supports:

* `region` - Indicates the region.

* `project_id` - Indicates the project ID.

* `project_name` - Indicates the project name.

* `instance_id` - Indicates the instance ID.
