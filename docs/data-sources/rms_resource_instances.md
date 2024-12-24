---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_instances"
description: |-
  Use this data source to get a list of resource instances that are filtered by tags.
---

# huaweicloud_rms_resource_instances

Use this data source to get a list of resource instances that are filtered by tags.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_rms_resource_instances" "test"
{
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  The valid value can be **config:policyAssignments**.

* `without_any_tag` - (Optional, Bool) Specifies if the resource has no tags.
  + If this parameter is set to true, all resources that do not have the specified tags will be returned.
    The `tags` parameter is ignored.
  + If this parameter is set to false or is not specified, it does not take effect.

* `tags` - (Optional, List) Specifies the tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.

* `values` - (Required, List) Specifies the tag values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The resource list.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `tags` - The tags.

  The [tags](#resources_tags_struct) structure is documented below.

* `resource_name` - The resource name.

* `resource_id` - The resource ID.

<a name="resources_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
