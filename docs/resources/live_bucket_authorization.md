---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_bucket_authorization"
description: ""
---

# huaweicloud_live_bucket_authorization

Manages a Live bucket authorization resource within HuaweiCloud.

## Example Usage

```hcl
variable "bucket" {}

resource "huaweicloud_live_bucket_authorization" "test"{
  bucket = var.bucket
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `bucket` - (Required, String, ForceNew) Specifies the bucket name of the OBS.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The live bucket authorize can be imported using the `bucket`, e.g.

```bash
$ terraform import huaweicloud_live_bucket_authorization.test <bucket>
```
