---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_appcode"
description: ""
---

# huaweicloud_apig_appcode

Manages an APPCODE in application resource within HuaweiCloud.

## Example Usage

### Auto generate APPCODE

```hcl
variable "instance_id" {}
variable "application_id" {}

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = var.instance_id
  application_id = var.application_id
}
```

### Manually configure APPCODE

```hcl
variable "instance_id" {}
variable "application_id" {}
variable "app_code" {}

resource "huaweicloud_apig_appcode" "test" {
  instance_id    = var.instance_id
  application_id = var.application_id
  value          = var.app_code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application and APPCODE are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the application
  and APPCODE belong.  
  Changing this will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the ID of application to which the APPCODE belongs.
  Changing this will create a new resource.

* `value` - (Optional, String, ForceNew) Specifies the APPCODE value (content).
  The value can contain `64` to `180` characters, starting with a letter, plus sign (+), or slash (/), or digit.
  Only letters, digit and the following special characters are allowed: `+_!@#$%/=`.
  If omitted, a random value will be generated.
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The APPCODE ID.

* `created_at` - The creation time of the APPCODE.

## Import

APPCODEs can be imported using related `instance_id`, `application_id` and their `id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_apig_appcode.test <instance_id>/<application_id>/<id>
```
