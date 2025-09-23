---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_flavors"
description: |-
  Use this data source to get flavor list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_flavors

Use this data source to get flavor list of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
variable "product_id" {}

data "huaweicloud_workspace_app_flavors" "test" {
  product_id = var.product_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `product_id` - (Optional, String) Specifies the product ID used to filter the app flavor list.

* `flavor_id` - (Optional, String) Specifies the flavor ID used to filter the app flavor list.

* `availability_zone` - (Optional, String) Specifies the availability zone used to filter the app flavor list.

* `os_type` - (Optional, String) Specifies the operating system type used to filter the app flavor list. The valid
  value is **Windows**.

* `charge_mode` - (Optional, String) Specifies the charge mode used to filter the app flavor list.
  + **1**. The billing method is pre-paid.
  + **0**. The billing method is post-paid.

* `architecture` - (Optional, String) Specifies the architecture type used to filter the app flavor list. The valid
  values are **x86** and **arm**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of app flavors that matched filter parameters.  
  The [flavors](#workspace_app_flavors_flavors) structure is documented below.

<a name="workspace_app_flavors_flavors"></a>
The `flavors` block supports:

* `product_id` - The product ID of app flavors that matched filter parameters.

* `id` - The flavor ID of app flavors that matched filter parameters.

* `architecture` - The flavor architecture of app flavors that matched filter parameters.

* `type` - The flavor type.

* `cpu` - The CPU core count.

* `memory` - The memory size in MB.

* `is_gpu` - Whether the flavor is GPU type.

* `system_disk_type` - The system disk type.

* `system_disk_size` - The system disk size.

* `descriptions` - The flavor description.

* `charge_mode` - The charge mode of app flavors that matched filter parameters.

* `contain_data_disk` - Whether the flavor includes data disk.

* `resource_type` - The resource type.

* `cloud_service_type` - The cloud service type.

* `volume_product_type` - The volume product type.

* `sessions` - The maximum number of sessions supported by the flavor.

* `status` - The flavor status.

* `cond_operation_az` - The flavor status in availability zones.

* `domain_ids` - The domain IDs that the flavor belongs to.

* `package_type` - The package type.
