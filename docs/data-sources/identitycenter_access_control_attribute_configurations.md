---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_access_control_attribute_configurations"
description: |-
  Use this data source to get the Identity Center access control attribute configurations.
---

# huaweicloud_identitycenter_access_control_attribute_configurations

Use this data source to get the Identity Center access control attribute configurations.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_identitycenter_access_control_attribute_configurations" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the IAM Identity Center instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `access_control_attributes` - The attributes configured for ABAC in the IAM Identity Center instance.

  The [access_control_attributes](#instance_access_control_attribute_configuration_struct) structure is documented below.

<a name="instance_access_control_attribute_configuration_struct"></a>
The `access_control_attributes` block supports:

* `value` - The value mapped to identity source from the specified attribute.

* `key` - The name of the attribute associated with the identity in the identity source.
