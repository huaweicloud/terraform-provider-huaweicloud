---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_project_tags"
description: |-
  Use this data source to query the project tags of GeminiDB within HuaweiCloud.
---

# huaweicloud_geminidb_project_tags

Use this data source to query the project tags of GeminiDB within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_geminidb_project_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the project tags.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of project tags.  
  The [tags](#tags) structure is documented below.

<a name="tags"></a>
The `tags` block supports:

* `key` - The tag key.

* `type` - The tag type.  
  The valid values are as follows:
  + **user**
  + **system**

* `values` - The list of tag values.
