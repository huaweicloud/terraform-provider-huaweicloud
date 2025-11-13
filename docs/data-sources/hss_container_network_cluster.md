---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_network_cluster"
description: |-
  Use this data source to get the list of HSS container network cluster information within HuaweiCloud.
---

# huaweicloud_hss_container_network_cluster

Use this data source to get the list of HSS container network cluster information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_network_cluster" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of cluster information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `cluster_type` - The cluster type.  
  The valid values are as follows:
  + **cce**: CCE cluster.
  + **k8s**: Native cluster.
  + **ali**: Alibaba Cloud cluster.
  + **tencent**: Tencent Cloud cluster.
  + **azure**: Microsoft Azure cluster.
  + **aws**: Amazon cluster.
  + **self_built_hw**: Huawei Cloud self-built cluster.
  + **self_built_idc**: IDC self-built cluster.

* `version` - The cluster version.

* `mode` - The network mode.  
  The valid values are as follows:
  + **overlay_l2**: Container tunnel network.
  + **vpc-router**: VPC network.
  + **eni**: Cloud native network 2.0.
  + **native-network**: K8S native network.

* `namespace_num` - The number of namespaces.

* `policy_num` - The number of policies.

* `protection_status` - The protection status. The valid values are **true** and **false**.
