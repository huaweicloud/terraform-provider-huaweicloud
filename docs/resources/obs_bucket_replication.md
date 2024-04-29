---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_replication"
description: ""
---

# huaweicloud_obs_bucket_replication

Manages an OBS bucket **Cross-Region Replication** resource within HuaweiCloud.

-> **NOTE:** When creating or updating the OBS bucket replication, the original bucket replication rules will be
overwritten. The source bucket and destination bucket must belong to the same account. More cross-Region replication
constraints see [Cross-Region replication](https://support.huaweicloud.com/intl/en-us/ugobs-obs/obs_41_0034.html).

## Example Usage

### Replicate all objects

```hcl
variable "bucket" {}
variable "destination_bucket" {}
variable "agency" {}

resource "huaweicloud_obs_bucket_replication" "test" {
  bucket             = var.bucket
  destination_bucket = var.destination_bucket
  agency             = var.agency
}
```

### Replicate objects matched by prefix

```hcl
variable "bucket" {}
variable "destination_bucket" {}
variable "agency" {}

resource "huaweicloud_obs_bucket_replication" "test" {
  bucket             = var.bucket
  destination_bucket = var.destination_bucket
  agency             = var.agency

  rule {
    prefix = "log"
  }

  rule {
    prefix          = "imgs/"
    storage_class   = "COLD"
    enabled         = true
    history_enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

  Changing this parameter will create a new resource.

* `bucket` - (Required, String, ForceNew) Specifies the name of the source bucket.

  Changing this parameter will create a new resource.

* `destination_bucket` - (Required, String) Specifies the name of the destination bucket.

  -> **NOTE:** The destination bucket cannot be in the region where the source bucket resides.
  Some regions do not support cross regional replication. More constraints information see:
  [Cross-Region replication](https://support.huaweicloud.com/intl/en-us/ugobs-obs/obs_41_0034.html)

* `agency` - (Required, String) Specifies the IAM agency applied to the cross-region replication.

  -> **NOTE:** The IAM agency is a cloud service agency of OBS. Which must has the **OBS Administrator** permission.

* `rule` - (Optional, List) Specifies the configurations of object cross-region replication management.
  The [rule_struct](#OBSBucketReplication_rule_struct) structure is documented below.

<a name="OBSBucketReplication_rule_struct"></a>
The `rule_struct` block supports:

* `prefix` - (Optional, String) Specifies the prefix of an object key name, applicable to one or more objects.
  The maximum length of a prefix is 1024 characters.
  Duplicated prefixes are not supported. If omitted, all objects in the bucket will be managed by the lifecycle rule.
  To copy a folder, end the prefix with a slash (/), for example, imgs/.

* `storage_class` - (Optional, String) Specifies the storage class for replicated objects. Valid values are `STANDARD`,
  `WARM` (Infrequent Access) and `COLD` (Archive).
  If omitted, the storage class of object copies is the same as that of objects in the source bucket.

* `enabled` - (Optional, Bool) Specifies cross-region replication rule status. Defaults to `true`.

* `history_enabled` - (Optional, Bool) Specifies cross-region replication history rule status. Defaults to `false`.
  If the value is `true`, historical objects meeting this rule are copied.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the bucket.
* `rule/id` - The ID of a rule in UUID format.

## Import

The obs bucket cross-region replication can be imported using the `bucket`, e.g.

```bash
$ terraform import huaweicloud_obs_bucket_replication.test <bucket-name>
```
