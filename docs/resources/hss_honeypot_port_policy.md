---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_honeypot_port_policy"
description: |-
  Manages a dynamic port honeypot policy resource within HuaweiCloud.
---

# huaweicloud_hss_honeypot_port_policy

Manages a dynamic port honeypot policy resource within HuaweiCloud.

-> Create the dynamic port honeypot policy resource need to meet the following conditions:
  <br/>1. The server that not bound the EIP.
  <br/>2. The HSS premium edition or higher has been enabled on the server.
  <br/>3. The server agent is online. The Windows agent version is 4.0.22 or later, and the Linux agent version
  is 3.2.10 or later.
  <br/>4. For more details, please refer to [document](https://support.huaweicloud.com/intl/en-us/usermanual-hss2.0/hss_01_0600.html).

## Example Usage

```hcl
variable "policy_name" {}
variable "os_type" {}
variable "white_ip" {
  type = list(string)
}
variable "host_id" {
  type = list(string)
}

resource "huaweicloud_hss_honeypot_port_policy" "test" {
  policy_name = var.policy_name
  os_type     = var.os_type

  ports_list {
    port     = 8006
    protocol = "tcp"
  }

  ports_list {
    port     = 8008
    protocol = "tcp"
  }

  white_ip = var.white_ip
  host_id  = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `policy_name` - (Required, String) Specifies the honeypot policy name.
  The name must be unique.

* `os_type` - (Required, String) Specifies the OS type.
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `ports_list` - (Required, List) Specifies the port and protocol list.
  The [ports_list](#honeypot_policy_ports_struct) structure is documented below.

* `white_ip` - (Required, List) Specifies the IP addresses whitelist.

* `host_id` - (Optional, List) Specifies the ID list of the hosts.

* `group_list` - (Optional, List) Specifies the server group ID list.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

<a name="honeypot_policy_ports_struct"></a>
The `ports_list` block supports:

* `port` - (Required, Int) Specifies the port number.

* `protocol` - (Required, String) Specifies the protocol type.
  The valid values are as follows:
  + **tcp**
  + **tcp6**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `host_list` - The host ID list.

* `port_list` - The port and protocol list.
  The [port_list](#honeypot_policy_port_struct_attr) structure is documented below.

<a name="honeypot_policy_port_struct_attr"></a>
The `port_list` block supports:

* `port` - The port number.

* `protocol` - The protocol type.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_hss_honeypot_port_policy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `ports_list`, `host_id`, `group_list`, `enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_honeypot_port_policy" "test" { 
  ...
  
  lifecycle {
    ignore_changes = [
      "ports_list", host_id, group_list, enterprise_project_id,
    ]
  }
}
```
