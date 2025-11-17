---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_owner_verification"
description: |-
  Use this data source to get the domain ownership verification information of CDN within HuaweiCloud.
---

# huaweicloud_cdn_domain_owner_verification

Use this data source to get the domain ownership verification information of CDN within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_cdn_domain_owner_verification" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the name of the accelerated domain.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dns_verify_type` - The DNS resolution type.

* `dns_verify_name` - The DNS resolution host record name.

* `file_verify_url` - The file verification URL address.

* `verify_domain_name` - The verification domain name.

* `file_verify_filename` - The file verification filename.

* `verify_content` - The verification value, which is the resolution value or file content.

* `file_verify_domains` - The list of file verification domain names.
