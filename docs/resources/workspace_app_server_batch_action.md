---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_batch_action"
description: |-
  Use this resource to batch operate APP servers within HuaweiCloud.
---

# huaweicloud_workspace_app_server_batch_action

Use this resource to batch operate APP servers within HuaweiCloud.

-> This resource is only a one-time action resource for batch operate APP servers. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Modify Server Image

```hcl
variable "operate_server_ids" {
  type = list(string)
}
variable "new_image_id" {}

resource "huaweicloud_workspace_app_server_batch_action" "test" {
  type    = "batch-change-image"
  content = jsonencode({
    server_ids          = var.operate_server_ids
    image_id            = var.new_image_id
    image_type          = "private"
    os_type             = "Windows"
    update_access_agent = true
  })
}
```

### Reinstall Server

```hcl
variable "operate_server_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_server_batch_action" "test" {
  type    = "batch-reinstall"
  content = jsonencode({
    server_ids          = var.operate_server_ids
    update_access_agent = false
  })
}
```

### Rejoin Domain

```hcl
variable "operate_server_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_server_batch_action" "rejoin_domain" {
  type    = "batch-rejoin-domain"
  content = jsonencode({
    items = var.operate_server_ids
  })
}
```

### Update Virtual Session IP Configuration

```hcl
variable "operate_server_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_server_batch_action" "update_tsvi" {
  type    = "batch-update-tsvi"
  content = jsonencode({
    items = [
      for o in var.operate_server_ids : {
        id     = o
        enable = true
      }
    ]
  })
}
```

### Mark Server Maintenance Status

```hcl
variable "operate_server_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_server_batch_action" "maintain_status" {
  type    = "batch-maint"
  content = jsonencode({
    items           = var.operate_server_ids
    maintain_status = true
  })
}
```

### Reboot Server

```hcl
variable "operate_server_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_server_batch_action" "reboot" {
  type    = "batch-reboot"
  content = jsonencode({
    items = var.operate_server_ids
    type  = "SOFT"
  })
}
```

### Stop Server

```hcl
variable "operate_server_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_server_batch_action" "stop" {
  type    = "batch-stop"
  content = jsonencode({
    items = var.operate_server_ids
    type  = "SOFT"
  })
}
```

### Start Server

```hcl
variable "operate_server_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_server_batch_action" "start" {
  type    = "batch-start"
  content = jsonencode({
    items = var.operate_server_ids
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the APP servers to be batch operated are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `type` - (Required, String, NonUpdatable) Specifies the batch operation (action) type for the APP servers.  
  The valid values are as follows:
  + **batch-change-image**: Modify server image.
  + **batch-reinstall**: Reinstall server.
  + **batch-rejoin-domain**: Rejoin AD domain.
  + **batch-update-tsvi**: Update virtual session IP configuration.
  + **batch-maint**: Mark server maintenance status.
  + **batch-reboot**: Reboot server.
  + **batch-start**: Start server.
  + **batch-stop**: Stop server.

* `content` - (Required, String, NonUpdatable) Specifies the JSON string content for the batch operation (action)
  request.

* `max_retries` - (Optional, Int) Specifies the maximum number of retries for the batch operation (action) when
  encountering 409 conflict errors.  
  The default value is **0**, which means no retry will be performed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
