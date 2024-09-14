---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_message_approve"
description: |-
  Use this resource to approve the API message for DataArts Data Service within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_message_approve

Use this resource to approve the API message for DataArts Data Service within HuaweiCloud.

-> 1. Only messages of the exclusive (and published) API can be approved.
   <br>2. This resource is a one-time action resource used only for approval message. Deleting this resource will not
   clear the corresponding request record, but will only remove the resource information from the tfstate file.
   <br>3. Before using this resource, please make sure that the current user has the approver permission.

## Example Usage

### Approve the message immediately

```hcl
variable "workspace_id" {}
variable "message_id" {}

resource "huaweicloud_dataarts_dataservice_message_approve" "test" {
  workspace_id = var.workspace_id
  message_id   = var.message_id
  action       = 0
}
```

### Approve the message on the hour after 1 day

```hcl
variable "workspace_id" {}
variable "message_id" {}

resource "huaweicloud_dataarts_dataservice_message_approve" "test" {
  workspace_id = var.workspace_id
  message_id   = var.message_id
  action       = 1
  time         = formatdate("YYYY-MM-DD'T'hh:mm:ss.000Z", timeadd(timestamp(), "24h"))
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the message (to be approved) is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID of the exclusive API to which the message
  (to be approved) belongs.
  Changing this parameter will create a new resource.

* `message_id` - (Required, String, ForceNew) Specifies the ID of the message (to be approved).  
  Changing this parameter will create a new resource.

* `action` - (Optional, Int, ForceNew) Specifies the approval action performed by the message.  
  The valid values are as follows:
  + **0**: Immediate approval.
  + **1**: Regular approval.

  Defaults to `0`. Changing this parameter will create a new resource.

* `time` - (Optional, String, ForceNew) Specifies the regular approval time.  
  The format is `YYYY-MM-DDThh:mm:ss.000Z`.  
  Required if the value of the parameter `action` is `1`.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.
