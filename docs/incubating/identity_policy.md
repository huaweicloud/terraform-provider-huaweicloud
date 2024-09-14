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

* `policy_document` - (Required, String, NonUpdatable) Specifies the policy document of the identity policy.
  It's a JSON string.

* `path` - (Optional, String, NonUpdatable) Specifies the resource path. It is made of several strings, each containing one
  or more English letters, digits, underscores (_), plus (+), equals (=), comma (,), dots (.), at (@) and hyphens (-),
  and must be ended with slash (/). Such as **foo/bar/**. It's a part of the uniform resource name. Default is empty.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the identity policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The identity policy ID.

* `urn` - The uniform resource name of the identity policy. Format is `iam::$accountID:policy:$path$policyName` where
  `$accountID` is IAM account ID, `$path` is `path`, `$policyName` is `name`.

* `policy_type` - The policy type.

* `default_version_id` - The default version ID of the policy.

* `attachment_count` - The attachment count.

* `created_at` - The time when the identity policy was created.

* `updated_at` - The time when the identity policy was updated.

## Import

Identity policies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identity_policy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`policy_document`. It is generally recommended running `terraform plan` after importing an identity policy.
You can then decide if changes should be applied to the policy, or the resource definition should be updated
to align with the policy. Also you can ignore changes as below.

```hcl
resource "huaweicloud_identity_policy" "test" {
    ...

  lifecycle {
    ignore_changes = [
      policy_document
    ]
  }
}
```
