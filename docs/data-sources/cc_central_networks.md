---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_networks"
description: |-
  Use this data source to get the list of CC central networks.
---

# huaweicloud_cc_central_networks

Use this data source to get the list of CC central networks.

## Example Usage

```hcl
variable "central_network_name" {}

data "huaweicloud_cc_central_networks" "test" {
  name = var.central_network_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `central_network_id` - (Optional, String) Specifies the ID of central network.

* `name` - (Optional, String) Specifies the name of the central network. The name supports fuzzy query.

* `state` - (Optional, String) Specifies the status of the central network.

* `enterprise_project_id` - (Optional, String) Specifies enterprise project ID to which the central network belongs.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the central network.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `central_networks` - The central network list.

  The [central_networks](#central_networks_struct) structure is documented below.

<a name="central_networks_struct"></a>
The `central_networks` block supports:

* `id` - The central network ID.

* `name` - The central network name.

* `description` - The central network description.

* `state` - The central network status.

* `created_at` - The creation time of central network.

* `updated_at` - The update time of central network.

* `default_plane_id` - The central network default plane ID.

* `enterprise_project_id` - The ID of the enterprise project that the central network belongs to.

* `tags` - The central network tags.
