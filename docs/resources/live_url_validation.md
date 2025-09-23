---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_url_validation"
description: |-
  Manages an URL validation resource within HuaweiCloud.
---

# huaweicloud_live_url_validation

Manages an URL validation resource within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "key" {}
variable "auth_type" {}
variable "timeout" {}

resource "huaweicloud_live_url_validation" "test" {
  domain_name = var.domain_name
  key         = var.key
  auth_type   = var.auth_type
  timeout     = var.timeout
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the domain name to which the URL validation belongs.
  Including the ingest domain name and streaming domain name.
  Changing this parameter will create a new resource.

* `key` - (Required, String) Specifies the URL validation key value.
  The valid length is `32` characters, only uppercase letters, lwercase letters and digits are allowed.
  The value cannot contain only digits or letters.

* `auth_type` - (Required, String) Specifies the signing method of the URL validation.
  The valid values are as follows:
  + **d_sha256**: Indicates signing method D, which uses the HMAC-SHA256 algorithm. This method is recommended.
  + **c_aes**: Indicates signing method C, which uses the symmetric encryption algorithm.
  + **b_md5**: Indicates signing method B, which uses the MD5 algorithm.
  + **a_md5**: Indicates signing method A, which uses the MD5 algorithm.

  -> The signing methods A, B and C have security risks. The signing method D is more secure and recommended.

* `timeout` - (Required, Int) Specifies the timeout interval of URL validation.
  The valid value ranges from `60` to `2,592,000`, in seconds.

  -> This parameter is used to check whether the Live ingest URL or streaming URL has expired.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, UUID format.

## Import

The resource can be imported using `domain_name`, e.g.

```bash
$ terraform import huaweicloud_live_url_validation.test <domain_name>
```
