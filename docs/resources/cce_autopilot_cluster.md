---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_cluster"
description: |-
  Manages a CCE Autopilot cluster resource within huaweicloud.
---

# huaweicloud_cce_autopilot_cluster

Manages a CCE Autopilot cluster resource within huaweicloud.

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
  vpc_id     = huaweicloud_vpc.myvpc.id
}

resource "huaweicloud_cce_autopilot_cluster" "mycluster" {
  name        = "cluster"
  flavor      = "cce.autopilot.cluster"
  description = "created by terraform"

  host_network {
    vpc    = huaweicloud_vpc.myvpc.id
    subnet = huaweicloud_vpc_subnet.mysubnet.id
  }

  container_network {
    mode = "eni"
  }

  eni_network {
    subnets {
      subnet_id = huaweicloud_vpc_subnet.mysubnet.ipv4_subnet_id
    }
  }

  tags = {
    "foo" = "bar"
    "key" = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE autopilot cluster resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new cluster resource.

* `name` - (Required, String, NonUpdatable) Specifies the cluster name. Enter 4 to 128 characters starting with a lowercase
  letter and not ending with a hyphen (-). Only lowercase letters, digits, and hyphens (-) are allowed.

* `flavor` - (Required, String, NonUpdatable) Specifies the cluster flavor. Only **cce.autopilot.cluster** is supported.

* `host_network` - (Required, List, NonUpdatable) Specifies the host network of the cluster.
  The [host_network](#autopilot_cluster_host_networks) structure is documented below.

* `container_network` - (Required, List, NonUpdatable) Specifies the container network of the cluster.
  The [container_network](#autopilot_cluster_container_network) structure is documented below.

* `eip_id` - (Optional, String) Specifies the EIP ID of the cluster.

* `alias` - (Optional, String) Specifies the alias of the cluster. Enter 4 to 128 characters starting
  with a lowercase letter and not ending with a hyphen (-). Only lowercase letters, digits, and hyphens (-) are allowed.
  If not specified, the alias is the same as the cluster name.

* `annotations` - (Optional, Map, NonUpdatable) Specifies the cluster annotations in the format of key-value pairs.

* `category` - (Optional, String, NonUpdatable) Specifies the cluster type. Only **Turbo** is supported.

* `type` - (Optional, String, NonUpdatable) Specifies the master node architecture. The value can be:
  + **VirtualMachine**: Indicates the master node is an x86 server.

* `description` - (Optional, String) Specifies the description of the cluster.

* `version` - (Optional, String, NonUpdatable) Specifies the version of the cluster.
  If not specified, a cluster of the latest version will be created.

* `custom_san` - (Optional, List) Specifies the custom SAN field in the API server certificate of the cluster.

* `enable_snat` - (Optional, Bool, NonUpdatable) Specifies whether SNAT is configured for the cluster.
  After this function is enabled, the cluster can access the Internet through a NAT gateway.
  By default, the existing NAT gateway in the selected VPC is used. Otherwise,
  the system automatically creates a NAT gateway of the default specifications,
  binds an EIP to the NAT gateway, and configures SNAT rules.

* `enable_swr_image_access` - (Optional, Bool, NonUpdatable) Specifies whether the cluster is interconnected with SWR.
  To ensure that your cluster nodes can pull images from SWR, the existing SWR and OBS endpoints in the selected
  VPC are used by default. If not, new SWR and OBS endpoints will be automatically created.

* `enable_autopilot` - (Optional, Bool, NonUpdatable) Specifies whether the cluster is an Autopilot cluster,
  defaults to **true**.

* `ipv6_enable` - (Optional, Bool, NonUpdatable) Specifies whether the cluster uses the IPv6 mode.

* `eni_network` - (Optional, List) Specifies the ENI network of the cluster.
  The [eni_network](#autopilot_cluster_eni_network) structure is documented below.

* `service_network` - (Optional, List, NonUpdatable) Specifies the service network of the cluster.
  The [service_network](#autopilot_cluster_service_network) structure is documented below.

* `authentication` - (Optional, List, NonUpdatable) Specifies the configurations of the cluster authentication mode.
  The [authentication](#autopilot_cluster_authentication) structure is documented below.

* `tags` - (Optional, Map) Specifies the cluster tags in the format of key-value pairs.

* `kube_proxy_mode` - (Optional, String, NonUpdatable) Specifies the kube proxy mode of the cluster.
  The value can be: **iptables**.

* `extend_param` - (Optional, List, NonUpdatable) Specifies the extend param of the cluster.
  The [extend_param](#autopilot_cluster_extend_param) structure is documented below.

* `configurations_override` - (Optional, List, NonUpdatable) Specifies the this parameter to override
  the default component configurations in the cluster.
  The [configurations_override](#autopilot_cluster_configurations_override) structure is documented below.

* `deletion_protection` - (Optional, Bool, NonUpdatable) Specifies whether to enable deletion protection for the cluster.

* `delete_efs` - (Optional, String) Specifies whether to delete the SFS Turbo volume.
  The value can be:
  + **true** or **block**: The system starts to delete the object. If the deletion fails, subsequent processes are blocked.

  + **try**: The system starts to delete the object. If the deletion fails, no deletion retry is performed,
    and subsequent processes will proceed.

  + **false** or **skip**: The deletion is skipped. This is the default option.

* `delete_eni` - (Optional, String) Specifies whether to delete the ENI port.
  The value can be:
  + **true** or **block**: The system starts to delete the object. If the deletion fails, subsequent processes are blocked.
    This is the default option.

  + **try**: The system starts to delete the object. If the deletion fails, no deletion retry is performed,
    and subsequent processes will proceed.

  + **false** or **skip**: The deletion is skipped.

* `delete_net` - (Optional, String) Specifies whether to delete the cluster service or ingress resources,
  such as a load balancer. The value can be:
  + **true** or **block**: The system starts to delete the object. If the deletion fails, subsequent processes are blocked.
    This is the default option.

  + **try**: The system starts to delete the object. If the deletion fails, no deletion retry is performed,
    and subsequent processes will proceed.

  + **false** or **skip**: The deletion is skipped.

* `delete_obs` - (Optional, String) Specifies whether to delete the OBS volume.
  The value can be:
  + **true** or **block**: The system starts to delete the object. If the deletion fails, subsequent processes are blocked.

  + **try**: The system starts to delete the object. If the deletion fails, no deletion retry is performed,
    and subsequent processes will proceed.

  + **false** or **skip**: The deletion is skipped. This is the default option.

* `delete_sfs30` - (Optional, String) Specifies whether to delete the SFS 3.0 volume.
  The value can be:
  + **true** or **block**: The system starts to delete the object. If the deletion fails, subsequent processes are blocked.

  + **try**: The system starts to delete the object. If the deletion fails, no deletion retry is performed,
    and subsequent processes will proceed.

  + **false** or **skip**: The deletion is skipped. This is the default option.

* `lts_reclaim_policy` - (Optional, String) Specifies whether to delete the LTS resource, such as a log group or
  a log stream. The value can be:
  + **Delete_Log_Group**: The system starts to delete a log group. If the deletion fails, no deletion retry is performed,
    and subsequent processes will proceed.

  + **Delete_Master_Log_Stream**: The system starts to delete a master log stream. If the deletion fails,
    no deletion retry is performed, and subsequent processes will proceed. This is the default option.

  + **Retain**: The deletion is skipped.

<a name="autopilot_cluster_host_networks"></a>
The `host_network` block supports:

* `vpc` - (Required, String, NonUpdatable) Specifies the ID of the VPC used to create a master node.

* `subnet` - (Required, String, NonUpdatable) Specifies ID of the subnet used to create a master node.

<a name="autopilot_cluster_container_network"></a>
The `container_network` block supports:

* `mode` - (Required, String, NonUpdatable) Specifies the container network type. The value can be: **eni**.

<a name="autopilot_cluster_eni_network"></a>
The `eni_network` block supports:

* `subnets` - (Required, List) Specifies the list of ENI subnets.
  The [subnets](#autopilot_cluster_eni_network_subnets) structure is documented below.

<a name="autopilot_cluster_eni_network_subnets"></a>
The `subnets` block supports:

* `subnet_id` - (Required, String) Specifies the IPv4 subnet ID of the subnet used to create control
  nodes and containers.

<a name="autopilot_cluster_service_network"></a>
The `service_network` block supports:

* `ipv4_cidr` - (Optional, String, NonUpdatable) Specifies the IPv4 CIDR of the service network.
  If not specified, the default value 10.247.0.0/16 will be used.

<a name="autopilot_cluster_authentication"></a>
The `authentication` block supports:

* `mode` - (Optional, String, NonUpdatable) Specifies the cluster authentication mode.
  The default value is **rbac**.

<a name="autopilot_cluster_extend_param"></a>
The `extend_param` block supports:

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project to which the
  cluster belongs.

<a name="autopilot_cluster_configurations_override"></a>
The `configurations_override` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the component name.

* `configurations` - (Optional, List, NonUpdatable) Specifies the component configuration items.
  The [configurations](#autopilot_cluster_configurations_override_configurations) structure is documented below.

<a name="autopilot_cluster_configurations_override_configurations"></a>
The `configurations` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the component configuration item name.

* `value` - (Optional, String, NonUpdatable) Specifies the component configuration item value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the cluster resource.
  
* `platform_version` - The cluster platform version.

* `created_at` - The time when the cluster was created.

* `updated_at` - The time when the cluster was updated.

* `az` - The AZ of the cluster.

* `status` - The status of the cluster.
  The [status](#autopilot_cluster_status) structure is documented below.

<a name="autopilot_cluster_status"></a>
The `status` block supports:

* `phase` - The phase of the cluster.

* `endpoints` - The access address of kube-apiserver in the cluster.
  The [endpoints](#autopilot_cluster_status_endpoints) structure is documented below.

<a name="autopilot_cluster_status_endpoints"></a>
The `endpoints` block supports:

* `url` - The phase of the cluster.

* `type` - The access address of kube-apiserver in the cluster.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The autopilot cluster can be imported using the cluster ID, e.g.

```bash
 $ terraform import huaweicloud_cce_autopilot_cluster.mycluster <cluster_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`enable_snat`, `enable_swr_image_access`, `eip_id`, `delete_efs`, `delete_eni`, `delete_net`, `delete_obs`, `delete_sfs30`
and `lts_reclaim_policy`. It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition should be updated to align
with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cce_autopilot_cluster" "mycluster" {
    ...

  lifecycle {
    ignore_changes = [
      enable_snat, delete_efs, delete_obs,
    ]
  }
}
```
