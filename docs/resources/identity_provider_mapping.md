---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_provider_mapping"
description: |-
  Manage the mapping rules of identity provider within HuaweiCloud IAM service.
---

# huaweicloud_identity_provider_mapping

Manage the mapping rules of identity provider within HuaweiCloud IAM service.

## Example Usage

```hcl
variable "provider_id" {}

resource "huaweicloud_identity_provider_mapping" "mapping" {
  provider_id = var.provider_id

  mapping_rules = <<RULES
    [
      {
        "local": [
          {
            "user": {
              "name": "{0}"
            }
          },
          {
            "group": {
              "name": "admin"
            }
          }
        ],
        "remote": [
          {
            "type": "UserName"
          },
          {
            "type": "Groups",
            "any_one_of": [
              ".*@mail.com$"
            ],
            "regex": true
          }
        ]
      }
    ]
  RULES
}
```

## Argument Reference

The following arguments are supported:

* `provider_id` - (Required, String, NonUpdatable) Specifies the ID of the identity provider used to manage the mapping rules.

* `mapping_rules` - (Required, String) Specifies the identity mapping rules in json string format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of mapping rules.

## Import

Identity provider mapping rules are imported using the `provider_id`, e.g.

```bash
$ terraform import huaweicloud_identity_provider_mapping.mapping <provider_id>
```
