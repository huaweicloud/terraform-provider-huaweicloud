---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_private_cas_by_tags"
description: |-
  Use this data source to get a list of CCM private CAs by tags.
---

# huaweicloud_ccm_private_cas_by_tags

Use this data source to get a list of CCM private CAs by tags.

## Example Usage

```hcl
data "huaweicloud_ccm_private_cas_by_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `tags` - (Optional, List) Specifies the containing tags.
  It can contain a maximum of `20` keys, with a maximum of `20` values ​​under each key. The value array for each key can
  be empty, but the structure must be complete. Keys cannot be duplicated, and values ​​for the same key cannot be
  duplicated. The result returns a list of resources containing all tags. Keys are related by AND, and values ​​in the
  key-value structure are related by OR. Without tag filtering, the full dataset is returned.

  The [tags](#tags_Tag) structure is documented below.

* `matches` - (Optional, List) Specifies the matches condition.
  The key is the field to be matched, such as **resource_name**. The value is the match value.
  The key is a fixed dictionary value and cannot contain duplicate keys or unsupported keys.

  The [matches](#matches_Match) structure is documented below.

<a name="tags_Tag"></a>
The `tags` block supports:

* `key` - (Optional, String) Specifies the key.

* `values` - (Optional, List) Specifies the value array.

<a name="matches_Match"></a>
The `matches` block supports:

* `key` - (Optional, String) Specifies the field to be matched.

* `value` - (Optional, String) Specifies the match value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of resources.

* `resources` - The CA details.

  The [resources](#resources_response_struct) structure is documented below.

<a name="resources_response_struct"></a>
The `resources` block supports:

* `resource_id` - The ID of the private CA.

* `tags` - The tags of the private CA.
  The [tags](#tags_response_struct) structure is documented below.

* `resource_name` - The name of the private CA.

* `resource_detail` - The details of the private CA in JSON string format.

<a name="tags_response_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
