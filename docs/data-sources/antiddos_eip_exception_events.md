---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_eip_exception_events"
description: |-
  Use this data source to query the exception events of an EIP protected by Anti-DDoS.
---

# huaweicloud_antiddos_eip_exception_events

Use this data source to query the exception events of an EIP protected by Anti-DDoS.

## Example Usage

```hcl
variable "floating_ip_id" {}
variable "eip_address" {}

data "huaweicloud_antiddos_eip_exception_events" "test" {
  floating_ip_id = var.floating_ip_id
  ip             = var.eip_address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `floating_ip_id` - (Required, String) Specifies the ID of the EIP.

* `ip` - (Optional, String) Specifies the EIP address.

* `sort_dir` - (Optional, String) Specifies the sort direction. The value can be **asc** or **desc**.
  Defaults to **desc**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `logs` - The list of exception events.
  The [logs](#logs_block) structure is documented below.

<a name="logs_block"></a>
The `logs` block supports:

* `start_time` - The start time of the exception event.

* `end_time` - The end time of the exception event.

* `status` - The protection status. The valid values are:
  + `1`: Cleaning.
  + `2`: Blackhole.

* `trigger_bps` - The traffic when the exception event is triggered, in bit/s.

* `trigger_pps` - The packet rate when the exception event is triggered, in pps.

* `trigger_http_pps` - The HTTP request rate when the exception event is triggered, in pps.
