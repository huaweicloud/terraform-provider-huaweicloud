---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_asset"
description: |-
  Manages a SecMaster asset resource within HuaweiCloud.
---

# huaweicloud_secmaster_asset

Manages a SecMaster asset resource within HuaweiCloud.
This resource allows you to manage various types of assets including ECS, VPC, EIP, RDS, and more in the SecMaster service.

## Example Usage

```hcl
variable "asset_id" {}
variable "workspace_id" {}

resource "huaweicloud_secmaster_asset" "test" {
  asset_id     = var.asset_id
  workspace_id = var.workspace_id

  data_object {
    checksum           = "testchecksum"
    created            = "2025-08-25T00:30:54.820Z+0800"
    id                 = var.asset_id
    level              = 1
    name               = "tf_test_zoj7e"
    provider           = "rds"
    provisioning_state = "temp-state"
    type               = "instances"

    department {
      id   = "XXX"
      name = "test-department-name"
    }

    environment {
      domain_id   = "XXX"
      ep_id       = "0"
      ep_name     = "default"
      idc_id      = "test-idc-id"
      idc_name    = "test-idc-name"
      project_id  = "XXX"
      region_id   = "XXX"
      vendor_name = "test-vendor-name"
      vendor_type = "CloudService"
    }

    governance_user {
      name = "llrds"
      type = "test-governance-user-type"
    }

    properties {
      hwc_rds {
        alias                 = "test-alias"
        associated_with_ddm   = false
        backup_used_space     = 0
        cpu                   = "2"
        created               = "2025-08-24T15:33:15+0000"
        db_user_name          = "root"
        disk_encryption_id    = "test-disk-encryption-id"
        enable_ssl            = false
        enterprise_project_id = "XXX"
        expiration_time       = "2025-08-24T15:33:15Z"
        flavor_ref            = "rds.mysql.x1.large.2"
        id                    = "XXX"
        maintenance_window    = "18:00-22:00"
        max_iops              = 0
        mem                   = "4"
        name                  = "tf_test_zoj7e"
        port                  = 3306
        private_dns_names = [
          "XXX",
        ]
        private_ips = [
          "XXX",
        ]
        project_id       = "XXX"
        protected_status = "CLOSE"
        public_ips = [
          "XXX",
        ]
        read_only_by_user  = false
        region             = "cn-north-4"
        security_group_id  = "XXX"
        status             = "ACTIVE"
        storage_used_space = 0
        subnet_id          = "XXX"
        switch_strategy    = "reliability"
        time_zone          = "UTC"
        type               = "Single"
        updated            = "2025-08-24T15:50:01+0000"
        vpc_id             = "XXX"

        backup_strategy {
          keep_days  = 7
          start_time = "02:00-03:00"
        }

        datastore {
          complete_version = "8.0.28.231003"
          type             = "MySQL"
          version          = "8.0"
        }

        ha {
          replication_mode = "async"
        }

        nodes {
          availability_zone = "cn-north-4a"
          id                = "XXX"
          name              = "tf_test_zoj7e_node0"
          role              = "master"
          status            = "ACTIVE"
        }

        related_instance {
          id   = "XXX"
          type = "replica_of"
        }

        tags {
          key = "test-key"
          values = [
            "test-value1",
            "test-value2",
          ]
        }

        volume {
          size = 40
          type = "CLOUDSSD"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the asset belongs.

* `asset_id` - (Required, String, NonUpdatable) Specifies the ID of the asset.

* `data_object` - (Required, List) Specifies the details of the asset.
  The [data_object](#data_object_struct) structure is documented below.

<a name="data_object_struct"></a>
The `data_object` block supports:

* `id` - (Optional, String) Specifies the ID of the asset.

* `name` - (Required, String) Specifies the name of the asset.

* `provider` - (Required, String) Specifies the provider of the asset.

* `type` - (Required, String) Specifies the type of the asset.
  Valid values are **ECS**, **VPC**, **EVS**, **IP**, **URL** and so on.

* `environment` - (Required, List) Specifies the environment of the asset.
  The [environment](#environment_struct) structure is documented below.

* `properties` - (Required, List) Specifies the properties of the asset.
  The [properties](#properties_struct) structure is documented below.

* `checksum` - (Optional, String) Specifies the checksum of the asset.

* `created` - (Optional, String) Specifies the creation time of the asset.

* `provisioning_state` - (Optional, String) Specifies the provisioning state of the asset.

* `department` - (Optional, List) Specifies the department of the asset.
  The [department](#department_struct) structure is documented below.

* `governance_user` - (Optional, List) Specifies the governance user of the asset.
  The [governance_user](#governance_user_struct) structure is documented below.

* `level` - (Optional, Int) Specifies the level of the asset. Valid values are:
  + **0**: Testing.
  + **1**: Normal.
  + **2**: Important.

<a name="environment_struct"></a>
The `environment` block supports:

* `vendor_type` - (Required, String) Specifies the environment vendor type.

* `domain_id` - (Required, String) Specifies the domain ID.

* `vendor_name` - (Required, String) Specifies the vendor name.

* `idc_name` - (Required, String) Specifies the IDC name.

* `region_id` - (Optional, String) Specifies the region ID.

* `project_id` - (Optional, String) Specifies the project ID.

* `ep_id` - (Optional, String) Specifies the enterprise project ID.

* `ep_name` - (Optional, String) Specifies the enterprise project name.

* `idc_id` - (Optional, String) Specifies the IDC ID.

<a name="properties_struct"></a>
The `properties` block supports:

* `hwc_ecs` - (Optional, List) Specifies the details of the ECS.
  The [hwc_ecs](#hwc_ecs_struct) structure is documented below.

* `hwc_eip` - (Optional, List) Specifies the details of the EIP.
  The [hwc_eip](#hwc_eip_struct) structure is documented below.

* `hwc_vpc` - (Optional, List) Specifies the details of the VPC.
  The [hwc_vpc](#hwc_vpc_struct) structure is documented below.

* `hwc_subnet` - (Optional, List) Specifies the details of the subnet.
  The [hwc_subnet](#hwc_subnet_struct) structure is documented below.

* `hwc_rds` - (Optional, List) Specifies the details of the RDS.
  The [hwc_rds](#hwc_rds_struct) structure is documented below.

* `hwc_domain` - (Optional, List) Specifies the details of the domain.
  The [hwc_domain](#hwc_domain_struct) structure is documented below.

* `website` - (Optional, List) Specifies the details of the website.
  The [website](#website_struct) structure is documented below.

* `oca_ip` - (Optional, List) Specifies the details of the cloud asset IP.
  The [oca_ip](#oca_ip_struct) structure is documented below.

<a name="hwc_ecs_struct"></a>
The `hwc_ecs` block supports:

* `id` - (Required, String) Specifies the ID of the ECS.

* `name` - (Required, String) Specifies the name of the ECS.

* `protected_status` - (Required, String) Specifies the protection status of the ECS.
  Valid values are **OPEN** and **CLOSE**.

* `description` - (Required, String) Specifies the description of the ECS.

* `status` - (Required, String) Specifies the status of the ECS.
  Valid values are **ACTIVE**, **BUILD**, **ERROR**, **HARD_REBOOT**, **MIGRATING**, **REBOOT**, **REBUILD**,
  **RESIZE**, **REVERT_RESIZE**, **SHUTOFF**, **VERIFY_RESIZE**, and **DELETED**.

* `locked` - (Required, Bool) Specifies whether the ECS is locked. Valid values are:
  + **true**: The ECS is locked.
  + **false**: The ECS is not locked.

* `user_id` - (Required, String) Specifies the user ID of the ECS.

* `project_id` - (Required, String) Specifies the project ID of the ECS.

* `host_id` - (Required, String) Specifies the host ID of the ECS.

* `host_name` - (Required, String) Specifies the host name of the ECS.

* `host_status` - (Required, String) Specifies the host status of the ECS. Valid values are:
  + **UP**: The host is normal.
  + **UNKNOWN**: The host status is unknown.
  + **DOWN**: The host is down.
  + **MAINTENANCE**: The host is under maintenance.

* `addresses` - (Required, List) Specifies the IP addresses of the ECS.
  The [addresses](#hwc_ecs_addresses_struct) structure is documented below.

* `security_groups` - (Required, List) Specifies the security groups of the ECS.
  The [security_groups](#hwc_ecs_security_groups_struct) structure is documented below.

* `availability_zone` - (Required, String) Specifies the availability zone of the ECS.

* `volumes_attached` - (Required, List) Specifies the volumes attached to the ECS.
  The [volumes_attached](#hwc_ecs_volumes_attached_struct) structure is documented below.

* `metadata` - (Required, List) Specifies the metadata of the ECS.
  The [metadata](#hwc_ecs_metadata_struct) structure is documented below.

* `updated` - (Optional, String) Specifies the update time of the ECS. Time format example: **2019-05-22T03:30:52Z**.

* `created` - (Optional, String) Specifies the creation time of the ECS. Time format example: **2019-05-22T03:19:19Z**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the ECS.

* `flavor` - (Optional, List) Specifies the flavor of the ECS.
  The [flavor](#hwc_ecs_flavor_struct) structure is documented below.

* `key_name` - (Optional, String) Specifies the key name of the ECS.

* `scheduler_hints` - (Optional, List) Specifies the scheduler hints of the ECS.
  The [scheduler_hints](#hwc_ecs_scheduler_hints_struct) structure is documented below.

* `hypervisor` - (Optional, List) Specifies the virtualization information of the ECS.
  The [hypervisor](#hwc_ecs_hypervisor_struct) structure is documented below.

<a name="hwc_ecs_volumes_attached_struct"></a>
The `volumes_attached` block supports:

* `id` - (Required, String) Specifies the ID of the disk.

* `delete_on_termination` - (Optional, String) Specifies whether to delete the disk when deleting the ECS.

* `boot_index` - (Optional, String) Specifies the boot order of the disk. Valid values:
  + **0**: The disk is system disk.
  + **1**: The disk is data disk.

* `device` - (Optional, String) Specifies the mount point of the disk.

<a name="hwc_ecs_metadata_struct"></a>
The `metadata` block supports:

* `image_id` - (Optional, String) Specifies the ID of the image.

* `image_type` - (Optional, String) Specifies the image type of the ECS. Valid values are:
  + **gold**: The gold image.
  + **private**: The private image.
  + **shared**: The shared image.

* `image_name` - (Optional, String) Specifies the image name of the ECS.

* `os_bit` - (Optional, String) Specifies the OS bit of the ECS. Valid values are:
  + **32**: The 32-bit OS.
  + **64**: The 64-bit OS.

* `os_type` - (Optional, String) Specifies the OS type of the ECS. Valid values are:
  + **Windows**: The Windows OS.
  + **Linux**: The Linux OS.

* `vpc_id` - (Optional, String) Specifies the VPC ID of the ECS.

* `resource_spec_code` - (Optional, String) Specifies the resource spec code of the ECS.

* `resource_type` - (Optional, String) Specifies the resource type of the ECS.

* `agency_name` - (Optional, String) Specifies the agency name of the ECS.

<a name="hwc_ecs_scheduler_hints_struct"></a>
The `scheduler_hints` block supports:

* `group` - (Optional, List) Specifies the cloud server group ID.

* `tenancy` - (Optional, List) Specifies the tenancy of the ECS. Valid values are:
  + **dedicated**: The dedicated ECS.
  + **shared**: The shared ECS.

* `dedicated_host_id` - (Optional, List) Specifies the dedicated host ID. This field is valid only when the ECS is
  a dedicated host.

<a name="hwc_ecs_hypervisor_struct"></a>
The `hypervisor` block supports:

* `hypervisor_type` - (Optional, String) Specifies the virtualization type.

<a name="hwc_ecs_addresses_struct"></a>
The `addresses` block supports:

* `addr` - (Required, String) Specifies the IP address.

* `version` - (Required, String) Specifies the IP version, **4** means IPv4, **6** means IPv6.

* `type` - (Required, String) Specifies the type of the IP address. Valid values:
  + **fixed**: The fixed IP address.
  + **floating**: The floating IP address.

* `mac_addr` - (Required, String) Specifies the MAC address.

* `port_id` - (Required, String) Specifies the ID of the port.

* `vpc_id` - (Required, String) Specifies the ID of the VPC.

<a name="hwc_ecs_flavor_struct"></a>
The `flavor` block supports:

* `id` - (Required, String) Specifies the ID of the flavor.

* `name` - (Required, String) Specifies the name of the flavor.

* `disk` - (Optional, String) Specifies the disk size of the flavor.

* `vcpus` - (Optional, String) Specifies the number of vCPUs.

* `ram` - (Optional, String) Specifies the memory size in MB.

<a name="hwc_ecs_security_groups_struct"></a>
The `security_groups` block supports:

* `id` - (Required, String) Specifies the ID of the security group.

* `name` - (Optional, String) Specifies the name of the security group.

<a name="hwc_eip_struct"></a>
The `hwc_eip` block supports:

* `id` - (Required, String) Specifies the ID of the EIP.

* `alias` - (Required, String) Specifies the name of the EIP.

* `protected_status` - (Required, String) Specifies the protection status of the EIP.
  Valid values are **OPEN** and **CLOSE**.

* `project_id` - (Required, String) Specifies the project ID of the EIP.

* `enterprise_project_id` - (Required, String) Specifies the enterprise project ID of the EIP.

* `ip_version` - (Required, Int) Specifies the IP version information. Valid values are `4` and `6`.

* `status` - (Required, String) Specifies the status of the EIP. Valid values are:
  + **FREEZED**: The EIP is frozen.
  + **BIND_ERROR**: The binding of the EIP fails.
  + **BINDING**: The EIP is binding.
  + **PENDING_DELETE**: The EIP is deleting.
  + **PENDING_CREATE**: The EIP is creating.
  + **NOTIFYING**: The EIP is notifying.
  + **NOTIFY_DELETE**: The EIP is notifying deleting.
  + **PENDING_UPDATE**: The EIP is updating.
  + **DOWN**: The EIP is down.
  + **ACTIVE**: The EIP is active.
  + **ELB**: The EIP is binding ELB.
  + **VPN**: The EIP is binding VPN.
  + **ERROR**: The EIP is error.

* `public_ip_address` - (Optional, String) Specifies the public IP address of the EIP.

* `public_ipv6_address` - (Optional, String) Specifies the public IPv6 address of the EIP.

* `publicip_pool_name` - (Optional, String) Specifies the public IP pool name of the EIP.

* `publicip_pool_id` - (Optional, String) Specifies the public IP pool ID of the EIP.

* `description` - (Optional, String) Specifies the description of the EIP.

* `tags` - (Optional, List) Specifies the tags of the EIP.

* `type` - (Optional, String) Specifies the type of the EIP. Valid values are **EIP**, **DUALSTACK**, and **DUALSTACK_SUBNET**.

* `vnic` - (Optional, List) Specifies the VNIC information of the EIP.
  The [vnic](#hwc_eip_vnic_struct) structure is documented below.

* `bandwidth` - (Optional, List) Specifies the bandwidth information of the EIP.
  The [bandwidth](#hwc_eip_bandwidth_struct) structure is documented below.

* `lock_status` - (Optional, String) Specifies the freeze status of the public IP. Valid values are **police** and **locked**.

* `associate_instance_type` - (Optional, String) Specifies the instance type of the public IP.
  Valid values are **PORT**, **NATGW**, **ELB**, **ELBV1**, **VPN**, and **null**.

* `associate_instance_id` - (Optional, String) Specifies the instance ID of the public IP.

* `allow_share_bandwidth_types` - (Optional, List) Specifies the list of shared bandwidth types that the public IP can join.
  If it is an empty list, it means that the public IP cannot join any shared bandwidth. Constraint: The public IP can
  only join the shared bandwidth with the same bandwidth type.

* `created_at` - (Optional, String) Specifies the creation UTC time of the public IP.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `updated_at` - (Optional, String) Specifies the update UTC time of the public IP.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `public_border_group` - (Optional, String) Specifies the center site asset or edge site asset. Value range: center,
  edge site name.

<a name="hwc_eip_vnic_struct"></a>
The `vnic` block supports:

* `private_ip_address` - (Optional, String) Specifies the private IP address of the public IP.

* `device_id` - (Optional, String) Specifies the device ID of the public IP.

* `device_owner` - (Optional, String) Specifies the device owner of the public IP. Valid values are **network:dhcp**,
  **network:VIP_PORT**, **network:router_interface_distributed**, and **network:router_centralized_snat**.

* `vpc_id` - (Optional, String) Specifies the virtual private cloud ID of the public IP.

* `port_id` - (Optional, String) Specifies the port ID of the public IP.

* `port_profile` - (Optional, String) Specifies the port profile information of the public IP.

* `mac` - (Optional, String) Specifies the MAC address of the public IP.

* `vtep` - (Optional, String) Specifies the VTEP IP of the public IP.

* `vni` - (Optional, String) Specifies the VXLAN ID of the public IP.

* `instance_id` - (Optional, String) Specifies the instance ID of the public IP.

* `instance_type` - (Optional, String) Specifies the instance type of the public IP.

<a name="hwc_eip_bandwidth_struct"></a>
The `bandwidth` block supports:

* `id` - (Required, String) Specifies the ID of the bandwidth.

* `name` - (Optional, String) Specifies the name of the bandwidth.

* `size` - (Optional, Int) Specifies the size of the bandwidth in Mbps. Ranges from `5` to `2000`.

* `share_type` - (Optional, String) Specifies the bandwidth type. Valid values are **PER** and **WHOLE**.

<a name="hwc_vpc_struct"></a>
The `hwc_vpc` block supports:

* `id` - (Required, String) Specifies the ID of the VPC.

* `name` - (Required, String) Specifies the name of the VPC.

* `protected_status` - (Required, String) Specifies the protected status of the VPC. Valid values are **OPEN** and **CLOSE**.

* `status` - (Required, String) Specifies the status of the VPC. Valid values are **PENDING** and **ACTIVE**.

* `project_id` - (Required, String) Specifies the project ID of the VPC.

* `created_at` - (Required, String) Specifies the created time of the VPC.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `updated_at` - (Required, String) Specifies the updated time of the VPC.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `description` - (Optional, String) Specifies the description of the VPC.

* `cidr` - (Optional, String) Specifies the cidr of the VPC. Values ranges are: **10.0.0.0/8~10.255.255.240/28**,
  **172.16.0.0/12 ~ 172.31.255.240/28**, and **192.168.0.0/16 ~ 192.168.255.240/28**.
  Constraint: Must be an ipv4 cidr format, for example: **192.168.0.0/16**.

* `extend_cidrs` - (Optional, List) Specifies the extend cidrs of the VPC. Currently only supports IPv4 cidr.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the VPC.

* `cloud_resources` - (Optional, List) Specifies the cloud resources of the VPC.
  The [cloud_resources](#hwc_vpc_cloud_resources_struct) structure is documented below.

* `tags` - (Optional, List) Specifies the tags of the VPC.
  The [tag](#hwc_vpc_tag_struct) structure is documented below.

<a name="hwc_vpc_cloud_resources_struct"></a>
The `cloud_resources` block supports:

* `resource_type` - (Optional, String) Specifies the type of the cloud resources.

* `resource_count` - (Optional, Int) Specifies the asset count of the cloud resources.

<a name="hwc_vpc_tag_struct"></a>
The `tag` block supports:

* `key` - (Required, String) Specifies the tag key. The maximum length is `128` Unicode characters. The key cannot be empty.

* `values` - (Required, List) Specifies the tag values. Each value has a maximum length of `255` Unicode characters.
  If values is an empty list, it means any_value (query any value).

<a name="hwc_subnet_struct"></a>
The `hwc_subnet` block supports:

* `id` - (Required, String) Specifies the ID of the subnet.

* `name` - (Required, String) Specifies the name of the subnet.

* `project_id` - (Required, String) Specifies the project ID to which the security group belongs.

* `created_at` - (Required, String) Specifies the security group creation time.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `updated_at` - (Required, String) Specifies the security group update time.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `description` - (Optional, String) Specifies the security group description.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the security group belongs.

* `security_group_rules` - (Optional, List) Specifies the security group rules.
  The [security_group_rules](#hwc_security_group_rules_struct) structure is documented below.

<a name="hwc_security_group_rules_struct"></a>
The `security_group_rules` block supports:

* `id` - (Required, String) Specifies the unique identifier of the security group rule.

* `security_group_id` - (Required, String) Specifies the security group ID to which the security group rule belongs.

* `direction` - (Required, String) Specifies the direction of the security group rule.
  Valid values are **ingress** and **egress**.

* `protocol` - (Required, String) Specifies the protocol type.
  Valid values are **icmp**, **tcp**, **udp**, **icmpv6** and IP protocol number.

* `ethertype` - (Required, String) Specifies the IP address protocol type.
  Valid values are **IPv4** and **IPv6**.

* `multiport` - (Required, String) Specifies the port range.
  Support single port (`80`), continuous ports (`1-30`), and non-continuous ports (`22, 3389, 80`).

* `action` - (Required, String) Specifies the security group rule action. Valid values are **allow** and **deny**.

* `priority` - (Required, Int) Specifies the priority of the security group rule.
  Valid values are `1` to `100`, `1` represents the highest priority.

* `created_at` - (Required, String) Specifies the security group rule creation time.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `updated_at` - (Required, String) Specifies the security group rule update time.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `project_id` - (Required, String) Specifies the project ID to which the security group rule belongs.

* `description` - (Optional, String) Specifies the security group rule description.

* `remote_group_id` - (Optional, String) Specifies the remote security group ID.
  Valid values are the security group ID of the tenant.
  It is mutually exclusive with `remote_ip_prefix` and `remote_address_group_id`.

* `remote_ip_prefix` - (Optional, String) Specifies the remote IP address.
  When `direction` is **egress**, it is the IP address of the virtual machine.
  When `direction` is **ingress**, it is the IP address of the virtual machine.
  It is mutually exclusive with `remote_group_id` and `remote_address_group_id`.

* `remote_address_group_id` - (Optional, String) Specifies the remote address group ID.
  Valid values are the address group ID of the tenant.
  It is mutually exclusive with `remote_group_id` and `remote_ip_prefix`.

<a name="hwc_rds_struct"></a>
The `hwc_rds` block supports:

* `id` - (Required, String) Specifies the ID of the RDS instance.

* `name` - (Required, String) Specifies the name of the RDS instance.

* `protected_status` - (Required, String) Specifies the DBSS opening status of the RDS.
  Valid values are **OPEN** and **CLOSE**.

* `status` - (Required, String) Specifies the status of the RDS instance.
  Valid valus are:
  + **BUILD**: Indicates that the instance is creating.
  + **ACTIVE**: Indicates that the instance is normal.
  + **FAILED**: Indicates that the instance is abnormal.
  + **FROZEN**: Indicates that the instance is frozen.
  + **MODIFYING**: Indicates that the instance is being scaled.
  + **REBOOTING**: Indicates that the instance is rebooting.
  + **RESTORING**: Indicates that the instance is restoring.
  + **MODIFYING INSTANCE TYPE**: Indicates that the instance is being converted to a primary-secondary architecture.
  + **SWITCHOVER**: Indicates that the instance is switching to a primary-secondary architecture.
  + **MIGRATING**: Indicates that the instance is migrating.
  + **BACKING UP**: Indicates that the instance is backing up.
  + **MODIFYING DATABASE PORT**: Indicates that the instance is modifying the database port.
  + **STORAGE FULL**: Indicates that the instance's disk space is full.

* `port` - (Required, Int) Specifies the database port of the RDS.
  The database port setting range for RDS for MySQL is `1,024`～`65,535` (`12,017` and `33,071` are occupied by the RDS
  system and cannot be set).
  The database port setting range for RDS for PostgreSQL is `2,100`～`9,500`.
  The database port setting range for RDS for SQL Server is `1,433` and ranges from `2,100` to `9,500` (where `5355`
  and `5985` are not settable). For `2017 EE`, `2017 SE`, and `2017 Web`, `5050`, `5353`, and `5986` are not settable.

* `enable_ssl` - (Required, Bool) Specifies the SSL flag of the instance.
  Valid values are **true** and **false**.

* `type` - (Required, String) Specifies the type of the RDS. Valid values are **Single**, **Ha**, **Replica**, and
  **Enterprise**.

* `ha` - (Required, List) Specifies the HA configuration.
  The [ha](#hwc_rds_ha_struct) structure is documented below.

* `region` - (Required, String) Specifies the region where the RDS is located.

* `datastore` - (Required, List) Specifies the database information.
  The [datastore](#hwc_rds_datastore_struct) structure is documented below.

* `created` - (Required, String) Specifies the creation time.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `vpc_id` - (Required, String) Specifies the VPC ID.

* `subnet_id` - (Required, String) Specifies the subnet ID of the RDS.

* `security_group_id` - (Required, String) Specifies the security group ID.

* `flavor_ref` - (Required, String) Specifies the flavor of the RDS.

* `cpu` - (Required, String) Specifies the CPU size of the RDS.

* `mem` - (Required, String) Specifies the memory size of the RDS.

* `volume` - (Required, List) Specifies the volume information.
  The [volume](#hwc_rds_volume_struct) structure is documented below.

* `project_id` - (Required, String) Specifies the project ID.

* `switch_strategy` - (Required, String) Specifies the database switch strategy respectively.

* `read_only_by_user` - (Required, Bool) Specifies the user set read-only status of the RDS. Only supports RDS for
  MySQL engine.

* `backup_strategy` - (Required, List) Specifies the backup policy.
  The [backup_strategy](#hwc_rds_backup_strategy_struct) structure is documented below.

* `related_instance` - (Required, List) Specifies the list of associated database instances.
  The [related_instance](#hwc_rds_related_instance_struct) structure is documented below.

* `time_zone` - (Required, String) Specifies the time zone of the RDS.

* `expiration_time` - (Required, String) Specifies the expiration time of the RDS.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `maintenance_window` - (Required, String) Specifies the maintenance window of the RDS.

* `storage_used_space` - (Required, Float) Specifies the disk space usage, unit is GB.

* `nodes` - (Required, List) Specifies the main and standby instance information.
  The [nodes](#hwc_rds_node_struct) structure is documented below.

* `associated_with_ddm` - (Required, Bool) Specifies whether the instance is associated with DDM.

* `max_iops` - (Required, Int) Specifies the maximum IOPS of the disk.

* `updated` - (Required, String) Specifies the update time of the RDS.
  Format: **yyyy-MM-ddTHH:mm:ssZ**.

* `alias` - (Optional, String) Specifies the alias of the RDS.

* `private_ips` - (Optional, List) Specifies the private IP addresses of the RDS.

* `private_dns_names` - (Optional, List) Specifies the private DNS names of the RDS.

* `public_ips` - (Optional, List) Specifies the list of public IP addresses of the instance.

* `db_user_name` - (Optional, String) Specifies the default username of the RDS.

* `tags` - (Optional, List) Specifies the tags.
  The [tag](#hwc_rds_tag_struct) structure is documented below.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the RDS belongs.

* `disk_encryption_id` - (Optional, String) Specifies the disk encryption ID.

* `backup_used_space` - (Optional, Float) Specifies the backup space usage of the RDS. Only supports RDS for SQL Server engine.

<a name="hwc_rds_ha_struct"></a>
The `ha` block supports:

* `replication_mode` - (Required, String) Specifies the replication mode.
  Valid values are **async** and **semisync** for RDS for MySQL.
  Valid values are **async** and **sync** for RDS for PostgreSQL.
  Valid values are **sync** for RDS for Microsoft SQL Server.

<a name="hwc_rds_datastore_struct"></a>
The `datastore` block supports:

* `type` - (Required, String) Specifies the database engine type.
  Valid values are **MySQL**, **PostgreSQL**, and **SQLServer**.

* `version` - (Required, String) Specifies the database version.

* `complete_version` - (Optional, String) Specifies the complete version number of the database.
  Only returns when the database engine is **PostgreSQL**.

<a name="hwc_rds_volume_struct"></a>
The `volume` block supports:

* `type` - (Required, String) Specifies the volume type.

* `size` - (Required, Int) Specifies the volume size in GB.

<a name="hwc_rds_backup_strategy_struct"></a>
The `backup_strategy` block supports:

* `start_time` - (Optional, String) Specifies the backup time period. Automatic backup will be triggered in this period.

* `keep_days` - (Optional, Int) Specifies the number of days that the generated backup files can be saved.
  Valid values are `0` to `732`. When the value is `0`, it means that the automatic backup policy is not set or the
  automatic backup policy is disabled. If you need to extend the retention period, contact the customer service for
  application. The maximum retention period for automatic backup is `2,562` days.

<a name="hwc_rds_node_struct"></a>
The `node` block supports:

* `id` - (Required, String) Specifies the node ID.

* `name` - (Required, String) Specifies the node name.

* `role` - (Required, String) Specifies the node role. Valid values are **master**, **slave**, and **readreplica**.

* `status` - (Required, String) Specifies the node status.

* `availability_zone` - (Required, String) Specifies the AZ where the node resides.

<a name="hwc_rds_related_instance_struct"></a>
The `related_instance` block supports:

* `id` - (Required, String) Specifies the ID of the associated instance.

* `type` - (Required, String) Specifies the type of the associated instance.
  Valid values are **replica_of** and **replica**.

<a name="hwc_rds_tag_struct"></a>
The `tag` block supports:

* `key` - (Required, String) Specifies the tag key. The maximum length is `128` unicode characters.
  Key cannot be empty. (Search does not validate the character set of this parameter), key cannot be empty or an empty string,
  cannot be a space, trim half-width spaces before and after the validation and usage.

* `values` - (Required, List) Specifies the tag values. Each value has a maximum length of `255` unicode characters.
  If values is an empty list, it means any_value (query any value).

<a name="hwc_domain_struct"></a>
The `hwc_domain` block supports:

* `domain_name` - (Required, String) Specifies the domain name.

* `expire_date` - (Required, String) Specifies the domain expiration date. eg：**2023-01-10**.

* `status` - (Required, String) Specifies the domain status.

* `audit_status` - (Required, String) Specifies the domain real-name authentication status. Valid values are:
  + **NONAUDIT**：Unauthenticated.
  + **SUCCEED**：Authenticated.
  + **FAILED**：Authentication failed.
  + **AUDITING**：Authentication in progress.

* `audit_unpass_reason` - (Required, String) Specifies the reason for domain real-name authentication failure.

* `reg_type` - (Required, String) Specifies the registration type. Valid values are:
  + **PERSONAL**：Personal.
  + **COMPANY**：Company.

* `privacy_protection` - (Required, String) Specifies whether privacy protection is enabled.

* `name_server` - (Required, List) Specifies the domain name server list.

* `credential_type` - (Required, String) Specifies the type of credential.

* `credential_id` - (Required, String) Specifies the credential ID.

* `registrar` - (Required, String) Specifies the registrar of the domain.

* `contact` - (Required, List) Specifies the contact information.
  The [contact](#hwc_domain_contact_struct) structure is documented below.

* `transfer_status` - (Optional, String) Specifies the domain transfer status.

<a name="hwc_domain_contact_struct"></a>
The `contact` block supports:

* `email` - (Required, String) Specifies the email address of the domain contact.

* `register` - (Required, String) Specifies the registrant information.

* `contact_name` - (Required, String) Specifies the contact name.

* `phone_num` - (Required, String) Specifies the contact phone number.

* `province` - (Required, String) Specifies the province of the contact.

* `city` - (Required, String) Specifies the city of the contact.

* `address` - (Required, String) Specifies the address of the contact.

* `zip_code` - (Required, String) Specifies the zip code of the contact.

<a name="website_struct"></a>
The `website` block supports:

* `value` - (Required, String) Specifies the website URL.

* `main_domain` - (Required, String) Specifies the main domain of the website.

* `protected_status` - (Required, String) Specifies the WAF status. Valid values are:
  + **OPEN**：Enabled.
  + **CLOSE**：Disabled.

* `is_public` - (Required, Bool) Specifies whether the website is public or private. Valid values are:
  + **true**：Public.
  + **false**：Private.

* `name_server` - (Required, List) Specifies the website server list.

* `remark` - (Optional, String) Specifies the website remark.

* `extend_propertites` - (Optional, List) Specifies the other properties.
  The [extend_propertites](#website_extend_propertites_struct) structure is documented below.

<a name="website_extend_propertites_struct"></a>
The `extend_propertites` block supports:

* `mac_addr` - (Optional, String) Specifies the MAC address of the website server.

<a name="oca_ip_struct"></a>
The `oca_ip` block supports:

* `value` - (Required, String) Specifies the asset value.

* `version` - (Required, String) Specifies the asset version. Valid values are:
  + **ipv4**：IPv4.
  + **ipv6**：IPv6.

* `network` - (Required, List) Specifies the network information.
  The [network](#oca_ip_network_struct) structure is documented below.

* `server_room` - (Required, String) Specifies the server room.

* `server_rack` - (Required, String) Specifies the server rack.

* `data_center` - (Required, List) Specifies the data center.
  The [data_center](#oca_ip_data_center_struct) structure is documented below.

* `remark` - (Optional, String) Specifies the asset remark.

* `name` - (Optional, String) Specifies the asset name, default value is asset value.

* `relative_value` - (Optional, String) Specifies the relative value, such as ipv6 if the asset is ipv4.

* `mac_addr` - (Optional, String) Specifies the MAC address.

* `important` - (Optional, Int) Specifies the importance level, `0`: not important, `1`: important.

* `extend_propertites` - (Optional, List) Specifies the other third-party attributes.
  The [extend_propertites](#oca_ip_extend_propertites_struct) structure is documented below.

<a name="oca_ip_network_struct"></a>
The `network` block supports:

* `is_public` - (Required, Bool) Specifies whether the IP is public or private. Valid values are:
  + **true**：Public.
  + **false**：Private.

* `partition` - (Optional, String) Specifies the network partition. Valid values are:
  + **OM**：OM.
  + **PSZ**：PSZ.
  + **DMZ**：DMZ.

* `plane` - (Optional, String) Specifies the network plane (offline has its own code).

* `vxlan_id` - (Optional, String) Specifies the virtual network ID.

<a name="oca_ip_data_center_struct"></a>
The `data_center` block supports:

* `city_code` - (Required, String) Specifies the city code.

* `country_code` - (Required, String) Specifies the country code.

* `latitude` - (Optional, Float) Specifies the latitude.

* `longitude` - (Optional, Float) Specifies the longitude.

<a name="oca_ip_extend_propertites_struct"></a>
The `extend_propertites` block supports:

* `device` - (Optional, List) Specifies the device information.
  The [device](#oca_ip_device_struct) structure is documented below.

* `system` - (Optional, List) Specifies the system information.
  The [system](#oca_ip_system_struct) structure is documented below.

* `services` - (Optional, List) Specifies the application information.
  The [service](#oca_ip_service_struct) structure is documented below.

<a name="oca_ip_device_struct"></a>
The `device` block supports:

* `type` - (Optional, String) Specifies the device type.

* `model` - (Optional, String) Specifies the device model.

* `vendor` - (Optional, List) Specifies the vendor information.
  The [vendor](#oca_ip_device_vendor_struct) structure is documented below.

<a name="oca_ip_device_vendor_struct"></a>
The `vendor` block supports:

* `name` - (Optional, String) Specifies the vendor name.

* `is_xc` - (Optional, Bool) Specifies whether the supplier is domestic or not.

<a name="oca_ip_system_struct"></a>
The `system` block supports:

* `family` - (Optional, String) Specifies the type of the system.

* `name` - (Optional, String) Specifies the name of the system.

* `version` - (Optional, String) Specifies the version of the system.

* `vendor` - (Optional, List) Specifies the vendor information.
  The [vendor](#oca_ip_system_vendor_struct) structure is documented below.

<a name="oca_ip_system_vendor_struct"></a>
The `vendor` block supports:

* `name` - (Required, String) Specifies the vendor name.

* `is_xc` - (Optional, Bool) Specifies whether the supplier is domestic or not.

<a name="oca_ip_service_struct"></a>
The `service` block supports:

* `port` - (Optional, Int) Specifies the service port.

* `protocol` - (Optional, String) Specifies the service protocol.

* `name` - (Optional, String) Specifies the service name.

* `version` - (Optional, String) Specifies the service version.

* `vendor` - (Optional, List) Specifies the vendor information.
  The [vendor](#oca_ip_service_vendor_struct) structure is documented below.

<a name="oca_ip_service_vendor_struct"></a>
The `service.vendor` block supports:

* `name` - (Required, String) Specifies the vendor name.

* `is_xc` - (Optional, Bool) Specifies whether the supplier is domestic or not.

<a name="department_struct"></a>
The `department` block supports:

* `name` - (Optional, String) Specifies the department name of the asset.

* `id` - (Optional, String) Specifies the department ID of the asset.

<a name="governance_user_struct"></a>
The `governance_user` block supports:

* `type` - (Optional, String) Specifies the governance user type of the asset.

* `name` - (Optional, String) Specifies the governance user name of the asset.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (the asset ID).

## Import

The asset can be imported using the workspace ID and the asset ID, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_asset.test <workspace_id>/<asset_id>
```
