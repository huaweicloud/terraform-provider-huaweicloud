---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket"
description: ""
---

# huaweicloud_obs_bucket

Provides an OBS bucket resource.

## Example Usage

### Private Bucket with Tags

```hcl
resource "huaweicloud_obs_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl    = "private"

  tags = {
    foo = "bar"
    Env = "Test"
  }
}
```

### Enable versioning

```hcl
resource "huaweicloud_obs_bucket" "b" {
  bucket     = "my-tf-test-bucket"
  acl        = "private"
  versioning = true
}
```

### Enable Logging

```hcl
variable "agency_name" {} # The agency must be an OBS cloud service agency and has the `PutObject` permission.

resource "huaweicloud_obs_bucket" "log_bucket" {
  bucket = "my-tf-log-bucket"
  acl    = "log-delivery-write"
}

resource "huaweicloud_obs_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl    = "private"

  logging {
    target_bucket = huaweicloud_obs_bucket.log_bucket.id
    target_prefix = "log/"
    agency        = var.agency_name
  }
}
```

### Static Website Hosting

```hcl
resource "huaweicloud_obs_bucket" "b" {
  bucket = "obs-website-test.hashicorp.com"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"

    routing_rules = <<EOF
[{
    "Condition": {
        "KeyPrefixEquals": "docs/"
    },
    "Redirect": {
        "ReplaceKeyPrefixWith": "documents/"
    }
}]
EOF
  }
}
```

### Using CORS

```hcl
resource "huaweicloud_obs_bucket" "b" {
  bucket = "obs-website-test.hashicorp.com"
  acl    = "public-read"

  cors_rule {
    allowed_origins = ["https://obs-website-test.hashicorp.com"]
    allowed_methods = ["PUT", "POST"]
    allowed_headers = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
```

### Using object lifecycle

```hcl
resource "huaweicloud_obs_bucket" "bucket" {
  bucket     = "my-bucket"
  acl        = "private"
  versioning = true

  lifecycle_rule {
    name    = "log"
    prefix  = "log/"
    enabled = true

    expiration {
      days = 365
    }
    transition {
      days          = 60
      storage_class = "WARM"
    }
    transition {
      days          = 180
      storage_class = "COLD"
    }
    abort_incomplete_multipart_upload {
      days = 360
    }
  }

  lifecycle_rule {
    name    = "tmp"
    enabled = true

    noncurrent_version_expiration {
      days = 180
    }
    noncurrent_version_transition {
      days          = 30
      storage_class = "WARM"
    }
    noncurrent_version_transition {
      days          = 60
      storage_class = "COLD"
    }
  }
}
```

### using encryption

