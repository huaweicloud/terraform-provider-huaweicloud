---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_node"
description: ""
---

# huaweicloud_cce_node

Add a node to a CCE cluster.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}
variable "node_name" {}
variable "keypair_name" {}

data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_compute_flavors" "myflavors" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_kps_keypair" "mykp" {
  name       = var.keypair_name
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_cce_node" "node" {
  cluster_id        = var.cluster_id
  name              = var.node_name
  flavor_id         = data.huaweicloud_compute_flavors.myflavors.ids[0]
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = huaweicloud_kps_keypair.mykp.name

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
```

### Node with Eip

```hcl
variable "cluster_id" {}
variable "node_name" {}
variable "keypair_name" {}

data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_kps_keypair" "mykp" {
  name       = var.keypair_name
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_cce_node" "mynode" {
  cluster_id        = var.cluster_id
  name              = var.node_name
  flavor_id         = data.huaweicloud_compute_flavors.myflavors.ids[0]
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = huaweicloud_kps_keypair.mykp.name

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }

  // Assign EIP
  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 100
}
```

### Node with Existing Eip

```hcl
variable "cluster_id" {}
variable "node_name" {}
variable "keypair_name" {}

data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_kps_keypair" "mykp" {
  name       = var.keypair_name
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_node" "mynode" {
  cluster_id        = var.cluster_id
  name              = var.node_name
  flavor_id         = data.huaweicloud_compute_flavors.myflavors.ids[0]
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = huaweicloud_kps_keypair.mykp.name

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }

  // Assign existing EIP
  eip_id = huaweicloud_vpc_eip.myeip.id
}
```

### Node with storage configuration

```hcl
variable "cluster_id" {}
variable "node_name" {}
variable "keypair_name" {}
variable "kms_key_name" {}

data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_kps_keypair" "mykp" {
  name       = var.keypair_name
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_kms_key" "mykey" {
  key_alias    = var.kms_key_name
  pending_days = "7"
}

resource "huaweicloud_cce_node" "mynode" {
  cluster_id        = var.cluster_id
  name              = var.node_name
  flavor_id         = data.huaweicloud_compute_flavors.myflavors.ids[0]
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = huaweicloud_kps_keypair.mykp.name

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
    kms_key_id = huaweicloud_kms_key.mykey.id
  }

  // Storage configuration
  storage {
    selectors {
      name              = "cceUse"
      type              = "evs"
      match_label_size  = "100"
      match_label_count = 1
    }

    selectors {
      name                           = "user"
      type                           = "evs"
      match_label_size               = "100"
      match_label_metadata_encrypted = "1"
      match_label_metadata_cmkid     = huaweicloud_kms_key.mykey.id
      match_label_count              = "1"
    }

    groups {
      name           = "vgpaas"
      selector_names = ["cceUse"]
      cce_managed    = true

      virtual_spaces {
        name        = "kubernetes"
        size        = "10%"
        lvm_lv_type = "linear"
      }

      virtual_spaces {
        name        = "runtime"
        size        = "90%"
      }
    }

    groups {
      name           = "vguser"
      selector_names = ["user"]

      virtual_spaces {
        name        = "user"
        size        = "100%"
        lvm_lv_type = "linear"
        lvm_path    = "/workspace"
      }
    }
  }
}
```

### Spot Node

```hcl
variable "cluster_id" {}
variable "node_name" {}
variable "keypair_name" {}

