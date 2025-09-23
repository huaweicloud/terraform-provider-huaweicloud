---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_permission_set"
description: ""
---

# huaweicloud_identitycenter_permission_set

Manages an Identity Center permission set resource within HuaweiCloud.  

## Example Usage

```hcl
data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_permission_set" "demo" {
  instance_id      = data.huaweicloud_identitycenter_instance.system.id
  name             = "demo"
  session_duration = "PT8H"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center instance.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the permission set.
  Changing this parameter will create a new resource.

* `session_duration` - (Required, String) Specifies the length of time that the user sessions are valid in the
  ISO-8601 standard, e.g. **PT4H**.

* `relay_state` - (Optional, String) Specifies the relay state URL used to redirect users within the application during
  the federation authentication process.

* `description` - (Optional, String) Specifies the description of the permission set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - The Uniform Resource Name of the permission set.

* `created_at` - The date the permission set was created in RFC3339 format.

* `account_ids` - The array of one or more account IDs bound to the permission set.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the permission set.

## Import

The Identity Center permission set can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_permission_set.demo <instance_id>/<id>
```
