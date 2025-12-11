---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protected_instances"
description: |-
  Use this data source to query SDRS protected instances within HuaweiCloud.
---

# huaweicloud_sdrs_protected_instances

Use this data source to query SDRS protected instances within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_sdrs_protected_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_group_id` - (Optional, String) Specifies the ID of the protection group, in which all protected instances are
  queried.
  The value of this parameter can query from datasource `huaweicloud_sdrs_protection_groups`.

* `server_group_ids` - (Optional, String) Specifies the protection group ID list. The value is in the following format:
  **[server_group_id1,server_group_id2,...,server_group_idx]**. Convert it using URL encoding.
  + All the protected instances with valid `server_group_id` in `server_group_ids` are returned.
  + The protected instances of a maximum of `30` `server_group_id` values can be queried.
  + If parameters `server_group_id` and `server_group_ids` are both specified in the request, `server_group_id` will be
    ignored.

* `protected_instance_ids` - (Optional, String) Specifies the protected instance ID list. The value is in the following
  format: **[protected_instance_id1,protected_instance_id2,...,protected_instance_idx]**. Convert it using URL encoding.
  + All the protected instances with valid `protected_instance_id` in `protected_instance_ids` are returned.
  + The protected instances of a maximum of `30` `protected_instance_id` values can be queried.
  + If parameter `server_group_id` or `server_group_ids` is specified in the request, `protected_instance_ids` will be
    ignored.

* `status` - (Optional, String) Specifies the status.
  For details, see [Protected Instance Status](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0126152931.html).

* `name` - (Optional, String) Specifies the name of a protected instance. Fuzzy search is supported.

* `query_type` - (Optional, String) Specifies the query type. Valid values are **status_abnormal** and **general**.

* `availability_zone` - (Optional, String) Specifies the current production site AZ of the protection group containing
  the protected instance. The value of this parameter can query from datasource `huaweicloud_sdrs_domain`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `protected_instances` - The information about protected instances.

  The [protected_instances](#protected_instances_struct) structure is documented below.

<a name="protected_instances_struct"></a>
The `protected_instances` block supports:

* `tags` - The tag list.

  The [tags](#protected_instances_tags_struct) structure is documented below.

* `progress` - The synchronization progress of a protected instance. Unit: %.

* `priority_station` - The current production site AZ of the protection group containing the protected instance.
  Valid values are:
  + **source**: Indicates that the current production site AZ is the source availability zone value.
  + **target**: Indicates that the current production site AZ is the target availability zone value.

* `server_group_id` - The ID of a protection group.

* `updated_at` - The time when a protected instance was updated.
  The default format is as follows: "yyyy-MM-dd HH:mm:ss.SSS", for example, **2019-04-01 12:00:00.000**.

* `metadata` - The metadata of a protected instance.

  The [metadata](#protected_instances_metadata_struct) structure is documented below.

* `attachment` - The attached replication pairs.

  The [attachment](#protected_instances_attachment_struct) structure is documented below.

* `id` - The ID of a protected instance.

* `name` - The name of a protected instance.

* `description` - The description of a protected instance.

* `created_at` - The time when a protected instance was created.
  The default format is as follows: "yyyy-MM-dd HH:mm:ss.SSS", for example, **2019-04-01 12:00:00.000**.

* `target_server` - The DR site server ID.

* `status` - The status of a protected instance.

* `source_server` - The production site server ID.

<a name="protected_instances_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.

<a name="protected_instances_metadata_struct"></a>
The `metadata` block supports:

* `_system_frozen` - Whether the resource is frozen. Value **true** indicates that the resource is frozen.
  Empty indicates that the resource is not frozen.

<a name="protected_instances_attachment_struct"></a>
The `attachment` block supports:

* `replication` - The ID of a replication pair.

* `device` - The device name.