```hcl
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "my-tf-encryption-bucket"
  storage_class = "STANDARD"
  acl           = "private"
  encryption    = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Specifies the name of the bucket. Changing this parameter will create a new
  resource. A bucket must be named according to the globally applied DNS naming regulations as follows:
    + The name must be globally unique in OBS.
    + The name must contain 3 to 63 characters. Only lowercase letters, digits, hyphens (-), and periods (.) are
      allowed.
    + The name cannot start or end with a period (.) or hyphen (-), and cannot contain two consecutive periods (.) or
      contain a period (.) and a hyphen (-) adjacent to each other.
    + The name cannot be an IP address.
    + If the name contains any periods (.), a security certificate verification message may appear when you access the
      bucket or its objects by entering a domain name.

* `storage_class` - (Optional, String) Specifies the storage class of the bucket. OBS provides three storage classes:
  "STANDARD", "WARM" (Infrequent Access) and "COLD" (Archive). Defaults to `STANDARD`.

* `acl` - (Optional, String) Specifies the ACL policy for a bucket. The predefined common policies are as follows:
  "private", "public-read", "public-read-write" and "log-delivery-write". Defaults to `private`.

* `policy` - (Optional, String) Specifies the text of the bucket policy in JSON format. For more information about obs
  format bucket policy,
  see the [Developer Guide](https://support.huaweicloud.com/intl/en-us/perms-cfg-obs/obs_40_0004.html).

* `policy_format` - (Optional, String) Specifies the policy format, the supported values are *obs* and *s3*. Defaults
  to *obs*.

* `tags` - (Optional, Map) A mapping of tags to assign to the bucket. Each tag is represented by one key-value pair.

* `versioning` - (Optional, Bool) Whether enable versioning. Once you version-enable a bucket, it can never return to an
  unversioned state. You can, however, suspend versioning on that bucket.

* `logging` - (Optional, List) A settings of bucket logging (documented below).

<!-- markdownlint-disable MD033 -->

* `quota` - (Optional, Int) Specifies bucket storage quota. Must be a positive integer in the unit of byte. The maximum
  storage quota is 2<sup>63</sup> – 1 bytes. The default bucket storage quota is `0`, indicating that the bucket storage
  quota is not limited.

<!-- markdownlint-enable MD033 -->

* `website` - (Optional, List) A website object (documented below).

* `cors_rule` - (Optional, List) A rule of Cross-Origin Resource Sharing (documented below).

* `lifecycle_rule` - (Optional, List) A configuration of object lifecycle management (documented below).

* `force_destroy` - (Optional, Bool) A boolean that indicates all objects should be deleted from the bucket, so that the
  bucket can be destroyed without error. Default to `false`.

* `region` - (Optional, String, ForceNew) Specifies the region where this bucket will be created. If not specified, used
  the region by the provider. Changing this will create a new bucket.

* `parallel_fs` - (Optional, Bool, ForceNew) Whether enable a bucket as a parallel file system. Changing this will
  create a new bucket.

* `multi_az` - (Optional, Bool, ForceNew) Whether enable the multi-AZ mode for the bucket. When the multi-AZ mode is
  enabled, data in the bucket is duplicated and stored in multiple AZs.

  -> **NOTE:** Once a bucket is created, you cannot enable or disable the multi-AZ mode. Changing this will create a new
  bucket, but the name of a deleted bucket can be reused for another bucket at least 30 minutes after the deletion.
  Exercise caution when changing this field.

* `encryption` - (Optional, Bool) Whether to enable default server-side encryption of the bucket.

* `sse_algorithm` - (Optional, String) Specifies the mode of encryption algorithm. The valid values are:
  + **kms**: Server-side encryption with keys hosted by KMS are used to encrypt your objects.
  + **AES256**: Server-side encryption with keys managed by OBS are used to encrypt your objects.

  Defaults to **kms**.

* `kms_key_id` - (Optional, String) Specifies the ID of a KMS key. If omitted, the default master key will be used. This
  field is used only when `sse_algorithm` value is **kms**.

* `kms_key_project_id` - (Optional, String) Specifies the project ID to which the KMS key belongs. This field is valid
  only when `kms_key_id` is specified.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the OBS bucket.
  Defaults to `0`.

* `user_domain_names` - (Optional, List) Specifies the user domain names. The restriction requirements for this field
  are as follows:
  + Each value must meet the domain name rules.
  + The maximum length of a domain name is 256 characters.
  + A maximum of 100 custom domain names can be set for a bucket.
  + A custom domain name can only be used by one bucket.
  + Ensure the domain name has been licensed by the Ministry of Industry and Information Technology.
  + The bound user domain names only support access over HTTP now.

  -> When creating or updating the OBS bucket user domain names, the original user domain names will be overwritten.

The `logging` object supports the following:

* `target_bucket` - (Required, String) The name of the bucket that will receive the log objects. The acl policy of the
  target bucket should be `log-delivery-write`.

* `target_prefix` - (Optional, String) To specify a key prefix for log objects.

* `agency` - (Required, String) Specifies the IAM agency of OBS cloud service.

  -> The IAM agency requires the `PutObject` permission for the target bucket.  If default encryption is enabled for the
  target bucket, the agency also requires the `KMS Administrator` permission in the region where the target bucket is
  located.

The `website` object supports the following:

* `index_document` - (Optional, String)  Unless using `redirect_all_requests_to`. Specifies the default homepage of the
  static website, only HTML web pages are supported. OBS only allows files such as `index.html` in the root directory of
  a bucket to function as the default homepage. That is to say, do not set the default homepage with a multi-level
  directory structure (for example, /page/index.html).

* `error_document` - (Optional, String) Specifies the error page returned when an error occurs during static website
  access. Only HTML, JPG, PNG, BMP, and WEBP files under the root directory are supported.

* `redirect_all_requests_to` - (Optional, String) A hostname to redirect all website requests for this bucket to.
  Hostname can optionally be prefixed with a protocol (`http://` or `https://`) to use when redirecting requests. The
  default is the protocol that is used in the original request.

* `routing_rules` - (Optional, String) A JSON or XML format containing routing rules describing redirect behavior and
  when redirects are applied. Each rule contains a `Condition` and a `Redirect` as shown in the following table:

  Parameter | Key
      --- | ---
  Condition | KeyPrefixEquals, HttpErrorCodeReturnedEquals
  Redirect | Protocol, HostName, ReplaceKeyPrefixWith, ReplaceKeyWith, HttpRedirectCode

The `cors_rule` object supports the following:

* `allowed_origins` - (Required, List) Requests from this origin can access the bucket. Multiple matching rules are
  allowed. One rule occupies one line, and allows one wildcard character (*) at most.

