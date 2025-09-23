---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_api_action"
description: |-
  Use this resource to operate published API for DataArts Data Service within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_api_action

Use this resource to operate published API for DataArts Data Service within HuaweiCloud.

-> 1. Only exclusive API can be unpublished, stopped or recovered.
   <br>2. This resource is only a one-time action resource for doing API action. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.
   <br>3. Repeated unpublishing is not supported when the API has been unpublished or in the pending approval status.
   <br>4. This resource will not remove the object API (in APIG service) during resource delete.
   <br>5. Before unpublishing or stopping, please make sure that the current user has the approver permission.
   <br>6. If the current user does not have the approver permission and doing recovery request, this request needs to
   be approved.
   <br>7. The unpublish/stop function will reserve 2 days for canceling the APP's authorization permissions or
   processing existing messages.

## Example Usage

```hcl
variable "workspace_id" {}
variable "published_api_id" {}
variable "exclusive_cluster_id" {}

resource "huaweicloud_dataarts_dataservice_api_action" "test" {
  workspace_id = var.workspace_id
  api_id       = var.published_api_id
  instance_id  = var.exclusive_cluster_id
  type         = "UNPUBLISH"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the API to be operated is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of workspace where the API is located.  
  Changing this parameter will create a new resource.

* `api_id` - (Required, String, ForceNew) Specifies the exclusive API ID, which in published status.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the exclusive cluster ID to which the published API belongs and
  on Data Service side.  
  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the action type.  
  The valid values are as follows:
  + **UNPUBLISH**
  + **STOP**
  + **RECOVER**

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
