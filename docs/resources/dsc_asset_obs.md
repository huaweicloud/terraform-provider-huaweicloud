---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_asset_obs"
description: |-
  Manages an OBS asset resource of DSC within HuaweiCloud.
---

# huaweicloud_dsc_asset_obs

Manages an OBS asset resource of DSC within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "bucket_name" {}

resource "huaweicloud_dsc_asset_obs" "test" {
  name          = var.name
  bucket_name   = var.bucket_name
  bucket_policy = "private"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of asset.

* `bucket_name` - (Required, String, ForceNew) Specifies the bucket name.

  Changing this parameter will create a new resource.

* `bucket_policy` - (Required, String, ForceNew) Specifies the bucket policy.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The OBS asset resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dsc_asset_obs.test <id>
```
