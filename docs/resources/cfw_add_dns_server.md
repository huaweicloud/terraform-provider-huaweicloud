---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_add_dns_server"
description: |-
  Manages a resource to add DNS server within HuaweiCloud.
---

# huaweicloud_cfw_add_dns_server

Manages a resource to add DNS server within HuaweiCloud.

-> This resource is a one-time action resource used to add DNS server. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "server_ip" {}

resource "huaweicloud_cfw_add_dns_server" "test" {
  fw_instance_id = var.fw_instance_id
  server_ip      = var.server_ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `server_ip` - (Required, String, NonUpdatable) Specifies the DNS server IP address.
  The DNS server IP can be obtained through the query DNS server list interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.
