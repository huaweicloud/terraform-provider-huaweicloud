---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_address_group_members"
description: |-
  Use this data source to get the list of CFW address group members.
---

# huaweicloud_cfw_address_group_members

Use this data source to get the list of CFW address group members.

## Example Usage

```hcl
variable "group_id" {}

data "huaweicloud_cfw_address_group_members" "test" {
  group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_id` - (Required, String) Specifies the ID of the IP address group.

* `key_word` - (Optional, String) Specifies the keyword.

* `address` - (Optional, String) Specifies the IP address

* `item_id` - (Optional, String) Specifies the address group member ID.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

* `query_address_set_type` - (Optional, String) Specifies the query address group type.
  + **0** means custom define address set.
  + **1** means predefined address set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The IP address group member list.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `item_id` - The ID of an address group member.

* `description` - The address group member description.

* `address_type` - The address type.

* `address` - The IP address.
