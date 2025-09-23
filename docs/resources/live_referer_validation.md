---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_referer_validation"
description: |-
  Manages a referer validation resource within HuaweiCloud.
---

# huaweicloud_live_referer_validation

Manages a referer validation resource within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "referer_config_empty" {}
variable "referer_white_list" {}
variable "referer_auth_list" {
  type = list(string)
}

resource "huaweicloud_live_referer_validation" "test" {
  domain_name          = var.domain_name
  referer_config_empty = var.referer_config_empty
  referer_white_list   = var.referer_white_list
  referer_auth_list    = var.referer_auth_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the streaming domain name to which the referer validation
  belongs.
  Changing this parameter will create a new resource.

* `referer_config_empty` - (Required, String) Specifies whether the referer header is included.
  The value can be **true** or **false**.

* `referer_white_list` - (Required, String) Specifies whether the referer is in the trustlist.
  The valid values are as follows:
  + **true**: Indicates referer whitelist.
  + **false**: Indicates referer blacklist.

* `referer_auth_list` - (Required, List) Specifies the domain name list.
  The maximum length is `100`.
  The domain name can be a specific domain name or a regular expression. e.g. `www.example.com`, `www.*com`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, UUID format.

## Import

The resource can be imported using `domain_name`, e.g.

```bash
$ terraform import huaweicloud_live_referer_validation.test <domain_name>
```
