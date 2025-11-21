---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_account_feature_status"
description: |-
  Use this data source to get the status of account in the Identity and Access Management V5 service.
---

# huaweicloud_identityv5_account_feature_status

Use this data source to get the status of account in the Identity and Access Management V5 service.

## Example Usage

```hcl
data "huaweicloud_identityv5_account_feature_status" "test" {
  feature_name = "test_name"
}
```

## Argument Reference

The following arguments are supported:

* `feature_name` - (Required, String) Specifies the name of the IAM account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `feature_status` - Indicates the status of account.
