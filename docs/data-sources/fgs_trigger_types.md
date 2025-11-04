---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_trigger_types"
description: |-
  Use this data source to query trigger type list within HuaweiCloud.
---

# huaweicloud_fgs_trigger_types

Use this data source to query trigger type list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_trigger_types" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the trigger types are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `types` - The list of trigger types.
