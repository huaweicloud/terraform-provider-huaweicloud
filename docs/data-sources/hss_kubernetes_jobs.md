---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_kubernetes_jobs"
description: |-
  Use this data source to get the list of HSS kubernetes jobs within HuaweiCloud.
---

# huaweicloud_hss_kubernetes_jobs

Use this data source to get the list of HSS kubernetes jobs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_kubernetes_jobs" "test" {}
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

* `job_name` - (Optional, String) Specifies the job name.

* `namespace_name` - (Optional, String) Specifies the namespace name.

* `cluster_name` - (Optional, String) Specifies the cluster name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of jobs.

* `last_update_time` - The last update time.

* `job_info_list` - The list of jobs information.

  The [job_info_list](#job_info_list_struct) structure is documented below.

<a name="job_info_list_struct"></a>
The `job_info_list` block supports:

* `name` - The job name.

* `namespace_name` - The namespace name.

* `cluster_name` - The cluster name.

* `status` - The status. Valid values are:
  + **Running**: Normal operation.
  + **Failed**: Abnormal.

* `pods_num` - The number of instances.

* `image_name` - The image name.

* `match_labels` - The labels.

  The [match_labels](#match_labels_struct) structure is documented below.

* `execute_time` - The execution time.

* `create_time` - The creation time.

<a name="match_labels_struct"></a>
The `match_labels` block supports:

* `key` - The label name.

* `val` - The label value.
