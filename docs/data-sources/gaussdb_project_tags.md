---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_project_tags"
description: |-
  Use this data source to query the tags of GaussDB instances within a project in HuaweiCloud.
---

# huaweicloud_gaussdb_project_tags

Use this data source to query the tags of GaussDB instances within a project in HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_project_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.

  The [tags](#gaussdb_project_tags_struct) structure is documented below.

<a name="gaussdb_project_tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - The list of tag values.
