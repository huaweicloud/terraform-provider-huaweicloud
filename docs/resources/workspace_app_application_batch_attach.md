---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_application_batch_attach"
description: |-
  Use this resource to attach applications to Workspace APP image instance in HuaweiCloud.
---

# huaweicloud_workspace_app_application_batch_attach

Use this resource to attach applications to Workspace APP image instance in HuaweiCloud.

-> This resource is a one-time action resource used to attach applications to specified image instance. Deleting resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "server_id" {}
variable "applications_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_application_batch_attach" "test" {
  server_id        = var.server_id
  applications_ids = var.applications_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Workspace APP is located.  
  Changing this creates a new resource.

* `server_id` - (Required, String) Specifies the ID of the image server instance.

* `applications_ids` - (Required, List) Specifies the list of application IDs to be attach.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `uri` - The URI of the application attachment.
