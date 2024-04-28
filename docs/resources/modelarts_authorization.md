---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_authorization"
description: ""
---

# huaweicloud_modelarts_authorization

Manages a ModelArts authorization resource within HuaweiCloud.  

## Example Usage

### Authorized to an IAM user

```hcl
variable "user_id" {}

resource "huaweicloud_modelarts_authorization" "test" {
  user_id     = var.user_id
  type        = "agency"
  agency_name = "ma_agency_userName"
}
```

### Authorized to all IAM users

```hcl
resource "huaweicloud_modelarts_authorization" "test" {
  user_id     = "all-users"
  type        = "agency"
  agency_name = "modelarts_agency"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `user_id` - (Required, String, ForceNew) User ID.  
  If user_id is set to **all-users**, all IAM users are authorized.  
  If this user has been authorized, the authorization setting will be updated.

  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Authorization type. The valid value is **agency**.  

  Changing this parameter will create a new resource.

* `agency_name` - (Required, String) Agency name.  
  If the agency does not exist, it will be created automatically,
  the agency name can be **modelarts_agency** or prefixed with **ma_agency_**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `user_id`.

* `user_name` - User Name.  
  The value is **all-users** if `user_id` is set to **all-users**.

## Import

The ModelArts authorization can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_authorization.test 0ce123456a00f2591fabc00385ff1234
```
