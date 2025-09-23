---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instances_filter"
description: |-
  Use this data source to query the dedicated instance list within HuaweiCloud.
---

# huaweicloud_apig_instances_filter

Use this data source to query the dedicated instance list within HuaweiCloud.

## Example Usage

### Query all instances

```hcl
data "huaweicloud_apig_instances_filter" "test" {}
```

### Query instance list by tags

```hcl
variable "tags" {
  type = list(object({
    key    = string
    values = optional(list(string), [])
  }))
}

data "huaweicloud_apig_instances_filter" "test" {
  dynamic "tags" {
    for_each = var.tags

    content {
      key    = tags.value.key
      values = tags.value.values
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `without_any_tag` - (Optional, Bool) Specifies whether to query resources without tags.  
  Defaults to **false**.
  + **true**: Only all resources without tags are queried.
  + **false**: All resources are queried.

* `tags` - (Optional, List) Specifies the list of the tags to be queried.  
  The [tags](#data_instances_filter_tags) structure is documented below.

* `matches` - (Optional, List) Specifies the fields to be queried.  
  The [matches](#data_instances_filter_matches) structure is documented below.

<a name="data_instances_filter_tags"></a>
The `tags` block supports:

* `key` - (Optional, String) Specifies the key of the tag.

* `values` - (Optional, List) Specifies the list of values of the tag.

<a name="data_instances_filter_matches"></a>
The `matches` block supports:

* `key` - (Optional, String) Specifies the key to be matched.  
  The valid values are as follows:
  + **resource_name**

* `value` - (Optional, String) Specifies the value of the matching field. Fuzzy match is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - All dedicated instances that match the filter parameters.

  The [instances](#data_instances_filter_struct) structure is documented below.

<a name="data_instances_filter_struct"></a>
The `instances` block supports:

* `resource_id` - The ID of the instance.

* `resource_name` - The name of the instance.

* `tags` - The tag list associated with the instance.

  The [tags](#data_instances_filter_tags_struct) structure is documented below.

<a name="data_instances_filter_tags_struct"></a>
The `tags` block supports:

* `key` - The key of the instance tag.

* `value` - The value of the instance tag.
