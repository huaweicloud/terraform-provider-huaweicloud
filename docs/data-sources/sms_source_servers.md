---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_source_servers"
description: ""
---

# huaweicloud_sms_source_servers

Use this data source to get a list of SMS source servers.

## Example Usage

```hcl
variable "server_name" {}

data "huaweicloud_sms_source_servers" "demo" {
  name = var.server_name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) Specifies the ID of the source server.

* `name` - (Optional, String) Specifies the name of the source server.

* `state` - (Optional, String) Specifies the status of the source server.

* `ip` - (Optional, String) Specifies the IP address of the source server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `servers` - An array of SMS source servers found. Structure is documented below.

The `servers` block supports:

* `name` - The name of the source server.

* `id` - The ID of the source server.

* `ip` - The IP address of the source server.

* `state` - The status of the source server.

* `connected` - Whether the source server is properly connected to SMS.

* `os_type` - The OS type of the source server. The value can be **WINDOWS** and **LINUX**.

* `os_version` - The OS version of the source server, for example, UBUNTU_20_4_64BIT.

* `registered_time` - The UTC time when the source server is registered.

* `agent_version` - The version of Agent installed on the source server.

* `vcpus` - The vcpus count of the source server.

* `memory` - The memory size in MB.

* `disks` - The disk information of the source server. Structure is documented below.

The `disks` blocks support:

* `name` - The disk name, for example, /dev/vda.

* `size` - The disk size in MB.

* `device_type` - The disk type. The value can be **BOOT**, **OS** and **NORMAL**.
