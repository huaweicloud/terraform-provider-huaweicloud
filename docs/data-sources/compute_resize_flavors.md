---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_resize_flavors"
description: |-
  Use this data source to get the list of target ECS flavors to which a flavor can be changed.
---

# huaweicloud_compute_resize_flavors

Use this data source to get the list of target ECS flavors to which a flavor can be changed.

## Example Usage

```hcl
variable "source_flavor_id" {}

data "huaweicloud_compute_resize_flavors" "this" {
  source_flavor_id = var.source_flavor_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_uuid` - (Optional, String) Specifies the target ECS ID in UUID format.

* `source_flavor_id` - (Optional, String) Specifies the source flavor ID.

* `source_flavor_name` - (Optional, String) Specifies the source flavor name.

* `sort_dir` - (Optional, String) Specifies the sorting of ECS flavors.
  Value options:
  + **asc**: indicates the ascending order.
  + **desc**: indicates the descending order.
  
  Defaults to **asc**.

* `sort_key` - (Optional, String) Specifies the field for sorting.
  Value options:
  + **flavorid**: indicates the flavor ID.
  + **name**: indicates the flavor name.
  + **memory_mb**: indicates the memory size.
  + **vcpus**: indicates the number of vCPUs.
  + **root_gb**: indicates the system disk size.
  
  Defaults to **flavorid**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Indicates the ECS flavors.

  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `id` - Indicates the ECS flavor ID.

* `name` - Indicates the ECS flavor name.

* `vcpus` - Indicates the number of vCPUs in the ECS flavor.

* `ram` - Indicates the memory size (MiB) in the ECS flavor.

* `disk` - Indicates the system disk size in the ECS flavor.
  This parameter has not been used. Its default value is **0**.

* `swap` - Indicates the swap partition size required by the ECS flavor.
  This parameter has not been used. Its default value is "".

* `os_flv_ext_data_ephemeral` - Indicates the temporary disk size. This is an extended attribute.
  This parameter has not been used. Its default value is **0**.

* `os_flv_disabled_disabled` - Indicates whether the ECS flavor has been disabled. This is an extended attribute.
  This parameter has not been used. Its default value is **false**.

* `rxtx_factor` - Indicates the ratio of the available network bandwidth to the network hardware bandwidth of the ECS.
  This parameter has not been used. Its default value is **1**.

* `rxtx_quota` - Indicates the software constraints of the network bandwidth that can be used by the ECS.
  This parameter has not been used. Its default value is **null**.

* `rxtx_cap` - Indicates the hardware constraints of the network bandwidth that can be used by the ECS.
  This parameter has not been used. Its default value is **null**.

* `os_flavor_access_is_public` - Indicates whether a flavor is available to all tenants. This is an extended attribute.
  Value options:
  + **true**: indicates that a flavor is available to all tenants.
  + **false**: indicates that a flavor is available only to certain tenants.

