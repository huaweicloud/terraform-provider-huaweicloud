---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_global_variable"
description: ""
---

# huaweicloud_dli_global_variable

Manages a DLI global variable resource within HuaweiCloud.  

## Example Usage

```hcl
  resource "huaweicloud_dli_global_variable" "test" {
    name  = "demo"
    value = "abc"
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The name of a Global variable.
  This parameter can contain a maximum of `128` characters, which may consist of digits, letters, and underscores (\_),
  but cannot start with an underscore (\_) or contain only digits.

  Changing this parameter will create a new resource.

* `value` - (Required, String) The value of Global variable.

* `is_sensitive` - (Optional, Bool, ForceNew) Whether to set a variable as a sensitive variable. The default value is **false**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the name.

## Import

The global variable can be imported using the `id` which equals the name, e.g.

```bash
$ terraform import huaweicloud_dli_global_variable.test demo_name
```
