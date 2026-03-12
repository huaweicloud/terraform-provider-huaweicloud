---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_host_nets"
description: |-
  Use this data source to get the list of host network metrics in DWS within HuaweiCloud.
---

# huaweicloud_dws_host_nets

Use this data source to get the list of host network metrics in DWS within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_host_nets" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region where the host nets are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Optional, String) The cluster ID to be queried.

* `instance_name` - (Optional, String) The instance name to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `host_nets` - The list of the host nets that matched filter parameters.  
  The [host_nets](#dws_host_nets_struct) structure is documented below.

<a name="dws_host_nets_struct"></a>
The `host_nets` block supports:

* `virtual_cluster_id` - The virtual cluster ID.

* `ctime` - The query timestamp in Unix milliseconds.

* `host_id` - The host ID.

* `host_name` - The host name.

* `instance_name` - The instance name.

* `interface_name` - The network interface name.

* `up` - Whether the network interface is up.

* `speed` - The network interface speed in Mbps.

* `recv_packets` - The received packets.

* `send_packets` - The sent packets.

* `recv_drop` - The dropped packets on receiving.

* `recv_rate` - The receiving rate in KB/s.

* `send_rate` - The sending rate in KB/s.

* `io_rate` - The network IO rate in KB/s.
