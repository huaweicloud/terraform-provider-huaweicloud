---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_geo_blocking"
description: |-
  Manages a Live geo-blocking resource within HuaweiCloud.
---

# huaweicloud_live_geo_blocking

Manages a Live geo-blocking resource within HuaweiCloud.

-> Destroying this resource means that there is no geo-blocking on the streaming domain name.

## Example Usage

```hcl
variable "domain_name" {}

resource "huaweicloud_live_geo_blocking" "test" {
  domain_name    = var.domain_name
  app_name       = "live"
  area_whitelist = ["CN-IN", "CN-HK"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the streaming domain name.

  Changing this will create a new resource.

* `app_name` - (Required, String) Specifies the application name.

* `area_whitelist` - (Required, List) Specifies the list of supported areas.
  The values of all region codes, except that of China, contain two uppercase letters.
  For the code format, see [ISO 3166-1alpha-2](https://www.iso.org/obp/ui/#search/code/).
  Some options are as follows:
  + **CN-IN**: Chinese mainland.
  + **CN-HK**: Hong Kong (China).
  + **CN-MO**: Macao (China).
  + **CN-TW**: Taiwan (China).
  + **BR**: Brazil.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The Live geo-blocking resource can be imported using `domain_name` and `app_name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_live_geo_blocking.test <domain_name>/<app_name>
```
