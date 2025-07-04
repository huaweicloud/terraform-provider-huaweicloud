---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_vaults_by_tags"
description: |-
  Use this data source to get CBR vaults filtered by tags within HuaweiCloud.
---

# huaweicloud_cbr_vaults_by_tags

Use this data source to get CBR vaults filtered by tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbr_vaults_by_tags" "test" {
  action = "filter"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the action name. Possible values are **count** and **filter**.
  + **count**: querying count of data filtered by tags.
  + **filter**: querying details of data filtered by tags.

* `enterprise_project_id` - (Optional, String) Specifies the ID of enterprise project.

* `without_any_tag` - (Optional, Bool) Specifies whether ignore tags params.
  If this parameter is set to **true**, all resources without tags are queried.
  In this case, the `tag`, `not_tags`, `tags_any`, and `not_tags_any` fields are ignored.

* `tags` - (Optional, List) Specifies the list of included tags. Backups with these tags will be filtered.
  The [tags](#tags_struct) structure is documented below.

* `tags_any` - (Optional, List) Specifies the list of tags. Backups with any tags in this list will be filtered.
  The [tags_any](#tags_struct) structure is documented below.

* `not_tags` - (Optional, List) Specifies the list of excluded tags. Backups without these tags will be filtered.
  The [not_tags](#tags_struct) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the list of tags. Backups without any tags in this list will be filtered.
  The [not_tags_any](#tags_struct) structure is documented below.

-> For arguments above, include `tags`, `tags_any`, `not_tags`, `not_tags_any` have limits as follows:
  <br/>1. This list cannot be an empty list.
  <br/>2. The list can contain up to `10` keys.
  <br/>3. Keys in this list must be unique.
  <br/>4. If no tag filtering condition is specified, full data is returned.

* `sys_tags` - (Optional, List) Specifies the system tags.
  The [sys_tags](#sys_tags_struct) structure is documented below.

  -> The sys_tags has limits as follows:
    <br/>1. Only users with the op_service permission can obtain this field.
    <br/>2. Field `sys_tags` and tag filter conditions (`tags`, `tags_any`, `not_tags`, `not_tags_any`)
    cannot  be used at the same time.
    <br/>3. If no `sys_tags` exists, use other tag APIs for filtering. If no tag filter is specified, full data is returned.
    <br/>4. This list cannot be an empty list.

* `matches` - (Optional, List) Specifies the matches supported by resources. Keys in this list must be unique.
  Only one key is supported currently. Multiple-key support will be available later.
  The [matches](#matches_struct) structure is documented below.

* `cloud_type` - (Optional, String) Specifies cloud type of the instances. Possible values are:
  + **public**: public cloud.
  + **hybrid**: hybrid cloud.

* `object_type` - (Optional, String) Specifies resource type of the instances. Possible values are:
  + **server**: elastic cloud server.
  + **disk**: elastic volume server.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the resource tag. It contains a maximum of `127` unicode characters.
  A tag key cannot be an empty string. Spaces before and after a key will be deprecated.

* `values` - (Required, List) Specifies the list of values corresponding to the key.

  -> The field has the following restrictions:
    <br/>1. The list can contain up to `10` values.
    <br/>2. A tag value contains up to `255` unicode characters. Spaces before and after a key will be deprecated.
    <br/>3. Values in this list must be unique.
    <br/>4. Values in this list are in an OR relationship.
    <br/>5. This list can be empty and each value can be an empty character string.
    <br/>6. If this list is left blank, it indicates that all values are included.
    <br/>7. The asterisk (*) is a reserved character in the system.
    If the value starts with (*), it indicates that fuzzy match is performed based on the value following (*).
    The value cannot contain only asterisks.

<a name="sys_tags_struct"></a>
The `sys_tags` block supports:

* `key` - (Required, String) Specifies the key of the system tag,
  which is obtained from the whitelist and cannot be defined  randomly.
  Currently, only the **_sys_enterprise_project_id** field is supported,
  and the corresponding value indicates the enterprise project ID.

* `values` - (Required, List) Specifies the list of values. Currently, only the enterprise project ID is used.
  The default enterprise project ID is `0`.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key of the resource tag.
  A key can only be set to **resource_name**, indicating the resource name.

* `value` - (Required, String) Specifies the value of the resource tag.
  A value consists of up to `255` characters.
  If key is **resource_name**, an empty string indicates exact match and any non-empty string indicates fuzzy match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of matched resources.

* `resources` - List of matched resources.
  The [resources](#cbr_vaults_resources) structure is documented below.

<a name="cbr_vaults_resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_detail` - The detail of the matched resources.
  The [resource_detail](#cbr_vaults_resources_detail) structure is documented below.

* `tags` - The tag list.
  The [tags](#cbr_vaults_tags) structure is documented below.

* `sys_tags` - The system tags. Only users with the op_service permission can obtain this field.
  Currently, only resource_tag is included. key: **_sys_enterprise_project_id** value: enterprise project ID.
  Value `0` indicates the default enterprise project. This field is not returned in the non-op_service scenario.
  The [sys_tags](#cbr_vaults_sys_tags) structure is documented below.

<a name="cbr_vaults_sys_tags"></a>
The `sys_tags` block supports:

* `key` - Only supports **_sys_enterprise_project_id**, and correspond value indicates the enterprise project ID.

* `value` - Only the enterprise project ID is used. The default enterprise project ID is `0`.

<a name="cbr_vaults_resources_detail"></a>
The `resource_detail` block supports:

* `vault` - The attributes of all vault.
  The [vault](#cbr_vaults_resource_detail_vault) structure is documented below.

<a name="cbr_vaults_resource_detail_vault"></a>
The `vault` block supports:

* `id` - The vault ID.

* `name` - The vault name.

* `description` - The vault description.

* `provider_id` - The ID of the vault resource type.

* `resources` - The attributes of all vault resources.
  The [resources](#cbr_vaults_resource_detail_vault_resources) structure is documented below.

* `tags` - The list of all tags for resources.
  The [tags](#cbr_vaults_tags) structure is documented below.

* `enterprise_project_id` - The enterprise project ID.

* `auto_bind` - Whether automatic association is enabled. Defaults to **false**.

* `bind_rules` - The association rule.
  The [bind_rules](#cbr_vaults_resource_detail_vault_bind_rules) structure is documented below.

* `user_id` - The user ID.

* `created_at` - The creation time.

* `auto_expand` - Whether to enable auto capacity expansion for the vault.
  Only pay-per-use vaults support auto capacity expansion.

* `smn_notify` - Whether enabled SMN message notification.

* `threshold` - The vault capacity threshold. If the vault capacity usage exceeds this threshold, an exception
  notification is sent.

* `sys_lock_source_service` - The identifier for the SMB service.

* `locked` - Whether the vault is locked. A locked vault cannot be unlocked.

* `billing` - The operation info.
  The [billing](#cbr_vaults_resource_detail_vault_billing) structure is documented below.

<a id="cbr_vaults_resource_detail_vault_resources"></a>
The `resources` block supports:

* `id` - The ID of the resource to be backed up.

* `type` - The resource type. Possible values are **OS::Nova::Server**, **OS::Cinder::Volume**,
  **OS::Ironic::BareMetalServer**, **OS::Native::Server**, **OS::Sfs::Turbo** or **OS::Workspace::DesktopV2**.

* `name` - The name of the resource to be backed up.

* `protect_status` - The protection status.

* `size` - The allocated capacity for the associated resource, in GB.

* `backup_size` - The backup size, in GB.

* `backup_count` - The number of backups.

* `extra_info` - The extra information of the resource.
  The [extra_info](#cbr_vaults_resource_detail_vault_resources_extra_info) structure is documented below.

<a id="cbr_vaults_resource_detail_vault_bind_rules"></a>
The `bind_rules` block supports:

* `tags` - The tags to filter automatically associated resources by tag.
  Minimum length: `0` characters. Maximum length: `5` characters.
  The [tags](#cbr_vaults_resource_detail_vault_bind_rules_tags) structure is documented below.

<a id="cbr_vaults_resource_detail_vault_bind_rules_tags"></a>
The `tags` block supports:

* `key` - The key of the bind rules tag.

  -> The value of tags has limits as follows:
    <br/>1. It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/.
    <br/>2. It can contain only letters, digits, periods (.), hyphens (-), and underscores (_).

* `value` - The value of the bind rules tag.

  -> The value of tags has limits as follows:
    <br/>1. It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/.
    <br/>2. It can contain only letters, digits, periods (.), hyphens (-), and underscores (_).

<a id="cbr_vaults_resource_detail_vault_billing"></a>
The `billing` block supports:

* `allocated` - The allocated capacity, in GB.

* `charging_mode` - The billing mode. Possible values are:
  + **post_paid**: pay-per-use.
  + **pre_paid**: yearly/monthly.

* `cloud_type` - The cloud type. Possible values are:
  + **public**: public cloud.
  + **hybrid**: hybrid cloud.

* `consistent_level` - The consistency level. Possible values are:
  + **crash_consistent**: crash consistent backup.
  + **app_consistent**: application consistency backup.

* `object_type` - The object type. Possible values are: **server**, **disk**, **turbo**, **workspace**,
  **vmware**, **rds** and **file**.

* `order_id` - The order ID.

* `product_id` - The product ID.

* `protect_type` - The protection type. Possible values are: **backup** and **replication**

* `size` - The capacity, in GB.

* `spec_code` - The product specification code.

  -> Possible values as follows:
    <br/> Cloud server backup vault: **vault.backup.server.normal**.
    <br/> Disk backup vault: **vault.backup.volume.normal**.
    <br/> File system backup vault: **vault.backup.turbo.normal**.

* `status` - The vault status.

* `storage_unit` - The name of the bucket for the vault.

* `used` - The used capacity, in MB.

* `frozen_scene` - The scenario when an account is frozen.

* `is_multi_az` - Whether the vault is multi-AZ.

<a name="cbr_vaults_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

  -> The key of tags has limits as follows:
    <br/>1. It can contain a maximum of `36` characters.
    <br/>2. It cannot be an empty string.
    <br/>3. Spaces before and after a key will be discarded.
    <br/>4. It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
    <br/>5. It can contain only letters, digits, hyphens (-), and underscores (_).

* `value` - The value of the resource tag.

  -> The value of tags has limits as follows:
    <br/>1. It is mandatory when a tag is added and optional when a tag is deleted.
    <br/>2. It can contain a maximum of `43` characters.
    <br/>3. It can be an empty string.
    <br/>4. Spaces before and after a value will be discarded.
    <br/>5. It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
    <br/>6. It can contain only letters, digits, hyphens (-), underscores (_), and periods (.).

<a name="cbr_vaults_resource_detail_vault_resources_extra_info"></a>
The `extra_info` block supports:

* `exclude_volumes` - The IDs of the disks that will not be backed up. This parameter is used when servers are added to
  a vault, which include all server disks. But some disks do not need to be backed up. Or in case that a server was
  previously added and some disks on this server do not need to be backed up.
