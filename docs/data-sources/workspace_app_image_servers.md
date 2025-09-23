---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_image_servers"
description: |-
  Use this data source to get the list of the image servers within HuaweiCloud.
---

# huaweicloud_workspace_app_image_servers

Use this data source to get the list of the image servers within HuaweiCloud.

## Example Usage

### Query all image servers

```hcl
data "huaweicloud_workspace_app_image_servers" "test" {}
```

### Query image server with the specified image server ID

```hcl
variable "image_server_id" {}

data "huaweicloud_workspace_app_image_servers" "test" {
  server_id = var.image_server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specified the name of the image server.
  Fuzzy search is supported.

* `server_id` - (Optional, String) Specified the ID of the image server.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the image server belongs.
  This parameter is only valid for enterprise users, if omitted, all enterprise project IDs will be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `servers` - All image servers that match the filter parameters.

  The [servers](#app_image_servers) structure is documented below.

<a name="app_image_servers"></a>
The `servers` block supports:

* `id` - The ID of the image server.

* `name` - The name of the image server.

* `description` - The description of the image server.

* `image_generated_product_id` - The ID of the generated image product.

* `image_id` - The ID of the basic image to which the image server belongs.

* `image_type` - The type of the basic image to which the image server belongs.
   + **gold**: The market image.
   + **public**: The public image.
   + **private**: The private image.
   + **shared**: The shared image.
   + **other**

* `spce_code` - The specification code of the basic image to which the image server belongs.

* `aps_server_group_id` - The ID of the APS server group associated with the image server.

* `aps_server_id` - The ID of the APS server associated with the image server.

* `app_group_id` - The ID of the application group associated with the image server.

* `status` - The current status of the image server.
  + **ACTIVE**
  + **BUILT**: Image task is finished.
  + **ERROR**

* `authorize_accounts` - The list of authorized users of the application group associated with the image server.

  The [authorize_accounts](#servers_authorize_accounts_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID to which the image server belongs.

* `created_at` - The creation time of the image server, in RFC3339 format.

* `updated_at` - The latest update time of the image server, in RFC3339 format.

<a name="servers_authorize_accounts_struct"></a>
The `authorize_accounts` block supports:

* `account` - The name of the account.

* `type` - The type of the account.

* `domain` - The domain name of the Workspace service.
