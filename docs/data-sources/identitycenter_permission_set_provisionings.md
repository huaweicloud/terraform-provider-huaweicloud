---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_permission_set_provisionings"
description: |-
  Use this data source to get the Identity Center permission set provisionings.
---

# huaweicloud_identitycenter_permission_set_provisionings

Use this data source to get the Identity Center permission set provisionings.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_identitycenter_permission_set_provisionings" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of an IAM Identity Center instance.

* `status` - (Optional, String) Specifies the status of the permission set provisioning process.
  The valid values are as follows:
  + **IN_PROGRESS**
  + **SUCCEEDED**
  + **FAILED**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `provisionings` - The authorization status of a permission set.

  The [provisionings](#provisionings_struct) structure is documented below.

<a name="provisionings_struct"></a>
The `provisionings` block supports:

* `created_at` - The date when a permission set was created.

* `request_id` - The unique ID of a request.

* `status` - The authorization status of a permission set.
