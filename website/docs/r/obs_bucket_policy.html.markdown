---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_policy"
sidebar_current: "docs-huaweicloud-resource-obs-bucket-policy"
description: |-
  Attaches a policy to an OBS bucket resource.
---

# huaweicloud\_obs\_bucket\_policy

Attaches a policy to an OBS bucket resource.

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
resource "huaweicloud_obs_bucket" "b" {
  bucket = "my-test-bucket"
}

resource "huaweicloud_obs_bucket_policy" "s3_policy" {
  bucket = huaweicloud_obs_bucket.bucket.id
  policy_format = "s3"
  policy = <<POLICY
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

* `bucket` - (Required) Specifies the name of the bucket to which to apply the policy.
* `policy` - (Required) Specifies the text of the bucket policy in JSON format. For more information about
  obs format bucket policy, see the [Developer Guide](https://support.huaweicloud.com/intl/en-us/devg-obs/obs_06_0048.html).
* `policy_format` - (Optional) Specifies the policy format, the supported values are *obs* and *s3*. Defaults to *obs* .

## Attributes Reference

The following attributes are exported:

* `bucket` - See Argument Reference above.
* `policy` - See Argument Reference above.
* `policy_format` - See Argument Reference above.

## Import

OBS format bucket policy can be imported using the `<bucket>`, e.g.

```
$ terraform import huaweicloud_obs_bucket_policy.policy <bucket-name>
```

S3 foramt bucket policy can be imported using the `<bucket>` and "s3" by a slash, e.g.

```
$ terraform import huaweicloud_obs_bucket_policy.s3_policy <bucket-name>/s3
```
