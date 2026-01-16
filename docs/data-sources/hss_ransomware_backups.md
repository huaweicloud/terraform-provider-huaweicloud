---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ransomware_backups"
description: |-
  Use this data source to get the list of HSS ransomware backups within HuaweiCloud.
---

# huaweicloud_hss_ransomware_backups

Use this data source to get the list of HSS ransomware backups within HuaweiCloud.

## Example Usage

```hcl
variable "host_id" {}

data "huaweicloud_hss_ransomware_backups" "test" {
  host_id = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the host ID to query ransomware backups.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `status` - (Optional, String) Specifies the backup status. Valid values are:
  + **available**: Available.
  + **protecting**: Protecting.
  + **deleting**: Deleting.
  + **restoring**: Restoring.
  + **error**: Error.
  + **waiting_protect**: Waiting protect.
  + **waiting_delete**: Waiting delete.
  + **waiting_restore**: Waiting restore.

* `name` - (Optional, String) Specifies the backup name.

* `last_days` - (Optional, Int) Specifies the query time range in days. Valid values are `1` - `30`.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of ransomware backup records.

* `data_list` - The list of ransomware backup data.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `backup_id` - The backup ID.

* `backup_name` - The backup name.

* `backup_status` - The backup status.

* `create_time` - The backup creation time.

* `os_images_data` - The list of backup registered image IDs.

  The [os_images_data](#os_images_data_struct) structure is documented below.

* `backup_tag` - The backup tag. Valid values are:
  + **0**: Scheduled backup.
  + **1**: Ransomware encryption backup.

<a name="os_images_data_struct"></a>
The `os_images_data` block supports:

* `image_id` - The backup registration image ID.
