---
subcategory: "Dedicated Host (DeH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_deh_instance_servers"
description: |-
  Use this data source to get the list of DeH instance servers.
---

# huaweicloud_deh_instance_servers

Use this data source to get the list of DeH instance servers.

## Example Usage

```hcl
variable "dedicated_host_id" {}

data "huaweicloud_deh_instance_servers" "test" {
  dedicated_host_id = var.dedicated_host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `dedicated_host_id` - (Required, String) Specifies the dedicated host ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `servers` - Indicates the servers.

  The [servers](#servers_struct) structure is documented below.

<a name="servers_struct"></a>
The `servers` block supports:

* `id` - Indicates the server ID.

* `name` - Indicates the server name.

* `flavor` - Indicates the flavor infos.

  The [flavor](#servers_flavor_struct) structure is documented below.

* `image` - Indicates the server image infos.

  The [image](#servers_image_struct) structure is documented below.

* `addresses` - Indicates the server network address infos.
  Key is VPC ID and value is network infos.

* `status` - Indicates the server status.

* `task_state` - Indicates the task state.

* `metadata` - Indicates the server metadata.

  The [metadata](#servers_metadata_struct) structure is documented below.

* `user_id` - Indicates the server user ID.

* `created` - Indicates the server create time.

* `updated` - Indicates the server update time.

<a name="servers_flavor_struct"></a>
The `flavor` block supports:

* `id` - Indicates the flavor ID.

<a name="servers_image_struct"></a>
The `image` block supports:

* `id` - Indicates the server image ID

<a name="servers_metadata_struct"></a>
The `metadata` block supports:

* `os_type` - Indicates the server OS type.
