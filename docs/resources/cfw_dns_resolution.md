---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_dns_resolution"
description: ""
---

# huaweicloud_cfw_dns_resolution

Manages a CFW DNS resolution resource within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "default_dns_servers" {}
variable "custom_dns_servers" {}
variable "health_check_domain_name" {}

resource "huaweicloud_cfw_dns_resolution" "test" {
  fw_instance_id           = var.fw_instance_id
  default_dns_servers      = var.default_dns_servers
  custom_dns_servers       = var.custom_dns_servers
  health_check_domain_name = var.health_check_domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `fw_instance_id` - (Required, String, ForceNew) The ID of the firewall.
  Changing this creates a new resource.

* `default_dns_servers` - (Optional, List) The default DNS servers.

* `custom_dns_servers` - (Optional, List) The custom DNS servers.
  Currently, only two custom DNS server addresses can be specified.

* `health_check_domain_name` - (Optional, String) The health check domain name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the firewall instance ID.

## Import

The DNS resolution resource can be imported using the firewall instance ID, e.g.

```bash
$ terraform import huaweicloud_cfw_dns_resolution.test <fw_instance_id>
```
