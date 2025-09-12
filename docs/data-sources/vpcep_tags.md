---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_tags"
description: |-
  Use this data source to query the tag list of all resources of the same type within HuaweiCloud.
---

# huaweicloud_vpcep_tags

Use this data source to query the tag list of all resources of the same type within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_vpcep_tags" "test" {
  resource_type = "endpoint"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource tags.  
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type to which the tags belong that to be queried.  
  The value can be **endpoint_service** or **endpoint**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of all tags for resources of the same type.  
  The [tags](#vpcep_project_tags) structure is documented below.

<a name="vpcep_project_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

* `values` - All values corresponding to the key.
