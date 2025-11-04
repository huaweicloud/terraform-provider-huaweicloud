---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cluster_protect_info"
description: |-
  Use this data source to get the list of HSS cluster protect information within HuaweiCloud.
---

# huaweicloud_hss_cluster_protect_info

Use this data source to get the list of HSS cluster protect information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_cluster_protect_info" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of cluster protect information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `cluster_name` - The cluster name.

* `cluster_id` - The cluster ID.

* `cluster_version` - The cluster version.

* `protect_status` - The protection status.  
  The valid values are as follows:
  + **unprotected**: Unprotected.
  + **plugin error**: Plugin error.
  + **protected with policy**: Protected with policy.
  + **deploy policy failed**: Deploy policy failed.
  + **protected without policy**: Protected without policy.
  + **uninstall failed**: Uninstall failed.
  + **uninstall**: Uninstall.

* `policy_num` - The number of policies.

* `cluster_status` - The cluster running status.  
  The valid values are as follows:
  + **Available**: Available.
  + **Unavailable**: Unavailable.

* `cluster_type` - The cluster type.  
  The valid values are as follows:
  + **k8s**: Native cluster.
  + **cce**: CCE cluster.
  + **ali**: Alibaba Cloud cluster.
  + **tencent**: Tencent Cloud cluster.
  + **azure**: Microsoft Azure cluster.
  + **aws**: Amazon cluster.
  + **self_built_hw**: Huawei Cloud self-built cluster.
  + **self_built_idc**: IDC self-built cluster.
