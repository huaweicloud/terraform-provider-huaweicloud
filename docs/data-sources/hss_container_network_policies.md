---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_network_policies"
description: |-
  Use this data source to get the list of container network policies within HuaweiCloud.
---
# huaweicloud_hss_container_network_policies

Use this data source to get the list of container network policies within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_hss_container_network_policies" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `namespace` - (Optional, String) Specifies the namespace to filter the network policies.

* `keyword` - (Optional, String) Specifies the keyword to filter the network policies by name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of network policies.

* `last_update_time` - The last update time of the network policies in milliseconds.

* `data_list` - The list of network policies.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `policy_id` - The ID of the network policy.

* `name` - The name of the network policy.

* `namespace` - The namespace to which the network policy belongs.

* `policy_content` - The content of the network policy in JSON format.

* `create_time` - The creation time of the network policy.

* `deploy_status` - The deployment status of the network policy.
