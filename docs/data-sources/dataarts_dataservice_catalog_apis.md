---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_catalog_apis"
description: |-
  Use this data source to get the list of Data Service APIs under a specified catalog within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_catalog_apis

Use this data source to get the list of Data Service APIs under a specified catalog within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "catalog_id" {}

data "huaweicloud_dataarts_dataservice_catalog_apis" "test" {
  workspace_id = var.workspace_id
  dlm_type     = "EXCLUSIVE"
  catalog_id   = var.catalog_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of workspace where the APIs are located.

* `dlm_type` - (Optional, String) Specifies the type of DLM engine.  
  The valid values are as follows:
  + **SHARED**: The shared data service.
  + **EXCLUSIVE**: The exclusive data service.

  Defaults to **SHARED**.

* `catalog_id` - (Required, String) Specifies the ID of the catalog to which the APIs belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apis` - All API summaries that under the specified catalog.  
  The [apis](#dataservice_api_summaries_under_catalog) structure is documented below.

<a name="dataservice_api_summaries_under_catalog"></a>
The `apis` block supports:

* `id` - The API ID, in UUID format.

* `name` - The API name.

* `description` - The description of the API.

* `group_id` - The ID of the group to which the shared API belongs.

* `status` - The status of the shared API.
  + **API_STATUS_CREATED**: API has been created.
  + **API_STATUS_PUBLISH_WAIT_REVIEW**: Release review status.
  + **API_STATUS_PUBLISH_REJECT**: Rejection status.
  + **API_STATUS_PUBLISHED**: Released status.
  + **API_STATUS_WAITING_STOP**: Disable review status.
  + **API_STATUS_STOPPED**: Disabled status.
  + **API_STATUS_RECOVER_WAIT_REVIEW**: Recover review status.
  + **API_STATUS_WAITING_OFFLINE**: Offline review status.
  + **API_STATUS_OFFLINE**: Already offline.
  + **API_STATUS_OFFLINE_WAIT_REVIEW**: Offline pending review status.

* `debug_status` - The debug status of the shared API.
  + **API_DEBUG_WAITING**: Waiting for debugging.
  + **API_DEBUG_FAILED**: Debugging failed.
  + **API_DEBUG_SUCCESS**: Debugging successful.

* `publish_messages` - All publish messages of the exclusive API.
  The [publish_messages](#datasource_api_publish_messages_under_catalog) structure is documented below.

* `type` - The API type.
  + **API_SPECIFIC_TYPE_CONFIGURATION**: Configuration API
  + **API_SPECIFIC_TYPE_SCRIPT**: Script API
  + **API_SPECIFIC_TYPE_REGISTER**: Registration API

* `manager` - The API reviewer.

* `created_at` - The creation time of the API, in RFC3339 format.

* `authorization_status` - The authorization status of the API.  
  + **NO_AUTHORIZATION_REQUIRED**: No authorization required.
  + **UNAUTHORIZED**: Unauthorized.
  + **AUTHORIZED**: Authorized.

<a name="datasource_api_publish_messages_under_catalog"></a>
The `publish_messages` block supports:

* `id` - The publish ID, in UUID format.

* `instance_id` - The ID of the instance used to publish the exclusive API.

* `instance_name` - The name of the instance used to publish the exclusive API.

* `api_status` - The publish status of the exclusive API.  
  + **API_STATUS_CREATED**: API has been created.
  + **API_STATUS_PUBLISH_WAIT_REVIEW**: Release review status.
  + **API_STATUS_PUBLISH_REJECT**: Rejection status.
  + **API_STATUS_PUBLISHED**: Released status.
  + **API_STATUS_WAITING_STOP**: Disable review status.
  + **API_STATUS_STOPPED**: Disabled status.
  + **API_STATUS_RECOVER_WAIT_REVIEW**: Recover review status.
  + **API_STATUS_WAITING_OFFLINE**: Offline review status.
  + **API_STATUS_OFFLINE**: Already offline.
  + **API_STATUS_OFFLINE_WAIT_REVIEW**: Offline pending review status.

* `api_debug` - The debug status of the exclusive API.
  + **API_DEBUG_WAITING**: Waiting for debugging.
  + **API_DEBUG_FAILED**: Debugging failed.
  + **API_DEBUG_SUCCESS**: Debugging successful.
