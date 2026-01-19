---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ransomware_backup_detail"
description: |-
  Use this data source to get the detail of HSS ransomware backup within HuaweiCloud.
---

# huaweicloud_hss_ransomware_backup_detail

Use this data source to get the detail of HSS ransomware backup within HuaweiCloud.

## Example Usage

```hcl
variable "backup_id" {}

data "huaweicloud_hss_ransomware_backup_detail" "test" {
  backup_id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `backup_id` - (Required, String) Specifies the ID of the backup.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, same as `backup_id`.

* `name` - The backup name.

* `image_type` - The backup type.

* `vault_id` - The repository ID where the backup is located.

* `created_at` - The creation Time.

* `status` - The status of the backup.

* `resource_size` - The resource size.

* `resource_id` - The resource ID, corresponds to the host ID.

* `resource_type` - The resource type.

* `resource_name` - The resource name, corresponds to the host name.

* `children` - The sub backup. It is a volume backup information.
  
  The [children](#children_struct) structure is documented below.

<a name="children_struct"></a>
The `children` block supports:

* `id` - The volume backup ID.

* `name` - The disk backup name.

* `image_type` - The backup type.

* `vault_id` - The repository ID where the backup is located.

* `status` - The status of the backup.

* `resource_size` - The resource size.

* `resource_id` - The resource ID, corresponds to the host ID.

* `resource_type` - The resource type.

* `resource_name` - The resource name, corresponds to the host name.
