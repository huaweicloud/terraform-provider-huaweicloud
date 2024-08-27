---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_policy"
description: ""
---

# huaweicloud_obs_bucket_policy

Attaches a policy to an OBS bucket resource.

-> **NOTE:** When creating or updating the OBS bucket policy, the original policy will be overwritten.

## Example Usage

### Policy with OBS format

```hcl
resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "my-test-bucket"
}

resource "huaweicloud_obs_bucket_policy" "policy" {
  bucket = huaweicloud_obs_bucket.bucket.id
  policy = <<POLICY
{
  "Statement": [
    {
      "Sid": "AddPerm",
      "Effect": "Allow",
      "Principal": {"ID": "*"},
      "Action": ["GetObject"],
      "Resource": "my-test-bucket/*"
    }
  ]
}
POLICY
}
```

### Policy with S3 format

```hcl
resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "my-test-bucket"
}

resource "huaweicloud_obs_bucket_policy" "s3_policy" {
  bucket        = huaweicloud_obs_bucket.bucket.id
  policy_format = "s3"
  policy        = <<POLICY
{
  "Version": "2008-10-17",
  "Id": "MYBUCKETPOLICY",
  "Statement": [
    {
      "Sid": "IPAllow",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::my-test-bucket/*",
      "Condition": {
        "IpAddress": {"aws:SourceIp": "8.8.8.8/32"}
      }
    }
  ]
}
POLICY
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the OBS bucket policy resource. If omitted, the
  provider-level region will be used. Changing this creates a new OBS bucket policy resource.

* `bucket` - (Required, String, ForceNew) Specifies the name of the bucket to which to apply the policy.

* `policy` - (Required, String) Specifies the text of the bucket policy in JSON format. For more information about obs
  format bucket policy,
  see the [Developer Guide](https://support.huaweicloud.com/intl/en-us/perms-cfg-obs/obs_40_0004.html).

* `policy_format` - (Optional, String) Specifies the policy format, the supported values are *obs* and *s3*. Defaults
  to *obs* .

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

## Import

OBS format bucket policy can be imported using the `<bucket>`, e.g.

```bash
$ terraform import huaweicloud_obs_bucket_policy.policy <bucket-name>
```

S3 format bucket policy can be imported using the `<bucket>` and "s3" by a slash, e.g.

```bash
$ terraform import huaweicloud_obs_bucket_policy.s3_policy <bucket-name>/s3
```
