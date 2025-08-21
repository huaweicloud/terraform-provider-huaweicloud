---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function_tracing_configuration"
description: |-
  Manages a FunctionGraph function tracing configuration resource within HuaweiCloud.
---

# huaweicloud_fgs_function_tracing_configuration

Manages a FunctionGraph function tracing configuration resource within HuaweiCloud.

-> This feature is only supported by version `v2` and runtime `JAVA`, the function memory cannot be less than `512`.

## Example Usage

```hcl
variable "function_urn" {}
variable "tracing_access_key" {}
variable "tracing_secret_key" {}

resource "huaweicloud_fgs_function_tracing_configuration" "test" {
  function_urn = var.function_urn
  tracing_ak   = var.tracing_access_key
  tracing_sk   = var.tracing_secret_key
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region where the function tracing configuration is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `function_urn` - (Required, String, NonUpdatable) Specifies the URN of the function to which the tracing configuration
  belongs.  
  Changing this parameter will create a new resource.

* `tracing_ak` - (Required, String) Specifies the APM access key for tracing configuration.

* `tracing_sk` - (Required, String) Specifies the APM secret key for tracing configuration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the function URN.

## Import

Tracing configuration can be imported by `id` (also the functin URN), e.g.

```bash
$ terraform import huaweicloud_fgs_function_tracing_configuration.test <id>
```

Note that the imported state may not be identical to your resource definition, due to security reason.
The missing attribute is `tracing_sk`. It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the resource, or the script definition should be updated to align
with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_fgs_function_tracing_configuration" "test" {
  ...

  lifecycle {
    ignore_changes = [
      tracing_sk,
    ]
  }
}
```
