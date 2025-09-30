---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_cluster_risks"
description: |-
  Use this data source to get the list of HSS container cluster risks within HuaweiCloud.
---

# huaweicloud_hss_container_cluster_risks

Use this data source to get the list of HSS container cluster risks within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_cluster_risks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `risk_type` - (Optional, String) Specifies the risk type. Valid values are:
  + **risk_assessment**: Risk assessment.
  + **benchmark**: Security compliance.

* `risk_status` - (Optional, String) Specifies the risk status. Valid values are:
  + **risky**: Risky.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `cluster_name` - (Optional, String) Specifies the cluster name.

* `risk_name` - (Optional, String) Specifies the risk name.

* `risk_level` - (Optional, String) Specifies the risk level. Valid values are:
  + **high**: High risk.
  + **medium**: Medium risk.
  + **low**: Low risk.
  + **tips**: Tips.

* `risk_category` - (Optional, String) Specifies the risk category. Valid values are:
  + **control_plane**: Control plane.
  + **access_control**: Access control.
  + **network**: Network.
  + **workload**: Workload.
  + **secrets**: Secrets management.
  + **node_escape**: Node escape.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of cluster risk information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `risk_id` - The risk ID.

* `risk_name` - The risk name.

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `risk_level` - The risk level.

* `risk_category` - The risk category.

* `risk_num` - The number of risks.

* `last_scan_time` - The last scan time.

* `description` - The risk description.

* `remediation` - The risk remediation suggestion.
