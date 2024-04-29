---
subcategory: "Cloud Data Migration (CDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_flavors"
description: ""
---

# huaweicloud_cdm_flavors

Use this data source to get available HuaweiCloud CDM flavors.

## Example Usage

```hcl
data "huaweicloud_cdm_flavors" "flavor" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the CDM flavors.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `version` - The version of the CDM cluster.

* `flavors` - Indicates the flavors information. Structure is documented below.
  The [flavors](#block-flavors) structure is documented below.

<a name="block-flavors"></a>
The `flavors` block supports:

* `id` - The ID of the CDM flavor.

* `name` - The name of the CDM flavor. Format is `cdm.<flavor_type>`.

* `cpu` - The numbers of CDM cluster vCPUs.

* `memory` - The memory size in GB.
