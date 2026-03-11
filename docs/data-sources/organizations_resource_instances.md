---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_resource_instances"
description: |-
  Use this data source to get the list of resource instances within HuaweiCloud.
---

# huaweicloud_organizations_resource_instances

Use this data source to get the list of resource instances within HuaweiCloud.

## Example Usage

### Query all resources under the specified resource type

```hcl
data "huaweicloud_organizations_resource_instances" "test"{
  resource_type = "organizations:accounts"
}
```

### Query resources by resource name

```hcl
variable "resource_value" {}

data "huaweicloud_organizations_resource_instances" "test" {
  resource_type = "organizations:accounts"

  matches {
    key   = "resource_name"
    value = var.resource_value
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String) Specifies the resource type.  
  The valid values are as follows:
  + **organizations:policies**
  + **organizations:ous**
  + **organizations:accounts**
  + **organizations:roots**

* `without_any_tag` - (Optional, Bool) Specifies whether only get the resources without tags.  
  Defaults to **false**, meaning all resources will be retrieved.

* `tags` - (Optional, List) Specifies the list of tags to be queried.  
  A maximum of `10` keys can be queried at a time, and each key can contain a maximum of `10` values.  
  The tag key cannot be left blank or be an empty string.  
  Each key must be unique, and each value for a key must be unique.  
  Resources that contain all keys and one or multiple values listed in tags will be found and returned.

  The [tags](#resource_instances_tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the fields to be queried.  
  The [matches](#resource_instances_matches_struct) structure is documented below.

<a name="resource_instances_tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.

* `values` - (Required, List) Specifies the list of values of the tag.

<a name="resource_instances_matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the field name, and must be unique.  
  Currently, only **resource_name** is supported.

* `value` - (Required, String) Specifies the value corresponding to the field name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of resources that match the filter parameters.  
  The [resources](#resource_instances_resources_struct) structure is documented below.

<a name="resource_instances_resources_struct"></a>
The `resources` block supports:

* `resource_id` - Indicates the resource ID.

* `resource_name` - Indicates the resource name.

* `tags` - Indicates the list of resource tags.  
  The [tags](#resource_instances_resource_tags_struct) structure is documented below.

<a name="resource_instances_resource_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `value` - Indicates the value of the tag.
