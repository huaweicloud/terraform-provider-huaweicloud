---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_flavor_specifications"
description: |-
  Use this data source to get the list of available cluster flavors.
---

# huaweicloud_cce_flavor_specifications

Use this data source to get the list of available cluster flavors.

## Example Usage

```hcl
variable "cluster_type" {}

data "huaweicloud_cce_flavor_specifications" "test" {
  cluster_type = var.cluster_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_type` - (Required, String) Specifies the type of cluster. Value options:
  + **VirtualMachine**: CCE cluster
  + **ARM64**: Kunpeng cluster

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_flavor_specs` - The cluster flavors for sale.

  The [cluster_flavor_specs](#cluster_flavor_specs_struct) structure is documented below.

<a name="cluster_flavor_specs_struct"></a>
The `cluster_flavor_specs` block supports:

* `name` - The flavor name.

* `node_capacity` - The number of nodes in a cluster.

* `is_sold_out` - Whether the cluster flavors are sold out.

* `is_support_multi_az` - Whether the control plane nodes in a cluster can be deployed in different AZs.

* `available_master_flavors` - The control plane node details.

  The [available_master_flavors](#available_master_flavors_struct) structure is documented below.

<a name="available_master_flavors_struct"></a>
The `available_master_flavors` block supports:

* `name` - The control plane node flavor name.

* `azs` - The AZs supported by control plane nodes.

* `az_fault_domains` - The fault domains supported by the AZs where control plane nodes reside.