* `allowed_methods` - (Required, List) Specifies the acceptable operation type of buckets and objects. The methods
  include `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.

* `allowed_headers` - (Optional, List) Specifies the allowed header of cross-origin requests. Only CORS requests
  matching the allowed header are valid.

* `expose_headers` - (Optional, List) Specifies the exposed header in CORS responses, providing additional information
  for clients.

* `max_age_seconds` - (Optional, Int) Specifies the duration that your browser can cache CORS responses, expressed in
  seconds. The default value is 100.

The `lifecycle_rule` object supports the following:

* `name` - (Required, String) Unique identifier for lifecycle rules. The Rule Name contains a maximum of 255 characters.

* `enabled` - (Required, Bool) Specifies lifecycle rule status.

* `prefix` - (Optional, String) Object key prefix identifying one or more objects to which the rule applies. If omitted,
  all objects in the bucket will be managed by the lifecycle rule. The prefix cannot start or end with a slash (/),
  cannot have consecutive slashes (/), and cannot contain the following special characters: \:*?"<>|.
  When configuring multiple `lifecycle_rule`, field `prefix` in multiple `lifecycle_rule` cannot have an inclusive
  relationship.

* `expiration` - (Optional, List) Specifies a period when objects that have been last updated are automatically
  deleted. (documented below).
* `transition` - (Optional, List) Specifies a period when objects that have been last updated are automatically
  transitioned to `WARM` or `COLD` storage class (documented below).
* `noncurrent_version_expiration` - (Optional, List) Specifies a period when noncurrent object versions are
  automatically deleted. (documented below).
* `noncurrent_version_transition` - (Optional, List) Specifies a period when noncurrent object versions are
  automatically transitioned to `WARM` or `COLD` storage class (documented below).
* `abort_incomplete_multipart_upload` - (Optional, List) Specifies a period when the not merged parts (fragments) in an
  incomplete upload are automatically deleted. (documented below).

At least one of `expiration`, `transition`, `noncurrent_version_expiration`, `noncurrent_version_transition`,
`abort_incomplete_multipart_upload` must be specified. The parameter `versioning` must be set to **true** before using
`noncurrent_version_expiration` or `noncurrent_version_transition`.

The `expiration` object supports the following

* `days` - (Required, Int) Specifies the number of days when objects that have been last updated are automatically
  deleted. The expiration time must be greater than the transition times.

The `transition` object supports the following

* `days` - (Required, Int) Specifies the number of days when objects that have been last updated are automatically
  transitioned to the specified storage class.
* `storage_class` - (Required, String) The class of storage used to store the object. Only `WARM` and `COLD` are
  supported.

The `noncurrent_version_expiration` object supports the following

* `days` - (Required, Int) Specifies the number of days when noncurrent object versions are automatically deleted.

The `noncurrent_version_transition` object supports the following

* `days` - (Required, Int) Specifies the number of days when noncurrent object versions are automatically transitioned
  to the specified storage class.
* `storage_class` - (Required, String) The class of storage used to store the object. Only `WARM` and `COLD` are
  supported.

The `abort_incomplete_multipart_upload` object supports the following

* `days` - (Required, Int) Specifies the number of days since the initiation of an incomplete multipart upload that OBS
  will wait before deleting the not merged parts (fragments) of the upload.
  The valid value ranges from 1 to 2,147,483,647.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the bucket.
* `bucket_domain_name` - The bucket domain name. Will be of format `bucketname.obs.region.myhuaweicloud.com`.
* `bucket_version` - The OBS version of the bucket.
* `region` - The region where this bucket resides in.
* `storage_info` - The OBS storage info of the bucket.
  The [object](#bucket_storage_info_attr) structure is documented below.

<a name="bucket_storage_info_attr"></a>
The `storage_info` block supports:

* `size` - The stored size of the bucket.
* `object_number` - The number of objects stored in the bucket.

## Import

OBS bucket can be imported using the `bucket`, e.g.

```bash
$ terraform import huaweicloud_obs_bucket.bucket <bucket-name>
```

OBS bucket with S3 format bucket policy can be imported using the `bucket` and "s3" by a slash, e.g.

```bash
$ terraform import huaweicloud_obs_bucket.bucket_with_s3_policy <bucket-name>/s3
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include `acl` and `force_destroy`. It is generally recommended
running `terraform plan` after importing an OBS bucket. Also you can ignore changes as below.

```hcl
resource "huaweicloud_obs_bucket" "bucket" {
    ...

  lifecycle {
    ignore_changes = [
      acl, force_destroy,
    ]
  }
}
```
