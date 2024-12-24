---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_notification_configuration"
description: |-
  Manages a notification configuration resource within HuaweiCloud.
---

# huaweicloud_live_notification_configuration

Manages a notification configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "url" {}

resource "huaweicloud_live_notification_configuration" "test" {
  domain_name = var.domain_name
  url         = var.url
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the ingest domain name to which the notification configuration
  belongs.
  Changing this parameter will create a new resource.

* `url` - (Required, String) Specifies the callback URL, which must start with `http://` or `https://`.

* `auth_sign_key` - (Optional, String) Specifies the authentication key.
  The valid length is `32` to `128` characters.

* `call_back_area` - (Optional, String) Specifies the region where the server that receives callback notifications
  is located.
  The valid vaules are as follows:
  + **mainland_china**: Indicates Chinese mainland.
  + **outside_mainland_china**: Indicates outside the Chinese mainland.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, UUID format.

## Import

The resource can be imported using `domain_name`, e.g.

```bash
$ terraform import huaweicloud_live_notification_configuration.test <domain_name>
```
