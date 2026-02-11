---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_provider_conversion"
description: |-
  Manages an IAM conversion resource within HuaweiCloud.
---

# huaweicloud_identity_provider_conversion

Manages an IAM conversion resource within HuaweiCloud.

## Example Usage

```hcl
variable "provider_id" {}

resource "huaweicloud_identity_provider_conversion" "test" {
  provider_id = var.provider_id

  conversion_rules {
    local {
      username = "Tom"
    }
    remote {
      attribute = "Tom"
    }
  }

  conversion_rules {
    local {
      username = "FederationUser"
    }
    remote {
      attribute = "username"
      condition = "any_one_of"
      value     = ["Tom", "Jerry"]
    }
  }
}
```

<!--markdownlint-disable MD033-->

## Argument Reference

The following arguments are supported:

* `provider_id` - (Required, String, ForceNew) Specifies the ID of the identity provider used to manage the
  conversion rules.  
  Changing this parameter will create a new resource.

* `conversion_rules` - (Required, List) Specifies the identity conversion rules of the identity provider.
  You can use identity conversion rules to map the identities of existing users to Huawei Cloud and manage their access
  to cloud resources.
  The [conversion_rules](#identity_provider_conversion_rules) structure is documented below.

<a name="identity_provider_conversion_rules"></a>
The `conversion_rules` block supports:

* `local` - (Required, List) Specifies the federated user information on the cloud platform.  
  The [local](#identity_provider_conversion_local) structure is documented below.

* `remote` - (Required, List) Specifies Federated user information in the IDP system.  
  The [remote](#identity_provider_conversion_remote) structure is documented below.

  -> **NOTE:** 1.If the protocol of identity provider is **SAML**, this field is an expression consisting of assertion
  attributes and operators.
  <br/>2.If the protocol of identity provider is **OIDC**, the value of this field is determined by the ID token.

<a name="identity_provider_conversion_local"></a>
The `local` block supports:

* `username` - (Optional, String) Specifies the name of a federated user on the cloud platform.

* `group` - (Optional, String) Specifies the user group to which the federated user belongs on the cloud platform.

<a name="identity_provider_conversion_remote"></a>
The `remote` block supports:

* `attribute` - (Required, String) Specifies the attribute in the IDP assertion.

* `condition` - (Optional, String) Specifies the condition of conversion rule.  
  The valid values are as follows:
  + **any_one_of**: The rule is matched only if the specified strings appear in the attribute type.
  + **not_any_of**: The rule is matched only if the specified strings do not appear in the attribute type.

  -> **NOTE:** 1. The condition result is Boolean rather than the argument that is passed as input.
  <br/>2. In a remote array, `any_one_of` and `not_any_of` are mutually exclusive and cannot be set at the same time.

* `value` - (Optional, List) Specifies the rule is matched only if the specified strings appear in the attribute type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of conversion rules.

## Import

Identity provider conversion rules are imported using the `provider_id`, e.g.

```bash
$ terraform import huaweicloud_identity_provider_conversion.test <id>
```
