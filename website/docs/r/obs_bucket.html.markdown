---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket"
sidebar_current: "docs-huaweicloud-resource-obs-bucket"
description: |-
  Provides an OBS bucket resource.
---

# huaweicloud\_obs\_bucket

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
      days = 60
      storage_class = "WARM"
    }
    transition {
      days = 180
      storage_class = "COLD"
    }
  }

  lifecycle_rule {
    name    = "tmp"
    prefix  = "tmp/"
    enabled = true

    noncurrent_version_expiration {
      days = 180
    }
    noncurrent_version_transition {
      days = 30
      storage_class = "WARM"
    }
    noncurrent_version_transition {
      days = 60
      storage_class = "COLD"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) Specifies the name of the bucket. Changing this parameter will create a new resource.
  A bucket must be named according to the globally applied DNS naming regulations as follows:
	* The name must be globally unique in OBS.
	* The name must contain 3 to 63 characters. Only lowercase letters, digits, hyphens (-), and periods (.) are allowed.
	* The name cannot start or end with a period (.) or hyphen (-), and cannot contain two consecutive periods (.) or contain a period (.) and a hyphen (-) adjacent to each other.
	* The name cannot be an IP address.
	* If the name contains any periods (.), a security certificate verification message may appear when you access the bucket or its objects by entering a domain name.

* `storage_class` - (Optional) Specifies the storage class of the bucket. OBS provides three storage classes: "STANDARD", "WARM" (Infrequent Access) and "COLD" (Archive). Defaults to `STANDARD`.

* `acl` - (Optional) Specifies the ACL policy for a bucket. The predefined common policies are as follows: "private", "public-read", "public-read-write" and "log-delivery-write". Defaults to `private`.

* `tags` - (Optional) A mapping of tags to assign to the bucket. Each tag is represented by one key-value pair.

* `versioning` - (Optional) Whether enable versioning. Once you version-enable a bucket, it can never return to an unversioned state.
  You can, however, suspend versioning on that bucket.

* `logging` - (Optional) A settings of bucket logging (documented below).
* `website` - (Optional) A website object (documented below).
* `cors_rule` - (Optional) A rule of Cross-Origin Resource Sharing (documented below).
* `lifecycle_rule` - (Optional) A configuration of object lifecycle management (documented below).

* `force_destroy` - (Optional) A boolean that indicates all objects should be deleted from the bucket so that the bucket can be destroyed without error. Default to `false`.

* `region` - (Optional) If specified, the region this bucket should reside in. Otherwise, the region used by the provider.

The `logging` object supports the following:

* `target_bucket` - (Required) The name of the bucket that will receive the log objects.
  The acl policy of the target bucket should be `log-delivery-write`.
* `target_prefix` - (Optional) To specify a key prefix for log objects.

The `website` object supports the following:

* `index_document` - (Required, unless using `redirect_all_requests_to`) Specifies the default homepage of the static website, only HTML web pages are supported.
  OBS only allows files such as `index.html` in the root directory of a bucket to function as the default homepage.
  That is to say, do not set the default homepage with a multi-level directory structure (for example, /page/index.html).

* `error_document` - (Optional) Specifies the error page returned when an error occurs during static website access.
  Only HTML, JPG, PNG, BMP, and WEBP files under the root directory are supported.

* `redirect_all_requests_to` - (Optional) A hostname to redirect all website requests for this bucket to. Hostname can optionally be prefixed with a protocol (`http://` or `https://`) to use when redirecting requests. The default is the protocol that is used in the original request.

* `routing_rules` - (Optional) A JSON or XML format containing routing rules describing redirect behavior and when redirects are applied.
  Each rule contains a `Condition` and a `Redirect` as shown in the following table:

Parameter | Key
-|-
Condition | KeyPrefixEquals, HttpErrorCodeReturnedEquals
Redirect | Protocol, HostName, ReplaceKeyPrefixWith, ReplaceKeyWith, HttpRedirectCode

The `cors_rule` object supports the following:

* `allowed_origins` (Required) Requests from this origin can access the bucket. Multiple matching rules are allowed.
  One rule occupies one line, and allows one wildcard character (*) at most.

* `allowed_methods` (Required) Specifies the acceptable operation type of buckets and objects.
  The methods include `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.

* `allowed_headers` (Optional) Specifies the allowed header of cross-origin requests.
  Only CORS requests matching the allowed header are valid.

* `expose_headers` (Optional) Specifies the exposed header in CORS responses, providing additional information for clients.

* `max_age_seconds` (Optional) Specifies the duration that your browser can cache CORS responses, expressed in seconds.
  The default value is 100.

The `lifecycle_rule` object supports the following:

* `name` - (Required) Unique identifier for lifecycle rules. The Rule Name contains a maximum of 255 characters.

* `enabled` - (Required) Specifies lifecycle rule status.

* `prefix` - (Optional) Object key prefix identifying one or more objects to which the rule applies.
  If omitted, all objects in the bucket will be managed by the lifecycle rule.
  The prefix cannot start or end with a slash (/), cannot have consecutive slashes (/), and cannot contain the following special characters: \:*?"<>|.

* `expiration` - (Optional) Specifies a period when objects that have been last updated are automatically deleted. (documented below).
* `transition` - (Optional) Specifies a period when objects that have been last updated are automatically transitioned to `WARM` or `COLD` storage class (documented below).
* `noncurrent_version_expiration` - (Optional) Specifies a period when noncurrent object versions are automatically deleted. (documented below).
* `noncurrent_version_transition` - (Optional) Specifies a period when noncurrent object versions are automatically transitioned to `WARM` or `COLD` storage class (documented below).

At least one of `expiration`, `transition`, `noncurrent_version_expiration`, `noncurrent_version_transition` must be specified.

The `expiration` object supports the following

* `days` (Required) Specifies the number of days when objects that have been last updated are automatically deleted.
  The expiration time must be greater than the transition times.

The `transition` object supports the following

* `days` (Required) Specifies the number of days when objects that have been last updated are automatically transitioned to the specified storage class.
* `storage_class` - (Required) The class of storage used to store the object. Only `WARM` and `COLD` are supported.

The `noncurrent_version_expiration` object supports the following

* `days` (Required) Specifies the number of days when noncurrent object versions are automatically deleted.

The `noncurrent_version_transition` object supports the following

* `days` (Required) Specifies the number of days when noncurrent object versions are automatically transitioned to the specified storage class.
* `storage_class` - (Required) The class of storage used to store the object. Only `WARM` and `COLD` are supported.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the bucket.
* `bucket_domain_name` - The bucket domain name. Will be of format `bucketname.obs.region.myhuaweicloud.com`.
* `region` - The region this bucket resides in.

## Import

OBS bucket can be imported using the `bucket`, e.g.

```
$ terraform import huaweicloud_obs_bucket.bucket bucket-name
```
