---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_encode_servers"
description: |-
  Use this data source to get encode server list of CPH server.
---

# huaweicloud_cph_encode_servers

Use this data source to get encode server list of CPH server.

## Example Usage

```hcl
data "huaweicloud_cph_encode_servers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the encode server type. The valid value can be **0** (server), **1** (container).

* `status` - (Optional, String) Specifies the encode server status.
  + **1**: Running
  + **2**: Abnormal
  + **3**: Restarting
  + **4**: Freeze
  + **5**: Shut down
  + **100**, **1014**, **0**: Creating

* `server_id` - (Optional, String) Specifies the CPH server ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `encode_servers` - The encode server list.

  The [encode_servers](#encode_servers_struct) structure is documented below.

<a name="encode_servers_struct"></a>
The `encode_servers` block supports:

* `keypair_name` - The encode server keypair name.

* `encode_server_id` - The encode server ID.

* `server_id` - The CPH server ID.

* `type` - The encode server type.

* `status` - The encode server status.

* `access_infos` - The encode server access list.

  The [access_infos](#encode_servers_access_infos_struct) structure is documented below.

* `encode_server_ipv6` - The server IPv6 of the encode server.

* `encode_server_name` - The server name of the encode server.

* `encode_server_ip` - The encode server IP.

<a name="encode_servers_access_infos_struct"></a>
The `access_infos` block supports:

* `listen_port` - The listen port of the encode server access.

* `public_ip` - The public IP of the encode server access.

* `server_ip` - The server IP of the encode server access.

* `access_port` - The access port of the encode server access.

* `type` - The type of the encode server access.

* `server_ipv6` - The server IPv6 of the encode server access.
