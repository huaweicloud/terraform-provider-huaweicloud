---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_upgrade_versions"
description: |-
  Use this data source to get the versions that a GaussDb OpenGauss instance can be upgraded to.
---

# huaweicloud_gaussdb_opengauss_upgrade_versions

Use this data source to get the versions that a GaussDb OpenGauss instance can be upgraded to.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_upgrade_versions" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `upgrade_type_list` - Indicates the list of upgrade types.

  The [upgrade_type_list](#upgrade_type_list_struct) structure is documented below.

* `rollback_enabled` - Indicates whether rollback is supported.

* `source_version` - Indicates the source instance version.

* `target_version` - Indicates the target version.
  The target version is only returned when the instance is in the rolling upgrade phase, or no information is returned.

* `roll_upgrade_progress` - Indicates the DN or AZ information during the rolling upgrade.

  The [roll_upgrade_progress](#roll_upgrade_progress_struct) structure is documented below.

* `upgrade_candidate_versions` - Indicates the versions that can be upgraded to, including minor and major versions.
  An empty array is returned during a rolling upgrade.

* `hotfix_upgrade_candidate_versions` - Indicates the hot patch versions that can be updated.
  An empty array is returned during a rolling upgrade.

* `hotfix_rollback_candidate_versions` - Indicates the hot patch versions that can be rolled back.
  An empty array is returned during a rolling upgrade.

* `hotfix_upgrade_infos` - Indicates the information about hot patch versions that can be installed.

  The [hotfix_upgrade_infos](#hotfix_upgrade_infos_struct) structure is documented below.

* `hotfix_rollback_infos` - Indicates the information about hot patch versions that can be rolled back.

  The [hotfix_rollback_infos](#hotfix_rollback_infos_struct) structure is documented below.

<a name="upgrade_type_list_struct"></a>
The `upgrade_type_list` block supports:

* `upgrade_type` - Indicates the upgrade type.
  The value can be:
  + **grey**: Gray upgrade
  + **inplace**: In-place upgrade
  + **hotfix**: Hot patch update

* `enable` - Indicates whether the upgrade type is available.

* `is_parallel_upgrade` - Indicates whether intra-AZ parallel upgrade is supported.
  The value can be:
  + **true**: The current instance is in the rolling upgrade phase of the gray upgrade. The intra-AZ parallel
  upgrade is supported. Once this parameter is configured, it cannot be changed later.
  + **false**: The current instance is being upgraded. The intra-AZ parallel upgrade is not supported. Once
  this parameter is configured, it cannot be changed later.
  + **null**: The current instance is not in the upgrade process.

* `upgrade_action_list` - Indicates the upgrade actions.

  The [upgrade_action_list](#upgrade_type_list_upgrade_action_list_struct) structure is documented below.

<a name="upgrade_type_list_upgrade_action_list_struct"></a>
The `upgrade_action_list` block supports:

* `upgrade_action` - Indicates the upgrade action.
  The value can be:
  + **upgrade**: Rolling upgrade
  + **upgradeAutoCommit**: Auto-commit
  + **commit**: Commit
  + **rollback**: Rollback

* `enable` - Indicates whether the upgrade action is available.

<a name="roll_upgrade_progress_struct"></a>
The `roll_upgrade_progress` block supports:

* `upgraded_dn_group_numbers` - Indicates the number of shards that have been upgraded.

* `total_dn_group_numbers` - Indicates the total number of shards.

* `not_fully_upgraded_az` - Indicates the AZs that have not been upgraded.
  Multiple AZs are separated by commas (,). For instances in the independent deployment, null is returned.

* `already_upgraded_az` - Indicates the AZs that have upgraded.
  Multiple AZs are separated by commas (,). For instances in the independent deployment, null is returned.

* `az_description_map` - Indicates the AZ description.

<a name="hotfix_upgrade_infos_struct"></a>
The `hotfix_upgrade_infos` block supports:

* `version` - Indicates the hot patch version.

* `common_patch` - Indicates the patch type.
  The value can be:
  + **common**: common patch
  + **certain**: custom patch

* `backup_sensitive` - Indicates whether the patch is related to backups.

* `descripition` - Indicates the description of the patch.

<a name="hotfix_rollback_infos_struct"></a>
The `hotfix_rollback_infos` block supports:

* `version` - Indicates the hot patch version.

* `common_patch` - Indicates the patch type.
  The value can be:
  + **common**: common patch
  + **certain**: custom patch

* `backup_sensitive` - Indicates whether the patch is related to backups.

* `descripition` - Indicates the description of the patch.
