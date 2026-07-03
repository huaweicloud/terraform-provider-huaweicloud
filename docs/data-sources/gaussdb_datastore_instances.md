---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_datastore_instances"
description: |-
  Use this data source to get the list instances by engine version.
---

# huaweicloud_gaussdb_datastore_instances

Use this data source to get the list instances by engine version.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_datastore_instances" "this" {}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the data source.

* `engine_instance_details` - Indicates the instance engine details.

  The [engine_instance_details](#engine_instance_details_struct) structure is documented below.

<a name="engine_instance_details_struct"></a>
The `engine_instance_details` block supports:

* `engine_version` - Indicates the engine version.

* `instances` - Indicates the instance details.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `status` - Indicates the instance status.
  + **BUILD**: The instance is being created.
  + **BUILD_FAILED**: The instance failed to be created.
  + **ACTIVE**: The instance is normal.
  + **FAILED**: The instance is abnormal.
  + **FROZEN**: The instance is frozen.
  + **MODIFYING**: The storage is being scaled up, or instance specifications are being changed.
  + **EXPANDING**: Read replicas, CNs, or DN shards are being added to the instance.
  + **REBOOTING**: The instance is being rebooted.
  + **REDUCING**: Read replicas are being deleted.
  + **UPGRADING**: The instance is being upgraded.
  + **RESTORING**: The instance is being restored.
  + **SWITCHOVER**: A primary/standby switchover is being performed.
  + **MIGRATING**: The instance is being migrated.
  + **BACKING UP**: The instance is being backed up.
  + **UPGRADE TO BE OBSERVED**: The instance upgrade is in the observation period.
  + **MODIFYING DATABASE PORT**: The database port is being changed.
  + **STORAGE FULL**: The instance storage is full.
  + **REPAIRING**: The instance is being repaired.
  + **SHUTDOWN**: The instance is stopped.

* `type` - Indicates the Instance type.
  + **enterprise**: Distributed.
  + **centralization_standard**: Centralized.

* `solution` - Indicates the instance deployment model.
  + **Ha**: Centralized deployment.
  + **Independent**: Independent deployment.
  + **Combined**: Combined deployment.

* `hotfix_versions` - Indicates the updated hot patch version.
