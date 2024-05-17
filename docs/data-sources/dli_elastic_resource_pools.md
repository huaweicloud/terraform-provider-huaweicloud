---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_elastic_resource_pools"
description: |-
  Use this data source to get the list of DLI elastic resource pools within HuaweiCloud.
---

# huaweicloud_dli_elastic_resource_pools

Use this data source to get the list of DLI elastic resource pools within HuaweiCloud.

## Example Usage

```hcl
variable "resoure_pool_name" {}

data "huaweicloud_dli_elastic_resource_pools" "test" {
  name = var.resoure_pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the elastic resource pool.

* `status` - (Optional, String) Specifies the status of the elastic resource pool.
  The valid values are as follows:
  + **available**
  + **failed**

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the elastic resource pool.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `elastic_resource_pools` - All elastic resource pools that match the filter parameters.

  The [elastic_resource_pools](#elastic_resource_pools_struct) structure is documented below.

<a name="elastic_resource_pools_struct"></a>
The `elastic_resource_pools` block supports:

* `id` - The elastic resource pool ID.

* `name` - The elastic resource pool name.

* `max_cu` - The maximum CUs number of the elastic resource pool.

* `min_cu` - The minimum CUs number of the elastic resource pool.

* `current_cu` - The current CUs number of the elastic resource pool.

* `actual_cu` - The actual CUs number of the elastic resource pool.

* `cidr` - The CIDR block of network to associate with the elastic resource pool.

* `resource_id` - The resource ID of the elastic resource pool.

* `enterprise_project_id` - The enterprise project ID corresponding to the elastic resource pool.

* `queues` - The list of queues association with the elastic resource pool.

* `description` - The description of the elastic resource pool.

* `status` - The current status of the elastic resource pool.

* `owner` - The account name for creating elastic resource pool.

* `manager` - The type of the elastic resource pool.

* `fail_reason` - The reason of elastic resource pool creation failed.

* `created_at` - The creation time of the elastic resource pool.
