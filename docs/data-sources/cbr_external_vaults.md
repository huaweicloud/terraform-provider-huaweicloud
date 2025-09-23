---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_external_vaults"
description: |-
  Use this data source to get available CBR external vaults within HuaweiCloud.
---

# huaweicloud_cbr_external_vaults

Use this data source to get available CBR external vaults within HuaweiCloud.

## Example Usage

```hcl
variable "region_id" {}
variable "external_project_id" {}

data "huaweicloud_cbr_external_vaults" "test" {
  region_id           = var.region_id
  external_project_id = var.external_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `external_project_id` - (Required, String) Specifies project ID of other region.

* `region_id` - (Required, String) Specifies the region ID.

* `cloud_type` - (Optional, String) Specifies cloud type of the instances. The value can be **public** or **hybrid**.

* `protect_type` - (Optional, String) Specifies the protection type. The value can be **backup**, **replication**,
  or **hybrid**.

* `vault_id` - (Optional, String) Specifies vault ID. If the vault ID is specified,
  other filtering criteria do not take effect.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vaults` - List of vault details. The [vaults](#cbr_external_vaults) structure is documented below.

<a name="cbr_external_vaults"></a>
The `vaults` block supports:

* `id` - Vault ID.

* `name` - Vault name.

* `description` - User-defined vault description.

* `provider_id` - ID of the vault resource type.

* `project_id` - Project ID.

* `enterprise_project_id` - Enterprise project ID. Its default value is **0**.

* `created_at` - Creation time.

* `auto_bind` - Indicates whether automatic association is enabled. Its default value is **false** (not enabled).

* `bind_rules` - Association rule. The [bind_rules](#cbr_external_vaults_bind_rules) structure is documented below.

* `user_id` - User ID.

* `auto_expand` - Whether to enable auto capacity expansion for the vault. Only pay-per-use vaults support auto
  capacity expansion.

* `smn_notify` - Exception notification function.

* `threshold` - Vault capacity threshold. If the vault capacity usage exceeds this threshold, an exception notification
  is sent.

* `sys_lock_source_service` - Used to identify the SMB service. You can set it to SMB or leave it empty.

* `locked` - Whether the vault is locked. A locked vault cannot be unlocked.

* `tags` - The list of tags.
  The [tags](#cbr_external_vaults_tags) structure is documented below.

* `resources` - The list of resources.
  The [resources](#cbr_external_vaults_resources) structure is documented below.

* `billing` - Capacity and billing info of the vault.
  The [billing](#cbr_external_vaults_billing) structure is documented below.

<a name="cbr_external_vaults_bind_rules"></a>
The `bind_rules` block supports:

* `tags` - Filters automatically associated resources by tag.
  Minimum length: `0` characters. Maximum length: `5` characters.
  The [tags](#cbr_external_vaults_bind_rules_tags) structure is documented below.

<a name="cbr_external_vaults_bind_rules_tags"></a>
The `tags` block supports:

* `key` - Key. It can contain a maximum of `36` characters.
  The key cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
  The key can contain only letters, digits, hyphens (-), and underscores (_).

* `value` - Value.
  The value cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
  The value can contain only letters, digits, periods (.), hyphens (-), and underscores (_).

<a name="cbr_external_vaults_tags"></a>
The `tags` block supports:

* `key` - Key. It can contain a maximum of `36` characters. It cannot be an empty string.
  Spaces before and after a key will be discarded. It cannot contain non-printable ASCII characters (`0`–`31`) and the
  following characters: =*<>,|/ It can contain only letters, digits, hyphens (-), and underscores (_).

* `value` - Value. It is mandatory when a tag is added and optional when a tag is deleted.
  It can contain a maximum of `43` characters. It can be an empty string.
  Spaces before and after a value will be discarded.
  It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
  It can contain only letters, digits, hyphens (-), underscores (_), and periods (.).

<a name="cbr_external_vaults_resources"></a>
The `resources` block supports:

* `id` - ID of the resource to be backed up.

* `type` - The resource type. Possible values are **OS::Nova::Server**, **OS::Cinder::Volume**,
  **OS::Ironic::BareMetalServer**, **OS::Native::Server**, **OS::Sfs::Turbo** or **OS::Workspace::DesktopV2**.

* `name` - Name of the resource to be backed up.

* `protect_status` - Protection status. Possible values are **available**, **error**, **protecting**, **restoring** or
  **removing**.

* `size` - Allocated capacity for the associated resource, in GB.

* `backup_size` - Backup size.

* `backup_count` - Number of backups.

* `extra_info` - Extra information of the resource.
  The [extra_info](#cbr_external_vaults_resources_extra_info) structure is documented below.

<a name="cbr_external_vaults_resources_extra_info"></a>
The `extra_info` block supports:

* `exclude_volumes` - IDs of the disks that will not be backed up. This parameter is used when servers are
  added to a vault, which include all server disks. But some disks do not need to be backed up.
  Or in case that a server was previously added and some disks on this server do not need to be backed up.

<a name="cbr_external_vaults_billing"></a>
The `billing` block supports:

* `allocated` - Allocated capacity, in GB.

* `used` - Used capacity, in MB.

* `size` - Capacity, in GB.

* `status` - Vault status.

* `charging_mode` - The billing mode. Possible values are:
    + **post_paid**: pay-per-use.
    + **pre_paid**: yearly/monthly.

* `cloud_type` - The cloud type. Possible values are:
    + **public**: public cloud.
    + **hybrid**: hybrid cloud.

* `consistent_level` - The vault specification. Possible values are:
    + **crash_consistent**: crash consistent backup.
    + **app_consistent**: application consistency backup.

* `protect_type` - The protection type. Possible values are: **backup** and **replication**

* `object_type` - The object type. Possible values are: **server**, **disk**, **turbo**, **workspace**,
  **vmware**, **rds** and **file**.

* `spec_code` - The product specification code.

  Possible values as follows:
  + Specification codeServer backup vault: **vault**, **backup**, **server**, **normal**.
  + Disk backup vault: **vault**, **backup**, **volume**, **normal**.
  + File system backup vault: **vault**, **backup**, **turbo**, **normal**.

* `order_id` - Order ID.

* `product_id` - Product ID.

* `storage_unit` - Name of the bucket for the vault.

* `frozen_scene` - Scenario when an account is frozen.

* `is_multi_az` - Multi-AZ attribute of a vault.
