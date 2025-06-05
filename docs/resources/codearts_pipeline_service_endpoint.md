---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_service_endpoint"
description: |-
  Manages a CodeArts pipeline service endpoint resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_service_endpoint

Manages a CodeArts pipeline service endpoint resource within HuaweiCloud.

-> Destroying resource does not delete the service endpoint.

## Example Usage

```hcl
variable "project_id" {}
variable "module_id" {}
variable "url" {}
variable "name" {}
variable "authorization_scheme" {}

resource "huaweicloud_codearts_pipeline_service_endpoint" "test" {
  project_id = var.project_id
  module_id  = var.module_id
  url        = var.url
  name       = var.name

  authorization {
    scheme     = var.authorization_scheme
    parameters = jsonencode({
      "username":"test",
      "password":"test"
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.

* `authorization` - (Optional, List, NonUpdatable) Specifies the permission information.
  The [authorization](#block--authorization) structure is documented below.

* `module_id` - (Optional, String, NonUpdatable) Specifies the module ID.

* `name` - (Optional, String, NonUpdatable) Specifies the endpoint name.

* `url` - (Optional, String, NonUpdatable) Specifies the URL.

<a name="block--authorization"></a>
The `authorization` block supports:

* `parameters` - (Optional, String, NonUpdatable) Specifies the authentication parameter.

* `scheme` - (Optional, String, NonUpdatable) Specifies the authentication mode.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_by` - Indicates the permission information.
  The [created_by](#attrblock--created_by) structure is documented below.

<a name="attrblock--created_by"></a>
The `created_by` block supports:

* `user_id` - Indicates the user ID.

* `user_name` - Indicates the user name.

## Import

The service endpoint can be imported using `project_id` and `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_service_endpoint.test <project_id>/<id>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `authorization`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the service endpoint, or the resource definition should be updated to
align with the service endpoint. Also you can ignore changes as below.

```hcl
resource "huaweicloud_codearts_pipeline_service_endpoint" "test" {
  ...

  lifecycle {
    ignore_changes = [
      authorization,
    ]
  }
}
```
