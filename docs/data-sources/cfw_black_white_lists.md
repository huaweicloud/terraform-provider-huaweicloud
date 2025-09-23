---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_black_white_lists"
description: |-
   Use this data source to get the list of CFW blacklists and whitelists.
---

# huaweicloud_cfw_black_white_lists

Use this data source to get the list of CFW blacklists and whitelists.

## Example Usage

```hcl
variable "object_id" {}
variable "list_type" {}

data "huaweicloud_cfw_black_white_lists" "test" {
  object_id = var.object_id
  list_type = var.list_type      
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

* `list_type` - (Required, Int) Specifies the blacklist/whitelist type.
  The options are `4` (blacklist) and `5` (whitelist).

* `address_type` - (Optional, String) Specifies the IP address type.
  The valid value can be **0** (IPv4).

* `list_id` - (Optional, String) Specifies the blacklist/whitelist ID.

* `address` - (Optional, String) Specifies the IP address.

* `direction` - (Optional, String) Specifies the direction of a black or white address.
  The options are as follows:
  + **0**: source address;
  + **1**: destination address;

* `port` - (Optional, String) Specifies the port.

* `protocol` - (Optional, Int) Specifies The protocol type.
  The options are as follows:
  + **6**: TCP;
  + **17**: UDP;
  + **1**: ICMP;
  + **-1**: any protocol;

* `description` - (Optional, String) Specifies the description.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The blacklist and whitelist records.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `list_id` - The blacklist/whitelist ID.

* `direction` - The direction of a black or white address.

* `address_type` - The IP address type.

* `address` - The IP address.

* `protocol` - The protocol type.

* `port` - The port.

* `description` - The description.
