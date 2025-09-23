---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_parameters_modify"
description: |-
  Manages a DDS instance parameters modify resource within HuaweiCloud.
---

# huaweicloud_dds_instance_parameters_modify

Manages a DDS instance parameters modify resource within HuaweiCloud.

-> **NOTE:** Deleting instance parameters modify is not supported. If you destroy a resource of instance parameters
  modify, it is only removed from the state, but still remains in the cloud. And the instance doesn't return to the
  state before modifying.

## Example Usage

### Modify Replica Set instance parameters

```hcl
variable "instance_id" {}
variable "parameter_name" {}
variable "parameter_value" {}

resource "huaweicloud_dds_instance_parameters_modify" "test" {
  instance_id = var.instance_id

  parameters {
    name  = var.parameter_name
    value = var.parameter_value
  }
}
```

### Modify Cluster Community instance's entity parameters

```hcl
variable "instance_id" {}
variable "entity_id" {}
variable "parameter_name" {}
variable "parameter_value" {}

resource "huaweicloud_dds_instance_parameters_modify" "test" {
  instance_id = var.instance_id
  entity_id   = var.entity_id

  parameters {
    name  = var.parameter_name
    value = var.parameter_value
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDS instance.
  Changing this creates a new resource.

* `parameters` - (Required, List) Specify an array of one or more parameters to be set to the DDS instance or entity.
  You can check on console to see which parameters supported.
  The [parameters](#block--parameters) structure is documented below.

* `entity_id` - (Optional, String, ForceNew) Specifies the ID of a DDS instance entity.
  + If the DB instance type is cluster and the shard or config parameter template is to be changed, the value is the
  group ID. If the parameter template of the mongos node is to be changed, the value is the node ID.
  + If the DB instance to be changed is a replica set instance, the value should be empty.

  Changing this creates a new resource.

<a name="block--parameters"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name. Some of them needs a restart of instance to take effect.

* `value` - (Required, String) Specifies the parameter value.

## Attribute Reference

* `id` - Indicates the resource ID. Same as `<instance_id>` or `<instance_id>/<entity_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `update` - Default is 30 minutes.
