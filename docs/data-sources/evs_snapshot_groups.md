---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshot_groups"
description: |-
  Use this data source to query the list of EVS snapshot groups within HuaweiCloud.
---

# huaweicloud_evs_snapshot_groups

Use this data source to query the list of EVS snapshot groups within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evs_snapshot_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `id` - (Optional, String) Specifies the snapshot group ID.

* `name` - (Optional, String) Specifies the snapshot group name.

* `status` - (Optional, String) Specifies the snapshot group status.

* `tag_key` - (Optional, String) Specifies the tag name used to filter results.

* `tags` - (Optional, String) Specifies the key/value pairs used to filter results. The value is in the following
  format: **[{"key":"key1","value":"value1"}]**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID for filtering.

* `sort_key` - (Optional, String) Specifies the keyword based on which the returned results are sorted.
  The value can be **id**, **status**, or **created_at**, and the default value is **created_at**.

* `sort_dir` - (Optional, String) Specifies the result sorting order. The default value is **desc**.
    + **desc**: The descending order.
    + **asc**: The ascending order.

* `server_id` - (Optional, String) Specifies the server ID to which the snapshot group are attached.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `snapshot_groups` - The snapshot group list.
  The [snapshot_groups](#snapshot_groups_structure) structure is documented below.

<a name="snapshot_groups_structure"></a>
The `snapshot_groups` block supports:

* `id` - The snapshot group ID.

* `created_at` - The time when the snapshot group was created.

* `status` - The snapshot group status.

* `updated_at` - The time when the snapshot group was updated.

* `name` - The snapshot group name.

* `description` - The snapshot group description.

* `enterprise_project_id` - The ID of the enterprise project associated with the snapshot.

* `tags` - The tags of the snapshot group.

* `server_id` - The server ID to which the snapshot group are attached.
