---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_resource_instances"
description: |-
  Use this data source to get the list of VPN resource instances.
---

# huaweicloud_vpn_resource_instances

Use this data source to get the list of VPN resource instances.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_vpn_resource_instances" "test" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  Valid values are **vpn-gateway**, **customer-gateway**, **vpn-connection**, **p2c-vpn-gateways**.

* `without_any_tag` - (Optional, Bool) Specifies whether to filter instances.
  + If this parameter is set to **true**, all resources without tags are queried. The `tags` field is ignored.
  + If this parameter is set to **false**, all resources are queried or resources are filtered by `tags` or `matches`.

  Defaults to **false**.

* `tags` - (Optional, List) Specifies the tag list.
  A maximum of **20** tags can be specified.

  The [tags](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the search field, including a key and a value.

  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.
  The value is a string of **1** to **128** characters.

* `values` - (Required, List) Specifies the value list of the tag.
  If values is an empty list, it indicates any value. The relationship between values is **OR**.
  The value is a sting of **0** to **255** characters. A maximum of **20** values can be specified.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the match key.
  The value can be **resource_name**.

* `value` - (Required, String) Specifies the match value.
  The value is a sting of **0** to **255** characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates the resource object list.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - Indicates the resource ID.

* `resource_name` - Indicates the resource name.

* `tags` - Indicates the tag list.

  The [tags](#resources_tags_struct) structure is documented below.

<a name="resources_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `value` - Indicates the tag value.
