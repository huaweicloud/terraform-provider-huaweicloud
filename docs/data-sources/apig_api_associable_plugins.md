---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_associable_plugins"
description: |-
  Use this data source to get the list of plugins that can be bound to the specified and published API within
  HuaweiCloud.
---

# huaweicloud_apig_api_associable_plugins

Use this data source to get the list of plugins that can be bound to the specified and published API within HuaweiCloud.

## Example Usage

### Query all plugins that can be bound to the specified API in an environment

```hcl
variable "instance_id" {}
variable "api_id" {}
variable "published_environment_id" {}

data "huaweicloud_apig_api_associable_plugins" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id
  env_id      = var.published_environment_id
}
```

### Query all plugins that can be bound to the specified API and filter by plugin type

```hcl
variable "instance_id" {}
variable "api_id" {}
variable "published_environment_id" {}

data "huaweicloud_apig_api_associable_plugins" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id
  env_id      = var.published_environment_id
  plugin_type = "cors"
}
```

### Query all plugins that can be bound to the specified API and filter by plugin name

```hcl
variable "instance_id" {}
variable "api_id" {}
variable "published_environment_id" {}
variable "plugin_name" {}

data "huaweicloud_apig_api_associable_plugins" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id
  env_id      = var.published_environment_id
  plugin_name = var.plugin_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the API and the associable plugins are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the API and the associable
  plugins belong.

* `api_id` - (Required, String) Specifies the ID of the API to be queried.

* `env_id` - (Optional, String) Specifies the ID of the environment where the API is published.

* `plugin_name` - (Optional, String) Specifies the name of the associable plugins to be queried.  
  The fuzzy search is supported.

* `plugin_type` - (Optional, String) Specifies the type of the associable plugins to be queried.  
  The valid values are as follows:
  + **cors** - Cross-Origin Resource Sharing
  + **set_resp_headers** - HTTP Response Header Management
  + **kafka_log** - Kafka Log Push
  + **breaker** - Circuit Breaker
  + **rate_limit** - Rate Limiting
  + **third_auth** - Third-party Authentication
  + **proxy_cache** - Response Cache
  + **proxy_mirror** - Request Mirror
  + **oidc_auth** - OIDC Authentication
  + **jwt_auth** - JWT Authentication

* `plugin_id` - (Optional, String) Specifies the ID of the specified associable plugin to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - The list of plugins that can be bound to the specified API.
  The [plugins](#apig_api_associable_plugins) structure is documented below.

<a name="apig_api_associable_plugins"></a>
The `plugins` block supports:

* `id` - The ID of the associable plugin.

* `name` - The name of the associable plugin.

* `type` - The type of the associable plugin.

* `scope` - The scope of the associable plugin.  
  The valid values are as follows:
  + **global** - Global scope

* `content` - The content of the associable plugin configuration.

* `description` - The description of the associable plugin.

* `create_time` - The creation time of the associable plugin, in RFC3339 format.

* `update_time` - The update time of the associable plugin, in RFC3339 format.
