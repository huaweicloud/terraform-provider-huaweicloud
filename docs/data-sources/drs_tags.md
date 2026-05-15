---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_tags"
description: |-
  Use this data source to get the tags of specified DRS resource type within HuaweiCloud.
---

# huaweicloud_drs_tags

Use this data source to get the tags of specified DRS resource type within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_drs_tags" "test" { 
  resource_type = "migration"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  The valid values are as follows:
  + **migration**: Real-time migration.
  + **sync**: Real-time synchronization.
  + **cloudDataGuard**: Real-time disaster recovery.
  + **subscription**: Data subscription.
  + **backupMigration**: Backup migration.
  + **replay**: Playback.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.

The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `values` - The list of tag values.
