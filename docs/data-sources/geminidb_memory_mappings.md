---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_memory_mappings"
description: |-
  Use this data source to get the list of memory acceleration mappings.
---

# huaweicloud_geminidb_memory_mappings

Use this data source to get the list of memory acceleration mappings.

## Example Usage

### Query all mappings

```hcl
data "huaweicloud_geminidb_memory_mappings" "test" {}
```

### Query mapping by mapping ID

```hcl
variable "mapping_id" {}

data "huaweicloud_geminidb_memory_mappings" "test" {
  mapping_id = var.mapping_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `mapping_id` - (Optional, String) Specifies the memory acceleration mapping ID.

* `name` - (Optional, String) Specifies the memory acceleration mapping name.
  If the name starts with `*`, fuzzy search results are returned. Otherwise, an exact result is returned.

* `source_instance_id` - (Optional, String) Specifies the source instance ID.

* `source_instance_name` - (Optional, String) Specifies the source instance name.
  If the name starts with `*`, fuzzy search results are returned. Otherwise, an exact result is returned.

* `target_instance_id` - (Optional, String) Specifies the target instance ID.

* `target_instance_name` - (Optional, String) Specifies the target instance name.
  If the name starts with `*`, fuzzy search results are returned. Otherwise, an exact result is returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dbcache_mappings` - The list of memory acceleration mappings.
  The [dbcache_mappings](#dbcache_mappings_struct) structure is documented below.

<a name="dbcache_mappings_struct"></a>
The `dbcache_mappings` block supports:

* `id` - The memory acceleration mapping ID.

* `name` - The memory acceleration mapping name.

* `source_instance_id` - The source instance ID.

* `source_instance_name` - The source instance name.

* `target_instance_id` - The target instance ID.

* `target_instance_name` - The target instance name.

* `status` - The memory acceleration mapping status.
  + **normal**: A memory mapping is normal.
  + **creating**: A memory mapping is being created.
  + **createfail**: A memory mapping failed to be created.
  + **deleting**: A memory mapping is being deleted.
  + **stopped**: A memory mapping is stopped.
  + **deleted**: A memory mapping is deleted.

* `created` - The creation time.

* `updated` - The update time.

* `rule_count` - The number of rules in the memory mapping.
