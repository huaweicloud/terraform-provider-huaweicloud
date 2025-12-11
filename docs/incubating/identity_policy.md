---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_policy"
description: |-
  Manages an identity policy resource within HuaweiCloud.
---
# huaweicloud_identity_policy

Manages an identity policy resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

```hcl
variable "policy_name" {}

resource "huaweicloud_identity_policy" "test" {
  name            = var.policy_name
  description     = "created by terraform"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Allow"
        },
      ]
      Version = "5.0"
    }
  )
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, NonUpdatable) Specifies the name of identity policy. The name is a string of `1` to `128`
  characters. Only English letters, digits, underscores (_), plus (+), equals (=), dots (.), ats (@) and
  hyphens (-) are allowed.

* `policy_document` - (Required, String) Specifies the policy document of the identity policy.
  It's a JSON string. If updated, a new version of policy will be created and set to default version.
  At most 5 versions of each policy are allowed, if there is no room for a new version, the earliest
  version will be deleted.

* `version_to_delete` - (Optional, String) Specifies the ID the policy version to be deleted, for example, **v3**.
  If specified, this version will be deleted instead of the earliest one when updating the `policy_document`.
  The value must be an existing version and can not be the default version.

* `path` - (Optional, String, NonUpdatable) Specifies the resource path. It is made of several strings, each containing
  one or more English letters, digits, underscores (_), plus (+), equals (=), comma (,), dots (.), at (@) and hyphens
  (-), and must be ended with slash (/). Such as **foo/bar/**. It's a part of the uniform resource name. Default is
  empty.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the identity policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The identity policy ID.

* `urn` - The uniform resource name of the identity policy. Format is `iam::$accountID:policy:$path$policyName` where
  `$accountID` is IAM account ID, `$path` is `path`, `$policyName` is `name`.

* `policy_type` - The policy type.

* `default_version_id` - The default version ID of the policy.

* `version_ids` - The version ID list of the policy.

* `attachment_count` - The attachment count.

* `created_at` - The time when the identity policy was created.

* `updated_at` - The time when the identity policy was updated.

## Import

Identity policies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identity_policy.test <id>
```
