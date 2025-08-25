---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_resource_tags"
description: |-
  Use this data source to get resource tag list of CTS within HuaweiCloud.
---

# huaweicloud_cts_resource_tags

Use this data source to get resource tag list of CTS within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_obs_bucket" "test" {
  bucket = "tf-test-bucket"
  acl    = "private"
}

resource "huaweicloud_cts_tracker" "test" {
  bucket_name = huaweicloud_obs_bucket.test.bucket
  file_prefix = "cts"

  tags = {
    foo1 = "bar1",
    foo2 = "bar2"
  }
}

data "huaweicloud_cts_resource_tags" "test" {
  resource_id = huaweicloud_cts_tracker.test.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.  
  If omitted, the provider-level region will be used.

* `resource_id` - (Required, String) Specifies the resource ID.

* `resource_type` - (Required, String) Specifies the resource type. The valid value is **cts-tracker**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `tags` - The list of tags.  
  The [tags](#cts_resource_tags_attr) structure is documented below.

<a name="cts_resource_tags_attr"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
