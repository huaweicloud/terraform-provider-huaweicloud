---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_api_auth"
description: |-
  Use this resource to authorize APP(s) to access API for DataArts Data Service within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_api_auth

Use this resource to authorize APP(s) to access API for DataArts Data Service within HuaweiCloud.

-> 1. Only exclusive API can used to authorize APPs.
   <br>2. This resource is only a one-time action resource for doing API authorization. Deleting this resource will not
   clear the corresponding request record, but will only remove the resource information from the tfstate file.
   <br>3. APP authentication APIs can only be authorized to APP-type applications.
   <br>4. Before using this resource, please make sure that the current user has the approver permission.

## Example Usage

```hcl
variable "workspace_id" {}
variable "auth_api_id" {} # The auth type is 'APP'
variable "exclusive_cluster_id" {}
variable "auth_expiration_time" {}
variable "authorized_app_ids" {
  type = list(string)
}

resource "huaweicloud_dataarts_dataservice_api_auth" "test" {
  workspace_id = var.workspace_id
  api_id       = var.auth_api_id
  instance_id  = var.exclusive_cluster_id
  expired_at   = var.auth_expiration_time
  app_ids      = var.authorized_app_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the API and APPs are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID to which the API and APPs belong.  
  Changing this parameter will create a new resource.

* `api_id` - (Required, String, ForceNew) Specifies the API ID that used to authorize the authentication.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the exclusive cluster ID.  
  Changing this parameter will create a new resource.

* `app_ids` - (Required, List, ForceNew) Specifies the list of authorized application IDs.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
