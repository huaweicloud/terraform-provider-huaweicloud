---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group_v3"
sidebar_current: "docs-huaweicloud-resource-identity-group-v3"
description: |-
  Manages a User Group resource within HuaweiCloud IAM service.
---

# huaweicloud\_identity\_group_v3

Manages a User Group resource within HuaweiCloud IAM service.

Note: You _must_ have admin privileges in your HuaweiCloud cloud to use
this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_group_v3" "group_1" {
  name = "group_1"
  description = "This is a test group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the group.The length is less than or equal 
     to 64 bytes 

* `description` - (Optional) A description of the group.

* `domain_id` - (Optional) The domain this group belongs to.

* `region` - (Optional) The region in which to obtain the V3 Keystone client.
    If omitted, the `region` argument of the provider is used. Changing this
    creates a new User Group.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.

## Import

Groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_identity_group_v3.group_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
