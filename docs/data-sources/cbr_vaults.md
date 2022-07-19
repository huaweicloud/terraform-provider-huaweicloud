---
subcategory: "Cloud Backup and Recovery (CBR)"
---

# huaweicloud_cbr_vaults

Use this data source to get available CBR vaults within Huaweicloud.

## Example Usage

### Get vaults for all server type

```hcl
data "huaweicloud_cbr_vaults" "test" {
  type = "server"
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the CBR vaults.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies a unique name of the CBR vault. This parameter can contain a maximum of 64
  characters, which may consist of letters, digits, underscores(_) and hyphens (-).

* `type` - (Optional, String) Specifies the object type of the CBR vault. The vaild values are as follows:
  + **server** (Cloud Servers)
  + **disk** (EVS Disks)
  + **turbo** (SFS Turbo file systems)

* `consistent_level` - (Optional, String) Specifies the backup specifications.
  The valid values are as follows:
  + **[crash_consistent](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0109.html)**
  + **[app_consistent](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0109.html)**

  Only server type vaults support application consistent.

* `protection_type` - (Optional, String) Specifies the protection type of the CBR vault.
  The valid values are **backup** and **replication**. Vaults of type **disk** don't support **replication**.

* `size` - (Optional, Int) Specifies the vault sapacity, in GB. The valid value range is `1` to `10,485,760`.

* `auto_expand_enabled` - (Optional, Bool) Specifies whether to enable automatic expansion of the backup protection
  type vault. Default to **false**.

* `enterprise_project_id` - (Optional, String) Specifies a unique ID in UUID format of enterprise project.

* `policy_id` - (Optional, String) Specifies a policy to associate with the CBR vault.
  The `policy_id` cannot be used with the vault of replicate protection type.

* `status` - (Optional, String) Specifies the CBR vault status, including **available**, **lock**, **frozen** and **error**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in hashcode format.

* `vaults` - List of CBR vault details. The object structure of each CBR vault is documented below.

The `vaults` block supports:

* `id` - The vault ID in UUID format.

* `name` - The CBR vault name.

* `type` - The object type of the CBR vault.

* `consistent_level` - The backup specifications.

* `protection_type` - The protection type of the CBR vault.

* `size` - The vault capacity, in GB.

* `auto_expand_enabled` - Whether to enable automatic expansion of the backup protection type vault.

* `enterprise_project_id` - The enterprise project ID.

* `policy_id` - The policy associated with the CBR vault.

* `allocated` - The allocated capacity of the vault, in GB.

* `used` - The used capacity, in GB.

* `spec_code` - The specification code.

* `status` - The vault status.

* `storage` - The name of the bucket for the vault.

* `tags` - The key/value pairs to associate with the CBR vault.

* `resources` - An array of one or more resources to attach to the CBR vault.
  The [object](#cbr_vault_resources) structure is documented below.

<a name="cbr_vault_resources"></a>
The `resources` block supports:

* `server_id` - The ID of the ECS instance to be backed up.

* `excludes` - An array of disk IDs which will be excluded in the backup.

* `includes` - An array of disk or SFS file system IDs which will be included in the backup.
