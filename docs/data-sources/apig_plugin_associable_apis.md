---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_plugin_associable_apis"
description: |-
  Use this data source to get the list of published APIs that can be bound to the current plugin within HuaweiCloud.
---

# huaweicloud_apig_plugin_associable_apis

Use this data source to get the list of published APIs that can be bound to the current plugin within HuaweiCloud.

## Example Usage

### Query all published APIs that can be bound to the specified plugin in an environment

```hcl
variable "instance_id" {}
variable "plugin_id" {}
variable "published_environment_id" {}

data "huaweicloud_apig_plugin_associable_apis" "test" {
  instance_id = var.instance_id
  plugin_id   = var.plugin_id
  env_id      = var.env_id
}
```

### Query all published APIs that can be bound to the specified plugin and filter by group ID

```hcl
variable "instance_id" {}
variable "plugin_id" {}
variable "published_environment_id" {}
variable "group_id" {}

data "huaweicloud_apig_plugin_associable_apis" "test" {
  instance_id = var.instance_id
  plugin_id   = var.plugin_id
  env_id      = var.published_environment_id
  group_id    = var.group_id
}
```

### Query all published APIs that can be bound to the specified plugin without any tags

```hcl
variable "instance_id" {}
variable "plugin_id" {}
variable "published_environment_id" {}

data "huaweicloud_apig_plugin_associable_apis" "test" {
  instance_id = var.instance_id
  plugin_id   = var.plugin_id
  env_id      = var.published_environment_id
  tags        = "#no_tags#"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the plugin and the associable APIs are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the plugin and the associable
  APIs belong.

* `plugin_id` - (Required, String) Specifies the ID of the plugin to be queried.

* `env_id` - (Required, String) Specifies the ID of the environment where the associable APIs are published.

* `api_name` - (Optional, String) Specifies the name of the associable APIs to be queried.  
  The fuzzy search is supported.

* `api_id` - (Optional, String) Specifies the ID of the specified associable API to be queried.

* `group_id` - (Optional, String) Specifies the ID of the API group to be queried to which the associable APIs belong.

* `req_method` - (Optional, String) Specifies the request method of the associable APIs to be queried.  
  The valid values are as follows:
  + **GET**
  + **POST**
  + **PUT**
  + **DELETE**
  + **HEAD**
  + **PATCH**
  + **OPTIONS**
  + **ANY**

* `tags` - (Optional, String) Specifies the tags of the associable APIs to be queried.
  The value `#no_tags#` means filtering APIs without tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apis` - The list of APIs that can be bound to the specified plugin.
  The [apis](#apig_plugin_associable_apis) structure is documented below.

<a name="apig_plugin_associable_apis"></a>
The `apis` block supports:

* `id` - The ID of the associable API.

* `name` - The name of the associable API.

* `type` - The type of the associable API.  
  The valid values are as follows:
  + **1** - Public API
  + **2** - Private API

* `req_protocol` - The request protocol of the associable API.  
  The valid values are as follows:
  + **HTTP**
  + **HTTPS**
  + **BOTH** - Both **HTTP** and **HTTPS**

* `req_method` - The request method of the associable API.

* `req_uri` - The request path of the associable API.

* `auth_type` - The authentication type of the associable API.  
  The valid values are as follows:
  + **NONE** - No authentication
  + **APP** - APP authentication
  + **IAM** - IAM authentication
  + **AUTHORIZER** - Custom authentication

* `match_mode` - The match mode of the associable API.  
  The valid values are as follows:
  + **SWA** - Prefix matching
  + **NORMAL** - Normal matching (exact matching)

* `description` - The description of the associable API.

* `group_id` - The ID of the API group to which the associable API belongs.

* `group_name` - The name of the API group to which the associable API belongs.

* `tags` - The tag list bound to the associable API.
