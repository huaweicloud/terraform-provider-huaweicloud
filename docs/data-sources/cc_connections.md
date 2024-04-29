---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_connections"
description: ""
---

# huaweicloud_cc_connections

Use this data source to get the list of cloud connections.

## Example Usage

```hcl
variable "id" {}
variable "name" {}

data "huaweicloud_cc_connections" "test" {
  connection_id = var.id
  name          = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `connection_id` - (Optional, String) Specifies the cloud connection ID.

* `name` - (Optional, String) Specifies the cloud connection name.

* `description` - (Optional, String) Specifies the cloud connection description.

* `status` - (Optional, String) Specifies the cloud connection status.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `used_scene` - (Optional, String) Specifies the application scenario.

* `tags` - (Optional, Map) Specifies the cloud connection tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - Cloud connection list.
  The [connections](#CloudConnections) structure is documented below.

<a name="CloudConnections"></a>
The `connections` block supports:

* `id` - Cloud connection ID.

* `name` - Cloud connection name.

* `description` - Cloud connection description.

* `domain_id` - ID of the account that the instance belongs to.

* `enterprise_project_id` - ID of the enterprise project that the cloud connection belongs to.

* `created_at` - Time when the cloud connection was created.

* `updated_at` - Time when the cloud connection was updated.

* `tags` - Cloud connection tags.

* `status` - Cloud connection status.

* `used_scene` - The application scenario of the cloud connection.

* `network_instance_number` - Number of the network instances loaded to the cloud connection.

* `bandwidth_package_number` - Number of the bandwidth packages bound to the cloud connection.

* `inter_region_bandwidth_number` - Number of the inter-region bandwidths configured for the cloud connection.
