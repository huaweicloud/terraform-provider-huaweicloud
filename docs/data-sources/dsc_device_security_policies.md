---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_device_security_policies"
description: |-
  Use this data source to get the list of device security policies within HuaweiCloud.
---

# huaweicloud_dsc_device_security_policies

Use this data source to get the list of device security policies within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_device_security_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the policy name for filtering.

* `type` - (Optional, String) Specifies the policy type for filtering.
  Valid values include **GDE** (database encryption), **GDE_DECRYPT** (database decryption),
  **DOM** (database O&M), **DBSS** (database security service),
  **DDM** (database dynamic masking), **SDM** (database static masking).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `policies` - The security policy list.
  The [policies](#dsc_policies_struct) structure is documented below.

<a name="dsc_policies_struct"></a>
The `policies` block supports:

* `id` - The policy ID.

* `name` - The policy name.

* `enabled` - Whether the policy is enabled.

* `status` - The policy status.

* `type` - The policy type.

* `update_time` - The update time.

* `related_datasource_id` - The related datasource ID.

* `related_datasource_name` - The related datasource name.

* `related_datasource_type` - The related datasource type.

* `related_device_id` - The related device ID.

* `related_device_name` - The related device name.

* `related_device_type` - The related device type.

* `target_datasource_id` - The target datasource ID.

* `target_datasource_name` - The target datasource name.

* `target_datasource_type` - The target datasource type.

* `ddm_config` - The dynamic masking configuration.
  The [ddm_config](#dsc_ddm_config_struct) structure is documented below.

* `ddm_policy_list` - The dynamic masking policy list.
  The [ddm_policy_list](#dsc_ddm_policy_struct) structure is documented below.

* `gde_config` - The encryption configuration.
  The [gde_config](#dsc_gde_config_struct) structure is documented below.

* `gde_policy` - The encryption policy.
  The [gde_policy](#dsc_gde_policy_struct) structure is documented below.

* `sdm_config` - The static masking configuration.
  The [sdm_config](#dsc_sdm_config_struct) structure is documented below.

* `sdm_policy_list` - The static masking policy list.
  The [sdm_policy_list](#dsc_sdm_policy_struct) structure is documented below.

* `resource` - The device resource information.
  The [resource](#dsc_resource_info_struct) structure is documented below.

* `target_resource` - The target device resource information.
  The [target_resource](#dsc_resource_info_struct) structure is documented below.

<a name="dsc_ddm_config_struct"></a>
The `ddm_config` block supports:

* `proxy_port` - The proxy port.

* `zk_election_port` - The ZK election port.

* `zk_port` - The ZK port.

<a name="dsc_ddm_policy_struct"></a>
The `ddm_policy_list` block supports:

* `namespace` - The namespace name.

* `table` - The table name.

* `columns` - The column information list.
  The [columns](#dsc_column_struct) structure is documented below.

<a name="dsc_gde_config_struct"></a>
The `gde_config` block supports:

* `enc_mode` - The encryption mode. 1 means encrypt, 2 means decrypt.

* `proxy_port` - The proxy port.

<a name="dsc_gde_policy_struct"></a>
The `gde_policy` block supports:

* `action` - The action. 1 means encrypt, 2 means decrypt.

* `alg` - The encryption algorithm.

* `table` - The table name.

* `columns` - The column information list.
  The [columns](#dsc_column_struct) structure is documented below.

<a name="dsc_sdm_config_struct"></a>
The `sdm_config` block supports:

* `auto_rebuild_target` - Whether to rebuild the target table.

* `clear_target` - Whether to clear the target table.

* `select_param` - The extraction parameter value.

* `select_type` - The extraction type.

* `skip_dirty_data` - Whether to skip dirty data.

<a name="dsc_sdm_policy_struct"></a>
The `sdm_policy_list` block supports:

* `table` - The table name.

* `namespace` - The namespace name.

* `do_mask` - Whether to mask data.

* `do_move` - Whether to move data.

* `columns` - The column information list.
  The [columns](#dsc_column_struct) structure is documented below.

<a name="dsc_column_struct"></a>
The `columns` block supports:

* `name` - The column name.

* `mask` - The masking algorithm name or ID.

<a name="dsc_resource_info_struct"></a>
The `resource` and `target_resource` block supports:

* `account` - The account name.

* `address` - The address.

* `address_type` - The address type.

* `case_sensitive` - Whether case sensitive.

* `database_name` - The database name.

* `extra_params` - The extra parameters.

* `password` - The password.

* `port` - The port.

* `res_id` - The resource ID.

* `res_type` - The resource type.

* `res_version` - The resource version.
