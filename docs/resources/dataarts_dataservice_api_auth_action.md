---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_api_auth_action"
description: |-
  Use this resource to authorize an API for DataArts Data Service within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_api_auth_action

Use this resource to operate (authorized) APP for DataArts Data Service within HuaweiCloud.

-> 1. Only exclusive (and published) API can used to authorize APPs.
   <br>2. This resource is only a one-time action resource for doing authorize operation. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.
   <br>3. Before using this resource, please make sure that the current user has the approver permission.

## Example Usage

```hcl
variable "workspace_id" {}
variable "operate_api_id" {}
variable "exclusive_cluster_id" {}
variable "operate_app_id" {}

resource "huaweicloud_dataarts_dataservice_api_auth_action" "test" {
  workspace_id = var.workspace_id
  api_id       = var.operate_api_id
  instance_id  = var.exclusive_cluster_id
  app_id       = var.operate_app_id
  type         = "APPLY_TYPE_APP_CANCEL_AUTHORIZE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the published API and APP to be operated are
  located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID to which the published API and APP to be
  operated belong.
  Changing this parameter will create a new resource.

* `api_id` - (Required, String, ForceNew) Specifies the ID of the published API that to be operated.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the exclusive cluster ID to which the published API belongs.  
  Changing this parameter will create a new resource.

* `app_id` - (Required, String, ForceNew) Specifies the ID of the APP to be operated.  
  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the operation type of the authorization.  
  The valid values are as follows:
  + **APPLY_TYPE_APPLY**: Apply for APP access to API.
  + **APPLY_TYPE_AUTHORIZE**: Apply for APP access to API (fast approval).
  + **APPLY_TYPE_API_CANCEL_AUTHORIZE**: Cancel the API authorization.
  + **APPLY_TYPE_APP_CANCEL_AUTHORIZE**: Cancel the APP authorization (fast approval).
  + **APPLY_TYPE_RENEW**: Renew API.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
