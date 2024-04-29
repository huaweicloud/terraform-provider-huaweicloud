---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_global_connection_bandwidth_line_levels"
description: ""
---

# huaweicloud_cc_global_connection_bandwidth_line_levels

Use this data source to get the list of CC line levels of global connection bandwidths.

## Example Usage

```hcl
variable "line_id" {}

data "huaweicloud_cc_global_connection_bandwidth_line_levels" "test" {
  line_id = var.line_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `line_id` - (Optional, String) Line ID.

* `local_area` - (Optional, String) Local access point code included in the line specification.

* `remote_area` - (Optional, String) Remote access point code included in the line specification.

* `levels` - (Optional, String) Line grade.
  + **Pt**: Platinum.
  + **Ag**: Silver.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `line_levels` - The line grade list.

  The [line_levels](#line_levels_struct) structure is documented below.

<a name="line_levels_struct"></a>
The `line_levels` block supports:

* `id` - The line ID.

* `created_at` - Time when the line was created.

* `updated_at` - Time when the line was updated.

* `local_area` - Local access point.

* `remote_area` - Remote access point.

* `levels` - Line grade.
