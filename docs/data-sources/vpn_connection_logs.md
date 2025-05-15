---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_connection_logs"
description: |-
  Use this data source to get the list of VPN connection logs.
---

# huaweicloud_vpn_connection_logs

Use this data source to get the list of VPN connection logs.

## Example Usage

```hcl
variable "connection_id" {}

data "huaweicloud_vpn_connection_logs" "test" {
  connection_id = var.connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `vpn_connection_id` - (Required, String) Specifies the VPN connection ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logs` - Indicates the logs list.

  The [logs](#logs_struct) structure is documented below.

<a name="logs_struct"></a>
The `logs` block supports:

* `raw_message` - Indicates the log information.

* `time` - Indicates the time stamp of log.
