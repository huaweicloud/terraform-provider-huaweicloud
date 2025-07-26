---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_export_task"
description: |-
  Manages an HSS container export task resource within HuaweiCloud.
---

# huaweicloud_hss_container_export_task

Manages an HSS container export task resource within HuaweiCloud.

-> This resource is only a one-time action resource for HSS container export. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_hss_container_export_task" "test" {
  export_headers = [
    ["container_name", "Container Name"],
    ["cluster_name", "Cluster Name"],
    ["status", "Status"]
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `export_headers` - (Required, List, NonUpdatable) Specifies the header information list for exporting container
  data. The type of this field is an `Array<Array<string>>`. Please refer to Example Usage for the format of valid values.
  Valid key values and their corresponding table header names (table header names can be customized): **container_id**,
  **container_name**, **image_name**, **pod_name**, **cluster_name**, **cluster_type**, **status**, **risky**, **low_risk**,
  **medium_risk**, **high_risk**, **fatal_risk**, **create_time**, **restart_count**, **cpu_limit**, and **memory_limit**.

* `cluster_container` - (Optional, String, NonUpdatable) Specifies whether the container is in a cluster. The valid
  values are:
  + **true**: Only containers in a cluster are exported.
  + **false**: Only non-cluster containers are exported.

  This field does not involve a default value.

* `cluster_type` - (Optional, String, NonUpdatable) Specifies the cluster type. Options:
  + **cce**: CCE cluster.
  + **ali**: Alibaba Cloud cluster.
  + **tencent**: Tencent Cloud cluster.
  + **azure**: Microsoft Azure Cloud cluster.
  + **aws**: AWS Cloud cluster.
  + **self_built_hw**: Customer-built cluster on Huawei Cloud.
  + **self_built_idc**: IDC on-premises cluster.

* `cluster_name` - (Optional, String, NonUpdatable) Specifies the name of the cluster to which the container belongs.
  The value contains `1` to `255` characters.

* `container_name` - (Optional, String, NonUpdatable) Specifies the name of the container to export.
  The value contains `1` to `255` characters.

* `pod_name` - (Optional, String, NonUpdatable) Specifies the name of the pod to which the container belongs.
  The value can contain `1` to `512` characters.

* `image_name` - (Optional, String, NonUpdatable) Specifies the name of the container image.
  The value contains `1` to `255` characters.

* `status` - (Optional, String, NonUpdatable) Specifies the container status. Valid values are **Running**, **Waiting**,
  **Terminated**, **Isolated**, and **Paused**.

* `risky` - (Optional, String, NonUpdatable) Specifies whether the container has security risks. Options:
  + **true**: Only containers with security risks are to be exported.
  + **false**: Only containers without security risks are to be exported.

* `create_time` - (Optional, List, NonUpdatable) Specifies the time range for filtering containers.
  The [create_time](#create_time_object) structure is documented below.

* `cpu_limit` - (Optional, String, NonUpdatable) Specifies the CPU limit for the container.
  You can enter `0` to `64` characters. The unit is m, for example, **100m**.

* `memory_limit` - (Optional, String, NonUpdatable) Specifies the memory limit for the container.
  You can enter `0` to `64` characters. The unit is Mi or Gi, for example, **300Mi**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project that the server
  belongs to. The value **0** indicates the default enterprise project. To query servers in all enterprise projects,
  set this parameter to **all_granted_eps**. If you have only the permission on an enterprise project, you need to
  transfer the enterprise project ID to query the server in the enterprise project.
  Otherwise, an error is reported due to insufficient permission.

  -> An enterprise project can be configured only after the enterprise project function is enabled.

* `export_size` - (Optional, Int, NonUpdatable) Specifies the number of containers to export.
  The value ranges from `1` to `100,000`.

<a name="create_time_object"></a>
The `create_time` block supports:

* `start_time` - (Optional, Int) Specifies the start time for filtering containers.

* `end_time` - (Optional, Int) Specifies the end time for filtering containers.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID same as `task_id`.

* `task_id` - The export task ID.

* `task_name` - The export task name.

* `task_status` - The export task status. The valid values are **success**, **failure**, and **running**.

* `file_id` - The export file ID.

* `file_name` - The export file name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
