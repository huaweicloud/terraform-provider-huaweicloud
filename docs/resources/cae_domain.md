---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_domain"
description: |-
  Manage a CAE domain name resource within HuaweiCloud.
---

# huaweicloud_cae_domain

Manage a CAE domain name resource within HuaweiCloud.

## Example Usage

```hcl
variable "environment_id" {}
variable "domain_name" {}

resource "huaweicloud_cae_domain" "test" {
  environment_id = var.environment_id
  name           = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the CAE environment.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the domain name to be associated with the CAE environment.
  Changing this creates a new resource.  
  The maximum length of the domain name is `254` characters.  
  The domain name consists of multiple strings separated by dots (.), and the maximum length of a single string is `63` characters.
  Only letters, digits, and hyphens (-) allowed, and must start with a letter or a digit.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  domain name belongs.  
  Changing this creates a new resource.

-> If the `environment_id` belongs to the non-default enterprise project, this parameter is required and is
   only valid for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `created_at` - The time when the domain name is associated, in RFC3339 format.

## Import

The resource can be imported using `environment_id` and `name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_cae_domain.test <environment_id>/<name>
```

For the domain with the non-default enterprise project ID, its enterprise project ID need to be specified
additionanlly when importing. All fields are separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_cae_domain.test <environment_id>/<name>/<enterprise_project_id>
```
