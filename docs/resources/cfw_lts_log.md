---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_lts_log"
description: ""
---

# huaweicloud_cfw_lts_log

Manages a CFW lts log resource within HuaweiCloud.

## Example Usage

```hcl
variable "lts_attack_log_stream_id" {}
variable "lts_flow_log_stream_id" {}
variable "lts_access_log_stream_id" {}
variable "fw_instance_id" {}
variable "lts_log_group_id" {}

resource "huaweicloud_cfw_lts_log" "test" {
  fw_instance_id               = var.fw_instance_id
  lts_log_group_id             = var.lts_log_group_id
  lts_attack_log_stream_enable = 1
  lts_access_log_stream_enable = 1
  lts_flow_log_stream_enable   = 1
  lts_attack_log_stream_id     = var.lts_attack_log_stream_id
  lts_access_log_stream_id     = var.lts_access_log_stream_id
  lts_flow_log_stream_id       = var.lts_flow_log_stream_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `fw_instance_id` - (Required, String, ForceNew) Specifies the ID of the firewall.
  Changing this creates a new resource.

* `lts_log_group_id` - (Required, String) Specifies the LTS log group ID.

* `lts_attack_log_stream_enable` - (Required, Int) Specifies whether to enable the attack log stream.
  The valid values are `0` and `1`, where `0` means disable and `1` means enable.

* `lts_access_log_stream_enable` - (Required, Int) Specifies whether to enable the access log stream.
  The valid values are `0` and `1`, where `0` means disable and `1` means enable.

* `lts_flow_log_stream_enable` - (Required, Int) Specifies whether to enable the flow log stream.
  The valid values are `0` and `1`, where `0` means disable and `1` means enable.

* `lts_attack_log_stream_id` - (Optional, String) Specifies the attack log stream ID.

* `lts_access_log_stream_id` - (Optional, String) Specifies the access log stream ID.

* `lts_flow_log_stream_id` - (Optional, String) Specifies the flow log stream ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the firewall instance ID.

## Import

The lts log resource can be imported using the firewall instance ID, e.g.

```bash
$ terraform import huaweicloud_cfw_lts_log.test <fw_instance_id>
```
