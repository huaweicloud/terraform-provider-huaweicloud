---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster"
description: ""
---

# huaweicloud_cce_cluster

Provides a CCE cluster resource.

## Example Usage

### Basic Usage

```hcl
resource "huaweicloud_vpc" "myvpc" {
  name = "vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "mysubnet" {
  name       = "subnet"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = huaweicloud_vpc.myvpc.id
}

resource "huaweicloud_cce_cluster" "cluster" {
  name                   = "cluster"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.myvpc.id
  subnet_id              = huaweicloud_vpc_subnet.mysubnet.id
  container_network_type = "overlay_l2"
}
```

### Cluster With EIP

```hcl
resource "huaweicloud_vpc" "myvpc" {
  name = "vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "mysubnet" {
  name       = "subnet"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = huaweicloud_vpc.myvpc.id
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

resource "huaweicloud_cce_cluster" "cluster" {
  name                   = "cluster"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.myvpc.id
  subnet_id              = huaweicloud_vpc_subnet.mysubnet.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = huaweicloud_vpc_eip.myeip.address
}
```

### CCE Turbo Cluster

```hcl
resource "huaweicloud_vpc" "myvpc" {
  name = "vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "mysubnet" {
  name       = "subnet"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = huaweicloud_vpc.myvpc.id
}

resource "huaweicloud_vpc_subnet" "eni_test_1" {
  name          = "subnet-eni-1"
  cidr          = "192.168.2.0/24"
  gateway_ip    = "192.168.2.1"
  vpc_id        = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_subnet" "eni_test_2" {
  name          = "subnet-eni-2"
  cidr          = "192.168.3.0/24"
  gateway_ip    = "192.168.3.1"
  vpc_id        = huaweicloud_vpc.test.id
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "cluster"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.myvpc.id
  subnet_id              = huaweicloud_vpc_subnet.mysubnet.id
  container_network_type = "eni"
  eni_subnet_id          = join(",", [
    huaweicloud_vpc_subnet.eni_test_1.ipv4_subnet_id,
    huaweicloud_vpc_subnet.eni_test_2.ipv4_subnet_id,
  ])
}
```

### CCE HA Cluster

```hcl
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_cce_cluster" "cluster" {
  name                   = "cluster"
  flavor_id              = "cce.s2.small"
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  container_network_type = "overlay_l2"

  masters {
    availability_zone = "cn-north-4a"
  }
  masters {
    availability_zone = "cn-north-4b"
  }
  masters {
    availability_zone = "cn-north-4c"
  }
}
```

### Cluster with Component Configurations

