---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_kubernetes_deployments"
description: |-
  Use this data source to get the list of HSS kubernetes deployments within HuaweiCloud.
---

# huaweicloud_hss_kubernetes_deployments

Use this data source to get the list of HSS kubernetes deployments within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_kubernetes_deployments" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `deployment_name` - (Optional, String) Specifies the deployment name.

* `namespace_name` - (Optional, String) Specifies the namespace name.

* `cluster_name` - (Optional, String) Specifies the cluster name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `last_update_time` - The last update time.

* `type` - The resource type.  
  The valid values are as follows:
  + **deploy**: Stateless load.
  + **sts**: Conditional load.
  + **job**: Regular job.
  + **cronjob**: Scheduled job.

* `resources_info_list` - The basic information list of resources.

  The [resources_info_list](#resources_info_list_struct) structure is documented below.

<a name="resources_info_list_struct"></a>
The `resources_info_list` block supports:

* `name` - The deployment name.

* `namespace_name` - The namespace name.

* `cluster_name` - The cluster name.

* `status` - The deployment status.  
  The valid values are as follows:
  + **Running**: Normal running.
  + **Failed**: There are exceptions.

* `protect_status` - The protection status.  
  The valid values are as follows:
  + **closed**: Closed.
  + **opened**: Opened.

* `pods_num` - The total number of instances.

* `image_name` - The image name.

* `match_labels` - The labels.

  The [match_labels](#match_labels_struct) structure is documented below.

* `create_time` - The creation time.

* `agent_installed_num` - The number of installed agent instances under the workload.

* `agent_install_failed_num` - The number of failed instances of agent installation under workload.

* `agent_not_install_num` - The number of uninstalled agent instances under the workload.

<a name="match_labels_struct"></a>
The `match_labels` block supports:

* `key` - The label name.

* `val` - The label value.
