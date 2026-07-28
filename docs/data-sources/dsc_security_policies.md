---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_security_policies"
description: |-
  Use this data source to get the list of security policies within HuaweiCloud.
---

# huaweicloud_dsc_security_policies

Use this data source to get the list of security policies within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_security_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the policy name used for filtering. The `name` uses fuzzy matching.

* `type` - (Optional, String) Specifies the policy type used for filtering.
  The valid values are as follows:
  + **GDE**: database encryption.
  + **GDE_DECRYPT**: database decryption.
  + **DOM**: database operations and maintenance.
  + **DBSS**: database security service.
  + **DDM**: database dynamic masking.
  + **SDM**: database static masking.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policy_list` - The security policy list.

  The [policy_list](#policy_list_struct) structure is documented below.

<a name="policy_list_struct"></a>
The `policy_list` block supports:

* `dbss_policy` - The database audit policy information.

  The [dbss_policy](#dbss_policy_struct) structure is documented below.

* `ddm_config` - The dynamic masking policy configuration.

  The [ddm_config](#ddm_config_struct) structure is documented below.

* `ddm_policy_list` - The dynamic masking policy list.

  The [ddm_policy_list](#ddm_policy_list_struct) structure is documented below.

* `dom_config` - The database operation and maintenance policy configuration.

  The [dom_config](#dom_config_struct) structure is documented below.

* `dom_policy` - The database operation and maintenance policy information.

  The [dom_policy](#dom_policy_struct) structure is documented below.

* `enabled` - Whether the policy is enabled.

* `gde_config` - The database encryption policy configuration.

  The [gde_config](#gde_config_struct) structure is documented below.

* `gde_policy` - The encryption policy information.

  The [gde_policy](#gde_policy_struct) structure is documented below.

* `id` - The policy ID.

* `name` - The policy name.

* `related_datasource_id` - The related datasource ID.

* `related_datasource_name` - The related datasource name.

* `related_datasource_type` - The related datasource type.

* `related_instance_id` - The related instance ID.

* `related_instance_name` - The related instance name.

* `related_instance_type` - The related instance type.

* `resource` - The datasource information.

  The [resource](#resource_struct) structure is documented below.

* `status` - The policy status.

* `type` - The policy type.

* `update_time` - The update time.

<a name="dbss_policy_struct"></a>
The `dbss_policy` block supports:

* `data_mask` - Whether to mask privacy data.

* `show_result` - Whether to show the result set.

<a name="ddm_config_struct"></a>
The `ddm_config` block supports:

* `proxy_port` - The proxy port.

* `zk_election_port` - The custom ZK election port.

* `zk_port` - The custom ZK port.

<a name="ddm_policy_list_struct"></a>
The `ddm_policy_list` block supports:

* `columns` - The column information list.

  The [columns](#columns_struct) structure is documented below.

* `namespace` - The namespace name.

* `table` - The table name.

<a name="dom_config_struct"></a>
The `dom_config` block supports:

* `deploy_mode` - The deployment mode.

<a name="dom_policy_struct"></a>
The `dom_policy` block supports:

* `custom_policy` - Whether to enable the custom policy.

* `data_audit` - Whether to perform data audit.

* `default_action` - The default policy action.

* `intelligent_protection_baseline` - Whether to enable the intelligent protection baseline.

* `virtual_patch` - Whether to enable the virtual patch.

<a name="gde_config_struct"></a>
The `gde_config` block supports:

* `enc_mode` - The encryption mode. The value `1` indicates encryption and `2` indicates decryption.

* `proxy_port` - The proxy port.

<a name="gde_policy_struct"></a>
The `gde_policy` block supports:

* `action` - The action. The value `1` indicates encryption and `2` indicates decryption.

* `alg` - The encryption algorithm.

* `columns` - The column information list.

  The [columns](#columns_struct) structure is documented below.

* `table` - The table name.

<a name="columns_struct"></a>
The `columns` block supports:

* `mask` - The masking algorithm name or ID.

* `name` - The column name.

<a name="resource_struct"></a>
The `resource` block supports:

* `account` - The database username.

* `address` - The database address.

* `address_type` - The address type.

* `case_sensitive` - Whether case-sensitive.

* `database_name` - The database name.

* `extra_params` - The extra parameters.

* `password` - The database password.

* `port` - The database port.

* `res_id` - The database ID.

* `res_type` - The database type.

* `res_version` - The database version.
