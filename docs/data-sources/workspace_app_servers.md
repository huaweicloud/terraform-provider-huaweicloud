---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_servers"
description: |-
  Use this data source to get server list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_servers

Use this data source to get server list of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
variable "server_group_id" {}

data "huaweicloud_workspace_app_servers" "test" {
  server_group_id = var.server_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_group_id` - (Optional, String) Specifies the ID of the server group.

* `server_name` - (Optional, String) Specifies the name of the server.

* `machine_name` - (Optional, String) Specifies the machine name of the server.

* `ip_addr` - (Optional, String) Specifies the IP address of the server.

* `server_id` - (Optional, String) Specifies the ID of the server.

* `maintain_status` - (Optional, Bool) Specifies whether the server is in maintenance status.
  + **true** : Instances in maintenance state.
  + **false**: Instances in non-maintenance state.

* `scaling_auto_create` - (Optional, Bool) Specifies whether the server is created by auto-scaling.
  + **true** : Created through elastic scaling.
  + **false**: Not created through elastic scaling.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `server_count` - The total number of servers.

* `items` - The list of servers.  
  The [items](#app_servers_items) structure is documented below.

<a name="app_servers_items"></a>
The `items` block supports:

* `id` - The ID of the server.

* `name` - The name of the server.

* `machine_name` - The machine name of the server group that matched filter parameters.

* `description` - The description of the server.

* `server_group_id` - The ID of the server group that matched filter parameters.

* `server_group_name` - The name of the server group.

* `status` - The status of the server.
  + **UNREGISTER**
  + **REGISTERED**
  + **MAINTAINING**
  + **FREEZE**
  + **STOPPED**
  + **NONE**

* `create_time` - The creation time of the server.

* `update_time` - The update time of the server.

* `image_id` - The ID of the image.

* `availability_zone` - The availability zone of the server.

* `domain` - The domain of the server.

* `ou_name` - The organization unit name.

* `sid` - The SID of the instance.

* `instance_id` - The ID of the instance.

* `os_version` - The version of the operating system.

* `os_type` - The type of the operating system. Currently, only **Windows** is supported.

* `order_id` - The ID of the order.

* `maintain_status` - The maintenance status of the server group that matched filter parameters.

* `scaling_auto_create` - Whether the server is created by auto-scaling.

* `job_id` - The ID of the last executed job.

* `job_type` - The type of the job.
  + **CREATE_SERVER**
  + **DELETE_SERVER**
  + **UPDATE_FREEZE_STATUS**
  + **CREATE_SERVER_IMAGE**
  + **REINSTALL_OS**
  + **CHANGE_SERVER_IMAGE**
  + **REJOIN_DOMAIN**
  + **MIGRATE_SERVER**
  + **UPGRADE_ACCESS_AGENT**
  + **UPDATE_SERVER_TSVI**
  + **SCHEDULED_TASK**
  + **COLLECT_HDA_LOG**
  + **COLLECT_APS_LOG**
  + **CREATE_SERVER_SNAPSHOT**
  + **DELETE_SERVER_SNAPSHOT**
  + **RESTORE_SERVER_SNAPSHOT**
  + **BATCH_INSTALL_APP**

* `job_status` - The status of the job.
  + **WAITING**
  + **RUNNING**
  + **SUCCESS**
  + **FAILED**

* `job_time` - The execution time of the last job.

* `resource_pool_id` - The ID of the resource pool.

* `resource_pool_type` - The type of the resource pool.
  + **private**: Private resource pool.
  + **public**: Public resource pool.

* `host_id` - The ID of the dedicated host.

* `session_count` - The number of sessions.

* `vm_status` - The steady state of a server, the stable state in which a certain operation is completed.
  + **BUILD**: Creating an APS instance, the state of the APS instance before it enters operation.
  + **BUILD_FAIL**: Failed to create an APS instance.
  + **REBOOT**: The instance is undergoing a reboot operation.
  + **HARD_REBOOT**: The instance is undergoing a forced reboot operation.
  + **REBUILD**: The instance is being rebuilt.
  + **REBUILD_FAIL**: The instance rebuilding failed.
  + **MIGRATING**: The instance is in the process of hot migration.
  + **RESIZE**: The instance has received a change request and is beginning the change operation.
  + **ACTIVE**: The instance is in a normal operating state.
  + **SHUTOFF**: The instance has been normally stopped.
  + **REVERT_RESIZE**: The instance is reverting the configuration of the changed specifications.
  + **VERIFY_RESIZE**: The instance is verifying the configuration after the change is completed.
  + **ERROR**: The instance is in an abnormal state.
  + **DELETING**: The instance is being deleted.
  + **FREEZE**: The instance is frozen.
  + **BUILD_IMAGE**: Creating an image of the instance.
  + **BUILD_SNAPSHOT**: Creating a snapshot of the instance.
  + **RESTORE_SNAPSHOT**: Restoring a snapshot of the instance.
  + **NULL**: Not set.

* `task_status` - The task status of the server.
  + **scheduling**: Instance is being scheduled during creation.
  + **block_device_mapping**: Instance is preparing disks during creation.
  + **networking**: Instance is preparing network during creation.
  + **spawning**: Instance is being internally created.
  + **rebooting**: Instance is rebooting.
  + **reboot_pending**: Instance is rebooting, restart command is being issued.
  + **reboot_started**: Instance has started internal reboot.
  + **rebooting_hard**: Instance is performing a hard reboot.
  + **reboot_pending_hard**: Instance is performing a hard reboot, restart command is being issued.
  + **reboot_started_hard**: Instance has started internal hard reboot.
  + **rebuilding**: Instance is being rebuilt.
  + **rebuild_fail**: Instance rebuild failed.
  + **updating_tsvi**: Instance is updating virtual session IP.
  + **updating_tsvi_failed**: Instance virtual session IP update failed.
  + **rebuild_block_device_mapping**: Instance is preparing disks during rebuild.
  + **rebuild_spawning**: Instance is internally rebuilding.
  + **migrating**: Instance is hot migrating.
  + **resize_prep**: Instance is in preparation phase of resizing.
  + **resize_migrating**: Instance is in migration phase of resizing.
  + **resize_migrated**: Instance has completed migration during resizing.
  + **resize_finish**: Instance is finalizing resizing.
  + **resize_reverting**: Instance is reverting resizing changes.
  + **powering-off**: Instance is shutting down.
  + **powering-on**: Instance is starting up.
  + **deleting**: Instance is being deleted.
  + **source_locking**: Resource is being locked.
  + **rejoining_domain**: Instance is rejoining domain.
  + **delete_failed**: Instance deletion failed.
  + **upgrading_access_agent**: Instance is upgrading AccessAgent.
  + **upgrad_access_agent_fail**: Instance AccessAgent upgrade failed.
  + **upgrad_access_agent_success**: Instance AccessAgent upgrade succeeded.
  + **updating_sid**: Instance is waiting to update SID during creation.
  + **migrate_failed**: Instance migration failed.
  + **build_image**: Image is being generated.
  + **build_snapshot**: Snapshot is being generated.
  + **restore_snapshot**: Snapshot is being restored.
  + **installing_app**: Application is being silently installed.
  + **install_app_failed**: Application installation failed.
  + **null**: Not set.

* `enterprise_project_id` - The enterprise project ID of the server.

* `metadata` - The metadata of the server.
  + **charging_mode**: The billing type for cloud servers.
    - **0**: Post-paid.
    - **1**: Pre-paid.
    - **2**: Bidding instance billing.
  + **metering.order_id**: The order ID corresponding to the cloud server charged on a yearly/monthly basis.
  + **metering.product_id**: The product ID corresponding to the cloud server charged on a yearly/monthly basis.
  + **vpc_id**: The VPC ID of the service.
  + **EcmResStatus**: The frozen state of the cloud server.
    - **Normal**: The cloud server is in a normal state (not frozen).
    - **Freeze**: The cloud server has been frozen.
  + **metering.image_id**: The image ID corresponding to the cloud server operating system.
  + **metering.imagetype**: The image type corresponding to the cloud server operating system.
    - **gold**
    - **private**
    - **shared**
  + **metering.resourcespeccode**: The resource specifications corresponding to the cloud server.
  + **image_name**: The image name corresponding to the cloud server.
  + **os_bit**: The number of bits in the operating system is usually set to **32** or **64**.
  + **lockCheckEndpoint**: Callback URL, used to check whether the locking of the ECS instance is effective.
  + **lockSource**: Which service does the ECS instance come from. It is order locking if the value is **ORDER**.
  + **lockSourceId**: Which ID does the locking of the ECS instance come from. When `lockSource` is **ORDER**,
    `lockSourceId` is the order ID.
  + **lockScene**: The locking types of ECS instance.
    - **TO_PERIOD_LOCK**: Transfer post-paid to pre-paid.
  + **virtual_env_type**: The virtual environment type.
    - Create a virtual machine from an IOS image. It has `"virtual_env_type": "IsoImage"` property.
    - Create a virtual machine with a non-IOS image. Virtual machines created after version **19.5.0** will not
      have the `virtual_env_type` property added. However, virtual machines created with versions prior to this may
      return the property `"virtual_env_type": "FusionCompute"`.
  + **metering.resourcetype**: The resource type corresponding to the ECS instance.
  + **os_type**: Operating system type.
    - **Linux**
    - **Windows**
  + **cascaded.instance_extrainfo**: The internal virtual machine extension information of the system.
  + **__support_agent_list**: The list of agents supported by ECS instance.
    - **hss**
    - **ces**
  + **agency_name**: The name of the agency.

  For more details, you can [submit a service ticket](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/ticket_02_0003.html).
  
* `flavor` - The flavor information of the server.  
  The [flavor](#app_servers_flavor) structure is documented below.

* `product_info` - The product information of the server.  
  The [product_info](#app_servers_product_info) structure is documented below.

* `freeze` - The freeze information of the server.  
  The [freeze](#app_servers_freeze) structure is documented below.

* `host_address` - The network information of the server.  
  The [host_address](#app_servers_host_address) structure is documented below.

* `tags` - The tags of the server.  
  The [tags](#app_servers_tags) structure is documented below.

<a name="app_servers_flavor"></a>
The `flavor` block supports:

* `id` - The ID of the flavor.

* `links` - The quick link information for relevant tags corresponding to server specifications.  
  The [links](#app_servers_links) structure is documented below.

<a name="app_servers_links"></a>
The `links` block supports:

* `rel` - The shortcut link tag name.

* `href` - The corresponding shortcut link.

<a name="app_servers_product_info"></a>
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

<a name="app_servers_freeze"></a>
The `freeze` block supports:

* `effect` - The effect of the freeze operation.
  + **1**: (Implement/Remove) Freeze can be released. The resource can be manually deleted or released after freezing.
  + **2**: (Implement/Remove) Freeze cannot be released. The resource cannot be manually deleted or released after
    freezing, and cannot be changed, equivalent to the resource being sealed. After thawing, customers can delete
    and modify the data.
  + **3**: (Implement/Remove) Non-renewable after freezing. The resource cannot initiate renewal operations after
    freezing; this can be done after thawing.

-> **NOTE:** The `effect` field is used in conjunction with the above `status` field (1 for freeze, 0 for thaw) to
indicate the freezing effect initiated by the freeze/thaw command. If empty, it defaults to `effect=1` (the cloud
service needs to be compatible and handle it as `effect=1` by default).

-> **NOTE:**: Cloud services restrict operations/APIs based on the `status` and `effect`. The following `scene` field
is used for customer experience tips on the Console page and API error codes, not for restricting
cloud service operations/APIs.

* `scene` - The scene of the service status update. Default to **ARREAR**.
  + **ARREAR**: Arrears scenario; for normal operational business scenarios, including expired periodic resources,
    failed on-demand resource billing.
  + **POLICE**: Police freeze scenario.
  + **ILLEGAL**: Illegal freeze scenario.
  + **VERIFY**: Customer under-authentication freeze scenario.
  + **PARTNER**: Partner freeze (partner freezes sub-customer resources).

<a name="app_servers_host_address"></a>
The `host_address` block supports:

* `addr` - The IP address.

* `version` - The IP address version.
  + **4**: IPV4.
  + **6**: IPV6.

* `mac_addr` - The MAC address.

* `type` - The IP address allocation type. Strings are case-insensitive formats.
  + **fixed**: Represents a private IP address.
  + **floating**: Represents a floating IP address.

* `port_id` - The port ID of the IP address.

* `vpc_id` - The ID of the VPC.

* `subnet_id` - The ID of the subnet.

* `tenant_type` - The type of the tenant.
  + **tenant**
  + **resource_tenant**

<a name="app_servers_tags"></a>
The `tags` block supports:

* `key` - The key of the tag. The maximum length is 128 Unicode characters.

* `value` - The value of the tag. The maximum length is 256 Unicode characters.
