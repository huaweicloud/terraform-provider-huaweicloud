---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_security_compliance"
description: |-
  Use this data source to get details of the user domain's security compliance.
---

# huaweicloud_identity_security_compliance

Use this data source to get details of the user domain's security compliance.

## Example Usage

```hcl
data "huaweicloud_identity_security_compliance" "security_compliance1" {}

data "huaweicloud_identity_security_compliance" "security_compliance2" {
  option = "password_regex"
}

data "huaweicloud_identity_security_compliance" "security_compliance3" {
  option = "password_regex_description"
}
```

## Argument Reference

* `option` - (Optional, String) Specifies the query type.
  The valid values are **password_regex**, **password_regex_description**.

## Attribute Reference

* `password_regex` - Indicates the regular expression for password strength policy.

* `password_regex_description` - Indicates the description of the password strength policy.
