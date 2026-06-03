---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_training_job_flavors"
description: |-
  Use this data source to get public resource pool flavor list supported by ModelArts training job within HuaweiCloud.
---

# huaweicloud_modelarts_training_job_flavors

Use this data source to get public resource pool flavor list supported by ModelArts training job within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_modelarts_training_job_flavors" "test" {}
```

### Filter by flavor type

```hcl
data "huaweicloud_modelarts_training_job_flavors" "test" {
  flavor_type = "CPU"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the training job flavors are located.  
  If omitted, the provider-level region will be used.

* `flavor_type` - (Optional, String) Specifies the type of the flavor.  
  The valid values are as follows:
  + **CPU**
  + **GPU**
  + **Ascend**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of training job flavors that match the filter parameters.  
  The [flavors](#modelarts_training_job_flavors_attr) structure is documented below.

<a name="modelarts_training_job_flavors_attr"></a>
The `flavors` block supports:

* `flavor_id` - The ID of the flavor.

* `flavor_name` - The name of the flavor.

* `flavor_type` - The type of the flavor.

* `billing` - The billing information of the flavor.  
  The [billing](#modelarts_training_job_flavors_billing) structure is documented below.

* `flavor_info` - The detailed information of the flavor.  
  The [flavor_info](#modelarts_training_job_flavors_flavor_info) structure is documented below.

* `attributes` - The other attributes of the flavor.

* `support_engines` - The engines supported by the flavor.

<a name="modelarts_training_job_flavors_billing"></a>
The `billing` block supports:

* `code` - The billing code.

* `unit_num` - The billing unit.

<a name="modelarts_training_job_flavors_flavor_info"></a>
The `flavor_info` block supports:

* `max_num` - The maximum number of nodes that can be selected.  
  `1` means that the flavor does not support distributed.

* `cpu` - The CPU information of the flavor.  
  The [cpu](#modelarts_training_job_flavors_cpu) structure is documented below.

* `gpu` - The GPU information of the flavor.  
  The [gpu](#modelarts_training_job_flavors_gpu) structure is documented below.

* `npu` - The Ascend information of the flavor.  
  The [npu](#modelarts_training_job_flavors_npu) structure is documented below.

* `memory` - The memory information of the flavor.  
  The [memory](#modelarts_training_job_flavors_memory) structure is documented below.

* `disk` - The disk information of the flavor.  
  The [disk](#modelarts_training_job_flavors_disk) structure is documented below.

<a name="modelarts_training_job_flavors_cpu"></a>
The `cpu` block supports:

* `arch` - The CPU architecture.

* `core_num` - The number of CPU cores.

<a name="modelarts_training_job_flavors_gpu"></a>
The `gpu` block supports:

* `unit_num` - The number of GPUs.

* `product_name` - The GPU product name.

* `memory` - The GPU memory.

<a name="modelarts_training_job_flavors_npu"></a>
The `npu` block supports:

* `unit_num` - The number of NPUs.

* `product_name` - The NPU product name.

* `memory` - The NPU memory.

<a name="modelarts_training_job_flavors_memory"></a>
The `memory` block supports:

* `size` - The memory size.

* `unit` - The memory unit.

<a name="modelarts_training_job_flavors_disk"></a>
The `disk` block supports:

* `size` - The disk size.

* `unit` - The disk unit.
