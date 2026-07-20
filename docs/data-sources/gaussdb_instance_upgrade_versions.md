---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_upgrade_versions"
description: |-
  Use this data source to query the available upgrade types and versions that GaussDB instances can be upgraded to in batches within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_upgrade_versions

Use this data source to query the available upgrade types and versions that GaussDB instances can be
upgraded to in batches within HuaweiCloud.

## Example Usage

```hcl
variable "instance_ids" {
  type = list(string)
}

data "huaweicloud_gaussdb_instance_upgrade_versions" "test" {
  instance_ids = var.instance_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the upgrade versions.
  If omitted, the provider-level region will be used.

* `instance_ids` - (Optional, List) Specifies the GaussDB instance IDs for batch query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `upgrade_type_list` - The upgrade types.
  The [upgrade_type_list](#upgrade_type_list_object) structure is documented below.

* `target_version` - The target version. The target version is returned when the instance is in the
  rolling upgrade phase, or no information is returned.

* `upgrade_candidate_versions` - The versions that can be upgraded, including major and minor versions.

* `hotfix_upgrade_infos` - The information about hot patch versions that can be installed.
  The [hotfix_infos](#hotfix_infos_object) structure is documented below.

* `hotfix_rollback_infos` - The information about hot patch versions that can be rolled back.
  The [hotfix_infos](#hotfix_infos_object) structure is documented below.

<a name="upgrade_type_list_object"></a>
The `upgrade_type_list` block supports:

* `upgrade_type` - The upgrade type.
  The valid values are as follows:
  + **grey**: Gray upgrade.
  + **inplace**: In-place upgrade.
  + **hotfix**: Hot patch installation.

* `enable` - Whether the upgrade type is available.

* `upgrade_action_list` - The upgrade actions.
  The [upgrade_action_list](#upgrade_action_list_object) structure is documented below.

* `is_parallel_upgrade` - Whether intra-AZ parallel upgrade is supported.
  + **true**: The current instance is in the rolling upgrade phase of the gray upgrade. The intra-AZ
     parallel upgrade is supported. Once this parameter is configured, it cannot be changed later.
  + **false**: The current instance is being upgraded. The intra-AZ parallel upgrade is not supported.
       Once this parameter is configured, it cannot be changed later.
  + **null**: The current instance is not in the upgrade process.

<a name="upgrade_action_list_object"></a>
The `upgrade_action_list` block supports:

* `upgrade_action` - The upgrade action.
  The valid values are as follows:
  + **upgrade**: Upgrade.
  + **upgradeAutoCommit**: Auto-commit.
  + **commit**: Commit.
  + **rollback**: Rollback.

* `enable` - Whether the upgrade action is available.

<a name="hotfix_infos_object"></a>
The `hotfix_upgrade_infos` and `hotfix_rollback_infos` blocks support:

* `version` - The hot patch version.

* `common_patch` - The patch type.
  The valid values are as follows:
  + **common**: Common patch.
  + **certain**: Custom patch.

* `backup_sensitive` - Whether the patch is related to backups.
  + **true**: The patch is related to backups.
  + **false**: The patch is not related to backups.

* `descripition` - The description of the patch.

* `default_upgrade` - The default upgrade.
  + **true**: The patch is default upgrade.
  + **false**: The patch is not default upgrade.
