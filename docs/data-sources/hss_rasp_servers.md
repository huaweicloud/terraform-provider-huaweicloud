---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_rasp_servers"
description: |-
  Use this data source to get the list of protection servers.
---

# huaweicloud_hss_rasp_servers

Use this data source to get the list of protection servers.

## Example Usage

```hcl
variable "app_status"

data "huaweicloud_hss_rasp_servers" "test" {
  app_status = var.app_status
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `app_status` - (Required, String) Specifies the application protection status.
  The valid values are as follows:
  + **closed**: Protection disabled.
  + **opened**: Protection enabled.

* `host_name` - (Optional, String) Specifies the host name.

* `os_type` - (Optional, String) Specifies the operating system type.
  The valid values are as follows:
  + **linux**
  + **windows**

* `host_ip` - (Optional, String) Specifies the host IP address.

* `app_type` - (Optional, String) Specifies the application type.
  The valid values are as follows:
  + **java**: Java application protection.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - All protection servers that match the filter parameters.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `agent_id` - The agent ID.

* `agent_version` - The agent version.

* `host_name` - The host name.

* `public_ip` - The elastic IP address.

* `private_ip` - The private IP address.

* `os_type` - The operating system type.

* `rasp_status` - The application protection status.
  The valid values are as follows:
  + **0**: Protection being enabled.
  + **2**: Protection successful.
  + **4**: Protection failed (installation failed).
  + **6**: Not protected.
  + **8**: Partially protected.
  + **9**: Protection failed.

* `policy_name` - The protection policy name.

* `is_friendly_user` - Whether the user is a friendly user.

* `agent_support_auto_attach` - Whether the agent supports dynamic loading.

* `agent_status` - The agent status.

* `auto_attach` - Whether dynamic loading is enabled.

* `protect_status` - The protection status.

* `group_id` - The server group ID.

* `group_name` - The server group name.

* `protect_event_num` - The number of protection events.

* `rasp_port` - The RASP port number.
