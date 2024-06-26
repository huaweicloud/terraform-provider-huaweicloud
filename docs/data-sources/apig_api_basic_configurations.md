---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_basic_configurations"
description: |-
  Use this data source to get the basic configuration list of the APIs within HuaweiCloud.
---

# huaweicloud_apig_api_basic_configurations

Use this data source to get the basic configuration list of the APIs within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_apig_api_basic_configurations" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the APIs belong.

* `api_id` - (Optional, String) Specifies the ID of the API.

* `name` - (Optional, String) Specifies the name of the API. Fuzzy search is supported.

* `group_id` - (Optional, String) Specifies the ID of the API group to which the APIs belong.

* `type` - (Optional, String) Specifies the type of the API.  
  The valid values are **Public** and **Private**.

* `request_method` - (Optional, String) Specifies the request method of the API.  
  The valid values are **GET**, **POST**, **PUT**, **DELETE**, **HEAD**, **PATCH**, **OPTIONS** and **ANY**.

* `request_path` - (Optional, String) Specifies the request address of the API. Fuzzy search is supported.

* `request_protocol` - (Optional, String) Specifies the request protocol of the API.  
  The valid values are **HTTP**, **HTTPS**, **BOTH** and **GRPCS**.

* `security_authentication` - (Optional, String) Specifies the security authentication mode of the API request.  
  The valid values are **NONE**, **APP**, **IAM** and **AUTHORIZER**.

* `vpc_channel_name` - (Optional, String) Specifies the name of the VPC channel. Fuzzy search is supported.

* `precise_search` - (Optional, String) Specifies the parameter name that you want to match exactly.  
  The valid values are as follows:
  + **name**: API name, corresponding to the field `name` in this data source arguments.
  + **req_uri**: Request path, corresponding to the field `request_path` in this data source arguments.

  This parameter can also be set to multiple enumerated values and separated by a comma (,), e.g. `name,req_uri`.  
  This parameter takes effect only after the corresponding parameter(s) is(are) set.

* `env_id` - (Optional, String) Specifies the ID of the environment where the API is published.

* `env_name` - (Optional, String) Specifies the name of the environment where the API is published.

* `backend_type` - (Optional, String) Specifies the backend type of the API.  
  The valid values are **HTTP**, **FUNCTION**, **MOCK** and **GRPC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - All API configurations that match the filter parameters.
  The [configurations](#basic_configurations) structure is documented below.

<a name="basic_configurations"></a>
The `configurations` block supports:

* `id` - The ID of the API.

* `name` - The name of the API.

* `type` - The type of the API.

* `request_method` - The request method of the API.

* `request_path` - The request address of the API.

* `request_protocol` - The request protocol of the API.

* `security_authentication` - The security authentication mode of the API request.

* `simple_authentication` - Whether the authentication of the application code is enabled.

* `group_id` - The ID of group corresponding to the API.

* `group_name` - The name of group corresponding to the API.

* `group_version` - The version of group corresponding to the API.

* `env_id` - The ID of the environment where the API is published.

* `env_name` - The name of the environment where the API is published.

* `authorizer_id` - The ID of the authorizer to which the API request used.

* `publish_id` - The ID of publish corresponding to the API.

* `backend_type` - The backend type of the API.

* `cors` - Whether CORS is supported.

* `matching` - The matching mode of the API.  
  + **Exact**
  + **Prefix**

* `description` - The description of the API.

* `tags` - The list of tags configuration.

* `registered_at` - The registered time of the API, in RFC3339 format.

* `updated_at` - The latest update time of the API, in RFC3339 format.

* `published_at` - The published time of the API, in RFC3339 format.
