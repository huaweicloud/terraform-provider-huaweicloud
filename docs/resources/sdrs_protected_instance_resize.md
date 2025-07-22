---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protected_instance_resize"
description: |-
  Using this resource to resize a protected instance in SDRS within HuaweiCloud.
---

# huaweicloud_sdrs_protected_instance_resize

Using this resource to resize a protected instance in SDRS within HuaweiCloud.

-> This is a one-time action resource to resize a protected instance. Deleting this resource will
not change the current configuration, but will only remove the resource information from the tfstate file. You can use
this resource in the following scenarios:
<br/>1. Modify the specifications of both the production and DR site servers.
<br/>2. Modify the specifications of only the production site server.
<br/>3. Modify the specifications of only the DR site server.
<br/>4. Running this resource may cause unexpected change to the `flavor_id` field of the `huaweicloud_compute_instance`.
Please using `lifecycle` to control the `flavor_id` field of the `huaweicloud_compute_instance`.

-> Servers of different specifications have different performance, which may affect applications running on the servers.
To ensure the server performance after a planned failover or failover, you are recommended to use servers of
specifications (CPU and memory) same or higher than the specifications of the production site servers at the DR site.

-> Before using this resource, please note the following restrictions:
<br/>1. The status of the protected group must be **available** or **protected**.
<br/>2. The status of the protected instance must be **available** or **protected** or **error-resizing**.
<br/>3. The ECS to be resized must be in the shutdown state.

## Example Usage

```hcl
variable "protected_instance_id" {}
variable "flavor_ref" {}

resource "huaweicloud_sdrs_protected_instance_resize" "test" {
  protected_instance_id = var.protected_instance_id
  flavor_ref            = var.flavor_ref
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted,
  the provider-level region will be used. Changing this will create a new resource.

* `protected_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the protected instance to resize.
  You can obtain this value by calling the datasource `huaweicloud_sdrs_protected_instances`.

* `flavor_ref` - (Optional, String, NonUpdatable) Specifies the flavor ID of the production and DR site servers after
  the modification.

  -> If you specify this parameter, the system modifies the specifications of both the production and DR site servers.
  After the modification, the production site server and DR site server use the same specifications.

* `production_flavor_ref` - (Optional, String, NonUpdatable) Specifies the flavor ID of the production site server after
  the modification.

  -> 1. If you specify this parameter, the system modifies the specifications of only the production site server.
  <br/>2. If `flavor_ref` is specified, `production_flavor_ref` does not take effect.

* `dr_flavor_ref` - (Optional, String, NonUpdatable) Specifies the flavor ID of the DR site server after the modification.

  -> 1. If you specify this parameter, the system modifies the specifications of only the DR site server.
  <br/>2. If `flavor_ref` is specified, `dr_flavor_ref` does not take effect.

-> You can obtain the values of fields `flavor_ref`, `production_flavor_ref`, and `dr_flavor_ref` by calling the API
  [Querying the Target ECS Flavors to Which a Flavor Can Be Changed](https://support.huaweicloud.com/intl/en-us/api-ecs/ecs_02_0402.html).

* `production_dedicated_host_id` - (Optional, String, NonUpdatable) Specifies the new DeH ID for the production site.

  -> 1. If the production site server is created on a DeH, this parameter must be specified when you modify the specifications
  of the production site server.
  <br/>2. You can set this parameter to the ID of the DeH where the production site server is currently located or the ID
  of another DeH.

* `dr_dedicated_host_id` - (Optional, String, NonUpdatable) Specifies the new DeH ID for the DR site.

  -> 1. If the DR site server is created on a DeH, this parameter must be specified when you modify the specifications of
  the DR site server.
  <br/>2. You can set this parameter to the ID of the DeH where the DR site server is currently located or the ID of
  another DeH.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
