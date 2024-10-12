---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_server_flavors"
description: |-
  Use this data source to get available flavors of CPH server.
---

# huaweicloud_cph_server_flavors

Use this data source to get available flavors of CPH server.

## Example Usage

```hcl
data "huaweicloud_cph_server_flavors" "flavor" {
  type = "0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) The type of the CPH server flavor.  
  The options are as follows:
  + **0**: Cloud phone servers are designed for app hosting and multi-platform live streaming.
  + **1**: Cloud mobile gaming servers, GPU hardware acceleration and graphics interfaces
    allow you to run mobile games on the cloud.

* `vcpus` - (Optional, Int) The vcpus of the CPH server.

* `memory` - (Optional, Int) The ram of the CPH server in GB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of flavor detail.
  The [Flavors](#serverFlavors_Flavors) structure is documented below.

<a name="serverFlavors_Flavors"></a>
The `Flavors` block supports:

* `flavor_id` - The name of the flavor.

* `vcpus` - The vcpus of the CPH server.

* `memory` - The ram of the CPH server in GB.

* `type` - The type of the CPH server flavor.  
  The options are as follows:
  + **0**: Cloud phone servers are designed for app hosting and multi-platform live streaming.
  + **1**: Cloud mobile gaming servers, GPU hardware acceleration and graphics interfaces
    allow you to run mobile games on the cloud.

* `extend_spec` - The extended attribute description.
  The [ExtendSpec](#serverFlavors_FlavorsExtendSpec) structure is documented below.

<a name="serverFlavors_FlavorsExtendSpec"></a>
The `FlavorsExtendSpec` block supports:

* `vcpus` - The extended description of the vcpus.

* `memory` - The extended description of the ram.

* `disk` - The extended description of the disk.

* `network_interface` - The extended description of the network interface.

* `gpu` - The extended description of the gpu.

* `bms_flavor` - The extended description of the bms flavor.

* `gpu_count` - The gpu count.

* `numa_count` - The numa count.
