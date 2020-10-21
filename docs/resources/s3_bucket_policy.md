---
subcategory: "Deprecated"
---

# huaweicloud\_s3\_bucket\_policy

!> **Warning:** It has been deprecated, use `huaweicloud_obs_bucket_policy` instead.

Attaches a policy to an S3 bucket resource.

## Example Usage

### Basic Usage

```hcl
resource "huaweicloud_s3_bucket" "b" {
  bucket = "my-tf-test-bucket"
}

resource "huaweicloud_s3_bucket_policy" "policy" {
  bucket = huaweicloud_s3_bucket.b.id
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
      "Resource": "arn:aws:s3:::my-tf-test-bucket/*",
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

* `bucket` - (Required) The name of the bucket to which to apply the policy.
* `policy` - (Required) The text of the policy.
