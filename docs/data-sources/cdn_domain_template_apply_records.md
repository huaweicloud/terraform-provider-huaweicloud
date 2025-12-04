---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_template_apply_records"
description: |-
  Use this data source to get a list of domain template apply records within HuaweiCloud.
---

# huaweicloud_cdn_domain_template_apply_records

Use this data source to get a list of domain template apply records within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_domain_template_apply_records" "test" {}
```

## Argument Reference

The following arguments are supported:

* `template_id` - (Optional, String) The ID of the domain template.

* `template_name` - (Optional, String) The name of the domain template.

* `operator_id` - (Optional, String) The operation ID of the domain template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of domain template apply records that matched the filter parameters.  
  The [records](#cdn_domain_template_apply_records_attr) structure is documented below.

<a name="cdn_domain_template_apply_records_attr"></a>
The `records` block supports:

* `operator_id` - The operation ID.

* `status` - The result of applying the template.
  + **success** - The template was applied successfully.
  + **fail** - The template application failed.

* `template_id` - The ID of the domain template.

* `template_name` - The name of the domain template.

* `description` - The description of the domain template.

* `apply_time` - The time when the domain template was applied, in RFC3339 format.

* `type` - The type of the domain template.
  + **1** - System preset template.
  + **2** - Tenant custom template.

* `account_id` - The account ID.

* `resources` - The list of resources to which the template was applied.  
  The [resources](#cdn_domain_template_apply_records_resources_attr) structure is documented below.

* `configs` - The configuration of the domain template, in JSON format.

<a name="cdn_domain_template_apply_records_resources_attr"></a>
The `resources` block supports:

* `status` - The status of applying the template to the domain.
  + **success** - The template was applied successfully to the domain.
  + **fail** - The template application to the domain failed.

* `domain_name` - The domain name.

* `error_msg` - The error message.
