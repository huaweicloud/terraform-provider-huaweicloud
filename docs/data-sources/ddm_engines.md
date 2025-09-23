---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_engines"
description: ""
---

# huaweicloud_ddm_engines

Use this data source to get the list of DDM engines.

## Example Usage

```hcl
data "huaweicloud_ddm_engines" "test" {
  version = "3.0.8.5"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `version` - (Optional, String) Specifies the engine version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `engines` - Indicates the list of DDM engine.
  The [Engine](#DdmEngines_Engine) structure is documented below.

<a name="DdmEngines_Engine"></a>
The `Engine` block supports:

* `id` - Indicates the ID of the engine.

* `version` - Indicates the engine version.
