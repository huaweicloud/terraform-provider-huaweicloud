---
subcategory: "Image Management Service (IMS)"
---

# huaweicloud_images_image_copy

Manages a IMS image copy resources within HuaweiCloud.

## Example Usage

```hcl
variable "source_image_id" {}
variable "name" {}
variable "target_region" {}
variable "agency_name" {}

resource "huaweicloud_images_image_copy" "test" {
  source_image_id = var.source_image_id
  name            = var.name
  target_region   = var.target_region
  agency_name     = var.agency_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `source_image_id` - (Required, String, ForceNew) Specifies the ID of the copied image.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the copy image.

* `target_region` - (Required, String) Specifies the target region name..

* `description` - (Optional, String, ForceNew) Specifies the description of the copy image.

  Changing this parameter will create a new resource.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the master key used for encrypting an image.
  Only copying scene within a region is supported.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the image.

  Changing this parameter will create a new resource.

* `agency_name` - (Optional, String, ForceNew) Specifies the agency name. It is required in the cross-region scene. 

  Changing this parameter will create a new resource.

* `vault_id` - (Optional, String, ForceNew) Specifies the ID of the vault. It is used in the cross-region scene,
  and it is mandatory if you are replicating a full-ECS image.

  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the copy image.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The IMS image copy can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_images_image_copy.test 7886e623-f1b3-473e-b882-67ba1c35887f
```
