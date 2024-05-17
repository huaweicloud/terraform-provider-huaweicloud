---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_advanced_query_schemas"
description: |-
  Use this data source to get the list of RMS advanced query schemas.
---

# huaweicloud_rms_advanced_query_schemas

Use this data source to get the list of RMS advanced query schemas.

## Example Usage

```hcl
data "huaweicloud_rms_advanced_query_schemas" "test" {
  type = "aad.instances"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Optional, String) Specifies the type of the schema.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schemas` - The list of schema.

  The [schemas](#schemas_struct) structure is documented below.

<a name="schemas_struct"></a>
The `schemas` block supports:

* `type` - The schema type.

* `schema` - The schema detail.
