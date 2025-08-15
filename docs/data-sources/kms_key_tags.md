---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_key_tags"
description: |-
  Use this data source to query the list of KMS key tags within HuaweiCloud.
---

# huaweicloud_kms_key_tags

Use this data source to query the list of KMS key tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_kms_key_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - the list of key tags.  
  The [tags](#kms_project_tags) structure is documented below.

<a name="kms_project_tags"></a>
The `tags` block supports:

* `key` - The tag key.

* `values` - The list of tag values.
