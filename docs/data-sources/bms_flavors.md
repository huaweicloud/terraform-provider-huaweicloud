---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_flavors"
description: ""
---

# huaweicloud_bms_flavors

Use this data source to get available BMS flavors.

## Example Usage

```hcl
data "huaweicloud_bms_flavors" "demo" {
  availability_zone = "cn-north-1a"
  vcpus             = 48
}

# Create BMS instance with the matched flavor
resource "huaweicloud_bms_instance" "instance" {
  flavor_id = data.huaweicloud_bms_flavors.demo.flavors[0].id

  # Other properties...
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the AZ name.

* `vcpus` - (Optional, Int) Specifies the number of vCPUs in the BMS flavor.

* `memory` - (Optional, Int) Specifies the memory size(GB) in the BMS flavor.

* `cpu_arch` - (Optional, String) Specifies the CPU architecture of the BMS flavor.
  The value can be x86_64 and aarch64, defaults to **x86_64**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `flavors` - Indicates the flavors information. Structure is documented below.

The `flavors` block contains:

* `id` - The id or name of the BMS flavor.
* `vcpus` - The number of vCPUs.
* `memory` - The memory size in GB.
* `cpu_arch` - The CPU architecture of the BMS flavor.
* `operation` - The operation status of the BMS flavor in an each AZs.
