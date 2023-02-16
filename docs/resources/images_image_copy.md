---
subcategory: "Image Management Service (IMS)"
---

# huaweicloud_images_image_copy

Manages IMS image copy resources within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_images_image_copy" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `image_id` - (Required, String, ForceNew) Special the ID of the copied image.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the image.

* `description` - (Optional, String, ForceNew) Specifies the description of the image.

  Changing this parameter will create a new resource.

* `cmk_id` - (Optional, String, ForceNew) Specifies the master key used for encrypting an image.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the image.

  Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The ims image copy can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_images_image_copy.test 7886e623-f1b3-473e-b882-67ba1c35887f
```
