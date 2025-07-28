---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_migration_tasks"
description: |-
  Use this data source to get the list of the migration tasks.
---

# huaweicloud_dcs_migration_tasks

Use this data source to get the list of the migration tasks.

## Example Usage

```hcl
data "huaweicloud_dcs_migration_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the migration task name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `migration_tasks` - Indicates the migration task list.

  The [migration_tasks](#migration_tasks_struct) structure is documented below.

<a name="migration_tasks_struct"></a>
The `migration_tasks` block supports:

* `task_id` - Indicates the migration task ID.

* `task_name` - Indicates the migration task name.

* `description` - Indicates the description of a migration task.

* `status` - Indicates the migration task status.
  The value can be **SUCCESS**, **FAILED**, **MIGRATING**, or **TERMINATED**.

* `migration_type` - Indicates the mode of the migration, which can be backup file import or online migration.

* `migration_method` - Indicates the type of the migration, which can be **full migration** or **incremental migration**.

* `ecs_tenant_private_ip` - Indicates the private IP address of the migration ECS on the tenant side.

* `data_source` - Indicates the source Redis address, which is **ip:port** or a bucket name.

* `source_instance_id` - Indicates the ID of the source instance.
  If the source Redis is self-hosted, this parameter is left blank.

* `source_instance_name` - Indicates the name of the source instance.
  If the source Redis is self-hosted, this parameter is left blank.

* `source_instance_subnet_id` - Indicates the subnet ID of the source instance.
  If the source Redis is self-hosted, this parameter is left blank.

* `source_instance_spec_code` - Indicates the source instance specification code.
  If the source Redis is self-hosted, this parameter is left blank.

* `source_instance_status` - Indicates the status of the source instance.
  If the source Redis is self-hosted, this parameter is left blank.

* `target_instance_id` - Indicates the target instance ID.

* `target_instance_name` - Indicates the target instance name.

* `target_instance_status` - Indicates the status of the target instance.

* `target_instance_subnet_id` - Indicates the subnet ID of the target instance.

* `target_instance_addrs` - Indicates the target Redis address. The format is **ip:port**.

* `target_instance_spec_code` - Indicates the ID of the target instance flavor.

* `version` - Indicates the version.

* `resume_mode` - Indicates the operation mode, which can be **auto** or **manual**.

* `supported_features` - Indicates the supported features.

* `error_message` - Indicates the error information.

* `created_at` - Indicates the time when the migration task is created.

* `released_at` - Indicates the time when the migration server is released.
