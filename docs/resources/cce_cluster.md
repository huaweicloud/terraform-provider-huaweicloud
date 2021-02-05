---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud\_cce\_cluster

Provides a CCE cluster resource.
This is an alternative to `huaweicloud_cce_cluster_v3`

## Basic Usage

```hcl
resource "huaweicloud_vpc" "myvpc" {
  name = "vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "mysubnet" {
  name          = "subnet"
  cidr          = "192.168.0.0/16"
  gateway_ip    = "192.168.0.1"

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

## Cluster With Eip

```hcl
resource "huaweicloud_vpc" "myvpc" {
  name = "vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "mysubnet" {
  name          = "subnet"
  cidr          = "192.168.0.0/16"
  gateway_ip    = "192.168.0.1"

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

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the cce cluster resource. If omitted, the provider-level region will be used. Changing this creates a new cce cluster resource.

* `name` - (Required, String, ForceNew) Cluster name. Changing this parameter will create a new cluster resource.

* `flavor_id` - (Required, String, ForceNew) Cluster specifications. Changing this parameter will create a new cluster resource. Possible values:

	* `cce.s1.small` - small-scale single cluster (up to 50 nodes).
	* `cce.s1.medium` - medium-scale single cluster (up to 200 nodes).
	* `cce.s1.large` - large-scale single cluster (up to 1000 nodes).
	* `cce.s2.small` - small-scale HA cluster (up to 50 nodes).
	* `cce.s2.medium` - medium-scale HA cluster (up to 200 nodes).
	* `cce.s2.large` - large-scale HA cluster (up to 1000 nodes).
	* `cce.t1.small` - small-scale single physical machine cluster (up to 10 nodes).
	* `cce.t1.medium` - medium-scale single physical machine cluster (up to 100 nodes).
	* `cce.t1.large` - large-scale single physical machine cluster (up to 500 nodes).
	* `cce.t2.small` - small-scale HA physical machine cluster (up to 10 nodes).
	* `cce.t2.medium` - medium-scale HA physical machine cluster (up to 100 nodes).
	* `cce.t2.large` - large-scale HA physical machine cluster (up to 500 nodes).

* `cluster_version` - (Optional, String, ForceNew) For the cluster version, defaults to the latest supported version. To learn which cluster
versions are available, choose Dashboard > Buy Cluster on the CCE console. Changing this parameter will create a new cluster resource.

* `cluster_type` - (Optional, String, ForceNew) Cluster Type, possible values are VirtualMachine, BareMetal and ARM64. Defaults to *VirtualMachine*.
  Changing this parameter will create a new cluster resource.

* `description` - (Optional, String) The Cluster description.

* `billing_mode` - (Optional, Int, ForceNew) Charging mode of the cluster, which is 0 (on demand). Changing this parameter will create a new cluster resource.

* `extend_param` - (Optional, Map, ForceNew) Extended parameter. Changing this parameter will create a new cluster resource.

* `vpc_id` - (Required, String, ForceNew) The ID of the VPC used to create the node. Changing this parameter will create a new cluster resource.

* `subnet_id` - (Required, String, ForceNew) The ID of the subnet used to create the node  which should be configured with a *DNS address*.
  Changing this parameter will create a new cluster resource.

* `highway_subnet_id` - (Optional, String, ForceNew) The ID of the high speed network used to create bare metal nodes. Changing this parameter will create a new cluster resource.

* `service_network_cidr` - (Optional, String, ForceNew) Service network segment. Changing this parameter will create a new cluster resource.

* `container_network_type` - (Required, String, ForceNew) Container network parameters. Possible values:

	* `overlay_l2` - An overlay_l2 network built for containers by using Open vSwitch(OVS)
	* `underlay_ipvlan` - An underlay_ipvlan network built for bare metal servers by using ipvlan.
	* `vpc-router` - An vpc-router network built for containers by using ipvlan and custom VPC routes.

* `container_network_cidr` - (Optional, String, ForceNew) Container network segment. Changing this parameter will create a new cluster resource.

* `authentication_mode` - (Optional, String, ForceNew) Authentication mode of the cluster, possible values are x509 and rbac. Defaults to *rbac*.
    Changing this parameter will create a new cluster resource.

* `authenticating_proxy_ca` - (Optional, String, ForceNew) CA root certificate provided in the authenticating_proxy mode. The CA root certificate
	is encoded to the Base64 format. Changing this parameter will create a new cluster resource.

* `multi_az` - (Optional, Bool, ForceNew) Enable multiple AZs for the cluster, only when using HA flavors. 
  Changing this parameter will create a new cluster resource. This parameter and `masters` are alternative

* `masters` - (Optional, List, ForceNew) Advanced configuration of master nodes. Changing this creates a new cluster.
  This parameter and `multi_az` are alternative.

* `eip` - (Optional, String, ForceNew) EIP address of the cluster. Changing this parameter will create a new cluster resource.

* `kube_proxy_mode` - (Optional, String, ForceNew) Service forwarding mode. Two modes are available:

  - iptables: Traditional kube-proxy uses iptables rules to implement service load balancing. In this mode, too many iptables rules will be generated when many services are deployed. In addition, non-incremental updates will cause a latency and even obvious performance issues in the case of heavy service traffic.
  - ipvs: Optimized kube-proxy mode with higher throughput and faster speed. This mode supports incremental updates and can keep connections uninterrupted during service updates. It is suitable for large-sized clusters.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the cce cluster. Changing this creates a new cluster.

The `masters` block supports:

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone of the master node. Changing this creates a new cluster.
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` -  Id of the cluster resource.

  * `status` -  Cluster status information.

  * `certificate_clusters/name` - The cluster name.

  * `certificate_clusters/server` - The server IP address.

  * `certificate_clusters/certificate_authority_data` - The certificate data.

  * `certificate_users/name` - The user name.

  * `certificate_users/client_certificate_data` - The client certificate data.

  * `certificate_users/client_key_data` - The client key data.

  * `security_group_id` - Security group ID of the cluster.

  * `kube_config_raw` - Raw Kubernetes config to be used by kubectl and other compatible tools.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.
- `delete` - Default is 30 minute.

## Import

 Cluster can be imported using the cluster id, e.g.
 ```
 $ terraform import huaweicloud_cce_cluster.cluster_1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d  
```

