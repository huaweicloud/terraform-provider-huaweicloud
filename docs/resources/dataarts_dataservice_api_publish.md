---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_api_publish"
description: |-
  Use this resource to publish an debugged API for DataArts Data Service within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_api_publish

Use this resource to publish an debugged API for DataArts Data Service within HuaweiCloud.

-> 1. The API must be debugged before it is published.
   <br>2. This resource is only a one-time action resource for publishing the API. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.
   <br>3. Repeated publishing is not supported when the API has been published or in the pending approval status.
   <br>4. Before using this resource, please make sure that the current user has the approver permission.
   <br>5. Please make sure that the network of the cluster and the target APIG instance (if used) are consistent.

## Example Usage

### Publish the API on Data Service side

```hcl
variable "workspace_id" {}
variable "publish_api_id" {}
variable "exclusive_cluster_id" {}

resource "huaweicloud_dataarts_dataservice_api_publish" "test" {
  workspace_id = var.workspace_id
  api_id       = var.publish_api_id
  instance_id  = var.exclusive_cluster_id
}
```

### Publish the API both on Data Service and APIG sides

```hcl
variable "workspace_id" {}
variable "publish_api_id" {}
variable "exclusive_cluster_id" {}
variable "apig_instance_id" {}
variable "apig_group_id" {}

resource "huaweicloud_dataarts_dataservice_api_publish" "test" {
  workspace_id     = var.workspace_id
  api_id           = var.publish_api_id
  instance_id      = var.exclusive_cluster_id
  apig_type        = "APIG"
  apig_instance_id = var.apig_instance_id
  apig_group_id    = var.apig_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the published API is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID to which the published API belongs.  
  Changing this parameter will create a new resource.

* `api_id` - (Required, String, ForceNew) Specifies the ID of the API to be published.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the exclusive cluster ID to which the published API belongs on
  Data Service side.  
  Changing this parameter will create a new resource.

* `apig_type` - (Optional, String, ForceNew) Specifies the type of the APIG object.  
  The valid values are as follows:
  + **APIG**
  + **APIGW**
  + **ROMA_APIC**

  Changing this parameter will create a new resource.

* `apig_instance_id` - (Optional, String, ForceNew) Specifies the APIG instance ID to which the API is published
  simultaneously in APIG service.  
  Changing this parameter will create a new resource.

* `apig_group_id` - (Optional, String, ForceNew) Specifies the APIG group ID to which the published API belongs.  
  Changing this parameter will create a new resource.

* `roma_app_id` - (Optional, String, ForceNew) Specifies the application ID for ROMA APIC.  
  Changing this parameter will create a new resource.

-> If `apig_type` and other optional parameters are omitted, the API will only be published on Data Service side.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