* `links` - Indicates the shortcut link of the ECS flavor.

  The [links](#flavors_links_struct) structure is documented below.

* `extra_specs` - Indicates the extended field of the ECS flavor.

  The [extra_specs](#flavors_extra_specs_struct) structure is documented below.

<a name="flavors_links_struct"></a>
The `links` block supports:

* `rel` - Indicates the shortcut link marker name.

* `href` - Indicates the provides the shortcut link.

* `type` - Indicates the shortcut link type.
  This parameter has not been used. Its default value is **null**.

<a name="flavors_extra_specs_struct"></a>
The `extra_specs` block supports:

* `hpet_support` - Indicates  whether to enable the high-precision clock on the ECS.
  The ECS specifications determine whether to return the parameter value.
  + **true** indicates to enable the function
  + **false** indicates to disable the function

* `resource_type` - Indicates the resource type.
  It is used to differentiate between the types of the physical servers accommodating ECSs.

* `cond_operation_az` - Indicates if an AZ is not configured in the `cond_operation_status` parameter.
  The value of this parameter is used by default. This parameter takes effect region-wide.
  This parameter is in the format of **az(xx)**. The value in parentheses is the flavor status in an AZ.
  If the parentheses are left blank, the configuration is invalid. The `cond_operation_az` options are the
  same as the `cond_operation_status` options.

* `quota_max_rate` - Indicates the maximum bandwidth.
  Unit: **Mbit/s**. If a bandwidth is in the unit of **Gbit/s**, it must be divided by **1,000**.

* `cond_operation_charge` - Indicates the billing type.
  + All the billing types are supported if this parameter is not set.
  + **period**: The billing type is yearly or monthly.
  + **demand**: The billing type is pay-per-use.

* `cond_storage` - Indicates the storage constraints.
  Disk features are supported. If this parameter is not set, the default configuration on the console is used.
  + **scsi**: indicates that SCSI is supported.
  + **localdisk**: indicates that local disks are supported.
  + **ib**: indicates that IB is supported.

* `info_gpu_name` - Indicates the number and names of GPUs.

* `info_cpu_name` - Indicates the CPU name.

* `quota_sub_network_interface_max_num` - Indicates the maximum number of auxiliary network interfaces that can be bound
  to an ECS.

* `quota_nvme_ssd` - Indicates the value of this parameter.
  The format is **{type}:{spec}:{num}:{size}:{safeFormat}**.
  + **type**: indicates the capacity of a single NVMe SSD disk attached to the ECS, which can only be 1.6 TB or 3.2 TB.
  + **spec: indicates the specifications of the NVMe SSD disk, which can be large (large specifications) or
    lvs (small specifications). If spec is set to large, only I series (for example, I3) is supported.
  + **num**: indicates the number of local disks.
  + **size**: indicates the capacity, in the unit of GiB, of the disk used by the guest user. If the spec value is large,
    the value of this parameter is the size of a single disk attached to the ECS. If the value of spec is lvs, the value
    of size is an integer multiple of 50.
  + **safeFormat**: indicates whether the local disks of the ECS are securely formatted.
    If safeFormat is set to True, only I series (for example, I3) is supported.

* `instance_vnic_type` - Indicates the NIC type.
  The value of this parameter is consistently **enhanced**, indicating that network enhancement ECSs are to be created.

* `instance_vnic_instance_bandwidth` - Indicates the maximum bandwidth in the unit of **Mbit/s**.
  The maximum value of this parameter is 10,000.

* `pci_passthrough_alias` - Indicates PCI passthrough device information, in the format of PCI device name: quantity.
  Multiple device information is separated by commas.

* `cond_network` - Indicates network constraints.
  Network features are supported. If this parameter is not set, the default configuration on the console is used.

* `quota_gpu` - Indicates the GPU name.

* `network_interface_traffic_mirroring_supported` - Indicates whether the flavor supports traffic mirroring.

* `info_asic_accelerators` - Indicates information about the accelerator.
  + **name**: accelerator name
  + **memory_mb**: accelerator memory
  + **count**: the number of accelerators
  + **alias_prefix**: internal alias of an accelerator

* `hw_numa_nodes` - Indicates the number of physical CPUs of the host.
  The ECS specifications determine whether to return the parameter value.

* `instance_vnic_max_count` - Indicates the maximum number of NICs.
  The maximum value of this parameter is 4.

* `extra_spec_io_persistent_grant` - Indicates whether persistence is supported.
  The value of this parameter is true.
  This parameter indicates that the ECS is persistently authorized to access the storage.

* `pci_passthrough_enable_gpu` - Indicates whether the graphics card is passed through.

* `cond_spot_operation_az` - Indicates  the AZ for the flavors in spot pricing billing mode.

* `quota_vif_max_num` - Indicates the maximum number of elastic network interfaces that can be bound to an ECS.

* `ecs_instance_architecture` - Indicates the CPU architecture corresponding to the flavor.
  This parameter is returned only for Kunpeng ECSs. The value arm64 indicates that the CPU architecture is Kunpeng.

* `security_enclave_supported` - Indicates whether the flavor supports QingTian Enclave.

* `info_gpus` - Indicates information about the GPU.
  + **name**: GPU name
  + **memory_mb**: GPU memory
  + **count**: the number of GPUs
  + **alias_prefix**: internal alias of a GPU

* `ecs_performancetype` - Indicates the ECS flavor type.
  + **normal**: general computing
  + **computingv3**: general computing-plus
  + **highmem**: memory-optimized
  + **cpuv1**: computing I
  + **cpuv2**: computing II
  + **highcpu**: high-performance computing
  + **diskintensive**: disk-intensive
  + **saphana**: large-memory
  + **kunpeng_highio**: Kunpeng ultra-high I/O
  + **kunpeng_accelerator**: Kunpeng application-accelerated
  + **advanced_smb**: general computing (providing resources for Huawei Cloud FlexusX)

* `quota_local_disk` - Indicates the value of this parameter is in format of **{type}:{count}:{size}:{safeFormat}**.
  + **type**: indicates the disk type, which can only be HDD.
  + **count**: indicates the number of local disks. The following types are supported:
      - For D1 ECSs, the value can be 3, 6, 12, or 24.
      - For D2 ECSs, the value can be 2, 4, 8, 12, 16, or 24.
      - For D3 ECSs, the value can be 2, 4, 8, 12, 16, 24, or 28.
  + **size**: indicates the capacity of a single disk, in GiB. Currently, only 1675 is supported.
    The actual disk size is 1800, and the available size after formatting is 1675.
  + **safeFormat**: indicates whether the local disks of the ECS are securely formatted. The following types are supported:
      - For D1 ECSs, the value is FALSE.
      - For D2 or D3 ECSs, the value is True.

* `cond_operation_status` - Indicates if an AZ is not configured in the `cond_operation_az` parameter.
  The value of this parameter is used by default. This parameter takes effect region-wide.
  If this parameter is not set or used, the meaning of normal applies. Options:
  + **normal**: indicates normal commercial use of the flavor.
  + **abandon**: indicates that the flavor has been taken offline (not displayed).
  + **sellout**: indicates that the flavor has been sold out.
  + **obt**: indicates that the flavor is under open beta testing (OBT).
  + **obt_sellout**: indicates that the OBT resources are sold out.
  + **promotion**: indicates that the flavor is recommended (for commercial use, which is similar to normal).

* `quota_min_rate` - Indicates the assured bandwidth.
  Unit: **Mbit/s**. If a bandwidth is in the unit of **Gbit/s**, it must be divided by **1,000**.

* `cond_operation_roles` - Indicates the allowed roles.
  Matched user tag (roles op_gatexxx), which is available to all users if this parameter is not set

* `cond_storage_type` - Indicates supported disk types.
  If you do not specify this parameter, the configuration on the console is used.
  + **SATA**: common I/O disks (sold out)
  + **SAS**: high I/O disks
  + **GPSSD**: General Purpose SSDs
  + **SSD**: ultra-I/O disks
  + **ESSD**: extreme SSDs
  + **GPSSD2**: General Purpose SSD V2
  + **ESSD2**: extreme SSD V2

* `ecs_virtualization_env_types` - Indicates a virtualization type.
  + If the parameter value is **FusionCompute**, the ECS uses Xen virtualization.
  + If the parameter value is **CloudCompute**, the ECS uses KVM virtualization.

* `pci_passthrough_gpu_specs` - Indicates  the technologies used by G1 and G2 cloud servers.
  It includes GPU virtualization and GPU pass-through.

* `quota_max_pps` - Indicates the maximum intranet PPS.
  Unit: number. If a value is in the unit of 10,000, it must be divided by 10,000.

* `cond_operation_charge_stop` - Indicates whether fees are billed for a stopped ECS.
  + No fees by default
  + **charge**
  + **free**

* `cond_spot_operation_status` - Indicates the status of a flavor in spot pricing billing mode.
  + Equivalent to abandon if this parameter is not set.
  + **normal**: indicates commercial use of the flavor.
  + **abandon**: indicates that the flavor has been taken offline.
  + **sellout**: indicates that the flavor has been sold out.
  + **obt**: indicates that the flavor is under OBT (not supported currently).
  + **private**: indicates that the flavor is private, which is available only to specified users (not supported currently).
  + **test**: indicates that the flavor is at free trial phase (not supported currently).
  + **promotion**: indicates that the flavor is recommended.

* `cond_compute_live_resizable` - Indicates computing constraints.
  + If the value of this parameter is true, online capacity expansion is supported.
  + If this parameter does not exist or its value is set to false, online capacity expansion is not supported.

* `cond_compute` - Indicates computing constraints.
  + **autorecovery**: indicates that automatic recovery is supported.
  + If this parameter does not exist, automatic recovery is not supported.

* `ecs_generation` - Indicates the generation of an ECS type.

* `info_features` - Indicates the features supported by the flavor.
