---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ransomware_protection_servers"
description: |-
  Use this data source to get the list of ransomware protection servers.
---

# huaweicloud_hss_ransomware_protection_servers

Use this data source to get the list of ransomware protection servers.

## Example Usage

```hcl
data "huaweicloud_hss_ransomware_protection_servers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_name` - (Optional, String) Specifies the server name.

* `host_id` - (Optional, String) Specifies the server ID.

* `os_type` - (Optional, String) Specifies the OS type.
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `host_ip` - (Optional, String) Specifies the server EIP.

* `private_ip` - (Optional, String) Specifies the server private IP address.

* `host_status` - (Optional, String) Specifies the handling status.
  The valid values are as follows:
  + **ACTIVE**
  + **SHUTOFF**

* `ransom_protection_status` - (Optional, String) Specifies the ransomware protection status.
  The valid values are as follows:
  + **closed**
  + **opened**
  + **opening**
  + **closing**
  + **protect_failed**
  + **protect_degraded**

* `protect_policy_name` - (Optional, String) Specifies the ransomware protection policy name.

* `policy_name` - (Optional, String) Specifies the policy name.

* `policy_id` - (Optional, String) Specifies the policy ID.

* `agent_status` - (Optional, String) Specifies the agent status.
  The valid values are as follows:
  + **installed**
  + **online**
  + **offline**
  + **install_failed**
  + **installing**
  + **not_installed**

  If you want to filter agents in all status except **online**, set this parameter to **not_online**.

* `group_id` - (Optional, String) Specifies the server group ID.

* `group_name` - (Optional, String) Specifies the server group name.

* `last_days` - (Optional, Int) Specifies the query time range.
  The valid value ranges from `1` to `30`.
  If this parameter is not specified, one day is queried by default.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of ransomware protection servers.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The server ID.

* `agent_id` - The agent ID.

* `host_name` - The server name.

* `host_ip` - The server EIP.

* `private_ip` - The server private IP address.

* `os_type` - The OS type.

* `os_name` - The OS name.

* `host_status` - The server status.

* `project_id` - The project ID.

* `enterprise_project_id` - The enterprise project ID.

* `ransom_protection_status` - The ransomware protection status.

* `ransom_protection_fail_reason` - The ransomware protection failure case.
  The valid values are as follows:
  + **driver_load_failed**: Driver loading failed.
  + **protect_interrupted**: Protection interrupted.
  + **decoy_deploy_totally_failed**: All honeypots failed to be deployed.
  + **decoy_deploy_partially_failed**: Some honeypots failed to be deployed.

* `failed_decoy_dir` - The directory where honeypot protection failed.

* `agent_version` - The agent version.

* `protect_status` - The protection status.
  The valid values are as follows:
  + **closed**
  + **opened**

* `group_id` - The server group ID.

* `group_name` - The server group name.

* `protect_policy_id` - The protection policy ID.

* `protect_policy_name` - The protection policy name.

* `backup_error` - The backup error message.

  The [backup_error](#data_list_backup_error_struct) structure is documented below.

* `backup_protection_status` - Whether to enable backup.
  The valid values are as follows:
  + **failed_to_turn_on_backup**
  + **closed**
  + **opened**

* `count_protect_event` - The number of protection events.

* `count_backuped` - The number of existing backups.

* `agent_status` - The agent status.

* `version` - The server protection version.
  The valid values are as follows:
  + **hss.version.null**
  + **hss.version.basic**
  + **hss.version.advanced**
  + **hss.version.enterprise**
  + **hss.version.premium**
  + **hss.version.wtp**
  + **hss.version.container.enterprise**

* `host_source` - The server type.
  The valid values are as follows:
  + **ecs**
  + **outside**
  + **workspace**

* `vault_id` - The vault ID.

* `vault_name` - The vault name.

* `vault_size` - The vault total capacity, in GB.

* `vault_used` - The vault used capacity, in MB.

* `vault_allocated` - The allocated bound server capacity, in GB.

* `vault_charging_mode` - The vault pay mode.
  The valid values are as follows:
  + **post_paid**: Pay-per-use.
  + **pre_paid**: Yearly/monthly.

* `vault_status` - The vault status.
  The valid values are as follows:
  + **available**
  + **lock**
  + **frozen**
  + **deleting**
  + **error**

* `backup_policy_id` - The backup policy ID.

* `backup_policy_name` - The backup policy name.

* `backup_policy_enabled` - Whether the backup policy is enabled.

* `resources_num` - The number of bound servers.

<a name="data_list_backup_error_struct"></a>
The `backup_error` block supports:

* `error_code` - The error code.
  The valid values are as follows:
  + **0**: No error information.
  + **1**: Backup cannot be enabled because another vault han been bound.
  + **2**: The number of backup vaults exceeds the upper limit.
  + **3**: An exception occurs when the CBR API is called.

* `error_description` - The error message.
