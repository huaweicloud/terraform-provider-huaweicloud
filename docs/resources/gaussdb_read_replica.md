---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_read_replica"
description: |-
  Manages a GaussDB read replica resource within HuaweiCloud.
---

# huaweicloud_gaussdb_read_replica

Manages a GaussDB read replica resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "configuration_id" {}

resource "huaweicloud_gaussdb_read_replica" "test" {
  instance_id       = var.instance_id
  availability_zone = "cn-north-4a"
  flavor_ref        = "gaussdb.bs.s6.xlarge.x864.ha"
  configuration_id  = var.configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the read replica.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance to which the read replica
  belongs.

* `availability_zone` - (Required, String, NonUpdatable) Specifies the availability zone where the read replica is located.

* `flavor_ref` - (Required, String, NonUpdatable) Specifies the specification code of the read replica.

* `configuration_id` - (Required, String, NonUpdatable) Specifies the parameter template ID for the read replica.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `name` - The name of the read replica.

* `status` - The status of the read replica. The valid values are as follows:
  + **BUILD**: The read replica is being created.
  + **ACTIVE**: The read replica is normal.
  + **FAILED**: The read replica failed to be created.
  + **DELETING**: The read replica is being deleted.

* `private_ip` - The private IP of the read replica.

* `data_ip` - The data IP of the read replica.

* `component_names` - The component names of the read replica.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

GaussDB read replica can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_read_replica.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `configuration_id`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the read replica, or the resource definition should be updated to
align with the read replica. Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_read_replica" "test" {
  ...

  lifecycle {
    ignore_changes = [
      configuration_id,
    ]
  }
}
```
