---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_object"
description: ""
---

# huaweicloud_obs_bucket_object

Use this data source to get info of special HuaweiCloud obs object.

```hcl
data "huaweicloud_obs_bucket_object" "object" {
  bucket = "my-test-bucket"
  key    = "new-key"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the OBS object. If omitted, the provider-level region will
  be used.

* `bucket` - (Required, String) The name of the bucket to put the file in.

* `key` - (Required, String) The name of the object once it is in the bucket.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - the `key` of the resource supplied above.

* `etag` - the ETag generated for the object (an MD5 sum of the object content). When the object is encrypted on the
  server side, the ETag value is not the MD5 value of the object, but the unique identifier calculated through the
  server-side encryption.

* `size` - the size of the object in bytes.

* `version_id` - a unique version ID value for the object, if bucket versioning is enabled.

* `storage_class` - specifies the storage class of the object.

* `content_type` - a standard MIME type describing the format of the object data, e.g. application/octet-stream. All
  Valid MIME Types are valid for this input.

* `body` - The content of an object which is available only for objects which have a human-readable Content-Type
  (text/* and application/json) and smaller than **64KB**. This is to prevent printing unsafe characters and
  potentially downloading large amount of data.
