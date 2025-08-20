---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_quotas"
description: |-
  Use this data source to get the quota list of Workspace APP server within HuaweiCloud.
---

# huaweicloud_workspace_app_server_quotas

Use this data source to get the quota list of Workspace APP server within HuaweiCloud.

## Example Usage

### Query all quotas under specified product ID

```hcl
variable "product_id" {}

data "huaweicloud_workspace_app_server_quotas" "test" {
  product_id       = var.product_id
  subscription_num =  1
  disk_size        =  80
  disk_num         =  2
}
```

### Query quotas under specified product ID by flavor ID

```hcl
variable "product_id" {}
variable "flavor_id" {}

data "huaweicloud_workspace_app_server_quotas" "advanced" {
  product_id       = var.product_id
  subscription_num = 1
  disk_size        = 80
  disk_num         = 2
  flavor_id        = var.flavor_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the Workspace APP server quotas are located.  
  If omitted, the provider-level region will be used.

* `product_id` - (Required, String) Specifies the ID of the product to be queried.

* `subscription_num` - (Required, Int) Specifies the number of server instances to be queried.

* `disk_size` - (Required, Int) Specifies the disk size of the single server instance to be queried.

* `disk_num` - (Required, Int) Specifies the number of disks for the single server instance to be queried.  
  The valid value ranges from `1` to `11`.

* `flavor_id` - (Optional, String) Specifies the ID of the flavor to be queried.

* `is_period` - (Optional, Bool) Specifies whether the instance is prepaid.  
  Defaults to **false**.

* `deh_id` - (Optional, String) Specifies the ID of the dedicated host.

* `cluster_id` - (Optional, String) Specifies the ID of the cloud dedicated distributed storage pool.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `is_enough` - Whether the quota is sufficient.

* `quotas` - The list of the quotas that match the filter parameters.  
  The [quotas](#data_app_server_quotas) structure is documented below.

<a name="data_app_server_quotas"></a>
The `quotas` block supports:

* `type` - The quota resource type.
  + **GPU_INSTANCES**: Number of GPU resource instances.
  + **INSTANCES**: Number of Normal instances.
  + **VOLUME_GIGABYTES**: Total disk capacity, in GB.
  + **VOLUMES**: Number of disks.
  + **CORES**: Number of CPUs.
  + **MEMORY**: Memory capacity, in MB.

* `remainder` - The remaining quota.

* `need` - The required quota.
