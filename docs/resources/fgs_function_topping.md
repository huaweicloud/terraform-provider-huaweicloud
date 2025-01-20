---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function_topping"
description: |-
  Using this function to top function within HuaweiCloud.
---

# huaweicloud_fgs_function_topping

Using this function to top function within HuaweiCloud.

## Example Usage

### Topping function

```hcl
variable "function_urn" {}

resource "huaweicloud_fgs_function_topping" "test" {
  function_urn = var.function_urn
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region where the function is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `function_urn` - (Required, String, ForceNew) Specifies the URN of the function to be topped.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
