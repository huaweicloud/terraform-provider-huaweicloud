---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_certificate"
description: |-
  Use this data source to get the certificate of WAF within HuaweiCloud.
---

# huaweicloud_waf_certificate

Use this data source to get the certificate of WAF within HuaweiCloud.

-> When multiple pieces of data are queried, the datasource will process the first piece of data and put it back.

## Example Usage

```hcl
variable enterprise_project_id {}

data "huaweicloud_waf_certificate" "test" {
  name                  = "test-name"
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the WAF. If omitted, the provider-level region
  will be used.

* `name` - (Optional, String) Specifies the name of certificate. The value is case-sensitive and supports fuzzy matching.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of WAF certificate.
  For enterprise users, if omitted, default enterprise project will be used.

* `expiration_status` - (Optional, String) Specifies the certificate expiration status. The options are as follows:
  + `0`: Not expired;
  + `1`: Expired;
  + `2`: Expired soon (The certificate will expire in one month.)

  -> If this field is not configured, all certificates that meet the expired status will be found.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The certificate ID in UUID format.

* `created_at` - Indicates the time when the certificate uploaded, in RFC3339 format.

* `expired_at` - Indicates the time when the certificate expires, in RFC3339 format.
