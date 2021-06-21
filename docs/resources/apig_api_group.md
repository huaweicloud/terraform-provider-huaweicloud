---
subcategory: "API Gateway (APIG)"
---

# huaweicloud_apig_group

Manages an APIG (API) group resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}
variable "description" {}

resource "huaweicloud_apig_group" "test" {
  instance_id = var.instance_id
  name        = var.group_name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the API group resource.
  If omitted, the provider-level region will be used.
  Changing this will create a new API group resource.

* `instance_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the
  API group belongs to.
  Changing this will create a new API group resource.

* `name` - (Required, String) Specifies the name of the API group.
  The API group name consists of 3 to 64 characters, starting with a letter.
  Only letters, digits and underscores (_) are allowed.
  Chinese characters must be in UTF-8 or Unicode format.

* `description` - (Optional, String) Specifies the description about the API group.
  The description contain a maximum of 255 characters and the angle brackets (< and >) are not allowed.
  Chinese characters must be in UTF-8 or Unicode format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the API group.
* `registraion_time` - Registration time, in RFC-3339 format.
* `update_time` - Time when the API group was last modified, in RFC-3339 format.

## Import

API groups of the APIG can be imported using their `id` and the ID of the APIG instance to which the group belongs,
separated by a slash, e.g.
```
$ terraform import huaweicloud_apig_group.test <instance id>/<id>
```
