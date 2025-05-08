---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_accelerators"
description: |-
  Use this data source to get the list of accelerators.
---

# huaweicloud_ga_accelerators

Use this data source to get the list of accelerators.

## Example Usage

```hcl
variable "accelerator_name" {}

data "huaweicloud_ga_accelerators" "test" {
  name = var.accelerator_name
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Optional, String) Specifies the ID of the accelerator.

* `name` - (Optional, String) Specifies the name of the accelerator.

* `status` - (Optional, String) Specifies the current status of the accelerator.
  The valid values are as follows:
  + **ACTIVE**: The status of the accelerator is normal operation.
  + **ERROR**: The status of the accelerator is error.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the accelerator
  belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accelerators` - The list of the accelerators.
  The [accelerators](#ga_accelerators) structure is documented below.

<a name="ga_accelerators"></a>
The `accelerators` block supports:

* `id` - The ID of the accelerator.

* `name` - The name of the accelerator.  

* `description` - The description of the accelerator.

* `status` - The current status of the accelerator.

* `ip_sets` - The IP information of the accelerator.
  The [ip_sets](#accelerator_ip_sets) structure is documented below.

* `flavor_id` - The ID of the flavor to which the accelerator belongs.

* `enterprise_project_id` - The ID of the enterprise project to which the accelerator belongs.

* `tags` - The key/value pairs to associate with the accelerator.

* `created_at` - The creation time of the accelerator.

* `updated_at` - The latest update time of the accelerator.

* `frozen_info` - The frozen details of cloud services or resources.
  The [frozen_info](#accelerators_frozen_info) structure is documented below.

<a name="accelerator_ip_sets"></a>
The `ip_sets` block supports:

* `ip_type` - The IP type of the accelerator.

* `ip_address` - The IP address of the accelerator.

* `area` - The acceleration zone of the accelerator.

<a name="accelerators_frozen_info"></a>
The `frozen_info` block supports:

* `status` - The status of a cloud service or resource.

* `effect` - The status of the resource after being forzen.

* `scene` - The service scenario.
