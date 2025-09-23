---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protected_instances_by_tags"
description: |-  
  Use this data source to query SDRS protected instances by tags within HuaweiCloud.
---

# huaweicloud_sdrs_protected_instances_by_tags

Use this data source to query SDRS protected instances by tags within HuaweiCloud.

## Example Usage

### Query SDRS protected instances list

```hcl
data "huaweicloud_sdrs_protected_instances_by_tags" "test" {
  action = "filter"

  tags {
    key    = "foo"
    values = ["bar"]
  }
}
```

### Query SDRS protected instances count

```hcl
data "huaweicloud_sdrs_protected_instances_by_tags" "test" {
  action = "count"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the datasource.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the operation to be performed. Valid values are:
  + **filter**: Returns all protected instances that match the search criteria.
    Only attribute `resources` will be returned in this case.
  + **count**: Returns the total number of protected instances that match the search criteria.
    Only attribute `total_count` will be returned in this case.

* `tags` - (Optional, List) Specifies the tags to query resource list which contain all the specified tags.
  The [tags](#tags_params_struct) structure is documented below.

  -> Each resource to be queried contains a maximum of `10` keys. Each tag key can have a maximum of `10` tag values.
   The tag value corresponding to each tag key can be an empty array but the structure cannot be missing.
   Each tag key must be unique, and each tag value in a tag must be unique. The response returns resources containing
   all tags in this list.
   Keys in this list are in an AND relationship while values in each key-value structure are in an OR relationship.
   If no tag filtering condition is specified, full data is returned.

* `tags_any` - (Optional, List) Specifies the tags to query resource list which contain any of the specified tags.
  The [tags_any](#tags_params_struct) structure is documented below.

  -> Each resource to be queried contains a maximum of `10` keys. Each tag key can have a maximum of `10` tag values.
   The tag value corresponding to each tag key can be an empty array but the structure cannot be missing.
   Each tag key must be unique, and each tag value in a tag must be unique. The response returns resources containing
   the tags in this list. Keys in this list are in an OR relationship and values in each key-value structure are also
   in an OR relationship. If no tag filtering condition is specified, full data is returned.

* `not_tags` - (Optional, List) Specifies the tags to query resource list which do not contain all the specified tags.
  The [not_tags](#tags_params_struct) structure is documented below.

  -> Each resource to be queried contains a maximum of `10` keys. Each tag key can have a maximum of `10` tag values.
   The tag value corresponding to each tag key can be an empty array but the structure cannot be missing.
   Each tag key must be unique, and each tag value in a tag must be unique. The response returns resources containing
   no tags in this list. Keys in this list are in an AND relationship while values in each key-value structure are
   in an OR relationship. If no tag filtering condition is specified, full data is returned.

* `not_tags_any` - (Optional, List) Specifies the tags to query resource list which do not contain any of the specified tags.
  The [not_tags_any](#tags_params_struct) structure is documented below.

  -> Each resource to be queried contains a maximum of `10` keys. Each tag key can have a maximum of `10` tag values.
   The tag value corresponding to each tag key can be an empty array but the structure cannot be missing.
   Each tag key must be unique, and each tag value in a tag must be unique. The response returns resources containing
   no tags in this list. Keys in this list are in an OR relationship and values in each key-value structure are also in
   an OR relationship. If no tag filtering condition is specified, full data is returned.

* `matches` - (Optional, List) Specifies the search field.
  The [matches](#matches_params_struct) structure is documented below.

  -> The tag key is the field to be matched, for example, **resource_name**. The tag value indicates the value to be matched.
   The key is a fixed dictionary value and cannot contain duplicate keys or unsupported keys.
   Determine whether fuzzy match is required based on the keys. For example, if key is **resource_name**,
   fuzzy search (case insensitive) is used by default. If value is an empty string, exact match is used.
   Currently, only **resource_name** for key is supported. Other key values will be available later.

<a name="tags_params_struct"></a>
The `tags`, `tags_any`, `not_tags`, `not_tags_any` block supports:

* `key` - (Required, String) Specifies the tag key. It contains a maximum of `127` Unicode characters.
  It cannot be left blank. Key cannot be empty, an empty string, or spaces.
  Before using key, delete spaces of single-byte character (SBC) before and after the value.

* `values` - (Required, List) Specifies the tag values. Each value contains a maximum of `255` Unicode characters.
  Before using values, delete SBC spaces before and after the value. The asterisk (*) is reserved for the system.
  If the value starts with (*), it indicates that fuzzy match is performed based on the value following (*).
  The value cannot contain only asterisks (*). If the values are null, it indicates any_value (querying any value).
  The resources containing one or more values listed in values will be found and displayed.

<a name="matches_params_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the search field key. Currently, only **resource_name** for key is supported.
  Other key values will be available later.

* `value` - (Required, String) Specifies the search value. Each value can contain a maximum of `255` unicode characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The information about protected instances that match the search criteria.
  The [resources](#resources_struct) structure is documented below.

  -> This attribute is only available when `action` is **filter**.

* `total_count` - The total number of protected instances that match the search criteria.

  -> This attribute is only available when `action` is **count**.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The ID of the protected instance.

* `resource_name` - The name of the protected instance. This attribute is blank by default if there is no name.

* `resource_detail` - The details of a protected instance.
  The [resource_detail](#resources_resource_detail_struct) structure is documented below.

* `tags` - The tags of the protected instance.
  The [tags](#resources_tags_struct) structure is documented below.

<a name="resources_resource_detail_struct"></a>
The `resource_detail` block supports:

* `id` - The ID of the protected instance.

* `name` - The name of the protected instance.

* `description` - The description of the protected instance.

* `status` - The status of the protected instance.
  For details, see [Protected Instance Status](https://support.huaweicloud.com/intl/en-us/api-sdrs/en-us_topic_0126152931.html).

* `source_server` - The ID of the production site server.

* `target_server` - The ID of the disaster recovery site server.

* `server_group_id` - The ID of the protection group.

* `created_at` - The creation time of the protected instance.

* `updated_at` - The last update time of the protected instance.

* `metadata` - The metadata of the protected instance.
  The [metadata](#resources_metadata_struct) structure is documented below.

* `attachment` - The attachments of the protected instance.
  The [attachment](#resources_attachments_struct) structure is documented below.

* `tags` - The tags of the protected instance.
  The [tags](#resources_tags_struct) structure is documented below.

* `progress` - The synchronization progress of the protected instance. Unit: %.

* `priority_station` - The priority station of the protected instance. Valid values are:
  + **source**: Indicates that the current production site AZ is the source availability zone value.
  + **target**: Indicates that the current production site AZ is the target availability zone value.

<a name="resources_metadata_struct"></a>
The `metadata` block supports:

* `__system__frozen` - Whether the resource is frozen. Value **true** indicates that the resource is frozen.
  Empty indicates that the resource is not frozen.

<a name="resources_attachments_struct"></a>
The `attachment` block supports:

* `replication` - The ID of the replication pair.

* `device` - The device name.

<a name="resources_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
