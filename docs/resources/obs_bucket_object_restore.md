---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_object_restore"
description: |-
  Use this resource to restore an archived OBS object within HuaweiCloud.
---

# huaweicloud_obs_bucket_object_restore

Use this resource to restore an archived OBS object within HuaweiCloud.

-> This resource is a one-time action resource for restoring an archived OBS object. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

-> After creating a resource, a new object of the corresponding type is generated, and no further creation is allowed
   within the specified time period.

## Example Usage

### Restore an archived object with standard tier

```hcl
variable "obs_bucket_name" {}
variable "obs_object_key" {}

resource "huaweicloud_obs_bucket_object_restore" "test" {
  bucket = var.obs_bucket_name
  key    = var.obs_object_key
  days   = 7
  tier   = "standard"
}
```

### Restore an archived object with expedited tier

```hcl
variable "obs_bucket_name" {}
variable "obs_object_key" {}

resource "huaweicloud_obs_bucket_object_restore" "test" {
  bucket = var.obs_bucket_name
  key    = var.obs_object_key
  days   = 3
  tier   = "expedited"
}
```

### Restore a specific version of an archived object

```hcl
variable "obs_bucket_name" {}
variable "obs_object_key" {}
variable "obs_version_id" {}

resource "huaweicloud_obs_bucket_object_restore" "test" {
  bucket     = var.obs_bucket_name
  key        = var.obs_object_key
  version_id = var.obs_version_id
  days       = 5
  tier       = "standard"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the archived object is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `bucket` - (Required, String, NonUpdatable) Specifies the name of the bucket to restore the object from.

* `key` - (Required, String, NonUpdatable) Specifies the name of the object to restore.

* `days` - (Required, Int, NonUpdatable) Specifies the number of days for which the restored object copy is valid.  
  The valid value is range from `1` to `20`.

* `tier` - (Required, String, NonUpdatable) Specifies the restore option.  
  The valid values are as follows:
  + **expedited**: Quick restore of the object.
  + **standard**: Standard restore of the object.

* `version_id` - (Optional, String, NonUpdatable) Specifies the version ID of the object to restore.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
