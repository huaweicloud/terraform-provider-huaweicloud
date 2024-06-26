---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_apis"
description: |-
  Use this data source to get the list of Data Service APIs within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_apis

Use this data source to get the list of Data Service APIs within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "api_name_to_be_queried" {}

data "huaweicloud_dataarts_dataservice_apis" "test" {
  workspace_id = var.workspace_id
  name         = var.api_name_to_be_queried
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

* `api_id` - (Optional, String) Specifies the API ID to be queried.

* `name` - (Optional, String) Specifies the API name to be fuzzy queried.  
  The valid length is limited from `3` to `64`, only Chinese and English characters, digits and underscores (_) are
  allowed.  
  The name must start with a Chinese or English character, and the Chinese characters must be in **UTF-8**
  or **Unicode** format.

* `type` - (Optional, String) Specifies the API type to be queried.  
  The valid values are as follows:
  + **API_SPECIFIC_TYPE_CONFIGURATION**
  + **API_SPECIFIC_TYPE_REGISTER**
  + **API_SPECIFIC_TYPE_ORCHESTRATE**
  + **API_SPECIFIC_TYPE_MYBATIS**
  + **API_SPECIFIC_TYPE_SCRIPT**
  + **API_SPECIFIC_TYPE_GROOVY**

* `description` - (Optional, String) Specifies the API description to be fuzzy queried.  
  Maximum of `255` characters are allowed.

* `create_user` - (Optional, String) Specifies the API creator to be queried.

* `datatable` - (Optional, String) Specifies the data table name used by API to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apis` - All APIs that match the filter parameters.  
  The [apis](#dataservice_apis_elem) structure is documented below.

<a name="dataservice_apis_elem"></a>
The `apis` block supports:

* `id` - The API ID, in UUID format.

* `name` - The name of the API.

* `type` - The API type.

* `description` - The description of the API.

* `protocol` - The request protocol of the API.

* `path` - The API request path.

* `request_type` - The request type of the API.

* `manager` - The API reviewer.

* `datasource_config` - The configuration of the API data source.  
  The [datasource_config](#dataservice_api_datasource_config_attr) structure is documented below.

* `request_params` - The parameters of the API request.  
  The [request_params](#dataservice_api_request_params_attr) structure is documented below.

* `backend_config` - The configuration of the API backend.  
  The [backend_config](#dataservice_api_backend_config_attr) structure is documented below.

* `create_user` - The creator name.

* `created_at` - The creation time of the API, in RFC3339 format.

* `updated_at` - The latest update time of the API, in RFC3339 format.

* `group_id` - The ID of the group to which the API belongs, for shared type.

* `status` - The API status, for shared type.

* `host` - The API host configuration, for shared type.

* `hosts` - The API host configuration, for exclusive type.

<a name="dataservice_api_datasource_config_attr"></a>
The `datasource_config` block supports:

* `type` - The type of the data source.

* `connection_name` - The name of the data connection for the DataArts Studio service.

* `connection_id` - The ID of the data connection for the DataArts Studio service.

* `database` - The name of the database.

* `datatable` - The name of the data table.

* `table_id` - The ID of the data table.

* `queue` - The ID of the DLI queue.

* `access_mode` - The access mode for the data.

* `sql` - The SQL statements in script access type.

* `backend_params` - The backend parameters of the API.  
  The [backend_params](#dataservice_api_datasource_config_backend_params) structure is documented below.

* `response_params` - The response parameters of the API.  
  The [response_params](#dataservice_api_datasource_config_response_params) structure is documented below.

* `order_params` - The order parameters of the API.  
  The [order_params](#dataservice_api_datasource_config_order_params) structure is documented below.

<a name="dataservice_api_datasource_config_backend_params"></a>
The `backend_params` block supports:

* `name` - The name of the backend parameter.

* `mapping` - The name of the mapping parameter.

* `condition` - The condition character.

<a name="dataservice_api_datasource_config_response_params"></a>
The `response_params` block supports:

* `name` - The name of the response parameter.

* `type` - The type of the response parameter.

* `field` - The bound table field for the response parameter.

* `description` - The description of the response parameter.

* `example_value` - The example value of the response parameter.

<a name="dataservice_api_datasource_config_order_params"></a>
The `order_params` block supports:

* `name` - The name of the order parameter.

* `field` - The corresponding parameter field for the order parameter.

* `optional` - Whether this order parameter is the optional parameter.

* `sort` - The sort type of the order parameter.

* `order` - The order of the sorting parameters.

<a name="dataservice_api_request_params_attr"></a>
The `request_params` block supports:

* `name` - The name of the request parameter.

* `position` - The position of the request parameter.

* `type` - The type of the request parameter.

* `description` - The description of the request parameter.

* `necessary` - Whether this parameter is the required parameter.

* `example_value` - The example value of the request parameter.

* `default_value` - The default value of the request parameter.

<a name="dataservice_api_backend_config_attr"></a>
The `backend_config` block supports:

* `type` - The type of the backend request.

* `protocol` - The protocol of the backend request.

* `host` - The backend host.

* `timeout` - The backend timeout.

* `path` - The backend path.

* `backend_params` - The backend parameters of the API.  
  The [backend_params](#dataservice_api_backend_config_attr_backend_params) structure is documented below.

* `constant_params` - The backend constant parameters of the API.  
  The [constant_params](#dataservice_api_backend_config_attr_constant_params) structure is documented below.

<a name="dataservice_api_backend_config_attr_backend_params"></a>
The `backend_params` block supports:

* `name` - The name of the request parameter.

* `position` - The position of the request parameter.

* `backend_param_name` - The name of the corresponding backend parameter.

<a name="dataservice_api_backend_config_attr_constant_params"></a>
The `constant_params` block supports:

* `name` - The name of the constant parameter.

* `type` - The type of the constant parameter.

* `position` - The position of the constant parameter.

* `value` - The value of the constant parameter.

* `description` - The description of the constant parameter.
