---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_dataservice_catalog_apis"
description: |-
  Use this data source to query the API list under a specified catalog within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_catalog_apis

Use this data source to query the API list under a specified catalog within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "catalog_id" {}

data "huaweicloud_dataarts_dataservice_catalog_apis" "test" {
  workspace_id = var.workspace_id
  catalog_id   = var.catalog_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the catalog APIs are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the catalog belongs.

* `catalog_id` - (Required, String) Specifies the ID of the catalog to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apis` - The list of the catalog APIs that matched filter parameters.
  The [apis](#catalog_apis) structure is documented below.

<a name="catalog_apis"></a>
The `apis` block supports:

* `id` - The ID of the API.

* `name` - The name of the API.

* `group_id` - The ID of the API group.

* `description` - The description of the API.

* `status` - The status of the API.

* `debug_status` - The debug status of the API.

* `publish_messages` - The publish information list of the API.
  The [publish_messages](#catalog_apis_publish_messages_attr) structure is documented below.

* `type` - The type of the API.

* `manager` - The manager of the API.

* `create_user` - The creator of the API.

* `create_time` - The creation time of the API, in RFC3339 format.

* `authorization_status` - The authorization status of the API.

<a name="catalog_apis_publish_messages_attr"></a>
The `publish_messages` block supports:

* `id` - The ID of the published message.

* `api_id` - The ID of the API.

* `instance_id` - The ID of the instance.

* `instance_name` - The name of the instance.

* `api_status` - The status of the API.

* `api_debug` - The debug status of the API.
