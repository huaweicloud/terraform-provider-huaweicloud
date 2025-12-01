---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_all_connections"
description: |-
  Use this data source to get the list of ESW all connections.
---

# huaweicloud_esw_all_connections

Use this data source to get the list of ESW all connections.

## Example Usage

```hcl
data "huaweicloud_esw_all_connections" "test" {}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the all instances. If omitted, the provider-level region
  will be used.

* `instance_id` - (Optional, String) Specifies the ID of the instance.

* `connection_id` - (Optional, String) Specifies the ID of the connection.

* `name` - (Optional, String) Specifies the name of the connection.

* `vpc_id` - (Optional, String) Specifies the vpc ID.

* `virsubnet_id` - (Optional, String) Specifies the subnet ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - Indicates the list of connections.
  The [connections](#connections_struct) structure is documented below.

<a name="connections_struct"></a>
The `connections` block supports:

* `id` - Indicates the ID of the connection.

* `instance_id` - Indicates the ID of the instance.

* `name` - Indicates the name of the connection.

* `project_id` - Indicates the project ID.

* `fixed_ips` - Indicates the downlink network port primary and standby IPs.

* `remote_infos` - Indicates the remote tunnel infos.
  The [remote_infos](#remote_infos_struct) structure is documented below.

* `vpc_id` - Indicates the vpc ID.

* `virsubnet_id` - Indicates the subnet ID.

* `status` - Indicates the status of the connection.

* `created_at` - Indicates the created time of the connection.

* `updated_at` - Indicates the updated time of the connection.

<a name="remote_infos_struct"></a>
The `remote_infos` block supports:

* `segmentation_id` - Indicates the tunnel number for the connection corresponds to the VXLAN network identifier (VNI).

* `tunnel_ip` - Indicates the remote tunnel IP of the ESW instance.

* `tunnel_port` - Indicates the remote tunnel port.

* `tunnel_type` - Indicates the remote tunnel type.
