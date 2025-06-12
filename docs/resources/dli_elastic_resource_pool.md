---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_elastic_resource_pool"
description: ""
---

# huaweicloud_dli_elastic_resource_pool

Manages an elastic resource pool within HuaweiCloud.

~> The elastic resource pool will regularly change the current number of CUs.
   The resource pool in the expansion status cannot be changed or deleted.
   Please refer to the [status](#dli_resource_pool_status) parameter.

## Example Usage

### Create a standard edition elastic resource pool under the default enterprise project

```hcl
variable "resoure_pool_name" {}

resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = var.resoure_pool_name
  description           = "Created by terraform script"
  min_cu                = 64
  max_cu                = 80
  cidr                  = "192.168.128.0/18"
  enterprise_project_id = "0"
}
```

### Create a basic elastic resource pool

```hcl
variable "resoure_pool_name" {}

resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name   = var.resoure_pool_name
  min_cu = 16
  max_cu = 64

  label = {
    spec = "basic"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the elastic resource pool is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the elastic resource pool.  
  The valid length is limited from `1` to `128`, only lowercase letters, digits and underscores (_) are allowed.  
  The name cannot contain only numbers or start with a number or an underscore.
  Changing this will create a new resource.

* `max_cu` - (Required, Int) Specifies the maximum number of CUs for elastic resource pool scaling.
  The interval is `16`.
  + For standard edition, the valid value ranges from `64` to `32,000`.
  + For basic edition, the valid value ranges from `16` to `64`.

* `min_cu` - (Required, Int) Specifies the minimum number of CUs for elastic resource pool scaling.
  The interval is `16`.
  + For standard edition, the valid value ranges from `64` to `32,000`.
  + For basic edition, the valid value ranges from `16` to `64`.

  ~> If the value needs to be updated, the `min_cu` value cannot be greater than the `current_cu` value.

* `description` - (Optional, String) Specifies the description of the elastic resource pool.  
  The valid length is limited from `1` to `128`.

* `cidr` - (Optional, String, ForceNew) Specifies the CIDR block of network to associate with the elastic resource pool.
  Defaults to `172.16.0.0/12`. Changing this will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the elastic resource
  pool belongs.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs to associate with the elastic resource pool.  
  Changing this will create a new resource.

* `label` - (Optional, Map, ForceNew) Specifies the attribute fields of the elastic resource pool.  
  Changing this will create a new resource.  
  If not specified, the default is the standard edition. The key/value corresponding to the basic edition is `spec = "basic"`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` -The resource ID in UUID format.

<a name="dli_resource_pool_status"></a>

* `status` - The current status of the elastic resource pool.

  + **AVAILABLE**: The resource pool can be used normally.
  + **SCALING**: The resource pool is changing CUs. There are not allow update and delete during this period.
  + **FAILED**: The resource pool is abnormal.

* `current_cu` - The current CU number of the elastic resource pool.

* `actual_cu` - The current CU number of the elastic resource pool.

* `created_at` - The creation time of the elastic resource pool.  
  The format is `YYYY-MM-DDThh:mm:ss{timezone}`, e.g. `2024-01-01T08:00:00+08:00`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

Elastic resource pools can be imported by their `name`, e.g.

```bash
$ terraform import huaweicloud_dli_elastic_resource_pool.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `tags` and `label`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dli_elastic_resource_pool" "test" {
  ...

  lifecycle {
    ignore_changes = [
      tags, label,
    ]
  }
}
```
