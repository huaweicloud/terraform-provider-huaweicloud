---
subcategory: "Cloud Data Migration (CDM)"
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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `version` - The version of the CDM cluster.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the CDM flavor.
* `id` - The id of the CDM flavor.
