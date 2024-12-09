---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instance_nodes"
description: |-
  Use this data source to get the list of nodes of the GaussDB OpenGauss instance.
---

# huaweicloud_gaussdb_opengauss_instance_nodes

Use this data source to get the list of nodes of the GaussDB OpenGauss instance.

## Example Usage

```hcl
varivable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_instance_nodes" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `component_type` - (Optional, String) Specifies the component type.
  The default value is **ALL**. Value options:
  + **ALL**: All component types are queried.
  + **CN**: CN component type is queried.
  + **DN**: DN component type is queried.
  + **CM**: CMS component type is queried.
  + **GTM**: GTM component type is queried.
  + **ETCD**: ETCD component type is queried.

* `availability_zone_id` - (Optional, String) Specifies the ID of the AZ where the primary component is located.
  + The default value is **ALL**, indicating that component information of nodes in all AZs of the instance is queried.
  + When you query the AZ where a primary DN is located, the information of all DNs in the same shard as the primary DN
    is displayed.
  + When you query the AZ where a CN is located, only the CN information in the AZ is displayed.
  + When you query the AZ where a component (except CNs or DNs) is located, information about all components of the same
    type is returned. If there is no such a component, no information is returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - Indicates the list of nodes.

  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `status` - Indicates the node status.

* `description` - Indicates the node description.

* `availability_zone_id` - Indicates the code of the AZ where the node is.

* `components` - Indicates the list of components.

  The [components](#nodes_components_struct) structure is documented below.

<a name="nodes_components_struct"></a>
The `components` block supports:

* `id` - Indicates the component ID.
  + **Global Transaction Manager (GTM)**: manages the status of transactions.
  + **Cluster Management Server (CMS)**: manages the instance status.
  + **Data node (DN)**: stores and queries table data.
  + **Coordinator node (CN)**: stores database metadata, distributes and executes query tasks, and then returns the query
    results from DNs to applications.
  + **Editable Text Configuration Daemon (ETCD)**: serves as a distributed key-value storage system used for configuration
    sharing and service discovery (registration and search).

* `role` - Indicates the node type.
  The value can be **master** or **slave**, indicating the primary node and standby node respectively.

* `status` - Indicates the component status.
  + **Primary**: The component is primary.
  + **Normal**: The component is normal.
  + **Down**: The component is down.
  + **Standby**: The component is standby.
  + **StateFollower**: The ETCD is standby.
  + **StateLeader**: The ETCD is primary.
  + **StateCandidate**: The ETCD is in arbitration.

* `distributed_id` - Indicates the Group ID, which is used to identify whether DNs are in the same shard.
  This parameter is suitable only for DNs. For other components, the value is an empty string.

* `type` - Indicates the node type.
  The value can be **DN**, **CN**, **GTM**, **CM**, **ETCD**.

* `detail` - Indicates the node details.
