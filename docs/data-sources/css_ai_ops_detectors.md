---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_ai_ops_detectors"
description: |-
  Use this data source to query intelligent operation and maintenance check items.
---

# huaweicloud_css_ai_ops_detectors

Use this data source to query intelligent operation and maintenance check items.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_ai_ops_detectors" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `full_detection` - The full check items.
  The [full_detection](#full_detection_struct) structure is documented below.

* `unavailability_detection` - The cluster unavailability check items.
  The [unavailability_detection](#unavailability_detection_struct) structure is documented below.

<a name="full_detection_struct"></a>
The `full_detection` block supports:

* `id` - The check item ID.

* `name` - The check item name.

* `desc` - The check item description.

<a name="unavailability_detection_struct"></a>
The `unavailability_detection` block supports:

* `id` - The check item ID.

* `name` - The check item name.

* `desc` - The check item description.
