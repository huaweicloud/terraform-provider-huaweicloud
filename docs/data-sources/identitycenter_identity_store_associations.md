---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_identity_store_associations"
description: |-
  Use this data source to get the Identity Center identity store associations.
---

# huaweicloud_identitycenter_identity_store_associations

Use this data source to get the Identity Center identity store associations.

## Example Usage

```hcl
variable "instance_id" {}
 
data "huaweicloud_identitycenter_identity_store_associations" "test"{
  instance_id = var.instance_id
}
```

## Argument Reference

* `instance_id` - (Required, String) Specifies the ID of an IAM Identity Center instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `identity_store_id` - Indicates the ID of the identity store that associated with IAM Identity Center.

* `identity_store_type` - Indicates the type of the identity store.

* `authentication_type` - Indicates the authentication type of the identity store.

* `provisioning_type` - Indicates the list of the identity store provisioning type.

* `status` - Indicates the status of the identity store.
