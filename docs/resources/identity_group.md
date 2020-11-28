---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud\_identity\_group

Manages a User Group resource within HuaweiCloud IAM service.
This is an alternative to `huaweicloud_identity_group_v3`

Note: You _must_ have admin privileges in your HuaweiCloud cloud to use
this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_group" "group_1" {
  name        = "group_1"
  description = "This is a test group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the group.The length is less than or equal to 64 bytes. 

* `description` - (Optional, String) A description of the group.

* `domain_id` - (Optional, String) The domain this group belongs to.

## Attributes Reference

The following attributes are exported:

## Import

Groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_identity_group_v3.group_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
