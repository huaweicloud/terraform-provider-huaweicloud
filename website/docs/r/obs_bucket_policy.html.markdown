---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_policy"
sidebar_current: "docs-huaweicloud-resource-obs-bucket-policy"
description: |-
  Attaches a policy to an OBS bucket resource.
---

# huaweicloud\_obs\_bucket\_policy

Attaches a policy to an OBS bucket resource.

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to which to apply the policy.
* `policy` - (Required) The text of the [bucket policy](https://support.huaweicloud.com/intl/en-us/devg-obs/obs_06_0048.html) in JSON format.

## Attributes Reference

The following attributes are exported:

* `bucket` - See Argument Reference above.
* `policy` - See Argument Reference above.

## Import

OBS bucket policy can be imported using the `bucket`, e.g.

```
$ terraform import huaweicloud_obs_bucket_policy.policy bucket-name
```
