---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_parameter_template_apply"
description: |-
  Manages a DDS parameter template apply resource within HuaweiCloud.
---

# huaweicloud_dds_parameter_template_apply

Manages a DDS parameter template apply resource within HuaweiCloud.

-> Please check whether the entities need to be restarted after applying parameter template.

## Example Usage

```hcl
variable "configuration_id" {}
variable "entity_ids" {}

resource "huaweicloud_dds_parameter_template_apply" "test" {
  configuration_id = var.configuration_id
  entity_ids       = var.entity_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `configuration_id` - (Required, String, ForceNew) Specifies the parameter template ID.
  Changing this creates a new resource.

* `entity_ids` - (Required, List, ForceNew) Specifies the entity IDs.
  + If the DB instance type is cluster and the shard or config parameter template is to be changed, the value is the
  group ID. If the parameter template of the mongos node is to be changed, the value is the node ID.
  + If the DB instance to be changed is a replica set instance, the value is the instance ID.

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
