---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_relationship"
description: |-
  Manages a GaussDB DR relationship resource within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_relationship

Manages a GaussDB DR relationship resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "dr_ip" {}
variable "dr_user_name" {}
variable "dr_user_password" {}

resource "huaweicloud_gaussdb_dr_relationship" "test" {
  instance_id      = var.instance_id
  disaster_type    = "stream"
  dr_ip            = var.dr_ip
  dr_user_name     = var.dr_user_name
  dr_user_password = var.dr_user_password
  dr_task_name     = "test-dr-task"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

* `disaster_type` - (Required, String, NonUpdatable) Specifies the disaster recovery type.
  The valid values are as follows:
  + **stream**: Stream disaster recovery.

* `dr_ip` - (Required, String, NonUpdatable) Specifies the data IP of the remote instance.

* `dr_user_name` - (Required, String, NonUpdatable) Specifies the account name of the remote instance.

* `dr_user_password` - (Required, String, NonUpdatable) Specifies the account password of the remote instance.

* `dr_task_name` - (Optional, String, NonUpdatable) Specifies the disaster recovery task name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.

* `dr_id` - The DR ID.

* `synchronization_id` - The ID of the DR relationship.

* `status` - The status of the disaster recovery record.
  The valid values are as follows:
  + **pending**: Task is being processed.
  + **normal**: Disaster recovery relationship is normal.
  + **failed**: Disaster recovery build failed.
  + **completed**: Disaster recovery relationship has been released.
  + **failover**: Disaster recovery has been promoted.
  + **simulation**: In simulation mode.
  + **dr_log_keep**: Log retention in progress.
  + **pre_checking**: Disaster recovery pre-check.
  + **pre_check_failed**: Disaster recovery pre-check failed.

* `precheck_failed_reason` - The reason for pre-check failure.

* `disaster_role` - The disaster recovery role.
  The valid values are as follows:
  + **master**: Primary instance.
  + **disaster**: Disaster recovery instance.

* `created` - The creation time.

* `updated` - The update time.

* `slave_region_instance_info` - The instance information in the DR region.
  The [slave_region_instance_info](#gaussdb_dr_relationship_region_instance_info) structure is documented below.

* `master_region_instance_info` - The instance information in the primary region.
  The [master_region_instance_info](#gaussdb_dr_relationship_region_instance_info) structure is documented below.

* `instance_name` - The instance name.

* `instance_status` - The instance status.

* `actions` - The list of actions currently being executed on the instance.

<a name="gaussdb_dr_relationship_region_instance_info"></a>
The `slave_region_instance_info` and `master_region_instance_info` block supports:

* `region_code` - The region code.

* `instance_id` - The instance ID.

* `project_id` - The project ID.

* `project_name` - The project name.

* `ip_address` - The data IP address list, separated by commas.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The GaussDB DR relationship can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_dr_relationship.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `dr_ip`, `dr_user_name`, `dr_user_password`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the resource. Also you can ignore
changes as below.

```hcl
resource "huaweicloud_gaussdb_dr_relationship" "test" {
  ...

  lifecycle {
    ignore_changes = [
      dr_ip, dr_user_name, dr_user_password,
    ]
  }
}
```
