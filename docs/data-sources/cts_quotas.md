---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_quotas"
description: |-
  Use this data source to get the list of CTS quotas.
---

# huaweicloud_cts_quotas

Use this data source to get the list of CTS quotas.

## Example Usage

```hcl
data "huaweicloud_cts_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The quota information.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `used` - The number of used resources.

* `quota` - The total number of resources.

* `type` - The resource type.
