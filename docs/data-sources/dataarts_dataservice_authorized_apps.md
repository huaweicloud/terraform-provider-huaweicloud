---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_authorized_apps"
description: |-
  Use this data source to get the list of authorized APPs under specified API within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_authorized_apps

Use this data source to get the list of authorized APPs under specified API within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "authorized_api_id" {}

data "huaweicloud_dataarts_dataservice_authorized_apps" "test" {
  workspace_id = var.workspace_id
  dlm_type     = "EXCLUSIVE"
  api_id       = var.authorized_api_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of workspace where the APIs are located.

* `dlm_type` - (Optional, String) Specifies the type of DLM engine.  
  The valid values are as follows:
  + **SHARED**: Shared data service.
  + **EXCLUSIVE**: The exclusive data service.

  Defaults to **SHARED**.

* `api_id` - (Required, String) Specifies the ID of the API used to authorize the APPs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apps` - All APPs authorized by API.  
  The [apis](#dataservice_authorized_apps_elem) structure is documented below.

<a name="dataservice_authorized_apps_elem"></a>
The `apps` block supports:

* `id` - The ID of the application that has authorization.

* `name` - The name of the application that has authorization.

* `instance_id` - The instance ID to which the authorized API belongs.

* `instance_name` - The instance name to which the authorized API belongs.

* `expired_at` - The expiration time, in RFC3339 format.

* `approved_at` - The approve time, in RFC3339 format.

* `relationship_type` - The relationship between the authorized API and the authorized APP list.
  + **LINK_WAITING_CHECK**: Pending to approval for authorize operation.
  + **LINKED**: Already authorized.
  + **OFFLINE_WAITING_CHECK**: Pending to approval for offline operation.
  + **RENEW_WAITING_CHECK**: Pending to approval for renew operation.

* `static_params` - The configuration of the static parameters.  
  The [static_params](#dataservice_authorized_app_static_param) structure is documented below.

<a name="dataservice_authorized_app_static_param"></a>
The `static_params` block supports:

* `name` - The name of the static parameter.

* `value` - The value of the static parameter.
