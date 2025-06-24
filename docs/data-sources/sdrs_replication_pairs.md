---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_replication_pairs"
description: |-
  Use this data source to query SDRS replication pairs within HuaweiCloud.
---

# huaweicloud_sdrs_replication_pairs

Use this data source to query SDRS replication pairs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_sdrs_replication_pairs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the SDRS replication pairs are located.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the AZ of the current production site of the protection group.

* `name` - (Optional, String) Specifies the name of the replication pair. Fuzzy search is supported.

* `protected_instance_id` - (Optional, String) Specifies the ID of the protected instance bound to the replication pair.

* `protected_instance_ids` - (Optional, String) Specifies the list of protected instance IDs (URL-encoded).
  The value is in the following format: **['protected_instance_id1','protected_instance_id2',...,'protected_instance_idx']**.
  Convert it using URL encoding.
  
  ->1. All the replication pairs with valid `protected_instance_id` in `protected_instance_ids` are returned.
  <br/>2. The replication pairs of a maximum of `30` protected instance ID values can be queried.
  <br/>3. If parameters `protected_instance_id` and `protected_instance_ids` are both specified in the request,
  `protected_instance_id` will be ignored.

* `query_type` - (Optional, String) Specifies the query type. Valid values are:
  + **status_abnormal**: Indicates to query replication pairs in the abnormal status.
  + **general**: Indicates to query all replication pairs.

* `server_group_id` - (Optional, String) Specifies the protection group ID.
  You can obtain this value through datasource `huaweicloud_sdrs_protection_groups`.

* `server_group_ids` - (Optional, String) Specifies the list of protection group IDs (URL-encoded).
  The value is in the following format: **['server_group_id1','server_group_id2',...,'server_group_idx']**.
  Convert it using URL encoding.

  -> 1. All the replication pairs with valid `server_group_id` in `server_group_ids` are returned.
  <br/>2. The replication pairs of a maximum of `30` server group ID values can be queried.
  <br/>3. If parameters `server_group_id` and `server_group_ids` are both specified in the request,
  `server_group_id` will be ignored.

* `status` - (Optional, String) Specifies the status of the replication pair.
  For details, see [Replication Pair Status](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0126152932.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `replication_pairs` - The information about replication pairs.

  The [replication_pairs](#replication_pairs_struct) structure is documented below.

<a name="replication_pairs_struct"></a>
The `replication_pairs` block supports:

* `id` - The ID of the replication pair.

* `name` - The name of the replication pair.

* `description` - The description of the replication pair.

* `status` - The status of the replication pair.
  For details, see [Replication Pair Status](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0126152932.html).

* `volume_ids` - The IDs of the volumes used by the replication pair.

* `attachment` - The attachment information.

  The [attachment](#attachment_struct) structure is documented below.

* `created_at` - The creation time of the replication pair.

* `updated_at` - The update time of the replication pair.

* `replication_model` - The replication model of the replication pair. The default value is **hypermetro**, indicating
  synchronous replication.

* `progress` - The synchronization progress, in percentage.

* `failure_detail` - The error code when the replication pair status is **error**.
  For details, see the returned value in [Error Codes](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0113127626.html).

* `record_metadata` - The record metadata.

  The [record_metadata](#record_metadata_struct) structure is documented below.

* `server_group_id` - The ID of the protection group.

* `fault_level` - The fault level of the replication pair. Valid values:
  + `0`: No fault occurs.
  + `2`: The disk of the current production site does not have read/write permissions. In this case, you are advised to
    perform a failover.
  + `5`: The replication link is disconnected. In this case, a failover is not allowed. Contact customer service.

* `priority_station` - The current production site AZ of the protection group containing the replication pair. Valid values:
  + **source**: Indicates that the current production site AZ is the `source_availability_zone` value.
  + **target**: Indicates that the current production site AZ is the `target_availability_zone` value.

* `replication_status` - The data synchronization status. Valid values:
  + **active**: Data has been synchronized.
  + **inactive**: Data is not synchronized.
  + **copying**: Data is being synchronized.
  + **active-stopped**: Data synchronization is stopped.

<a name="attachment_struct"></a>
The `attachment` block supports:

* `protected_instance` - The ID of the protected instance.

* `device` - The device name on the server.

<a name="record_metadata_struct"></a>
The `record_metadata` block supports:

* `multiattach` - Whether the volume supports multi-attach.

* `bootable` - Whether the volume is a system disk.

* `volume_size` - The size of the volume in GB.

* `volume_type` - The type of the volume.
