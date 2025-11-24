---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ransomware_backup_vaults"
description: |-
  Use this data source to get the list of HSS ransomware backup vaults within HuaweiCloud.
---

# huaweicloud_hss_ransomware_backup_vaults

Use this data source to get the list of HSS ransomware backup vaults within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_ransomware_backup_vaults" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `vault_name` - (Optional, String) Specifies the backup vault name.

* `vault_id` - (Optional, String) Specifies the backup vault ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of backup vaults.

* `data_list` - The list of backup vaults.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `vault_id` - The vault ID.

* `vault_name` - The vault name.

* `vault_size` - The total capacity of the vault in GB.

* `vault_used` - The used capacity in MB. Refers to the capacity occupied by existing backups, for example, if one host
  is bound and there are already two backups with a capacity of `60`G, then the used capacity is `60`G.

* `vault_allocated` - The allocated capacity in GB. Refers to the size of the bound server, for example, if one host is
  bound and the host size is `40`G, then the allocated capacity is `40`G.

* `vault_charging_mode` - The vault charging mode.  
  The valid values are as follows:
  + **post_paid**: Pay-per-use.
  + **pre_paid**: Yearly/Monthly.

* `vault_status` - The vault status.  
  The valid values are as follows:
  + **available**: Available.
  + **lock**: Locked.
  + **frozen**: Frozen.
  + **deleting**: Deleting.
  + **error**: Error.

* `backup_policy_id` - The backup policy ID. If it is empty, it is in an unbound state. If it is not empty,
  determine whether the policy is enabled through the `backup_policy_enabled` field.

* `backup_policy_name` - The backup policy name.

* `backup_policy_enabled` - Whether the backup policy is enabled.

* `resources_num` - The number of bound servers.
