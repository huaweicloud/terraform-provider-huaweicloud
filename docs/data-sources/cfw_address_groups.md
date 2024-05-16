---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_address_groups"
description: |-
  Use this data source to get the list of CFW address groups.
---

# huaweicloud_cfw_address_groups

Use this data source to get the list of CFW address groups.

## Example Usage

```hcl
variable "object_id" {}

data "huaweicloud_cfw_address_groups" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

* `name` - (Optional, String) Specifies the name of the address group.

* `query_address_set_type` - (Optional, Int) Specifies the address group type of the query.
   + **0:** indicates a custom IP address group.
   + **1:** indicates a predefined IP address group.

* `address_type` - (Optional, String) Specifies the IP address type.
  The value can be **0** (IPv4) or **1** (IPv6).

* `address` - (Optional, String) Specifies IP address of the IP address group.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id to which the IP address group belongs.

* `key_word` - (Optional, String) Specifies the keyword of the address group description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `address_groups` - The IP address group list.

  The [address_groups](#data_address_groups_struct) structure is documented below.

<a name="data_address_groups_struct"></a>
The `address_groups` block supports:

* `id` - The ID of the IP address group.

* `name` - The IP address group name.

* `ref_count` - The number of times this address group has been referenced.

* `description` - The address groups description.

* `object_id` - The protected object ID.

* `type` - The address group type.

* `address_type` - The address type.
