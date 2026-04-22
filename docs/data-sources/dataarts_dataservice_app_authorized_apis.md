---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_app_authorized_apis"
description: |-
  Use this data source to get the list of authorized APIs under specified APP within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_app_authorized_apis

Use this data source to get the list of authorized APIs under specified APP within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "app_id" {}

data "huaweicloud_dataarts_dataservice_app_authorized_apis" "test" {
  workspace_id = var.workspace_id
  app_id       = var.app_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the authorized APIs are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the APP belongs.

* `app_id` - (Required, String) Specifies the ID of the APP used to query authorized APIs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apis` - The list of APIs authorized to the APP.  
  The [apis](#dataarts_dataservice_app_authorized_apis_elem) structure is documented below.

<a name="dataarts_dataservice_app_authorized_apis_elem"></a>
The `apis` block supports:

* `id` - The ID of the API.

* `name` - The name of the API.

* `description` - The description of the API.

* `approval_time` - The approval time, in RFC3339 format.

* `manager` - The name of the API reviewer.

* `deadline` - The deadline for using the API, in RFC3339 format.

* `relationship_type` - The relationship between the authorized API and the APP.
  + **LINK_WAITING_CHECK**: Pending to approval for authorize operation.
  + **LINKED**: Already authorized.
  + **OFFLINE_WAITING_CHECK**: Pending to approval for offline operation.
  + **RENEW_WAITING_CHECK**: Pending to approval for renew operation.

* `static_params` - The configuration of the static parameters.  
  The [static_params](#dataarts_dataservice_app_authorized_api_static_param) structure is documented below.

<a name="dataarts_dataservice_app_authorized_api_static_param"></a>
The `static_params` block supports:

* `name` - The name of the static parameter.

* `value` - The value of the static parameter.
