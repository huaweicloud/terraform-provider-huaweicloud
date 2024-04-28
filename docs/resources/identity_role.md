---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_role"
description: ""
---

# huaweicloud_identity_role

Manages a **Custom Policy** resource within HuaweiCloud IAM service.

->**Note** You *must* have admin privileges to use this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_role" "role1" {
  name        = "test"
  description = "created by terraform"
  type        = "AX"
  policy      = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ],
      "Condition": {
        "StringStartWith": {
          "g:ProjectName": [
            "cn-north-4"
          ]
        }
      }
    }
  ]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the custom policy.

* `description` - (Required, String) Specifies the description of the custom policy.

* `type` - (Required, String) Specifies the display mode of the custom policy. Valid options are as follows:
  + **AX**: the global service project.
  + **XA**: region-specific projects.

* `policy` - (Required, String) Specifies the content of the custom policy in JSON format. For more details,
  please refer to the [official document](https://support.huaweicloud.com/intl/en-us/usermanual-iam/iam_01_0017.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The custom policy ID.

* `references` - The number of references.

## Import

IAM custom policies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identity_role.role1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
