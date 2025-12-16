---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_retention_policy"
description: ""
---

# huaweicloud_swr_image_retention_policy

Manages a SWR image retention policy within HuaweiCloud.

## Example Usage

```hcl
variable "organization_name" {}
variable "repository_name" {}

resource "huaweicloud_swr_image_retention_policy" "test"{
  organization = var.organization_name
  repository   = var.repository_name
  type         = "date_rule"
  number       = 20

  tag_selectors {
    kind    = "label"
    pattern = "abc"
  }

  tag_selectors {
    kind    = "regexp"
    pattern = "abc*"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization.

  Changing this parameter will create a new resource.

* `repository` - (Required, String, ForceNew) Specifies the name of the repository.

  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the retention policy type.
  Value options: **date_rule**, **tag_rule**.

  Changing this parameter will create a new resource.

* `number` - (Required, Int) Specifies the number of retention.
  + If type is set to `date_rule`, it represents the number of retention days.
  + If type is set to `tag_rule`, it represents the retention number.

* `tag_selectors` - (Optional, List) Specifies the image tags that are not counted in the retention policy
The [TagSelector](#SwrImageRetentionPolicy_TagSelector) structure is documented below.

<a name="SwrImageRetentionPolicy_TagSelector"></a>
The `TagSelector` block supports:

* `kind` - (Optional, String) Specifies the Matching rule. Value options: **label**, **regexp**.

* `pattern` - (Optional, String) Specifies the Matching pattern.
  + If kind is set to `label`, set this parameter to specific image tags.
  + If kind is set to `regexp`, set this parameter to a regular expression.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `retention_id` - The retention ID.

## Import

The swr image retention policy can be imported using the organization name, repository name
and retention ID separated by slashes or commas, e.g.:

Only when repository name is with no slashes, can use slashes to separate.

```bash
$ terraform import huaweicloud_swr_image_retention_policy.test <organization_name>/<repository_name>/<retention_id>
```

Using comma to separate is available for repository name with slashes or not.

```bash
$ terraform import huaweicloud_swr_image_retention_policy.test <organization_name>,<repository_name>,<retention_id>
```
