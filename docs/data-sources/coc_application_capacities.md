---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_application_capacities"
description: |-
  Use this data source to get the list of COC application capacities.
---

# huaweicloud_coc_application_capacities

Use this data source to get the list of COC application capacities.

## Example Usage

```hcl
variable "application_id" {}

data "huaweicloud_coc_application_capacities" "test" {
  application_id = var.application_id
  provider_obj {
    cloud_service_name = "ecs"
    type               = "cloudservers"
  }
}
```

## Argument Reference

The following arguments are supported:

* `provider_obj` - (Required, List) Specifies the resource object.
  The [provider_obj](#provider_obj_struct) structure is documented below.

* `group_id` - (Optional, String) Specifies the group ID.

* `component_id` - (Optional, String) Specifies the component ID.

* `application_id` - (Optional, String) Specifies the application ID.

-> Exactly one of `group_id`, `component_id` or `application_id` must be set.

<a name="provider_obj_struct"></a>
The `provider_obj` block supports:

* `cloud_service_name` - (Optional, String) Specifies the cloud service name.
  The value can be **ecs**, **cce**, **rds** and so on.

* `type` - (Optional, String) Specifies the resource type name.
  There are many resource types. Choose the resource type based on your business needs. Common resource types are as follows:
  + **cloudservers**: Elastic Cloud Servers.
  + **servers**: Bare Metal Servers.
  + **clusters**: Cloud Container Engines.
  + **instances**: Cloud Databases.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the application capacity list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `sum_size` - Indicates the total size of the hard disk.

* `sum_cpu` - Indicates the total amount of CPU allocated.

* `sum_mem` - Indicates the total amount of memory allocated.

* `cloud_service_name` - Indicates the cloud service name.

* `type` - Indicates the resource type name.
