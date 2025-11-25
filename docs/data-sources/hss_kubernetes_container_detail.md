---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_kubernetes_container_detail"
description: |-
  Use this data source to get the detail of HSS kubernetes container within HuaweiCloud.
---

# huaweicloud_hss_kubernetes_container_detail

Use this data source to get the detail of HSS kubernetes container within HuaweiCloud.

## Example Usage

```hcl
variable "container_id" {}

data "huaweicloud_hss_kubernetes_container_detail" "test" {
  container_id = var.container_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `container_id` - (Required, String) Specifies the container ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The container ID.

* `service_name` - The service name.

* `service_username` - The service username.

* `service_password` - The service password.

* `service_port_list` - The list of service ports.
  
  The [service_port_list](#service_port_list_struct) structure is documented below.

* `enable_simulate` - The data simulation is turned off by default. After activation, simulation data will be injected
  into the sandbox.  
  The valid values are as follows:
  + **open**: Enabled.
  + **close**: Disabled.

* `hosts` - The sandbox domain names, separated by commas between domains.

* `extra` - The additional configuration of sandbox.
  
  The [extra](#extra_struct) structure is documented below.

<a name="service_port_list_struct"></a>
The `service_port_list` block supports:

* `desc` - The service name.

* `type` - The type.  
  The valid values are as follows:
  + **http**: HTTP port.
  + **https**: HTTPS port.

* `protocol` - The protocol. The valid values are **tcp** and **udp**. Defaults to **tcp**.

* `user_port` - The user port.

* `port` - The container internal port.

<a name="extra_struct"></a>
The `extra` block supports:

* `openvpn` - The VPN drainage sandbox.
  
  The [openvpn](#extra_openvpn_struct) structure is documented below.

* `linux` - The Linux sandbox.
  
  The [linux](#extra_linux_struct) structure is documented below.

* `rdp` - The RDP sandbox.
  
  The [rdp](#extra_rdp_struct) structure is documented below.

* `mysql` - The MYSQL, MYSQLCHEAT sandbox.
  
  The [mysql](#extra_mysql_struct) structure is documented below.

<a name="extra_openvpn_struct"></a>
The `openvpn` block supports:

* `outside_ip` - The mapped IP.

* `outside_port` - The mapped port.

<a name="extra_linux_struct"></a>
The `linux` block supports:

* `os` - The operating system.  
  The valid values are as follows:
  + **ubt**: Ubuntu.
  + **centos**
  + **debian**
  + **redhat**
  + **opensuse**
  + **kylin**
  + **uos**
  + **euleros**

<a name="extra_rdp_struct"></a>
The `rdp` block supports:

* `proto_env` - The protocol type.  
  The valid values are as follows:
  + **0**: Protocol simulation.
  + **1**: System simulation.

* `system` - The system type. Used when simulating the system.  
  The valid values are as follows:
  + **win7**
  + **win8**
  + **win10**
  + **win-server2012**
  + **win-server2016**

<a name="extra_mysql_struct"></a>
The `mysql` block supports:

* `custom_path` - The custom countermeasure path.
