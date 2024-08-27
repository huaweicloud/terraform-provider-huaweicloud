---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_dedicated_instance"
description: ""
---

# huaweicloud_waf_dedicated_instance

Manages a WAF dedicated instance resource within HuaweiCloud.

## Example Usage

### Creating with common tenant

```hcl
variable az_name {}
variable ecs_flavor_id {}
variable vpc_id {}
variable subnet_id {}
variable security_group_id {}
variable enterprise_project_id {}

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name                  = "instance_1"
  available_zone        = var.az_name
  specification_code    = "waf.instance.professional"
  ecs_flavor            = var.ecs_flavor_id
  vpc_id                = var.vpc_id
  subnet_id             = var.subnet_id
  enterprise_project_id = var.enterprise_project_id

  security_group = [
    var.security_group_id
  ]
}
```

### Creating with resource tenant

```hcl
variable az_name {}
variable vpc_id {}
variable subnet_id {}
variable security_group_id {}
variable enterprise_project_id {}

resource "huaweicloud_waf_dedicated_instance" "instance_1" {
  name                  = "instance_1"
  available_zone        = var.az_name
  specification_code    = "waf.instance.professional"
  vpc_id                = var.vpc_id
  subnet_id             = var.subnet_id
  enterprise_project_id = var.enterprise_project_id
  res_tenant            = true

  security_group = [
    var.security_group_id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the WAF dedicated instance. If omitted, the
  provider-level region will be used. Changing this setting will create a new instance.

* `name` - (Required, String) The name of WAF dedicated instance. Duplicate names are allowed, we suggest to keeping the
  name unique.

* `available_zone` - (Required, String, ForceNew) The available zone names for the dedicated instances. It can be
  obtained through this data source `huaweicloud_availability_zones`. Changing this will create a new instance.

* `specification_code` - (Required, String, ForceNew) The specification code of instance. Different specifications have
  different throughput. Changing this will create a new instance. Values are:
  + `waf.instance.professional` - The professional edition, throughput: 100 Mbit/s; QPS: 2,000 (Reference only).
  + `waf.instance.enterprise` - The enterprise edition, throughput: 500 Mbit/s; QPS: 10,000 (Reference only).

* `vpc_id` - (Required, String, ForceNew) The VPC id of WAF dedicated instance. Changing this will create a new
  instance.

* `subnet_id` - (Required, String, ForceNew) The subnet id of WAF dedicated instance VPC. Changing this will create a
  new instance.

* `enterprise_project_id` - (Optional, String) The enterprise project ID of WAF dedicated instance. Changing this
  will migrate the WAF instance to a new enterprise project.

* `security_group` - (Required, List, ForceNew) The security group of the instance. This is an array of security group
  ids. Changing this will create a new instance.

* `cpu_architecture` - (Optional, String, ForceNew) The ECS cpu architecture of instance, Default value is `x86`.
  Changing this will create a new instance.

* `ecs_flavor` - (Optional, String, ForceNew) Specifies the flavor of the ECS used by the WAF instance. Flavors can be
  obtained through this data source `huaweicloud_compute_flavors`. Changing this will create a new instance.
  This field is valid and required only when `res_tenant` is set to **false**.

  -> **NOTE:** If the instance specification is the professional edition, the ECS specification should be 2U4G. If the
  instance specification is the enterprise edition, the ECS specification should be 8U16G.

* `group_id` - (Optional, String, ForceNew) The instance group ID used by the WAF dedicated instance in ELB mode.
  Changing this will create a new instance.

* `res_tenant` - (Optional, Bool, ForceNew) Specifies whether this is resource tenant.
  Changing this will create a new instance.
  + **false**: Common tenant.
  + **true**: Resource tenant.

  Defaults to **false**.

* `anti_affinity` - (Optional, Bool, ForceNew) Specifies whether to enable anti-affinity. This field is valid only
  when `res_tenant` is set to **true**. Changing this will create a new instance.

## Attribute Reference

The following attributes are exported:

* `id` - The id of the instance.

* `server_id` - The id of the instance server.

* `service_ip` - The ip of the instance service.

* `run_status` - The running status of the instance. Values are:
  + `0` - Instance is creating.
  + `1` - Instance has created.
  + `2` - Instance is deleting.
  + `3` - Instance has deleted.
  + `4` - Instance create failed.

* `access_status` - The access status of the instance. `0`: inaccessible, `1`: accessible.

* `upgradable` - The instance is to support upgrades. `0`: Cannot be upgraded, `1`: Can be upgraded.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 20 minutes.

## Import

There are two ways to import WAF dedicated instance state.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_dedicated_instance.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_dedicated_instance.test <id>/<enterprise_project_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `res_tenant`, `anti_affinity`. It is generally recommended running
`terraform plan` after importing the resource. You can then decide if changes should be applied to the resource,
or the resource definition should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_waf_dedicated_instance" "test" {
  ...

  lifecycle {
    ignore_changes = [
      res_tenant, anti_affinity,
    ]
  }
}
```
