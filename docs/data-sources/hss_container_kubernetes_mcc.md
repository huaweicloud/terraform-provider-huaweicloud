---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_mcc"
description: |-
  Use this data source to get the list of HSS container kubernetes multi-cloud clusters within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_mcc

Use this data source to get the list of HSS container kubernetes multi-cloud clusters within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_kubernetes_mcc" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of clusters.

* `data_list` - The list of multi-cloud cluster data.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `provider` - The cluster provider.

* `server` - The cluster apiserver address.

* `image_repo` - The image repository address.

* `status` - The anp-agent connection status.  
  The valid values are as follows:
  + **not_connect**
  + **connect_success**
  + **connect_fail**
  + **connect_success**
  + **connect_disruption**

* `version` - The anp-agent version.

* `current_expiration_date` - The current expiration date.

* `certificate_expiration_date` - The certificate expiration date.
