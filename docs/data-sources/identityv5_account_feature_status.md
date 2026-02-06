---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_account_feature_status"
description: |-
  Use this data source to get the status of the specified feature for the IAM account within HuaweiCloud.
---

# huaweicloud_identityv5_account_feature_status

Use this data source to get the status of the specified feature for the IAM account within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identityv5_account_feature_status" "test" {
  feature_name = "v5_console"
}
```

## Argument Reference

The following arguments are supported:

* `feature_name` - (Required, String) Specifies the name of the feature to be queried.  
  The valid values are as follows:
  + **v5_console**
  + **access_analyzer**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `feature_status` - The status of the feature.
  + **on**
  + **off**
