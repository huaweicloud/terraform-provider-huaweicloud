---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_clusters_risks"
description: |-
  Use this data source to get the list of HSS container kubernetes clusters risks within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_clusters_risks

Use this data source to get the list of HSS container kubernetes clusters risks within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id_list" {
  type = list(string)
}

data "huaweicloud_hss_container_kubernetes_clusters_risks" "test" {
  cluster_id_list = var.cluster_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id_list` - (Required, List) Specifies the cluster ID list.

* `detect_type` - (Optional, String) Specifies the detect type.  
  The valid values are as follows:
  + **image**: Image risk.
  + **baseline**: Baseline risk.
  + **vul**: Vulnerability risk.
  + **event**: Intrusion risk.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of clusters.

* `data_list` - The list of cluster risks.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `images_num` - The number of risky images.

* `baseline_risk_num` - The total number of baseline check risks.

* `vul_num` - The number of vulnerabilities.

* `event_num` - The number of alarm events.

* `protect_node_num` - The number of protected nodes in the cluster.

* `node_total_num` - The total number of nodes in the cluster.

* `cluster_id` - The cluster ID.

* `charging_mode` - The charging mode.  
  The valid values are as follows:
  + **on_demand**: Pay-per-use.
  + **free**: Free.
