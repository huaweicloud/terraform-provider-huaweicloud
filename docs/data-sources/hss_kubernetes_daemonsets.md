---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_kubernetes_daemonsets"
description: |-
  Use this data source to get the list of HSS kubernetes daemonsets within HuaweiCloud.
---

# huaweicloud_hss_kubernetes_daemonsets

Use this data source to get the list of HSS kubernetes daemonsets within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_kubernetes_daemonsets" "test" {}
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

* `daemonset_name` - (Optional, String) Specifies the daemonset name.

* `namespace_name` - (Optional, String) Specifies the namespace name.

* `cluster_name` - (Optional, String) Specifies the cluster name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of daemonsets.

* `data_list` - The list of daemonsets.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `name` - The daemonset name.

* `namespace_name` - The namespace name.

* `cluster_id` - The cluster ID.

* `cluster_type` - The cluster type. Valid values are:
  + **k8s**: Native cluster.
  + **cce**: CCE cluster.
  + **ali**: Alibaba Cloud cluster.
  + **tencent**: Tencent Cloud cluster.
  + **azure**: Microsoft Azure cluster.
  + **aws**: Amazon cluster.
  + **self_built_hw**: Huawei Cloud self-built cluster.
  + **self_built_idc**: IDC self-built cluster.

* `cluster_name` - The cluster name.

* `status` - The status. Valid values are:
  + **Running**: Normal operation.
  + **Failed**: Abnormal.

* `pods_num` - The number of instances.

* `image_name` - The image name.

* `match_labels` - The labels.

  The [match_labels](#match_labels_struct) structure is documented below.

* `create_time` - The creation time.

<a name="match_labels_struct"></a>
The `match_labels` block supports:

* `key` - The label name.

* `val` - The label value.
