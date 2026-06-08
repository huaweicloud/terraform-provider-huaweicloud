---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_mirror_back_to_source_rules"
description: |-
  Use this data source to query the mirror back to source rules of an OBS bucket within HuaweiCloud.
---

# huaweicloud_obs_bucket_mirror_back_to_source_rules

Use this data source to query the mirror back to source rules of an OBS bucket within HuaweiCloud.

## Example Usage

```hcl
variable "bucket_name" {}

data "huaweicloud_obs_bucket_mirror_back_to_source_rules" "test" {
  bucket = var.bucket_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the OBS bucket mirror back to source rules are located.  
  If omitted, the provider-level region will be used.

* `bucket` - (Required, String) Specifies the name of the OBS bucket.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The rules of the bucket mirror back to source, in JSON format.