data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_compute_flavors" "myflavors" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_kps_keypair" "mykp" {
  name       = var.keypair_name
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_cce_node" "node" {
  cluster_id        = var.cluster_id
  name              = var.node_name
  flavor_id         = data.huaweicloud_compute_flavors.myflavors.ids[0]
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = huaweicloud_kps_keypair.mykp.name

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }

  extend_params {
    market_type = "spot"
    spot_price  = "0.83"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE node resource.
  If omitted, the provider-level region will be used. Changing this creates a new CCE node resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster.

* `name` - (Optional, String) Specifies the node name.

* `flavor_id` - (Required, String, NonUpdatable) Specifies the flavor ID.
  resource.

* `availability_zone` - (Required, String, NonUpdatable) Specifies the name of the available partition (AZ).

* `os` - (Optional, String, NonUpdatable) Specifies the operating system of the node.
  The value can be **EulerOS 2.9** and **CentOS 7.6** e.g. For more details,
  please see [documentation](https://support.huaweicloud.com/intl/en-us/api-cce/node-os.html).
  This parameter is required when the `node_image_id` in `extend_params` is not specified.

* `key_pair` - (Optional, String) Specifies the key pair name when logging in to select the key pair mode.
  This parameter and `password` are alternative.

* `password` - (Optional, String) Specifies the root password when logging in to select the password mode.
  The password consists of 8 to 26 characters and must contain at least three of following: uppercase letters,
  lowercase letters, digits, special characters(!@$%^-_=+[{}]:,./?~#*).
  This parameter can be plain or salted and is alternative to `key_pair`.

  -> A new password is in plain text and takes effect after the node is started or restarted.

* `private_key` - (Optional, String) Specifies the private key of the in used `key_pair`. This parameter is mandatory
  when replacing or unbinding a keypair if the CCE node is in **Active** state.

* `root_volume` - (Required, List, NonUpdatable) Specifies the configuration of the system disk.

  + `size` - (Required, Int, NonUpdatable) Specifies the disk size in GB.

  + `volumetype` - (Required, String, NonUpdatable) Specifies the disk type.

  + `extend_params` - (Optional, Map, NonUpdatable) Specifies the disk expansion parameters.

  + `kms_key_id` - (Optional, String, NonUpdatable) Specifies the ID of a KMS key. This is used to encrypt the volume.

  + `dss_pool_id` - (Optional, String, NonUpdatable) Specifies the DSS pool ID. This field is used only for

  + `iops` - (Optional, Int, NonUpdatable) Specifies the iops of the disk,
    required when `volumetype` is **GPSSD2** or **ESSD2**.

  + `throughput` - (Optional, Int, NonUpdatable) Specifies the throughput of the disk in MiB/s,
    required when `volumetype` is **GPSSD2**.

* `data_volumes` - (Optional, List, NonUpdatable) Specifies the configurations of the data disk.

  + `size` - (Required, Int, NonUpdatable) Specifies the disk size in GB.

  + `volumetype` - (Required, String, NonUpdatable) Specifies the disk type.

  + `extend_params` - (Optional, Map, NonUpdatable) Specifies the disk expansion parameters.

  + `kms_key_id` - (Optional, String, NonUpdatable) Specifies the ID of a KMS key. This is used to encrypt the volume.

  + `dss_pool_id` - (Optional, String, NonUpdatable) Specifies the DSS pool ID. This field is used only for

  + `iops` - (Optional, Int, NonUpdatable) Specifies the iops of the disk,
    required when `volumetype` is **GPSSD2** or **ESSD2**.

  + `throughput` - (Optional, Int, NonUpdatable) Specifies the throughput of the disk in MiB/s,
    required when `volumetype` is **GPSSD2**.

    -> You need to create an agency (EVSAccessKMS) when disk encryption is used in the current project for the first
    time ever.

* `storage` - (Optional, List, NonUpdatable) Specifies the disk initialization management parameter.
  If omitted, disks are managed based on the DockerLVMConfigOverride parameter in extendParam.
  This parameter is supported for clusters of v1.15.11 and later.
  If the node has both local and EVS disks attached,
  this parameter must be specified, or it may result in unexpected disk partitions.
  If you want to change the value range of a data disk to **20** to **32768**, this parameter must be specified.
  If you want to use the shared disk space (with the runtime and Kubernetes partitions cancelled),
  this parameter must be specified.
  If you want to store system components in the system disk, this parameter must be specified.

  + `selectors` - (Required, List, NonUpdatable) Specifies the disk selection.
    Matched disks are managed according to match labels and storage type. Structure is documented below.

  + `groups` - (Required, List, NonUpdatable) Specifies the storage group consists of multiple storage devices.
    This is used to divide storage space. Structure is documented below.

* `subnet_id` - (Optional, String, NonUpdatable) Specifies the ID of the subnet to which the NIC belongs.

* `fixed_ip` - (Optional, String, NonUpdatable) Specifies the fixed IP of the NIC.

* `extension_nics` - (Optional, List, NonUpdatable) Specifies extension NICs of the node.
  The [object](#extension_nics) structure is documented below.

* `eip_id` - (Optional, String, NonUpdatable) Specifies the ID of the EIP.

-> **NOTE:** If the eip_id parameter is configured, you do not need to configure the bandwidth parameters:
`iptype`, `bandwidth_charge_mode`, `bandwidth_size` and `share_type`.

* `iptype` - (Optional, String, NonUpdatable) Specifies the elastic IP type.

* `bandwidth_charge_mode` - (Optional, String, NonUpdatable) Specifies the bandwidth billing type.

* `sharetype` - (Optional, String, NonUpdatable) Specifies the bandwidth sharing type.

* `bandwidth_size` - (Optional, Int, NonUpdatable) Specifies the bandwidth size.

* `ecs_group_id` - (Optional, String, NonUpdatable) Specifies the ECS group ID. If specified, the node will be created under
  the cloud server group.

* `charging_mode` - (Optional, String, NonUpdatable) Specifies the charging mode of the CCE node. Valid values are *prePaid*
  and *postPaid*, defaults to *postPaid*.

* `period_unit` - (Optional, String, NonUpdatable) Specifies the charging period unit of the CCE node.
  Valid values are *month* and *year*. This parameter is mandatory if `charging_mode` is set to *prePaid*.

* `period` - (Optional, Int, NonUpdatable) Specifies the charging period of the CCE node. If `period_unit` is set to *month*
  , the value ranges from 1 to 9. If `period_unit` is set to *year*, the value ranges from 1 to 3. This parameter is
  mandatory if `charging_mode` is set to *prePaid*.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled. Valid values are "true" and "false".

* `runtime` - (Optional, String, NonUpdatable) Specifies the runtime of the CCE node. Valid values are *docker* and
  *containerd*.

* `extend_params` - (Optional, List, NonUpdatable) Specifies the extended parameters.
  The [object](#extend_params) structure is documented below.

* `dedicated_host_id` - (Optional, String, NonUpdatable) Specifies the ID of the DeH to which the node is scheduled.

* `initialized_conditions` - (Optional, List, NonUpdatable) Specifies the custom initialization flags.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID of the CCE node.

* `labels` - (Optional, Map, NonUpdatable) Specifies the tags of a Kubernetes node, key/value pair format.

* `tags` - (Optional, Map) Specifies the tags of a VM node, key/value pair format.

* `taints` - (Optional, List, NonUpdatable) Specifies the taints configuration of the nodes to set anti-affinity.
  Each taint contains the following parameters:

  + `key` - (Required, String, NonUpdatable) A key must contain 1 to 63 characters starting with a letter or digit.
    Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed. A DNS subdomain name can be used
    as the prefix of a key.

  + `value` - (Optional, String, NonUpdatable) A value must start with a letter or digit and can contain a maximum of 63
    characters, including letters, digits, hyphens (-), underscores (_), and periods (.).

  + `effect` - (Required, String, NonUpdatable) Available options are NoSchedule, PreferNoSchedule, and NoExecute.

* `hostname_config` - (Optional, List, NonUpdatable) Specifies the hostname config of the kubernetes node,
  which is supported by clusters of v1.23.6-r0 to v1.25 or clusters of v1.25.2-r0 or later versions.
  The [object](#hostname_config) structure is documented below.

* `partition` - (Optional, String, NonUpdatable) Specifies the partition to which the node belongs. Value options:
  + **center**: center cloud.
  + The availability zone ID of the edge station.

<a name="extension_nics"></a>
The `extension_nics` block supports:

* `subnet_id` - (Required, String, NonUpdatable) Specifies the ID of the subnet to which the NIC belongs.

<a name="extend_params"></a>
The `extend_params` block supports:

* `max_pods` - (Optional, Int, NonUpdatable) Specifies the maximum number of instances a node is allowed to create.

* `docker_base_size` - (Optional, Int, NonUpdatable) Specifies the available disk space of a single container on a node,
  in GB.

* `preinstall` - (Optional, String, NonUpdatable) Specifies the script to be executed before installation.
  The input value can be a Base64 encoded string or not.

* `postinstall` - (Optional, String, NonUpdatable) Specifies the script to be executed after installation.
  The input value can be a Base64 encoded string or not.

* `node_image_id` - (Optional, String, NonUpdatable) Specifies the image ID to create the node.

* `node_multi_queue` - (Optional, String, NonUpdatable) Specifies the number of ENI queues.
  Example setting: **"[{\"queue\":4}]"**.

* `nic_threshold` - (Optional, String, NonUpdatable) Specifies the ENI pre-binding thresholds.
  Example setting: **"0.3:0.6"**.

* `agency_name` - (Optional, String, NonUpdatable) Specifies the agency name.

* `kube_reserved_mem` - (Optional, Int, NonUpdatable) Specifies the reserved node memory, which is reserved for
  Kubernetes-related components.

* `system_reserved_mem` - (Optional, Int, NonUpdatable) Specifies the reserved node memory, which is reserved
  value for system components.

* `security_reinforcement_type` - (Optional, String, NonUpdatable) Specifies the security reinforcement type.
  The value can be: **null** or **cybersecurity**.

* `market_type` - (Optional, String, NonUpdatable) Specifies the market type. When creating a spot node,
  this parameter should be set to **spot**.

* `spot_price` - (Optional, String, NonUpdatable) Specifies the highest price per hour a user accepts for a spot node.

The `selectors` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the selector name, used as the index of `selector_names`
  in storage group. The name of each selector must be unique.

* `type` - (Optional, String, NonUpdatable) Specifies the storage type. Currently, only **evs (EVS volumes)** is supported.
  The default value is **evs**.

* `match_label_size` - (Optional, String, NonUpdatable) Specifies the matched disk size. If omitted,
  the disk size is not limited. Example: 100.

* `match_label_volume_type` - (Optional, String, NonUpdatable) Specifies the EVS disk type. Currently,
  **SSD**, **GPSSD**, and **SAS** are supported. If omitted, the disk type is not limited.

* `match_label_metadata_encrypted` - (Optional, String, NonUpdatable) Specifies the disk encryption identifier.
  Values can be: **0** indicates that the disk is not encrypted and **1** indicates that the disk is encrypted.
  If omitted, whether the disk is encrypted is not limited.

* `match_label_metadata_cmkid` - (Optional, String, NonUpdatable) Specifies the customer master key ID of an encrypted
  disk.

* `match_label_count` - (Optional, String, NonUpdatable) Specifies the number of disks to be selected. If omitted,
  all disks of this type are selected.

The `groups` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of a virtual storage group. Each group name must be unique.

* `cce_managed` - (Optional, Bool, NonUpdatable) Specifies the whether the storage space is for **kubernetes** and
  **runtime** components. Only one group can be set to true. The default value is **false**.

* `selector_names` - (Required, List, NonUpdatable) Specifies the list of names of selectors to match.
  This parameter corresponds to name in `selectors`. A group can match multiple selectors,
  but a selector can match only one group.

* `virtual_spaces` - (Required, List, NonUpdatable) Specifies the detailed management of space configuration in a group.

  + `name` - (Required, String, NonUpdatable) Specifies the virtual space name. Currently, only **kubernetes**, **runtime**,
    and **user** are supported.

  + `size` - (Required, String, NonUpdatable) Specifies the size of a virtual space. Only an integer percentage is supported.
    Example: 90%. Note that the total percentage of all virtual spaces in a group cannot exceed 100%.

  + `lvm_lv_type` - (Optional, String, NonUpdatable) Specifies the LVM write mode, values can be **linear** and **striped**.
    This parameter takes effect only in **kubernetes** and **user** configuration.

  + `lvm_path` - (Optional, String, NonUpdatable) Specifies the absolute path to which the disk is attached.
    This parameter takes effect only in **user** configuration.

  + `runtime_lv_type` - (Optional, String, NonUpdatable) Specifies the LVM write mode, values can be **linear** and **striped**.
    This parameter takes effect only in **runtime** configuration.

<a name="hostname_config"></a>
The `hostname_config` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the hostname type of the kubernetes node.
  The value can be:
  + **privateIp**: The Kubernetes node is named after its IP address.
  + **cceNodeName**: The Kubernetes node is named after the CCE node.
  
  If `hostname_config` not specified, the default value is **privateIp**.

  ~>For a node which is configured using cceNodeName, the name is the same as the Kubernetes node name and the ECS name.
    The node name cannot be changed. If the ECS name is changed on the ECS console, the node name will retain unchanged
    after ECS synchronization. To avoid a conflict between Kubernetes nodes, the system automatically adds a suffix to
    each node name. The suffix is in the format of A hyphen (-) Five random characters. The value of the random
    characters is a lowercase letter or a digit ranging from 0 to 9.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `server_id` - ID of the ECS instance associated with the node.
* `private_ip` - Private IP of the CCE node.
* `public_ip` - Public IP of the CCE node.
* `status` - The status of the CCE node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

CCE node can be imported using the cluster ID and node ID separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_cce_node.my_node 5c20fdad-7288-11eb-b817-0255ac10158b/e9287dff-7288-11eb-b817-0255ac10158b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`password`, `private_key`, `storage`, `fixed_ip`, `extension_nics`, `eip_id`, `iptype`, `bandwidth_charge_mode`,
`bandwidth_size`, `share_type`, `extend_params`, `dedicated_host_id`, `initialized_conditions`, `labels`, `taints`
and arguments for pre-paid. It is generally recommended running `terraform plan` after importing a node.
You can then decide if changes should be applied to the node, or the resource definition should be updated to align
with the node. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cce_node" "my_node" {
    ...

  lifecycle {
    ignore_changes = [
      extend_params, labels,
    ]
  }
}
```
