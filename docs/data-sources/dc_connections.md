---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_connections"
description: |-
  Use this data source to get the list of connections.
---

# huaweicloud_dc_connections

Use this data source to get the list of connections.

## Example Usage

```hcl
variable connection_id {}

data "huaweicloud_dc_connections" {
  connection_id = var.connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `connection_id` - (Optional, String) Specifies the ID of the connection.

* `name` - (Optional, String) Specifies the name of the connection.

* `type` - (Optional, String) Specifies the type of the connection.
  The value can be **standard** (a standard connection), **hosting** (an operations connection) or
  **hosted** (a hosted connection).

* `status` - (Optional, String) Specifies the status of the connection.
  The valid values are as follows:
  + **ACTIVE**: The connection is in the normal state.
  + **DOWN**: The port for the connection is in the down state, which may cause line faults.
  + **BUILD**: Operations related to the connection are being performed.
  + **ERROR**: The connection configuration is incorrect. Contact customer service to rectify the fault.
  + **PENDING_DELETE**: The connection is being deleted.
  + **DELETED**: The connection has been deleted.
  + **APPLY**: A request for a connection is submitted.
  + **DENY**: A site survey is rejected because the customer fails to meet the requirements.
  + **PENDING_PAY**: The order for the connection is to be paid.
  + **PAID**: The order for the connection has been paid.
  + **PENDING_SURVEY**: A site survey is required for the connection.

* `hosting_id` - (Optional, String) Specifies operations connection ID by which hosted connections are filtered.

* `port_type` - (Optional, String) Specifies the type of the port used by the connection.
  The value can be **1G**, **10G**, **40G**, or **100G**.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the connections belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `direct_connects` - All connections that match the filter parameters.

  The [direct_connects](#direct_connects_struct) structure is documented below.

<a name="direct_connects_struct"></a>
The `direct_connects` block supports:

* `id` - The ID of the connection.

* `name` - The name of the connection.

* `type` - The type of the connection.

* `status` - The status of the connection.

* `description` - The description of the connection.

* `bandwidth` - The connection bandwidth, in Mbit/s.

* `location` - The access location information of the DC.

* `port_type` - The type of the port used by the connection.

* `provider` - The line carrier of the connection.

* `provider_status` - The status of the carrier's leased line.
  The value can be **ACTIVE** or **DOWN**.

* `support_feature` - Lists the features supported by the connection.

* `vgw_type` - The gateway type of the DC.
  The default value is **default**.

* `vlan` - The VLAN allocated to the hosted connection.

* `hosting_id` - The ID of the operations connection on which the hosted connection is created.

* `device_id` - The ID of the device connected to the connection.

* `lag_id` - The ID of the LAG to which the connection belongs.

* `ies_id` - The ID of an IES edge site.

* `charge_mode` - The billing mode.
  The value can be **prepayment**, **bandwidth**, or **traffic**.

* `peer_location` - The location of the on-premises facility at the other end of the connection.

* `peer_provider` - The carrier connected to the connection.

* `peer_port_type` - The peer port type.

* `public_border_group` - The public border group of the AZ, indicating whether the site is a HomeZones site.

* `email` - The customer email information.

* `onestopdc_status` - The status of a full-service connection.

* `modified_bandwidth` - The new bandwidth after the line bandwidth is changed.

* `change_mode` - The status of a renewal change.

* `ratio_95peak` - The percentage of the minimum bandwidth for 95th percentile billing.

* `enterprise_project_id` - The ID of the enterprise project to which the connection belongs.

* `tags` - The key/value pairs to associate with the connection.

* `apply_time` - The application time of the connection, in RFC3339 format.

* `created_at` - The creation time of the connection, in RFC3339 format.
