---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_image"
description: |-
  Manages a private image resource of Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_image

Manages a private image resource of Workspace APP within HuaweiCloud.

## Example Usage

```hcl
variable "image_server_id" {}
variable "generated_image_name" {}
variable "generated_image_description" {}

resource "huaweicloud_workspace_app_image" "test" {
  server_id   = var.image_server_id
  name        = var.generated_image_name
  description = var.generated_image_description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `server_id` - (Required, String, ForceNew) Specifies the image server ID for generating a private image.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the image.
  Changing this creates a new resource.  
  The name valid length is limited from `1` to `128` characters. Only Chinese and English characters, digits, spaces and
  special characters (-._) and are allowed, and the first and last characters cannot be spaces.

* `description` - (Optional, String, ForceNew) Specifies the description of the image.
  Changing this creates a new resource.  
  The description contain a maximum of `1,024` characters. The carriage return characters and angle brackets (< and >) are
  not allowed, and the first character cannot be space.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the image belongs.
  Changing this creates a new resource.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
