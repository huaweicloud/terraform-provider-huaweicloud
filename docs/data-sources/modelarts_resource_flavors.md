---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_resource_flavors"
description: ""
---

# huaweicloud_modelarts_resource_flavors

Use this data source to get resource flavors of ModelArts.

## Example Usage

```hcl
data "huaweicloud_modelarts_resource_flavors" "test" {
  type = "Dedicate"
  tag  ="os.modelarts/scope"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `tag` - (Optional, String) The tag key.

* `type` - (Optional, String) The type of resource flavor.  
  Value options are as follows:
    + **Dedicate**: physical resources.
    + **Logical**: logical resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of resource flavors.
  The [flavors](#ResourceFlavors_Flavors) structure is documented below.

<a name="ResourceFlavors_Flavors"></a>
The `flavors` block supports:

* `id` - Flavor ID.

* `tags` - The key/value pairs to associate with the flavor.

* `type` - The type of resource flavor.  
  Value options are as follows:
    + **Dedicate**: physical resources.
    + **Logical**: logical resources.

* `arch` - Computer architecture.  
  The value can be **x86** or **arm64**.

* `cpu` - Number of CPU cores.  

* `memory` - Memory size in GiB.  

* `gpu` - GPU information.
  The [gpu](#ResourceFlavors_FlavorsGpu) structure is documented below.

* `npu` - NPU information.
  The [npu](#ResourceFlavors_FlavorsNpu) structure is documented below.

* `volume` - The list of data disk information.
  The [volume](#ResourceFlavors_FlavorsVolume) structure is documented below.

* `billing_modes` - Billing mode supported by the flavor.  
  Value options are as follows:
    + **0**: pay-per-use.
    + **1**: yearly/monthly.

* `job_flavors` - Training job types supported by the resource flavor.  

* `az_status` - Sales status of a resource specification in each AZ. The value is (AZ, Status).  
  Status options are as follows:
    + **normal**: on-sales.
    + **soldout**: sold out.

<a name="ResourceFlavors_FlavorsGpu"></a>
The `gpu` block supports:

* `type` - GPU type.

* `size` - Number of GPUs.

<a name="ResourceFlavors_FlavorsNpu"></a>
The `npu` block supports:

* `type` - NPU type.

* `size` - Number of NPUs.

<a name="ResourceFlavors_FlavorsVolume"></a>
The `volume` block supports:

* `type` - Disk type.  
  Value options are as follows:
    + **SSD**: ultra-high I/O disk.
    + **GPSSD**: general-purpose SSD disk.
    + **SAS**: high I/O disk.
    + **SATA**: common disk.

* `size` - Disk size, in GiB.
