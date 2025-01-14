---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_phone_connections"
description: |-
  Use this data source to get a list of CPH phone connections.
---

# huaweicloud_cph_phone_connections

Use this data source to get a list of CPH phone connections.

## Example Usage

```hcl
variable "phone_ids" {}
variable "client_type" {}

data "huaweicloud_cph_phone_connections" "test" {
  phone_ids   = var.phone_ids
  client_type = var.client_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `phone_ids` - (Required, List) Specifies the phone ids.

* `client_type` - (Required, String) Specifies the client type applying for access.
  The values are as follows:
  + **ANDROID**
  + **WINDOWS**
  + **H5_MOBILE**
  + **H5_PC**
  + **IOS**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connect_infos` - The phone connection list.

  The [connect_infos](#connect_infos_struct) structure is documented below.

<a name="connect_infos_struct"></a>
The `connect_infos` block supports:

* `phone_id` - The phone ID.

* `access_info` - The phone access information.

  The [access_info](#connect_infos_access_info_struct) structure is documented below.

<a name="connect_infos_access_info_struct"></a>
The `access_info` block supports:

* `access_port` - The access port of cloud phone instance.

* `session_id` - The session ID of this access.

* `access_time` - The time of this access.

* `ticket` - The ticket of this access.

* `access_ip` - The IP of this access.

* `intranet_ip` - The intranet IP address of this access.

* `access_ipv6` - The IPv6 address of this access.