```hcl
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_cce_cluster" "cluster" {
  name                   = "cluster"
  flavor_id              = "cce.s1.small"
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  container_network_type = "overlay_l2"

  component_configurations {
    name           = "kube-apiserver"
    configurations = jsonencode([
      {
        name  = "default-not-ready-toleration-seconds"
        value = "100"
      }
    ])
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE cluster resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new cluster resource.

* `name` - (Required, String, ForceNew) Specifies the cluster name.
  Changing this parameter will create a new cluster resource.

* `flavor_id` - (Required, String) Specifies the cluster specifications.
  Possible values:
  + **cce.s1.small**: small-scale single cluster (up to 50 nodes).
  + **cce.s1.medium**: medium-scale single cluster (up to 200 nodes).
  + **cce.s2.small**: small-scale HA cluster (up to 50 nodes).
  + **cce.s2.medium**: medium-scale HA cluster (up to 200 nodes).
  + **cce.s2.large**: large-scale HA cluster (up to 1000 nodes).
  + **cce.s2.xlarge**: large-scale HA cluster (up to 2000 nodes).

  -> Changing the number of control nodes or reducing cluster flavor is not supported.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC used to create the node.
  Changing this parameter will create a new cluster resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the subnet used to create the node which should be
  configured with a *DNS address*. Changing this parameter will create a new cluster resource.

* `container_network_type` - (Required, String, ForceNew) Specifies the container network type.
  Changing this parameter will create a new cluster resource. Possible values:
  + **overlay_l2**: An overlay_l2 network built for containers by using Open vSwitch(OVS).
  + **vpc-router**: An vpc-router network built for containers by using ipvlan and custom VPC routes.
  + **eni**: A Yangtse network built for CCE Turbo cluster. The container network deeply integrates the native ENI
    capability of VPC, uses the VPC CIDR block to allocate container addresses, and supports direct connections between
    ELB and containers to provide high performance.

* `security_group_id` - (Optional, String) Specifies the default worker node security group ID of the cluster.
  If left empty, the system will automatically create a default worker node security group for you.
  The default worker node security group needs to allow access from certain ports to ensure normal communications.
  For details, see [documentation](https://support.huaweicloud.com/intl/en-us/cce_faq/cce_faq_00265.html).
  If updated, the modified security group will only be applied to nodes newly created or accepted.
  For existing nodes, you need to manually modify the security group rules for them.

* `cluster_version` - (Optional, String) Specifies the cluster version, defaults to the latest supported
  version. Changing this parameter will not upgrade the cluster. If you want to upgrade the cluster, please use
  resource `huaweicloud_cce_cluster_upgrade`. After upgrading cluster successfully, you can update this parameter
  to avoid unexpected changing plan.

* `cluster_type` - (Optional, String, ForceNew) Specifies the cluster Type, possible values are **VirtualMachine** and
  **ARM64**. Defaults to **VirtualMachine**. Changing this parameter will create a new cluster resource.

* `alias` - (Optional, String) Specifies the display name of a cluster. The value of `alias` cannot be the same as the `name`
  and display names of other clusters.

* `timezone` - (Optional, String, ForceNew) Specifies the time zone of a cluster. Changing this parameter will create a
  new cluster resource.

* `description` - (Optional, String) Specifies the cluster description.

* `container_network_cidr` - (Optional, String) Specifies the container network segments.
  In clusters of v1.21 and later, when the `container_network_type` is **vpc-router**, you can add multiple container
  segments, separated with comma (,). In other situations, only the first segment takes effect.

* `service_network_cidr` - (Optional, String, ForceNew) Specifies the service network segment.
  Changing this parameter will create a new cluster resource.

* `eni_subnet_id` - (Optional, String) Specifies the **IPv4 subnet ID** of the subnet where the ENI resides.
  Specified when creating a CCE Turbo cluster. You can add multiple IPv4 subnet ID, separated with comma (,).
  Only adding subnets is allowed, removing subnets is not allowed.

* `authentication_mode` - (Optional, String, ForceNew) Specifies the authentication mode of the cluster, possible values
  are **rbac** and **authenticating_proxy**. Defaults to **rbac**.
  Changing this parameter will create a new cluster resource.

* `authenticating_proxy_ca` - (Optional, String, ForceNew) Specifies the CA root certificate provided in the
  **authenticating_proxy** mode. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new cluster resource.

* `authenticating_proxy_cert` - (Optional, String, ForceNew) Specifies the Client certificate provided in the
  **authenticating_proxy** mode. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new cluster resource.

* `authenticating_proxy_private_key` - (Optional, String, ForceNew) Specifies the private key of the client certificate
  provided in the **authenticating_proxy** mode. The input value can be a Base64 encoded string or not.
  Changing this parameter will create a new cluster resource.

-> **Note:** For more detailed description of authenticating_proxy mode for authentication_mode see
[Enhanced authentication](https://github.com/huaweicloud/terraform-provider-huaweicloud/blob/master/examples/cce/basic/cce-cluster-enhanced-authentication.md).

* `multi_az` - (Optional, Bool, ForceNew) Specifies whether to enable multiple AZs for the cluster, only when using HA
  flavors. Changing this parameter will create a new cluster resource. This parameter and `masters` are alternative.

* `masters` - (Optional, List, ForceNew) Specifies the advanced configuration of master nodes.
  The [object](#cce_cluster_masters) structure is documented below.
  This parameter and `multi_az` are alternative. Changing this parameter will create a new cluster resource.

* `eip` - (Optional, String) Specifies the EIP address of the cluster.

* `kube_proxy_mode` - (Optional, String, ForceNew) Specifies the service forwarding mode.
  Changing this parameter will create a new cluster resource. Two modes are available:

  + **iptables**: Traditional kube-proxy uses iptables rules to implement service load balancing. In this mode, too many
    iptables rules will be generated when many services are deployed. In addition, non-incremental updates will cause a
    latency and even obvious performance issues in the case of heavy service traffic.
  + **ipvs**: Optimized kube-proxy mode with higher throughput and faster speed. This mode supports incremental updates
    and can keep connections uninterrupted during service updates. It is suitable for large-sized clusters.

* `custom_san` - (Optional, List) Specifies the custom san to add to certificate (array of string).

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether to enable IPv6 in the cluster.
  Changing this parameter will create a new cluster resource.

* `enable_distribute_management` - (Optional, Bool, ForceNew) Specifies whether to enable support for remote clouds.
  Changing this parameter will create a new cluster resource.

* `extend_params` - (Optional, List, ForceNew) Specifies the extended parameter.
  The [object](#cce_cluster_extend_params) structure is documented below.
  Changing this parameter will create a new cluster resource.

* `component_configurations` - (Optional, List) Specifies the kubernetes component configurations.
  For details, see [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-cce/cce_10_0213.html).
  The [object](#cce_cluster_component_configurations) structure is documented below.

* `encryption_config` - (Optional, List, ForceNew) Specifies the encryption configuration.
  The [object](#cce_cluster_encryption_config) structure is documented below.
  Changing this parameter will create a new cluster resource.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the CCE cluster.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.
  Changing this parameter will create a new cluster resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the CCE cluster.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new cluster resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the CCE cluster.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new cluster resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled. Valid values are **true** and **false**.

* `enterprise_project_id` - (Optional, String) The enterprise project ID of the CCE cluster.

* `tags` - (Optional, Map) Specifies the tags of the CCE cluster, key/value pair format.

* `delete_evs` - (Optional, String) Specified whether to delete associated EVS disks when deleting the CCE cluster.
  valid values are **true**, **try** and **false**. Default is **false**.

* `delete_obs` - (Optional, String) Specified whether to delete associated OBS buckets when deleting the CCE cluster.
  valid values are **true**, **try** and **false**. Default is **false**.

* `delete_sfs` - (Optional, String) Specified whether to delete associated SFS file systems when deleting the CCE
  cluster. valid values are **true**, **try** and **false**. Default is **false**.

* `delete_efs` - (Optional, String) Specified whether to unbind associated SFS Turbo file systems when deleting the CCE
  cluster. valid values are **true**, **try** and **false**. Default is **false**.

* `delete_all` - (Optional, String) Specified whether to delete all associated storage resources when deleting the CCE
  cluster. valid values are **true**, **try** and **false**. Default is **false**.

* `lts_reclaim_policy` - (Optional, String) Specified whether to delete LTS resources when deleting the CCE cluster.
  Valid values are:
  + **Delete_Log_Group**: Delete the log group, ignore it if it fails, and continue with the subsequent process.
  + **Delete_Master_Log_Stream**: Delete the the log stream, ignore it if it fails, and continue the subsequent process.
  The default option.
  + **Retain**: Skip the deletion process.

* `hibernate` - (Optional, Bool) Specifies whether to hibernate the CCE cluster. Defaults to **false**. After a cluster
  is hibernated, resources such as workloads cannot be created or managed in the cluster, and the cluster cannot be
  deleted.

<a name="cce_cluster_masters"></a>
The `masters` block supports:

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone of the master node.
  Changing this parameter will create a new cluster resource.

<a name="cce_cluster_extend_params"></a>
The `extend_params` block supports:

* `cluster_az` - (Optional, String, ForceNew) Specifies the AZ of master nodes in the cluster. The value can be:
  + **multi_az**: The cluster will span across AZs. This field is configurable only for high-availability clusters.
  + **AZ of the dedicated cloud computing pool**: The cluster will be deployed in the AZ of Dedicated Cloud (DeC).
  This parameter is mandatory for dedicated CCE clusters.

  Changing this parameter will create a new cluster resource.

* `dss_master_volumes` - (Optional, String, ForceNew) Specifies whether the system and data disks of a master node
  use dedicated distributed storage. If left unspecified, EVS disks are used by default.
  This parameter is mandatory for dedicated CCE clusters.
  It is in the following format:

  ```bash
  <rootVol.dssPoolID>.<rootVol.volType>;<dataVol.dssPoolID>.<dataVol.volType>
  ```

  Changing this parameter will create a new cluster resource.

* `fix_pool_mask` - (Optional, String, ForceNew) Specifies the number of mask bits of the fixed IP address pool
  of the container network model. This field can only be used when `container_network_type` is set to **vpc-router**.
  Changing this parameter will create a new cluster resource.

* `dec_master_flavor` - (Optional, String, ForceNew) Specifies the specifications of the master node
  in the dedicated hybrid cluster.
  Changing this parameter will create a new cluster resource.

* `docker_umask_mode` - (Optional, String, ForceNew) Specifies the default UmaskMode configuration of Docker in a
  cluster. The value can be **secure** or **normal**, defaults to normal.
  Changing this parameter will create a new cluster resource.

* `cpu_manager_policy` - (Optional, String, ForceNew) Specifies the cluster CPU management policy.
  The value can be:
  + **none**: CPU cores will not be exclusively allocated to workload pods.
    Select this value if you want a large pool of shareable CPU cores.
  + **static**: CPU cores can be exclusively allocated to workload pods.
    Select this value if your workload is sensitive to latency in CPU cache and scheduling.In a CCE Turbo cluster,
    this setting is valid only for nodes where common containers, not Kata containers, run.

  Defaults to none.  
  Changing this parameter will create a new cluster resource.

<a name="cce_cluster_component_configurations"></a>
The `component_configurations` block supports:

* `name` - (Required, String) Specifies the component name.

* `configurations` - (Optional, String) Specifies JSON string of the component configurations.

<a name="cce_cluster_encryption_config"></a>
The `encryption_config` block supports:

* `mode` - (Optional, String, ForceNew) Specifies the encryption mode. The value can be: **Default** and **KMS**.

* `kms_key_id` - (Optional, String, ForceNew) Specifies KMS key ID, required if `mode` is set to **KMS**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the cluster resource.

* `status` - Cluster status information.

* `category` - The category of the cluster. The value can be **CCE** and **Turbo**.

* `certificate_clusters` - The certificate clusters. Structure is documented below.

* `certificate_users` - The certificate users. Structure is documented below.

* `eni_subnet_cidr` - The ENI network segment. This value is valid when only one eni_subnet_id is specified.

* `kube_config_raw` - Raw Kubernetes config to be used by kubectl and other compatible tools.

* `support_istio` - Whether Istio is supported in the cluster.

The `certificate_clusters` block supports:

* `name` - The cluster name.

* `server` - The server IP address.

* `certificate_authority_data` - The certificate data.

The `certificate_users` block supports:

* `name` - The user name.

* `client_certificate_data` - The client certificate data.

* `client_key_data` - The client key data.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

Cluster can be imported using the cluster ID, e.g.

```
 $ terraform import huaweicloud_cce_cluster.cluster_1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`delete_efs`, `delete_eni`, `delete_evs`, `delete_net`, `delete_obs`, `delete_sfs` and `delete_all`. It is generally
recommended running `terraform plan` after importing an CCE cluster. You can then decide if changes should be applied to
the cluster, or the resource definition should be updated to align with the cluster. Also you can ignore changes as
below.

```hcl
resource "huaweicloud_cce_cluster" "cluster_1" {
    ...

  lifecycle {
    ignore_changes = [
      delete_efs, delete_obs,
    ]
  }
}
```
