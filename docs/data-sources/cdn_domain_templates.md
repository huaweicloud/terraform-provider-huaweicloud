---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_templates"
description: |-
  Use this data source to get a list of domain templates within HuaweiCloud.
---

# huaweicloud_cdn_domain_templates

Use this data source to get a list of domain templates within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_domain_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of domain templates.  
  The [templates](#cdn_domain_templates_attr) structure is documented below.

<a name="cdn_domain_templates_attr"></a>
The `templates` block supports:

* `id` - The ID of the domain template.

* `name` - The name of the domain template.

* `description` - The description of the domain template.

* `configs` - The configuration of the domain template, in JSON format.

* `type` - The type of the domain template.  
  Valid values are:
  + **1** - System preset template.
  + **2** - Tenant custom template.

* `account_id` - The account ID.

* `create_time` - The creation time of the domain template, in RFC3339 format.

* `modify_time` - The modification time of the domain template, in RFC3339 format.
