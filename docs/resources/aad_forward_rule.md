---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_forward_rule"
description: ""
---

# huaweicloud_aad_forward_rule

Manages a forward rule resource of Advanced Anti-DDos service within HuaweiCloud.

## Example Usage

```hcl
variable "aad_instance_id" {}
variable "aad_ip_address" {}

resource "huaweicloud_aad_forward_rule" "test" {
  instance_id         = var.aad_instance_id
  ip                  = var.aad_ip_address
  forward_protocol    = "udp"
  forward_port        = 808
  source_port         = 888
  source_ip_addresses = "1.1.1.1,2.2.2.2"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of advanced Anti-DDoS instance.
  Changing this will create a new rule resource.

* `ip` - (Required, String, ForceNew) Specifies the public IP address to which Advanced Anti-DDoS instance
  belongs. Changing this will create a new rule resource.

* `forward_protocol` - (Required, String) Specifies the forward protocol.
  The valid values are **tcp** and **udp**.

* `forward_port` - (Required, Int) Specifies the forward port.
  The valid value is range from `1` to `65,535`.

* `source_port` - (Required, Int) Specifies the source port.
  The valid value is range from `1` to `65,535`.

* `source_ip` - (Required, String) Specifies the source IP addresses, separated by commas (,).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of forward rule.

* `lb_method` - The LVS forward policy.

## Import

Rule can be imported using the `id` (combination of `instance_id`, `ip`, `forward_protocol` and `forward_port`),
separated by slashes (/), e.g.

```bash
terraform import huaweicloud_dds_database_user.test <instance_id>/<ip>/<forward_protocol>/<forward_port>
```
