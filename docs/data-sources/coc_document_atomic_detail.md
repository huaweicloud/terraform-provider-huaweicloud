---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document_atomic_detail"
description: |-
  Use this data source to get the document atomic detail.
---

# huaweicloud_coc_document_atomic_detail

Use this data source to get the document atomic detail.

## Example Usage

```hcl
variable "atomic_unique_key" {}

data "huaweicloud_coc_document_atomic_detail" "test" {
  atomic_unique_key = var.atomic_unique_key
}
```

## Argument Reference

The following arguments are supported:

* `atomic_unique_key` - (Required, String) Specifies the unique identifier of an atomic capability.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `atomic_name_zh` - Indicates the Chinese name.

* `atomic_name_en` - Indicates the English name.

* `tags` - Indicates the tag information.

* `inputs` - Indicates the atomic capability input.

  The [inputs](#inputs_struct) structure is documented below.

* `outputs` - Indicates the atomic capability output.

  The [outputs](#outputs_struct) structure is documented below.

<a name="inputs_struct"></a>
The `inputs` block supports:

* `param_key` - Indicates the parameter variable name.

* `param_name_zh` - Indicates the Chinese name of the parameter.

* `param_name_en` - Indicates the English name of the parameter.

* `required` - Indicates whether the field is required.

* `param_type` - Indicates the parameter type: constant/dictionary.

* `min` - Indicates the minimum value.

* `max` - Indicates the maximum value.

* `min_len` - Indicates the minimum length.

* `max_len` - Indicates the maximum length.

* `pattern` - Indicates the regular restriction expression.

<a name="outputs_struct"></a>
The `outputs` block supports:

* `supported` - Indicates whether output is supported.

* `type` - Indicates the output type.
