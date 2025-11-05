---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_runtime_types"
description: |-
  Use this data source to query runtime type list within HuaweiCloud.
---

# huaweicloud_fgs_runtime_types

Use this data source to query runtime type list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_runtime_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the runtime types are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `types` - The list of runtime types.  
  The [types](#fgs_runtime_types_attr) structure is documented below.

<a name="fgs_runtime_types_attr"></a>
The `types` block supports:

* `type` - The type of the runtime.

* `display_name` - The display name of the runtime.
