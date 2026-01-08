---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_public_zone_detection_status"
description: |-
  Use this data source to get the domain detection status of the public zone within HuaweiCloud.
---

# huaweicloud_dns_public_zone_detection_status

Use this data source to get the domain detection status of the public zone within HuaweiCloud.

## Example Usage

```hcl
variable "zone_id" {}
variable "domain_name" {}

data "huaweicloud_dns_public_zone_detection_status" "test" {
  zone_id     = var.zone_id
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Specifies the ID of the public zone.

* `domain_name` - (Required, String) Specifies the name of the record set to be detected.

* `type` - (Optional, String) Specifies the type of the record set to be detected.  
  The valid values are as follows:
  + **MX**
  + **CNAME**
  + **TXT**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The domain detection status.
  + **OK**: The resolution is successful.
  + **FAILED**: Whois query failed.
  + **NOT_REGISTERED**: The domain name is not registered.
  + **CANNOT_RESOLVE**: The domain name cannot be resolved.
  + **NOT_HWDNS**: The domain name is not hosted on Huawei Cloud.
  + **NO_WEBSITE_RECORD**: No website record set is configured.
  + **NO_EMAIL_RECORD**: No email record set is configured.
  + **NO_DEFAULT_VIEW**: No default record set is configured.
  + **NOT_EFFECT**: The authoritative server is not applied.
