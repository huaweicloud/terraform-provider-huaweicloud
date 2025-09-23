---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_histories"
description: |-
  Use this data source to get the list of RMS resource histories.
---

# huaweicloud_rms_resource_histories

Use this data source to get the list of RMS resource histories.

## Example Usage

```hcl
variable "resource_id" {}

data "huaweicloud_rms_resource_histories" "test" {
  resource_id = var.resource_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_id` - (Required, String) Specifies the resource ID.

* `earlier_time` - (Optional, String) Specifies the start time of the query.
  If you do not set this parameter, the action returns paginated results starting from the earliest history item.
  The time format is **YYYY-MM-DD hh:mm:ss**.

* `later_time` - (Optional, String) Specifies the end time of the query.
  If you do not set this parameter, the current time is used by default.
  The time format is **YYYY-MM-DD hh:mm:ss**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The resource history list.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `domain_id` - The user ID.

* `resource_id` - The resource ID.

* `resource_type` - The resource type.

* `capture_time` - The time when the resource is captured by Config.

* `status` - The resource status.

* `relations` - The list of resource relationships.

  The [relations](#items_relations_struct) structure is documented below.

* `resource` - The resource information.

  The [resource](#items_resource_struct) structure is documented below.

<a name="items_relations_struct"></a>
The `relations` block supports:

* `relation_type` - The relationship type.

* `from_resource_type` - The type of the source resource.

* `to_resource_type` - The type of the destination resource.

* `from_resource_id` - The source resource ID.

* `to_resource_id` - The destination resource ID.

<a name="items_resource_struct"></a>
The `resource` block supports:

* `id` - The resource ID.

* `name` - The resource name.

* `ep_id` - The enterprise project ID.

* `ep_name` - The enterprise project name.

* `checksum` - The resource checksum.

* `provider` - The provider name.

* `created` - The time when the resource is created.

* `updated` - The time when the resource is updated.

* `provisioning_state` - The status of the operation that causes the resource change.

* `type` - The resource type.

* `project_id` - The project ID.

* `region_id` - The region ID.

* `project_name` - The project name.

* `tags` - The resource tag.

* `properties` - The resource properties.
