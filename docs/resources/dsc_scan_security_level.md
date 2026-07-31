---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_scan_security_level"
description: |-
  Manages a DSC scan security level resource within HuaweiCloud.
---

# huaweicloud_dsc_scan_security_level

Manages a DSC scan security level resource within HuaweiCloud.

## Example Usage

```hcl
variable "security_level_name" {}

resource "huaweicloud_dsc_scan_security_level" "test" {
  security_level_name = var.security_level_name
  color_number        = 6
  security_level_desc = "Created by terraform script"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the security level is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `security_level_name` - (Required, String) Specifies the name of the security level.

* `color_number` - (Optional, Int) Specifies the color number of the security level displayed on the console.

* `security_level_desc` - (Optional, String) Specifies the description of the security level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the security level ID.

* `is_deleted` - Whether the security level is disabled.

* `category` - The category of the security level.
  The valid values are as follows:
  + **BUILT_IN**: The built-in security level.
  + **BUILT_IN_COPY**: The copy of the built-in security level.
  + **BUILT_SELF**: The custom security level.

* `used_count` - The number of the scan templates bound to the security level.

* `sort_weight` - The sort weight of the security level.

* `create_time` - The creation time of the security level.

* `update_time` - The latest update time of the security level.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dsc_scan_security_level.test <id>
```
