---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_template_apply"
description: |-
  Manages a CDN domain template apply resource within HuaweiCloud.
---

# huaweicloud_cdn_domain_template_apply

Manages a CDN domain template apply resource within HuaweiCloud.

-> This resource is a one-time action resource for applying CDN domain template to domains. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "template_id" {}
variable "apply_domain_names" {
  type = list(string)
}

resource "huaweicloud_cdn_domain_template_apply" "test" {
  template_id = var.template_id
  resources   = join(",", var.apply_domain_names)
}
```

## Argument Reference

The following arguments are supported:

* `template_id` - (Required, String) Specifies the ID of the domain template to apply.

* `resources` - (Required, String) Specifies the list of domain names to apply the template.  
  Multiple domain names are separated by commas (,).  
  A maximum of 50 domains can be applied in a single operation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
