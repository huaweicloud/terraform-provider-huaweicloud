---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_migrate_status"
description: |-
  Use this data source to get the migration status of CBR resources within HuaweiCloud.
---

# huaweicloud_cbr_migrate_status

Use this data source to get the migration status of CBR resources within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbr_migrate_status" "test" {
  all_regions = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the datasource.
  If omitted, the provider-level region will be used.

* `all_regions` - (Optional, Bool) Specifies whether to query the migration results in other regions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The overall migration status. Possible values are:
  + **migrating**: Migration processing.
  + **success**: Migration completed successfully.
  + **failed**: Migration failed.

* `project_status` - List of project migration status details.
  The [project_status](#cbr_migrate_status_project_status) structure is documented below.

<a name="cbr_migrate_status_project_status"></a>
The `project_status` block supports:

* `status` - The migration status of the project. Possible values are:
  + **migrating**: Migration processing.
  + **success**: Migration completed successfully.
  + **failed**: Migration failed.

* `project_id` - The project ID.

* `project_name` - The project name.

* `region_id` - The region ID.

* `progress` - The migration progress percentage. The value ranges from `0` to `100`.

* `fail_code` - The failure code when migration fails. This field is only present when status is **failed**.

* `fail_reason` - The failure reason when migration fails. This field is only present when status is **failed**.
