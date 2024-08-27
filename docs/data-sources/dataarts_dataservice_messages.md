---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_messages"
description: |-
  Use this data source to get the list of approval messages within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_messages

Use this data source to get the list of approval messages within HuaweiCloud.

-> Only exclusive messages are supported.

## Example Usage

### Query all exclusive approval messages

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_dataservice_messages" "test" {
  workspace_id = var.workspace_id
}
```

### Query all exclusive approval messages for a specified API

```hcl
variable "workspace_id" {}
variable "api_name" {}

data "huaweicloud_dataarts_dataservice_messages" "test" {
  workspace_id = var.workspace_id
  api_name     = var.api_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the approval messages are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID of the exclusive API to which the approval message
  belongs.

* `api_name` - (Optional, String) Specifies the name of the API to be approved.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `messages` - All approval messages that match the filter parameters.  
  The [messages](#dataservice_approve_messages_elem) structure is documented below.

<a name="dataservice_approve_messages_elem"></a>
The `messages` block supports:

* `id` - The ID of the approval message, in UUID format.

* `api_apply_status` - The apply status for API.
  + **STATUS_TYPE_PENDING_APPROVAL**: Pending review.
  + **STATUS_TYPE_REJECTED**: Rejected.
  + **STATUS_TYPE_PENDING_CHECK**: Pending confirmation.
  + **STATUS_TYPE_PENDING_EXECUTE**: Pending execution.
  + **STATUS_TYPE_SYNCHRONOUS_EXECUTE**: Synchronous execution.
  + **STATUS_TYPE_FORCED_CANCEL**: Forced cancellation.
  + **STATUS_TYPE_PASSED**: Passed.

* `api_apply_type` - The apply type.
  + **APPLY_TYPE_PUBLISH**: Release API.
  + **APPLY_TYPE_AUTHORIZE**: API active authorization.
  + **APPLY_TYPE_APPLY**: Review API.
  + **APPLY_TYPE_RENEW**: Apply for renewal of API.
  + **APPLY_TYPE_STOP**: Apply for suspension of API.
  + **APPLY_TYPE_RECOVER**: Apply for recovery of API.
  + **APPLY_TYPE_API_CANCEL_AUTHORIZE**: API cancellation of authorization.
  + **APPLY_TYPE_APP_CANCEL_AUTHORIZE**: APP cancellation of authorization.
  + **APPLY_TYPE_OFFLINE**: Apply for offline.

* `api_id` - The ID of the exclusive API to which the approval message belongs.

* `api_name` - The name of the exclusive API to which the approval message belongs.

* `api_using_time` - The expiration time used by the API, in RFC3339 format.

* `app_id` - The application ID of the API that has been bound (or is to be bound).

* `app_name` - The application name of the API that has been bound (or is to be bound).

* `apply_time` - The apply time, in RFC3339 format.

* `approval_time` - The approval time of the approval message, in RFC3339 format.

* `approver_name` - The approver name.

* `comment` - The approval comment.

* `user_name` - The name of applicant.
