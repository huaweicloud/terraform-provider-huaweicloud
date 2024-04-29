---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_access_policy"
description: ""
---

# huaweicloud_workspace_access_policy

Manages access policy resource within HuaweiCloud.

## Example Usage

### Create a private access policy

```hcl
variable "workspace_user_id" {}

resource "huaweicloud_workspace_access_policy" "test" {
  name           = "PRIVATE_ACCESS"
  blacklist_type = "INTERNET"

  blacklist {
    object_type = "USER"
    object_id   = var.workspace_user_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the access policy is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the access policy.
  + **PRIVATE_ACCESS**

  Changing this will create a new resource.

  -> Custom names are not supported, and only one resource can be created with each name at the same time.

* `blacklist_type` - (Required, String, ForceNew) Specifies the type of access policy blacklist.  
  The valid values are as follows:
  + **INTERNET**

  Changing this will create a new resource.

* `blacklist` - (Optional, List) Specifies the blacklist configuration to which the policy applies.  
  The [blacklist](#access_policy_blacklist_objects_args) structure is documented below.

<a name="access_policy_blacklist_objects_args"></a>
The `blacklist` block supports:

* `object_type` - (Required, String) Specifies the object type.
  The valid values are as follows:
  + **USER**
  + **USERGROUP**

* `object_id` - (Required, String) Specifies the object ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The access policy ID in UUID format.

* `blacklist` - The blacklist configuration to which the policy applies.  
  The [blacklist](#access_policy_blacklist_objects_attr) structure is documented below.

* `created_at` - The creation time of the access policy.

<a name="access_policy_blacklist_objects_attr"></a>
The `blacklist` block supports:

* `object_name` - The object name.

## Import

Access policies can be imported using their `name`, e.g.

```bash
$ terraform import huaweicloud_workspace_access_policy.test <name>
```
