---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_sync_mccs"
description: |-  
  Manages a container kubernetes cluster status synchronization task within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_sync_mccs

Manages a container kubernetes cluster status synchronization task within HuaweiCloud.

-> This resource is only a one-time action resource for HSS container kubernetes cluster status synchronization. Deleting
   this resource will not clear the synchronization status, but will only remove the resource information from the tf
   state file.

## Example Usage

```hcl
variable "cluster_id1" {}
variable "cluster_id2" {}

resource "huaweicloud_hss_container_kubernetes_sync_mccs" "test" {
  total_num             = 2
  enterprise_project_id = "0"

  data_list {
    cluster_id = var.cluster_id1
  }

  data_list {
    cluster_id = var.cluster_id2
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `total_num` - (Required, Int, NonUpdatable) Specifies the total number of clusters to synchronize.

* `data_list` - (Required, List, NonUpdatable) Specifies the list of clusters to synchronize.
  The [data_list](#data_list_object) structure is documented below.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project that the server
  belongs to. The value **0** indicates the default enterprise project. To query servers in all enterprise projects,
  set this parameter to **all_granted_eps**. If you have only the permission on an enterprise project, you need to
  transfer the enterprise project ID to query the server in the enterprise project.
  Otherwise, an error is reported due to insufficient permission.

  -> An enterprise project can be configured only after the enterprise project function is enabled.

<a name="data_list_object"></a>
The `data_list` block supports:

* `cluster_id` - (Required, String) Specifies the ID of the cluster to synchronize.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
