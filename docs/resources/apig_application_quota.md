---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_quota"
description: |-
  Manages an APIG application quota resource within HuaweiCloud.
---

# huaweicloud_apig_application_quota

Manages an APIG application quota resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}

resource "huaweicloud_apig_application_quota" "test" {
  instance_id   = var.instance_id
  name          = var.name
  time_unit     = "MINUTE"
  call_limits   = 200
  time_interval = 3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to
  which the application quota belongs. Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the application quota.  
  The value only Chinese and English letters, digits and underscores (_) are allowed,
  and must start with a Chinese or English letter. The valid value ranges from `3` to `255`.

* `time_unit` - (Required, String) Specifies the limited time unit of the application quota.  
  The valid values are as follows:
  + **SECOND**
  + **MINUTE**
  + **HOUR**
  + **DAY**

* `call_limits` - (Required, Int) Specifies the access limit of the application quota.  
  The valid value ranges from `1` to `2,147,483,647`.

* `time_interval` - (Required, Int) Specifies the limited time value for flow control of the application quota.  
  The valid value ranges from `1` to `2,147,483,647`.

* `description` - (Optional, String) Specifies the description of the application quota.  
  The description contain a maximum of `255` characters and the angle brackets (< and >) are not allowed.  
  Chinese characters must be in **UTF-8** or **Unicode** format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `bind_num` - The number of bound APPs.

* `created_at` - The creation time of the application quota, in RFC3339 format.

## Import

The application quota can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_apig_application_quota.test <instance_id>/<id>
```
