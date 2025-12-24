---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_local_hosts"
description: |-
  Use this data source to get the list of HSS local image hosts within HuaweiCloud.
---

# huaweicloud_hss_image_local_hosts

Use this data source to get the list of HSS local image hosts within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_local_hosts" "test" {}
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

* `image_digest` - (Optional, String) Specifies the image digest.

* `image_name` - (Optional, String) Specifies the image name.

* `image_version` - (Optional, String) Specifies the image version.

* `host_name` - (Optional, String) Specifies the host name.

* `agent_id` - (Optional, String) Specifies the agent ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of local image hosts.

* `data_list` - The list of local image hosts data.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `agent_status` - The agent status.

* `host_name` - The server name.

* `host_id` - The host ID.

* `version` - The server activated version.

* `private_ip` - The private IP address.

* `public_ip` - The elastic public IP address.

* `docker_name` - The docker name.

* `docker_type` - The docker type.
