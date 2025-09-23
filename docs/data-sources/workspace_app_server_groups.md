---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_groups"
description: |-
  Use this data source to get server group list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_server_groups

Use this data source to get server group list of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
variable "server_group_name" {}
        
data "huaweicloud_workspace_app_server_groups" "test" {
  server_group_name = var.server_group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_group_name` - (Optional, String) Specifies the name of the server group.

* `server_group_id` - (Optional, String) Specifies the ID of the server group.

* `app_type` - (Optional, String) Specifies the type of application group.
  + **SESSION_DESKTOP_APP**
  + **COMMON_APP**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `is_secondary_server_group` - (Optional, String) Specifies whether it is a secondary server group,
  default to **false**.

* `tags` - (Optional, String) Specifies the tag value to filter server groups.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `server_groups` - The list of server groups that matched filter parameters.  
  The [server_groups](#workspace_app_server_groups) structure is documented below.

<a name="workspace_app_server_groups"></a>
The `server_groups` block supports:

* `id` - The unique ID of the server group.

* `name` - The name of the server group.

* `description` - The description of the server group.

* `image_id` - The image ID used to create servers in this group.

* `os_type` - The type of the operating system. Currently, only **Windows** is supported.

* `product_id` - The product ID.

* `subnet_id` - The subnet ID for the network interface.

* `system_disk_type` - The type of the system disk.
  + **ESSD**
  + **SS**
  + **GPSSD**
  + **SAS**
  + **SATA**

* `system_disk_size` - The size of the system disk.

* `is_vdi` - Whether it is VDI single-session mode.

* `extra_session_type` - The paid session type.
  + **GPU**
  + **CPU**

* `extra_session_size` - The number of paid sessions.

* `app_type` - The application type of app server groups that matched filter parameters.
  + **SESSION_DESKTOP_APP**
  + **COMMON_APP**

* `create_time` - The creation time of the server group.

* `update_time` - The last update time of the server group.

* `storage_mount_policy` - The NAS storage directory mounting policy on APS.
  + **USER**
  + **SHARE**
  + **ANY**

* `enterprise_project_id` - The enterprise project ID.

* `primary_server_group_ids` - The list of primary server group IDs.

* `secondary_server_group_ids` - The list of secondary server group IDs.

* `server_group_status` - Whether the server group is enabled.

* `site_type` - The site type.
  + **CENTER**
  + **IES**

* `site_id` - The site ID.

* `app_server_flavor_count` - The total number of server configurations.

* `app_server_count` - The total number of servers.

* `app_group_count` - The total number of associated application groups.

* `image_name` - The image name of the group.

* `subnet_name` - The subnet name of the group.

* `ou_name` - The default organization name of the group.

* `product_info` - The product specification information.  
  The [product_info](#workspace_app_server_groups_product_info) structure is documented below.

* `scaling_policy` - The auto-scaling policy.  
  The [scaling_policy](#workspace_app_server_groups_scaling_policy) structure is documented below.

* `tags` - The tag information of app server groups that matched filter parameters.  
  The [tags](#workspace_app_server_groups_tags) structure is documented below.

<a name="workspace_app_server_groups_product_info"></a>
The `product_info` block supports:

* `product_id` - The ID of the product.

* `flavor_id` - The ID of the flavor.

* `type` - The type of the product.
  + **BASE**: Indicates the basic product package, which does not include any commercial software other than the
    operating system, and can only be used in private image scenarios.
  + **ADVANCED**: Indicates the advanced product package, which includes some commercial software in the package image.

* `architecture` - The architecture of the product. Currently, only **x86** is supported.

* `cpu` - The CPU information.

* `cpu_desc` - The CPU description.

* `memory` - The memory size in MB.

* `is_gpu` - Whether the flavor is GPU type.

* `system_disk_type` - The type of the system disk.

* `system_disk_size` - The size of the system disk.

* `gpu_desc` - The GPU description.

* `descriptions` - The product description.

* `charge_mode` - The charging mode.
  + **"1"**: Pre-paid.
  + **"0"**: Post-paid.

* `contain_data_disk` - Whether the package includes data disk.

* `resource_type` - The type of the resource.

* `cloud_service_type` - The type of the cloud service.

* `volume_product_type` - The type of the volume product.

* `sessions` - The maximum number of sessions supported by the package.

* `status` - The status of the product package in sales mode.
  + **normal**: Normal commercial use (Default).
  + **ababdon**: Discontinued (i.e., not displayed).
  + **sellout**: Sold out.
  + **obt**: Public testing.
  + **obs_sellout**: Public testing sold out.
  + **promotion**: Recommended (equivalent to normal, also commercial).

* `cond_operation_az` - The status of the product package in the availability zone.

* `sub_product_list` - The list of sub products.

* `domain_ids` - The list of domain IDs.

* `package_type` - The type of the package.
  + **general**: Indicates the general product package.
  + **dedicated**: Indicates the dedicated host product package.

* `expire_time` - The expiration time of the product package.

* `support_gpu_type` - The GPU type supported by the product package.

<a name="workspace_app_server_groups_scaling_policy"></a>
The `scaling_policy` block supports:

* `enable` - Whether to enable the policy.

* `max_scaling_amount` - The maximum scaling amount.

* `single_expansion_count` - The number of instances to add in a single scaling operation.

* `scaling_policy_by_session` - The session-based scaling policy.  
  The [scaling_policy_by_session](#workspace_app_server_groups_scaling_policy_by_session) structure is documented below.

<a name="workspace_app_server_groups_scaling_policy_by_session"></a>
The `scaling_policy_by_session` block supports:

* `session_usage_threshold` - The total session usage threshold of the group.

* `shrink_after_session_idle_minutes` - The release time for instances without session connections.

<a name="workspace_app_server_groups_tags"></a>
The `tags` block supports:

* `key` - The key of the tag. The maximum length is 128 Unicode characters.

* `value` - The value of the tag. The maximum length is 256 Unicode characters.
