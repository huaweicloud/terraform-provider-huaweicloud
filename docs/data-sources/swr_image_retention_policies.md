---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_retention_policies"
description: |-
  Use this data source to get the list of SWR image retention policies.
---

# huaweicloud_swr_image_retention_policies

Use this data source to get the list of SWR image retention policies.

## Example Usage

```hcl
variable "organization" {}
variable "repository" {}

data "huaweicloud_swr_image_retention_policies" "test" {
  organization = var.organization
  repository   = var.repository
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `organization` - (Required, String) Specifies the name of the organization to which the image belongs.

* `repository` - (Required, String) Specifies the name of the repository to which the image belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `retention_policies` - All retention policies that match the filter parameters.
  The [retention_policies](#swr_image_retention_policies) structure is documented below.

<a name="swr_image_retention_policies"></a>
The `retention_policies` block supports:

* `id` - The image retention policy ID.

* `algorithm` - The image retention policy matching rule.

* `rules` - The rules of the image retention policy.
  The [rules](#swr_image_retention_policies_rules) structure is documented below.

* `scope` - The reserved field.

<a name="swr_image_retention_policies_rules"></a>
The `rules` block supports:

* `template` - The template of the image retention policy. The value can be **date_rule** and **tag_rule**.

* `params` - The params of matching template.
  + If `template` is **date_rule**, the `params` will be **{"days": "xxx"}**.
  + If `template` is **tag_rule**, the `params` will be **{"num": "xxx"}**.

* `tag_selectors` - The exception images.
  The [tag_selectors](#swr_image_retention_policies_rules_tag_selectors) structure is documented below.

<a name="swr_image_retention_policies_rules_tag_selectors"></a>
The `tag_selectors` block supports:

* `kind` - The matching kind. The value can be **label** or **regexp**.

* `pattern` - The pattern of the matching kind.
  + If `kind` is **label**, this parameter will be the image tag.
  + If `kind` is **regexp**, this parameter will be a regular expression.
