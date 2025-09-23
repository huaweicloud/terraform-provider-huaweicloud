---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_service_group_members"
description: |-
  Use this data source to get the list of CFW service group members.
---

# huaweicloud_cfw_service_group_members

Use this data source to get the list of CFW service group members.

## Example Usage

```hcl
variable "group_id" {}

data "huaweicloud_cfw_service_group_members" "test" {
  group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_id` - (Required, String) Specifies the service group ID.

* `key_word` - (Optional, String) Specifies the key word.

* `group_type` - (Optional, String) Specifies the service group type.
  The value can be **0** (custom service group), **1** (predefined service group).

* `item_id` - (Optional, String) Specifies the service group member ID.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

* `protocol` - (Optional, Int) Specifies the protocol type.
  The options are as follows:
  + **6**: TCP;
  + **17**: UDP;
  + **1**: ICMP.

* `source_port` - (Optional, String) Specifies the source port.

* `dest_port` - (Optional, String) Specifies the destination port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The service group member list.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `item_id` - The service group member ID.

* `protocol` - The protocol type.

* `source_port` - The source port.

* `dest_port` - The destination port.

* `description` - The service group member description.
