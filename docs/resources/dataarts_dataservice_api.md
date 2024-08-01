---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_dataservice_api"
description: |-
  Use this resource to manage API under data service catalog within HuaweiCloud
---

# huaweicloud_dataarts_dataservice_api

Use this resource to manage API under data service catalog within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "catalog_id" {}
variable "random_resource_name" {}
variable "reviewer_name" {}

resource "huaweicloud_dataarts_dataservice_api" "test" {
  workspace_id = var.workspace_id
  dlm_type     = "SHARED"
  type         = "API_SPECIFIC_TYPE_CONFIGURATION"
  catalog_id   = var.catalog_id
  name         = var.random_resource_name
  description  = "Created by terraform script"
  auth_type    = "NONE"
  manager      = var.reviewer_name
  path         = "/terraform/auto/resource_create/{resource_type}/{resource_name}"
  protocol     = "PROTOCOL_TYPE_HTTP"
  request_type = "REQUEST_TYPE_POST"
  visibility   = "WORKSPACE"

  request_params {
    name          = "resource_type"
    position      = "REQUEST_PARAMETER_POSITION_PATH"
    type          = "REQUEST_PARAMETER_TYPE_STRING"
    description   = "The type of the terraform resource to be automatically created"
    necessary     = true
    example_value = "huaweicloud_vpc"
  }
  request_params {
    name          = "resource_name"
    position      = "REQUEST_PARAMETER_POSITION_PATH"
    type          = "REQUEST_PARAMETER_TYPE_STRING"
    description   = "The name of the terraform resource to be automatically created"
    necessary     = true
    example_value = "test"
  }
  request_params {
    name        = "configuration"
    position    = "REQUEST_PARAMETER_POSITION_BODY"
    type        = "REQUEST_PARAMETER_TYPE_STRING"
    description = "The configuration of the terraform resource, in JSON format"
    necessary   = true
  }
  request_params {
    name        = "resource_id"
    position    = "REQUEST_PARAMETER_POSITION_BODY"
    type        = "REQUEST_PARAMETER_TYPE_STRING"
    description = "The resource ID, in UUID format"
    necessary   = false
  }
  request_params {
    name          = "order"
    position      = "REQUEST_PARAMETER_POSITION_BODY"
    type          = "REQUEST_PARAMETER_TYPE_STRING"
    description   = "The filter parameter for resource configuration details"
    necessary     = false
    example_value = "asc"
    default_value = "desc"
  }

  datasource_config {
    type          = "DLI"
    connection_id = huaweicloud_dataarts_studio_data_connection.test.id
    database      = huaweicloud_dli_database.test.name
    datatable     = huaweicloud_dli_table.test.name
    queue         = huaweicloud_dli_queue.test.name
    access_mode   = "SQL"

    backend_params {
      name      = "configuration"
      mapping   = "configuration"
      condition = "CONDITION_TYPE_EQ"
    }

    response_params {
      name        = "resourceId"
      type        = "REQUEST_PARAMETER_TYPE_STRING"
      field       = "resource_id"
      description = "The resource ID, in UUID format"
    }

    order_params {
      name     = "bePlans"
      field    = "plans"
      optional = true
      sort     = "ASC"
      order    = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of workspace where the API is located.

  Changing this parameter will create a new resource.

* `dlm_type` - (Optional, String, ForceNew) Specifies the type of DLM engine.  
  The valid values are as follows:
  + **SHARED**: Shared data service.
  + **EXCLUSIVE**: The exclusive data service.

  Defaults to **SHARED**.  
  Changing this parameter will create a new resource.

* `catalog_id` - (Required, String) Specifies the ID of the catalog where the API is located.

* `name` - (Required, String) Specifies the name of the API.  
  The valid length is limited from `3` to `64`, only Chinese and English characters, digits and underscores (_) are
  allowed.  
  The name must start with a Chinese or English character, and the Chinese characters must be in **UTF-8**
  or **Unicode** format.

* `type` - (Required, String) Specifies the API type.  
  The valid values are as follows:
  + **API_SPECIFIC_TYPE_CONFIGURATION**
  + **API_SPECIFIC_TYPE_REGISTER**
  + **API_SPECIFIC_TYPE_ORCHESTRATE**
  + **API_SPECIFIC_TYPE_MYBATIS**
  + **API_SPECIFIC_TYPE_SCRIPT**
  + **API_SPECIFIC_TYPE_GROOVY**

* `auth_type` - (Required, String) Specifies the authentication type.  
  The valid values are as follows:
  + **APP**
  + **IAM**
  + **NONE**

* `protocol` - (Required, String) Specifies the request protocol of the API.  
  The valid values are as follows:
  + **PROTOCOL_TYPE_HTTP**
  + **PROTOCOL_TYPE_HTTPS**

* `path` (Required, String) Specifies the API request path.

* `request_type` - (Required, String) Specifies the request type of the API.  
  The valid values are as follows:
  + **REQUEST_TYPE_POST**
  + **REQUEST_TYPE_GET**

* `manager` - (Required, String) Specifies the API reviewer.

  -> Make sure the reviewer is already created in the data source reviewers management before the resource creation.

* `datasource_config` - (Required, List) Specifies the configuration of the API data source.  
  The [datasource_config](#dataservice_api_datasource_config) structure is documented below.

* `description` - (Optional, String) Specifies the description of API.  
  Maximum of `255` characters are allowed.

* `visibility` - (Optional, String) Specifies the visibility to the catalog of API.  
  The valid values are as follows:
  + **WORKSPACE**
  + **PROJECT**
  + **DOMAIN**

* `request_params` - (Optional, List) Specifies the parameters of the API request.  
  The [request_params](#dataservice_api_request_params) structure is documented below.

* `backend_config` - (Optional, List) Specifies the configuration of the API backend.  
  The [backend_config](#dataservice_api_backend_config) structure is documented below.

<a name="dataservice_api_datasource_config"></a>
The `datasource_config` block supports:

* `type` - (Required, String) Specifies the type of the data source.  
  The valid values are as follows:
  + **MYSQL**
  + **DLI**
  + **DWS**
  + **HIVE**
  + **HBASE**

* `connection_id` - (Required, String) Specifies the ID of the data connection for the DataArts Studio service.

* `database` - (Required, String) Specifies the name of the database.

* `datatable` - (Required, String) Specifies the name of the data table.

* `queue` - (Optional, String) Specifies the ID of the DLI queue.

* `access_mode` - (Optional, String) Specifies the access mode for the data.  
  The valid values are as follows:
  + **SQL**
  + **ROW_KEY**
  + **PREFIX_FILTER**

* `sql` - (Optional, String) Specifies the SQL statements in script access type.

* `backend_params` - (Optional, List) Specifies the backend parameters of the API.  
  The [backend_params](#dataservice_api_datasource_config_backend_params) structure is documented below.

* `response_params` - (Optional, List) Specifies the response parameters of the API.  
  The [response_params](#dataservice_api_datasource_config_response_params) structure is documented below.

* `order_params` - (Optional, List) Specifies the order parameters of the API.  
  The [order_params](#dataservice_api_datasource_config_order_params) structure is documented below.

-> All column names that appear in the data table must have corresponding parameter mappings.

<a name="dataservice_api_datasource_config_backend_params"></a>
The `backend_params` block supports:

* `name` - (Required, String) Specifies the name of the backend parameter.

* `mapping` - (Required, String) Specifies the name of the mapping parameter.

* `condition` - (Required, String) Specifies the condition character.  
  The valid values are as follows:
  + **CONDITION_TYPE_EQ**: =
  + **CONDITION_TYPE_NE**: <>
  + **CONDITION_TYPE_GT**: >
  + **CONDITION_TYPE_GE**: >=
  + **CONDITION_TYPE_LT**: <
  + **CONDITION_TYPE_LE**: <=
  + **CONDITION_TYPE_LIKE**: %like%
  + **CONDITION_TYPE_LIKE_L**: %like
  + **CONDITION_TYPE_LIKE_R**: like%

<a name="dataservice_api_datasource_config_response_params"></a>
The `response_params` block supports:

* `name` - (Required, String) Specifies the name of the response parameter.

* `type` - (Required, String) Specifies the type of the response parameter.  
  The valid values are as follows:
  + **REQUEST_PARAMETER_TYPE_NUMBER**
  + **REQUEST_PARAMETER_TYPE_STRING**

* `field` - (Required, String) Specifies the bound table field for the response parameter.

* `description` - (Optional, String) Specifies the description of the response parameter.

* `example_value` - (Optional, String) Specifies the example value of the response parameter.

<a name="dataservice_api_datasource_config_order_params"></a>
The `order_params` block supports:

* `name` - (Required, String) Specifies the name of the order parameter.

* `field` - (Required, String) Specifies the corresponding parameter field for the order parameter.

* `optional` - (Optional, Bool) Specifies whether this order parameter is the optional parameter.

* `sort` - (Optional, String) Specifies the sort type of the order parameter.  
  The valid values are as follows:
  + **ASC**
  + **DESC**
  + **CUSTOM**

* `order` - (Optional, Int) Specifies the order of the sorting parameters.

<a name="dataservice_api_request_params"></a>
The `request_params` block supports:

* `name` - (Required, String) Specifies the name of the request parameter.
  The valid length is limited from `1` to `32`, only letters, digits, dots (.), hyphens (-) and underscores (_) are
  allowed.  
  The name must start with a letter.  

  -> The parameter name cannot start with `x-apig-` or `x-sdk-`, and cannot be `x-stage`, `x-api-id`, `x-app-id`,
     `x-request-id`. When the parameter position is `HEADER`, the parameter name cannot be `Authorization` or
     `X-Auth-Token` (not case sensitive).

* `position` - (Required, String) Specifies the position of the request parameter.  
  The valid values are as follows:
  + **REQUEST_PARAMETER_POSITION_PATH**
  + **REQUEST_PARAMETER_POSITION_HEADER**
  + **REQUEST_PARAMETER_POSITION_QUERY**

* `type` - (Required, String) Specifies the type of the request parameter.  
  The valid values are as follows:
  + **REQUEST_PARAMETER_TYPE_NUMBER**
  + **REQUEST_PARAMETER_TYPE_STRING**

* `description` - (Optional, String) Specifies the description of the request parameter.  
  Maximum of `255` characters are allowed.

* `necessary` - (Optional, Bool) Specifies whether this parameter is the required parameter.

* `example_value` - (Optional, String) Specifies the example value of the request parameter.

* `default_value` - (Optional, String) Specifies the default value of the request parameter.

<a name="dataservice_api_backend_config"></a>
The `backend_config` block supports:

* `type` - (Required, String) Specifies the type of the backend request.  
  The valid values are as follows:
  + **REQUEST_TYPE_POST**
  + **REQUEST_TYPE_GET**

* `protocol` - (Required, String) Specifies the protocol of the backend request.  
  The valid values are as follows:
  + **PROTOCOL_TYPE_HTTP**
  + **PROTOCOL_TYPE_HTTPS**

* `host` - (Required, String) Specifies the backend host.

* `path` - (Required, String) Specifies the backend path.

* `timeout` - (Required, Int) Specifies the backend timeout.

* `backend_params` - (Optional, List) Specifies the backend parameters of the API.  
  The [backend_params](#dataservice_api_backend_config_backend_params) structure is documented below.

* `constant_params` - (Optional, List) Specifies the backend constant parameters of the API.  
  The [constant_params](#dataservice_api_backend_config_constant_params) structure is documented below.

<a name="dataservice_api_backend_config_backend_params"></a>
The `backend_params` block supports:

* `name` - (Required, String) Specifies the name of the request parameter.

* `position` - (Required, String) Specifies the position of the request parameter.  
  The valid values are as follows:
  + **REQUEST_PARAMETER_POSITION_PATH**
  + **REQUEST_PARAMETER_POSITION_HEADER**
  + **REQUEST_PARAMETER_POSITION_QUERY**

* `backend_param_name` - (Required, String) Specifies the name of the corresponding backend parameter.

<a name="dataservice_api_backend_config_constant_params"></a>
The `constant_params` block supports:

* `name` - (Required, String) Specifies the name of the constant parameter.

* `type` - (Required, String) Specifies the type of the constant parameter.  
  The valid values are as follows:
  + **REQUEST_PARAMETER_TYPE_NUMBER**
  + **REQUEST_PARAMETER_TYPE_STRING**

* `position` - (Required, String) Specifies the position of the constant parameter.  
  The valid values are as follows:
  + **REQUEST_PARAMETER_POSITION_PATH**
  + **REQUEST_PARAMETER_POSITION_HEADER**
  + **REQUEST_PARAMETER_POSITION_QUERY**

* `value` - (Required, String) Specifies the value of the constant parameter.

* `description` - (Optional, String) Specifies the description of the constant parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `create_user` - The creator name.

* `created_at` - The creation time of the API, in RFC3339 format.

* `updated_at` - The latest update time of the API, in RFC3339 format.

* `group_id` - The ID of the group to which the API belongs, for shared type.

* `status` - The API status, for shared type.

* `host` - The API host configuration, for shared type.

* `hosts` - The API host configuration, for exclusive type.  
  The [hosts](#dataservice_api_hosts_attr) structure is documented below.

<a name="dataservice_api_hosts_attr"></a>
The `hosts` block supports:

* `instance_id` - The cluster ID to which the API belongs.

* `instance_name` - The cluster name to which the API belongs.

* `intranet_host` - The intranet address.

* `external_host` - The exrernal address.

* `domains` - The list of gateway damains.

## Import

The API can be imported using `workspace_id`, `dlm_type` and `id` separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_dataarts_dataservice_api.test <workspace_id>/<dlm_type>/<id>
```

Also, you can omit `dlm_type` and provide just `workspace_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_dataservice_api.test <workspace_id>/<id>
```

~> This way only supports importing the API of the **SHARED** type, but does not support the API imported for
   **EXCLUSIVE** type. If an error is reported, please carefully check the `dlm_type` value to which imported API
   you want.

Note that the imported state may not be identical to your resource definition, because the attributes are missing in the
API response. The missing attributes includes: `auth_type`, `catalog_id` and `visibility`.
It is generally recommended running `terraform plan` after importing an resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_dataservice_api" "test" {
  ...

  lifecycle {
    ignore_changes = [
      auth_type, catalog_id, visibility,
    ]
  }
}
```
