---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_instances_by_tags"
description: |-
  Use this data source to get the list of DRS resource instances by tags.
---

# huaweicloud_drs_instances_by_tags

Use this data source to get the list of DRS resource instances by tags.

## Example Usage

```hcl
data "huaweicloud_drs_instances_by_tags" "test" {
  resource_type = "sync"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type of the DRS job.
  The valid values are:
  + **migration**: Online migration.
  + **sync**: Data synchronization.
  + **cloudDataGuard**: Disaster recovery.
  + **subscription**: Data subscription.
  + **backupMigration**: Backup migration.
  + **replay**: Replay.

* `without_any_tag` - (Optional, Bool) Specifies whether to query resources without any tags. If set to **true**,
  the tags parameter will be ignored, and only resources without tags will be returned.

* `tags` - (Optional, List) Specifies the list of tags to filter resources.

  The [tags](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the search criteria to filter resources.

  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.

* `values` - (Required, List) Specifies the list of tag values.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the field to match.

* `value` - (Required, String) Specifies the value to match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of DRS resource instances.
  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_detail` - The resource details.

* `tags` - The list of tags associated with the resource.
  The [tags](#tags_response_struct) structure is documented below.

<a name="tags_response_struct"></a>  
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
