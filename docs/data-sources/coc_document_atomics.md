---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document_atomics"
description: |-
  Use this data source to get the list of COC document atomics.
---

# huaweicloud_coc_document_atomics

Use this data source to get the list of COC document atomics.

## Example Usage

```hcl
data "huaweicloud_coc_document_atomics" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of atomic capabilities.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `atomic_unique_key` - Indicates the unique identifier of an atomic capability.

* `atomic_name_zh` - Indicates the atomic Chinese name.

* `atomic_name_en` - Indicates the atomic English name.

* `tags` - Indicates the tag information.
  The value can be **CLOUD_API**, **FUNCTION**.
