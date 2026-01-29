---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_resources_by_tags"
description: |-
  Use this data source to get the lis of VPCEP resources filtered by tags within HuaweiCloud.
---

# huaweicloud_vpcep_resources_by_tags

Use this data source to get the lis of VPCEP resources filtered by tags within HuaweiCloud.

## Example Usage

```hcl
variable "resource_type" {}
variable "action" {}

data "huaweicloud_vpcep_resources_by_tags" "test" {
  resource_type = var.resource_type
  action        = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  The value can be **endpoint_service** or **endpoint**.

* `action` - (Required, String) Specifies the action name.
  The valid values are **filter** and **count**.
  If `action` set to **filter**, indicates query the VPCEP resources based on tags filtering conditions.
  If `action` set to **count**, indicates only query the total number of VPCEP resources.

* `tags` - (Optional, List) Specifies the tags are included.
  The [tags](#tags_struct) structure is documented below.

* `tags_any` - (Optional, List) Specifies the any tags are included.
  The [tags_any](#tags_struct) structure is documented below.

* `not_tags` - (Optional, List) Specifies the tags are excluded.
  The [not_tags](#tags_struct) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the any tags are excluded.
  The [not_tags_any](#tags_struct) structure is documented below.

-> For arguments above, include `tags`, `tags_any`, `not_tags`, `not_tags_any` have limits as follows:
  <br/>1. A maximum of `10` tag keys are included, and each tag value can have a maximum of `10` values.
  <br/>2. Each tag value can be an empty array, but the tag structure cannot be missing.
  Tag keys must be unique. Values of the same tag key must be unique.
  <br/>3. Keys are in the AND relationship (`tags`,`not_tags`).
  <br/>4. Keys are in the OR relationship (`tags_any`,`not_tags_any`).
  <br/>5. Values in the key-value structure are in the OR relationship.

* `matches` - (Optional, List) Specifies the search field.
  The [matches](#matches_struct) structure is documented below.

* `without_any_tag` - (Optional, Bool) Specifies whether ignore tags parameters.
  The value can be **true** or **false** (Default value).
  When `without_any_tag` is set to **true**, ignore parameter verification on the `tags`, `not_tags`,
  `tags_any`, and `not_tags_any`.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.
  Each tag key can contain a maximum of `128` characters but cannot be left blank.
  key cannot be an empty string or spaces.

* `values` - (Required, List) Specifies the tag values.
  Each tag value contains a maximum of `255` characters.
  The tag value can be an empty array but cannot be left blank. If values are left blank,
  it indicates querying any value.
  Values are in the OR relationship.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the matches key.
  Only **resource_name** for key is supported.

* `value` - (Required, String) Specifies the matches value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of matched resources.

* `resources` - The resource details.
  The [resources](#vpcep_resources) structure is documented below.

<a name="vpcep_resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name. If the resource does not have a name, the ID is returned.

* `tags` - The tags list.
  The [tags](#vpcep_resources_tags) structure is documented below.

<a name="vpcep_resources_tags"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
