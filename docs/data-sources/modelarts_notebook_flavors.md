---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_notebook_flavors"
description: ""
---

# huaweicloud_modelarts_notebook_flavors

Use this data source to get available flavors of ModelArts notebook.

## Example Usage

```hcl
data "huaweicloud_modelarts_notebook_flavors" "flavors" {
  category = "CPU"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `category` - (Optional, String) Processor type. The valid values are: **CPU**, **GPU**, **ASCEND**.  

* `type` - (Optional, String) Cluster type.  
  The options are as follows:
    - **MANAGED**: Public cluster.
    - **DEDICATED**: Dedicated cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of flavors.
  The [Flavors](#notebookFlavors_Flavors) structure is documented below.

<a name="notebookFlavors_Flavors"></a>
The `Flavors` block supports:

* `id` - The ID of the flavor.

* `name` - The name of the flavor.

* `arch` - Architecture type. The valid values are **X86_64** and **AARCH64**.

* `category` - Processor type. The valid values are: **CPU**, **GPU**, **ASCEND**.

* `description` - Specification description.

* `feature` - Flavor type.  
  The options are as follows:
    - **DEFAULT**: CodeLab.
    - **NOTEBOOK**: notebook.

* `memory` - Memory size, in KB.

* `vcpus` - Number of vCPUs.

* `free` - Free flavor or not.

* `sold_out` - Whether resources are sold out.

* `billing` - Billing information.
  The [billing](#notebookFlavors_FlavorsBilling) structure is documented below.

* `gpu` - GPU information.
  The [gpu](#notebookFlavors_FlavorsGpu) structure is documented below.

<a name="notebookFlavors_FlavorsBilling"></a>
The `billing` block supports:

* `code` - Billing code.

* `unit_num` - Billing unit.

<a name="notebookFlavors_FlavorsGpu"></a>
The `gpu` block supports:

* `gpu` - Number of GPUs.

* `gpu_memory` - GPU memory, in GB.

* `type` - GPU type.
