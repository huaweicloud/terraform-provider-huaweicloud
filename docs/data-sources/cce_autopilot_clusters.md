---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_clusters"
description: |-
  Use this data source to get  the list of CCE Autopilot clusters within huaweicloud.
---

# huaweicloud_cce_autopilot_clusters

Use this data source to get  the list of CCE Autopilot clusters within huaweicloud.

## Example Usage

```hcl
data "huaweicloud_cce_autopilot_clusters" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the master node architecture.
  The value can be: **VirtualMachine**.

* `version` - (Optional, String) Specifies the version of the cluster.

* `status` - (Optional, String) Specifies the status of the cluster.
  The value can be: **Available**, **Unavailable**, **Creating**, **Deleting**, **Upgrading**,
  **RollingBack**, **RollbackFailed** and **Error**.

* `detail` - (Optional, String) Specifies whether to get the details of the cluster.
  The value can be **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clusters` - The list of the clusters.

  The [clusters](#clusters_struct) structure is documented below.

<a name="clusters_struct"></a>
The `clusters` block supports:

* `name` - The cluster name.

* `alias` - The alias of the cluster.

* `annotations` - The cluster annotations in the format of key-value pairs.

* `id` - The ID of the cluster.

* `created_at` - The time when the cluster was created.

* `updated_at` - The time when the cluster was updated.

* `flavor` - The cluster flavor.

* `host_network` - The host network of the cluster.

  The [host_network](#spec_host_network_struct) structure is documented below.

* `container_network` - The container network of the cluster.

  The [container_network](#spec_container_network_struct) structure is documented below.

* `category` - The cluster type. Only **Turbo** is supported.

* `type` - The master node architecture.

* `description` - The description of the cluster.

* `version` - The version of the cluster.

* `custom_san` - The custom SAN field in the API server certificate of the cluster.

* `enable_snat` - Whether SNAT is configured for the cluster.
  After this function is enabled, the cluster can access the Internet through a NAT gateway.
  By default, the existing NAT gateway in the selected VPC is used. Otherwise, the system
  automatically creates a NAT gateway of the default specifications, binds an EIP to the NAT
  gateway, and configures SNAT rules.

* `enable_swr_image_access` - Whether the cluster is interconnected with SWR.
  To ensure that your cluster nodes can pull images from SWR, the existing SWR and OBS
  endpoints in the selected VPC are used by default. If not, new SWR and OBS endpoints
  will be automatically created.

* `enable_autopilot` - Whether the cluster is an Autopilot cluster.

* `ipv6_enable` - Whether the cluster uses the IPv6 mode.

* `eni_network` - The ENI network of the cluster.

  The [eni_network](#spec_eni_network_struct) structure is documented below.

* `service_network` - The service network of the cluster.

  The [service_network](#spec_service_network_struct) structure is documented below.

* `authentication` - The configurations of the cluster authentication mode.

  The [authentication](#spec_authentication_struct) structure is documented below.

* `tags` - The cluster tags in the format of key-value pairs.

* `kube_proxy_mode` - The kube proxy mode of the cluster.

* `configurations_override` - The parameter to override the default component configurations in the cluster.

  The [configurations_override](#spec_configurations_override_struct) structure is documented below.

* `extend_param` - The extend param of the cluster.

  The [extend_param](#spec_extend_param_struct) structure is documented below.

* `az` - The AZ of the cluster.

* `platform_version` - The cluster platform version.

* `status` - The status of the cluster.

  The [status](#clusters_status_struct) structure is documented below.

<a name="spec_host_network_struct"></a>
The `host_network` block supports:

* `vpc` - The ID of the VPC used to create a master node.

* `subnet` - The ID of the subnet used to create a master node.

<a name="spec_container_network_struct"></a>
The `container_network` block supports:

* `mode` - The container network type.

<a name="spec_eni_network_struct"></a>
The `eni_network` block supports:

* `subnets` - The list of ENI subnets.

  The [subnets](#eni_network_subnets_struct) structure is documented below.

<a name="eni_network_subnets_struct"></a>
The `subnets` block supports:

* `subnet_id` - The IPv4 subnet ID of the subnet used to create control nodes and containers.

<a name="spec_service_network_struct"></a>
The `service_network` block supports:

* `ipv4_cidr` - The IPv4 CIDR of the service network.

<a name="spec_authentication_struct"></a>
The `authentication` block supports:

* `mode` - The cluster authentication mode.

<a name="spec_configurations_override_struct"></a>
The `configurations_override` block supports:

* `name` - The component name.

* `configurations` - The component configuration items.

  The [configurations](#configurations_override_configurations_struct) structure is documented below.

<a name="configurations_override_configurations_struct"></a>
The `configurations` block supports:

* `name` - The component configuration item name.

* `value` - The component configuration item value.

<a name="spec_extend_param_struct"></a>
The `extend_param` block supports:

* `enterprise_project_id` - The ID of the enterprise project to which the cluster belongs.

<a name="clusters_status_struct"></a>
The `status` block supports:

* `phase` - The phase of the cluster.

* `endpoints` - The access address of kube-apiserver in the cluster.

  The [endpoints](#status_endpoints_struct) structure is documented below.

<a name="status_endpoints_struct"></a>
The `endpoints` block supports:

* `url` - The URL of the endpoint.

* `type` - The type of the endpoint.
