---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_action"
description: |-
  Use this resource to operate APP server within HuaweiCloud.
---

# huaweicloud_workspace_app_server_action

Use this resource to operate APP server within HuaweiCloud.

-> This resource is only a one-time action resource for operate APP server. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Modify Server Image

```hcl
variable "operate_server_id" {}
variable "new_image_id" {}

resource "huaweicloud_workspace_app_server_action" "test" {
  type      = "change-image"
  server_id = "server-uuid"
  content   = jsonencode({
    image_id            = var.operate_server_id
    image_type          = var.new_image_id
    os_type             = "Windows"
    image_product_id    = "product-uuid"
    update_access_agent = true
  })
}
```

### Reinstall Server

```hcl
resource "huaweicloud_workspace_app_server_action" "test" {
  type      = "reinstall"
  server_id = "server-uuid"
  content   = jsonencode({
    update_access_agent = false
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the APP server to be operated is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `type` - (Required, String, NonUpdatable) Specifies the operation (action) type for the APP server.  
  The valid values are as follows:
  + **change-image**: Modify server image.
  + **reinstall**: Reinstall server.

* `server_id` - (Required, String, NonUpdatable) Specifies the ID of the server to be operated.

* `content` - (Required, String, NonUpdatable) Specifies the JSON string content for the operation (action)
  request.

* `max_retries` - (Optional, Int) Specifies the maximum number of retries for the operation (action) when
  encountering 409 conflict errors.  
  The default value is **0**, which means no retry will be performed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
