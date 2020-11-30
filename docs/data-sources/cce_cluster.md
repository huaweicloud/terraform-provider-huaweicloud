---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud\_cce\_cluster

Provides details about all clusters and obtains certificate for accessing cluster information.
This is an alternative to `huaweicloud_cce_cluster_v3`

## Example Usage

```hcl
variable "cluster_name" { }

data "huaweicloud_cce_cluster" "cluster" {
  name   = var.cluster_name
  status = "Available"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to obtain the cce clusters. If omitted, the provider-level region will be used.

* `name` -  (Optional, String)The Name of the cluster resource.
 
* `id` - (Optional, String) The ID of container cluster.

* `status` - (Optional, String) The state of the cluster.

* `cluster_type` - (Optional, String) Type of the cluster. Possible values: VirtualMachine, BareMetal.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `billingMode` - Charging mode of the cluster.

* `description` - Cluster description.

* `name` - The name of the cluster in string format.
  
* `flavor_id` - The cluster specification in string format.

* `cluster_version` - The version of cluster in string format.

* `container_network_cidr` - The container network segment.

* `container_network_type` - The container network type: overlay_l2 , underlay_ipvlan or vpc-router.
  
* `subnet_id` - The ID of the subnet used to create the node.

* `highway_subnet_id` - The ID of the high speed network used to create bare metal nodes.

**endpoints**

* `internal` - The address accessed within the user's subnet.

* `external` - Public network access address.

* `certificate_clusters/name` - The cluster name.

* `certificate_clusters/server` - The server IP address.

* `certificate_clusters/certificate_authority_data` - The certificate data.

* `certificate_users/name` - The user name.

* `certificate_users/client_certificate_data` - The client certificate data.

* `certificate_users/client_key_data` - The client key data.

* `kube_config_raw` - Raw Kubernetes config to be used by kubectl and other compatible tools.
