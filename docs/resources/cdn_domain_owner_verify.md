---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_owner_verify"
description: |-
  Use this resource to verify domain owner within HuaweiCloud.
---

# huaweicloud_cdn_domain_owner_verify

Use this resource to verify domain owner within HuaweiCloud.

-> This is resource only a one-time action resource for verify domain owner. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information
   from the tfstate file.

## Example Usage

### Verify domain ownership with default verification method

```hcl
variable "domain_name" {}

resource "huaweicloud_cdn_domain_owner_verify" "test" {
  domain_name = var.domain_name
}
```

### Verify domain ownership with DNS verification

```hcl
variable "domain_name" {}

resource "huaweicloud_cdn_domain_owner_verify" "test" {
  domain_name = var.domain_name
  verify_type = "dns"
}
```

### Verify domain ownership with file verification

```hcl
variable "domain_name" {}

resource "huaweicloud_cdn_domain_owner_verify" "test" {
  domain_name = var.domain_name
  verify_type = "file"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String, NonUpdatable) Specifies the domain name to be verified.

* `verify_type` - (Optional, String, NonUpdatable) Specifies the verification method.  
  The valid values are as follows:
  + **dns**: DNS resolution verification.
  + **file**: File verification.
  + **all**: Both DNS and file verification will be performed, and verification
    succeeds if either method passes.

  Defaults to **all**.
