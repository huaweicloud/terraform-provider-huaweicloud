---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud\_cce\_node
Add a node to a container cluster.
This is an alternative to `huaweicloud_cce_node_v3`

## Basic Usage

```hcl

data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_compute_keypair" "mykp" {
  name       = "mykp"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_cce_cluster" "mycluster" {
  name                   = "mycluster"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.myvpc.id
  subnet_id              = huaweicloud_vpc_subnet.mysubnet.id
  container_network_type = "overlay_l2"
}

resource "huaweicloud_cce_node" "node" {
  cluster_id        = huaweicloud_cce_cluster.mycluster.id
  name              = "node"
  flavor_id         = "s3.large.2"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = huaweicloud_compute_keypair.mykp.name

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

## Node with Eip

```hcl
resource "huaweicloud_cce_node" "mynode" {
  cluster_id        = huaweicloud_cce_cluster.mycluster.id
  name              = "mynode"
  flavor_id         = "s3.large.2"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = huaweicloud_compute_keypair.mykp.name

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

## Node with Existing Eip

```hcl
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
  cluster_id        = huaweicloud_cce_cluster.mycluster.id
  name              = "mynode"
  flavor_id         = "s3.large.2"
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  key_pair          = huaweicloud_compute_keypair.mykp.name

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

## Argument Reference
The following arguments are supported:

* `region` - (Optional) The region in which to obtain the cce node resource. If omitted, the provider-level region will work as default. Changing this creates a new cce node resource.

* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.

* `name` - (Optional) Node Name.

* `flavor_id` - (Required) Specifies the flavor id. Changing this parameter will create a new resource.

* `availability_zone` - (Required) specify the name of the available partition (AZ). Changing this parameter will create a new resource.

* `os` - (Optional) Operating System of the node. Changing this parameter will create a new resource.
    - For VM nodes, clusters of v1.13 and later support *EulerOS 2.5* and *CentOS 7.6*.
    - For BMS nodes purchased in the yearly/monthly billing mode, only *EulerOS 2.3* is supported.

* `key_pair` - (Optional) Key pair name when logging in to select the key pair mode. This parameter and `password` are alternative.
    Changing this parameter will create a new resource.

* `password` - (Optional) root password when logging in to select the password mode. This parameter must be salted and alternative to `key_pair`.
    Changing this parameter will create a new resource.

* `subnet_id` - (Optional) The ID of the subnet to which the NIC belongs. Changing this parameter will create a new resource.

* `eip_id` - (Optional) The ID of the EIP. Changing this parameter will create a new resource.

* `eip_ids` - (Deprecated) This has been deprecated, use eip_id instead. List of existing elastic IP IDs.
    Changing this parameter will create a new resource.

-> **Note:** If the eip_id parameter is configured, you do not need to configure the bandwidth parameters:
  `iptype`, `bandwidth_charge_mode`, `bandwidth_size` and `share_type`.

* `iptype` - (Optional) Elastic IP type. Changing this parameter will create a new resource.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Changing this parameter will create a new resource.

* `sharetype` - (Optional) Bandwidth sharing type. Changing this parameter will create a new resource.

* `bandwidth_size` - (Optional) Bandwidth size. Changing this parameter will create a new resource.

* `max_pods` - (Optional) The maximum number of instances a node is allowed to create. Changing this parameter will create a new cluster resource.

* `preinstall` - (Optional) Script required before installation. The input value can be a Base64 encoded string or not.
    Changing this parameter will create a new resource.

* `postinstall` - (Optional) Script required after installation. The input value can be a Base64 encoded string or not.
   Changing this parameter will create a new resource.

* `tags` - (Optional) Tags of a VM node, key/value pair format.

* `root_volume` - (Required) It corresponds to the system disk related configuration. Changing this parameter will create a new resource.

  * `size` - (Required) Disk size in GB.
  * `volumetype` - (Required) Disk type.
  * `extend_param` - (Optional) Disk expansion parameters.

* `data_volumes` - (Required) Represents the data disk to be created. Changing this parameter will create a new resource.

  * `size` - (Required) Disk size in GB.
  * `volumetype` - (Required) Disk type.
  * `extend_param` - (Optional) Disk expansion parameters.

* `taints` - (Optional) You can add taints to created nodes to configure anti-affinity. Each taint contains the following parameters:

  * `key` - (Required) A key must contain 1 to 63 characters starting with a letter or digit. Only letters, digits, hyphens (-),
    underscores (_), and periods (.) are allowed. A DNS subdomain name can be used as the prefix of a key.
  * `value` - (Required) A value must start with a letter or digit and can contain a maximum of 63 characters, including letters,
    digits, hyphens (-), underscores (_), and periods (.).
  * `effect` - (Required) Available options are NoSchedule, PreferNoSchedule, and NoExecute.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

 * `status` -  Node status information.

 * `server_id` - ID of the ECS instance associated with the node.

 * `private_ip` - Private IP of the CCE node.

 * `public_ip` - Public IP of the CCE node.
