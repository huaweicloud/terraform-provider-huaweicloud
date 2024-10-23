---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instances"
description: ""
---

# huaweicloud_ddm_instances

Use this data source to get the list of DDM instances.

## Example Usage

```hcl
variable "instance_name" {}

data "huaweicloud_ddm_instances" "test" {
  name = var.instance_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the instance.

* `status` - (Optional, String) Specifies the status of the instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id.

* `engine_version` - (Optional, String) Specifies the engine version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of DDM instance.
  The [Instance](#DdmInstances_Instance) structure is documented below.

<a name="DdmInstances_Instance"></a>
The `Instance` block supports:

* `status` - Indicates the status of the DDM instance.

* `name` - Indicates the name of the DDM instance.

* `availability_zones` - Indicates the list of availability zones.

* `vpc_id` - Indicates the ID of a VPC.

* `subnet_id` - Indicates the ID of a subnet.

* `security_group_id` - Indicates the ID of a security group.

* `node_num` - Indicates the number of nodes.

* `access_ip` - Indicates the address for accessing the DDM instance.

* `access_port` - Indicates the port for accessing the DDM instance.

* `enterprise_project_id` - Indicates the enterprise project id.

* `engine_version` - Indicates the engine version.
