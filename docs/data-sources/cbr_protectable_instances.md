---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_protectable_instances"
description: |-
  Use this data source to query the CBR protectable instances within HuaweiCloud.
---

# huaweicloud_cbr_protectable_instances

Use this data source to query the CBR protectable instances within HuaweiCloud.

## Example Usage

```hcl
variable "protectable_type" {}

data "huaweicloud_cbr_protectable_instances" "test" {
  protectable_type = var.protectable_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `protectable_type` - (Required, String) Specifies the object type. Valid values are:
  + **server**: Cloud servers.
  + **disk**: Cloud disks.
  + **turbo**: SFS Turbo file systems.
  + **workspace**: Workspace desktop.
  + **workspace_v2**: Workspace V2 desktops.

* `resource_id` - (Optional, String) Specifies the resource ID.

* `name` - (Optional, String) Specifies the resource name.

* `server_id` - (Optional, String) Specifies the server ID. This parameter is mandatory only for users who have enabled
  enterprise multi-project.

* `status` - (Optional, String) Specifies the resource status. Valid values are **available** and **error**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For enterprise users, if omitted, all enterprise project resources will be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The protectable instances.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `children` - The children resources. This field is a string in JSON format.

* `detail` - The resource detail. This field is a string in JSON format.

* `id` - The resource ID.

* `name` - The resource name.

* `protectable` - The backup information.

  The [protectable](#instances_protectable_struct) structure is documented below.

* `size` - The size of the resource, in GB.

* `status` - The resource status. Valid values are **active**, **deleted**, and **error**.

* `type` - The type of the resource to be backed up. Valid values are **OS::Nova::Server**, **OS::Cinder::Volume**,
  **OS::Ironic::BareMetalServer**, **OS::Native::Server**, **OS::Sfs::Turbo**, and **OS::Workspace::DesktopV2**.

<a name="instances_protectable_struct"></a>
The `protectable` block supports:

* `code` - The error code for unsupported backup.

* `reason` - The reason why backup is not supported.

* `result` - Whether backup is supported.

* `vault` - The associated vault.

  The [vault](#protectable_vault_struct) structure is documented below.

* `message` - The reason why the resource cannot be backed up. This field is returned only if the resource protectability
  check fails.

<a name="protectable_vault_struct"></a>
The `vault` block supports:

* `billing` - The operation information.

  The [billing](#vault_billing_struct) structure is documented below.

* `description` - The user-defined vault description.

* `id` - The vault ID.

* `name` - The vault name.

* `project_id` - The project ID.

* `provider_id` - The ID of the vault resource type.

* `resources` - The resources.

  The [resources](#vault_resources_struct) structure is documented below.

* `tags` - The vault tags.

  The [tags](#vault_tags_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID. Its default value is **0**.

* `auto_bind` - Whether automatic association is enabled. Its default value is **false** (not enabled).

* `bind_rules` - The association rules.

  The [bind_rules](#vault_bind_rules_struct) structure is documented below.

* `user_id` - The user ID.

* `created_at` - The creation time. For example, **2020-02-05T10:38:34.209782**.

* `auto_expand` - Whether to enable auto capacity expansion for the vault. Only pay-per-use vaults support this function.

* `smn_notify` - The SMN notification switch for the vault.

* `threshold` - The vault capacity threshold. If the vault capacity usage exceeds this threshold, an exception
  notification is sent.

* `sys_lock_source_service` - The identity of the SMB service.

* `locked` - Whether the vault is locked. A locked vault cannot be unlocked.

* `updated_at` - The latest update time. For example, **2020-02-05T10:38:34.209782**.

* `version` - The version.

<a name="vault_billing_struct"></a>
The `billing` block supports:

* `allocated` - The allocated capacity, in GB.

* `charging_mode` - The charging mode, which can be **post_paid** (pay-per-use) or **pre_paid** (yearly/monthly).
  The default value is **post_paid**.

* `cloud_type` - The cloud type. Which can be **public** or **hybrid**.

* `consistent_level` - The backup specifications. Which can be **crash_consistent** (crash consistent backup) or
  **app_consistent** (application consistency backup).

* `object_type` - The object type, which can be **server**, **disk**, **turbo**, **workspace**, **vmware**, **rds**,
  or **file**.

* `order_id` - The order ID.

* `product_id` - The product ID.

* `protect_type` - The protection type, which can be **backup** or **replication**.

* `size` - The capacity, in GB.

* `spec_code` - The specification code, which can be **vault.backup.server.normal** (server vault),
  **vault.backup.volume.normal** (disk vault), or **vault.backup.turbo.normal** (file system vault).

* `status` - The vault status. Valid values are **available**, **lock**, **frozen**, **deleting**, and **error**.

* `storage_unit` - The name of the bucket for the vault.

* `used` - The used capacity, in MB.

* `frozen_scene` - The scenario when an account is frozen.

* `is_multi_az` - The multi-AZ attribute of a vault.

<a name="vault_resources_struct"></a>
The `resources` block supports:

* `extra_info` - The additional information of the resource.

  The [extra_info](#resources_extra_info_struct) structure is documented below.

* `id` - The ID of the resource to be backed up.

* `name` - The name of the resource to be backed up.

* `protect_status` - The protection status.

* `size` - The allocated capacity for the associated resources, in GB.

* `type` - The type of the resource to be backed up, which can be **OS::Nova::Server**, **OS::Cinder::Volume**,
  **OS::Ironic::BareMetalServer**, **OS::Native::Server**, **OS::Sfs::Turbo**, or **OS::Workspace::DesktopV2**.

* `backup_size` - The backup size.

* `backup_count` - The number of backups.

<a name="resources_extra_info_struct"></a>
The `extra_info` block supports:

* `exclude_volumes` - The ID of the disk that will not be backed up. This parameter is used when servers are added to a
  vault, which include all server disks. But some disks do not need to be backed up. Or in case that a server was
  previously added and some disks on this server do not need to be backed up.

<a name="vault_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.

<a name="vault_bind_rules_struct"></a>
The `bind_rules` block supports:

* `tags` - The tags using to filter automatically associated resources.

  The [tags](#bind_rules_tags_struct) structure is documented below.

<a name="bind_rules_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
