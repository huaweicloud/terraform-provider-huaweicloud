---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_tags"
description: |-
  Use this data source to query the tag list of all resources of the same type within HuaweiCloud.
---

# huaweicloud_cbr_tags

Use this data source to query the tag list of all resources of the same type within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cbr_tags" "test" {
  resource_type = "vault"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource tags.  
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type to which the tags belong that to be queried.  
  The valid values are as follows:
  + **vault**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of all tags for resources of the same type.  
  The [tags](#cbr_project_tags) structure is documented below.

<a name="cbr_project_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

* `values` - All values corresponding to the key.
