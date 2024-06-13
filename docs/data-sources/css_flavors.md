---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_flavors"
description: ""
---

# huaweicloud_css_flavors

Use this data source to get available flavors of HuaweiCloud CSS node instance.

## Example Usage

```hcl
data "huaweicloud_css_flavors" "test" {
  type    = "ess"
  version = "7.9.3"
  vcpus   = 4
  memory  = 32
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the CSS flavors. If omitted, the
  provider-level region will be used.

* `type` - (Optional, String) Specifies the node instance type. The options are `ess`, `ess-cold`, `ess-master`
 and `ess-client`.

* `version` - (Optional, String) Specifies the engine version. The options are `5.5.1`, `6.2.3`, `6.5.4`, `7.1.1`,
 `7.6.2` and `7.9.3`.

* `name` - (Optional, String) Specifies the name of the CSS flavor.

* `vcpus` - (Optional, Int) Specifies the number of vCPUs in the CSS flavor.

* `memory` - (Optional, Int) Specifies the memory size(GB) in the CSS flavor.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID in UUID format.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `name` - The name of the CSS flavor. It is referenced by `node_config.flavor` in `huaweicloud_css_cluster`.

* `id` - The ID of CSS flavor.

* `region` - The region where the node resides.

* `type` - The node instance type.

* `version` - The engine version.

* `vcpus` - The number of vCPUs.

* `memory` - The memory size in GB.

* `disk_range` - The disk capacity range of an instance, in GB.

* `availability_zones` - The valid availability zones for current flavor.
