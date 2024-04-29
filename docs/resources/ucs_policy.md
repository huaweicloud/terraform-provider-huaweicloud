---
subcategory: "Ubiquitous Cloud Native Service (UCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ucs_policy"
description: ""
---

# huaweicloud_ucs_policy

Manages a UCS policy resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "user_id_1" {}
variable "user_id_2" {}

resource "huaweicloud_ucs_policy" "test" {
  name         = "policy-1"
  iam_user_ids = [var.user_id_1, var.user_id_2]
  type         = "admin"
  description  = "created by terraform"
}
```

### Custom Type Policy

```hcl
variable "user_id_1" {}
variable "user_id_2" {}

resource "huaweicloud_ucs_policy" "test" {
  name         = "policy-1"
  iam_user_ids = [var.user_id_1, var.user_id_2]
  type         = "custom"
  description  = "created by terraform"

  details {
    operations = ["*"]
    resources  = ["*"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the UCS policy. The name consists of 1 to 63 characters,
  including lowercase letters, digits and hyphens (-), must start with a letter and end with a letter or digit.

  Changing this parameter will create a new resource.

* `iam_user_ids` - (Required, List) Specifies the list of IAM user IDs to associate to the policy.

* `type` - (Required, String) Specifies the type of the UCS policy.
  The value can be: **readonly**, **develop**, **admin** and **custom**.

* `description` - (Optional, String) Specifies the description of the UCS policy.
  The description consists of 0 to 255 characters.

* `details` - (Optional, List) Specifies the details of the UCS policy.
  This only works when the type is **custom**.
  The [Details](#Policy_Details) structure is documented below.

<a name="Policy_Details"></a>
The `details` block supports:

* `operations` - (Optional, List) Specifies the list of operations.

* `resources` - (Optional, List) Specifies the list of resources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The created time.

* `updated_at` - The updated time.

## Import

The UCS policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ucs_policy.test 8b12072c-0c25-11ee-b6b2-0255ac1000de
```
