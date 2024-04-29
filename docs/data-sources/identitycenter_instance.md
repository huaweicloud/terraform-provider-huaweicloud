---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_instance"
description: ""
---

# huaweicloud_identitycenter_instance

Use this data source to get the Identity Center instance info.

## Example Usage

```hcl
data "huaweicloud_identitycenter_instance" "system"{}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the IAM Identity Center instance.

* `identity_store_id` - Indicates the ID of the identity store that associated with IAM Identity Center.

* `urn` - Indicates the uniform resource name of the instance.

* `alias` - Indicates the alias of the instance.
