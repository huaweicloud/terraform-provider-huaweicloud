---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_access_control_attribute_configuration"
description: |-
  Manages an Identity Center access control attribute configuration resource within HuaweiCloud.
---

# huaweicloud_identitycenter_access_control_attribute_configuration

Manages an Identity Center access control attribute configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_identitycenter_access_control_attribute_configuration" "test" {
  instance_id = var.instance_id

  access_control_attributes {
    key   = "test"
    value = ["$${user:email}"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center instance.
  Changing this creates a new resource.

* `access_control_attributes` - (Optional, List) Specifies the properties of ABAC configuration in IAM Identity Center instance.
  The [access_control_attributes](#access_control_attributes) structure is documented below.

<a name="access_control_attributes"></a>
The `access_control_attributes` block supports:

* `key` - (Required, String) Specifies the name of the attribute associated with the identity in your identity source.

* `value` - (Required, List) Specifies the value used to map the specified attribute to the identity source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
