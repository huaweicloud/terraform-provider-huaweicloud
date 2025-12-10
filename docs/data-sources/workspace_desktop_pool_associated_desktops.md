---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_pool_associated_desktops"
description: |-
  Use this data source to query the associated desktops under the desktop pool within HuaweiCloud.
---

# huaweicloud_workspace_desktop_pool_associated_desktops

Use this data source to query the associated desktops under the desktop pool within HuaweiCloud.

## Example Usage

```hcl
variable "desktop_pool_id" {}

data "huaweicloud_workspace_desktop_pool_associated_desktops" "test" {
  pool_id = var.desktop_pool_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the associated desktops are located.

* `pool_id` - (Required, String) Specifies the ID of the desktop pool to which the associated desktops belong.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `desktops` - The list of associated desktops.  
  The [desktops](#workspace_desktop_pool_associated_desktops_attr) structure is documented below.

<a name="workspace_desktop_pool_associated_desktops_attr"></a>
The `desktops` block supports:

* `desktop_id` - The ID of the desktop.

* `computer_name` - The name of the desktop.

* `os_host_name` - The OS host name of the desktop.

* `ip_addresses` - The list of IP addresses of the desktop.

* `ipv4` - The IPv4 address of the desktop.

* `ipv6` - The IPv6 address of the desktop.

* `desktop_type` - The type of the desktop.

* `status` - The status of the desktop.

* `in_maintenance_mode` - Whether the desktop is in maintenance mode.

* `created` - The creation time of the desktop.

* `login_status` - The login status of the desktop.

* `product_id` - The product ID of the desktop.

* `root_volume` - The root volume information of the desktop.  
  The [root_volume](#workspace_desktop_pool_associated_desktops_volume) structure is documented below.

* `data_volumes` - The list of data volumes of the desktop.  
  The [data_volumes](#workspace_desktop_pool_associated_desktops_volume) structure is documented below.

* `availability_zone` - The availability zone of the desktop.

* `site_type` - The site type of the desktop.

* `site_name` - The site name of the desktop.

* `product` - The product information of the desktop.  
  The [product](#workspace_desktop_pool_associated_desktops_product) structure is documented below.

* `os_version` - The OS version of the desktop.

* `sid` - The SID of the desktop.

* `tags` - The tags of the desktop.

* `is_support_internet` - Whether the desktop supports internet access.

* `is_attaching_eip` - Whether the desktop is attaching an EIP.

* `attach_state` - The attach state of the desktop.

* `enterprise_project_id` - The enterprise project ID of the desktop.

* `subnet_id` - The subnet ID of the desktop.

* `bill_resource_id` - The billing resource ID of the desktop.

<a name="workspace_desktop_pool_associated_desktops_volume"></a>
The `root_volume` and `data_volumes` block supports:

* `type` - The type of the volume.

* `size` - The size of the volume in GB.

* `device` - The device name of the volume.

* `id` - The ID of the volume.

* `volume_id` - The volume ID.

* `bill_resource_id` - The billing resource ID of the volume.

* `create_time` - The creation time of the volume.

* `display_name` - The display name of the volume.

* `resource_spec_code` - The resource specification code of the volume.

<a name="workspace_desktop_pool_associated_desktops_product"></a>
The `product` block supports:

* `product_id` - The product ID.

* `flavor_id` - The flavor ID.

* `type` - The product type.

* `cpu` - The CPU specification.

* `memory` - The memory specification.

* `descriptions` - The product description.

* `charge_mode` - The charging mode.

* `architecture` - The architecture of the product.

* `is_gpu` - Whether the product is GPU type.

* `package_type` - The package type of the product.

* `system_disk_type` - The system disk type.

* `system_disk_size` - The system disk size.

* `contain_data_disk` - Whether the product contains data disk.

* `resource_type` - The resource type.

* `cloud_service_type` - The cloud service type.

* `volume_product_type` - The volume product type.

* `status` - The status of the product.
