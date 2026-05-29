---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ai_component_detail"
description: |-
  Use this data source to get the AI component detail list of HSS within HuaweiCloud.
---

# huaweicloud_hss_ai_component_detail

Use this data source to get the AI component detail list of HSS within HuaweiCloud.

## Example Usage

```hcl
variable "category" {}
variable "catalogue" {}

data "huaweicloud_hss_ai_component_detail" "test" {
  category  = var.category
  catalogue = var.catalogue
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the asset category.  
  The valid values are as follows:
  + **host**: Host asset.
  + **container**: Container asset.

* `catalogue` - (Required, String) Specifies the AI component category.  
  The valid values are as follows:
  + **app**: Application.
  + **tool**: Tool.

* `server_name` - (Optional, String) Specifies the server name.  
  + When `category` is **host**, it represents the server name.
  + When `category` is **container**, it represents the node name.
  + When `category` is **serverless**, it represents the instance name.

* `server_ip` - (Optional, String) Specifies the server IP address.  
  + When `category` is **host**, it represents the server IP address.
  + When `category` is **container**, it represents the node IP address.
  + When `category` is **serverless**, it represents the instance IP address.

* `ai_application` - (Optional, String) Specifies the AI application name.

* `host_id` - (Optional, String) Specifies the server ID.

* `ai_tool` - (Optional, String) Specifies the AI tool name.

* `type` - (Optional, String) Specifies the AI application type.

* `version` - (Optional, String) Specifies the AI version.

* `installation_path` - (Optional, String) Specifies the installation path.

* `first_scan_time` - (Optional, String) Specifies the first scan time in milliseconds.

* `latest_scan_time` - (Optional, String) Specifies the latest scan time in milliseconds.

* `container_name` - (Optional, String) Specifies the container name.

* `container_id` - (Optional, String) Specifies the container ID.

* `image_name` - (Optional, String) Specifies the image name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The AI component detail list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `server_name` - The server name.

* `server_ip` - The server IP address.

* `ai_application` - The AI application name.

* `ai_tool` - The AI tool name.

* `type` - The AI application type.

* `version` - The version.

* `startup_path` - The application startup path.

* `startup_time` - The application startup time in milliseconds.

* `install_path` - The installation path.

* `cmdline` - The application startup command line.

* `first_scan_time` - The first scan time in milliseconds.

* `latest_scan_time` - The latest scan time in milliseconds.

* `container_name` - The container name.

* `container_id` - The container ID.

* `host_id` - The server ID.

* `pid` - The application process PID.

* `ppid` - The parent process PID of the application process.

* `user` - The user that runs the application.

* `net_info` - The network information listened by the application process.

  The [net_info](#net_info_struct) structure is documented below.

<a name="net_info_struct"></a>
The `net_info` block supports:

* `listen_ip` - The listening IP address of the application process.

* `listen_protocol` - The network protocol listened by the application process.  
  The valid values are as follows:
  + **tcp**: TCP protocol.
  + **udp**: UDP protocol.

* `listen_port` - The listening port of the application process.

* `listen_status` - The listening status of the application process.  
  The valid values are as follows:
  + **established**: Connection established.
  + **closed**: Connection closed.
  + **listening**: Listening.
  + **other**: Intermediate connection state.
