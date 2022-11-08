---
subcategory: "IoT Device Access (IoTDA)"
---

# huaweicloud_iotda_space

Manages an IoTDA resource space within HuaweiCloud.

A resource space is the basic unit of service management and provides independent device management and platform
configuration capabilities at the service layer. Resources (such as products and devices) must be created on
a resource space.

## Example Usage

```hcl
resource "huaweicloud_iotda_space" "space" {
  name = "first_space"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA resource space resource.
If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the space name. The name contains a maximum of 64 characters.
Only letters, digits, hyphens (-), underscore (_) and the following special characters are allowed: `?'#().,&%@!`.
Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `is_default` - Whether it is the default resource space. The IoT platform automatically creates and assigns
a default resource space (undeletable) to your account.

## Import

Spaces can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iotda_space.test 10022532f4f94f26b01daa1e424853e1
```
