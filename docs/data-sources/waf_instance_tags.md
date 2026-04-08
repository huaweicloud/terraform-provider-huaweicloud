---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_instance_tags"
description: |-
  Use this data source to get the tags of a WAF resource within HuaweiCloud.
---

# huaweicloud_waf_instance_tags

Use this data source to get the tags of a WAF resource within HuaweiCloud.

## Example Usage

```hcl
variable "resource_type" {} 
variable "resource_id" {}

data "huaweicloud_waf_instance_tags" "test" { 
  resource_type = var.resource_type 
  resource_id   = var.resource_id 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  The value can be **waf-instance** and **waf** .

* `resource_id` - (Required, String) Specifies the resource ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags associated with the resource.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
