---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_member_check_task"
description: |-
  Manages an ELB member check task resource within HuaweiCloud.
---

# huaweicloud_elb_member_check_task

Manages an ELB member check task resource within HuaweiCloud.

## Example Usage

```hcl
variable "member_id" {}
variable "listener_id" {}

resource "huaweicloud_elb_member_check_task" "test" {
  member_id   = var.member_id
  listener_id = var.listener_id
  subject     = "all"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `member_id` - (Required, String, NonUpdatable) Specifies the backend server ID.

* `listener_id` - (Required, String, NonUpdatable) Specifies the ID of a listener to query the status of the backend
  servers that are associated with this listener.

* `subject` - (Required, String, NonUpdatable) Specifies the check items. Value options:
  + **securityGroup**: security group checks
  + **networkAcl**: network ACL checks
  + **config**: health check port checks
  + **all**: all check items

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the project ID.

* `status` - Indicates the status of a backend server check task.

* `result` - Indicates the result of each check item.
  The [result](#result_struct) structure is documented below.

* `check_item_total_num` - Indicates the total number of check items.

* `check_item_finished_num` - Indicates the number of checked items.

* `created_at` - Indicates the task creation time.

* `updated_at` - Indicates the task update time.

<a name="result_struct"></a>
The `result` block supports:

* `config` - Indicates the configuration check.
  The [config](#result_group_struct) structure is documented below.

* `acl` - Indicates the network ACL rule check.
  The [acl](#result_group_struct) structure is documented below.

* `security_group` - Indicates the security group rule check.
  The [security_group](#result_group_struct) structure is documented below.

<a name="result_group_struct"></a>
The `config`, `acl`, `security_group` block supports:

* `check_result` - Indicates the check result. **true** indicates that the check is passed, and **false** indicates that
  the check is not passed.

* `check_items` - Indicates the summary of grouped check items.
  The [check_items](#check_items_struct) structure is documented below.

* `check_status` - Indicates the status of a backend server check task. The value can be **processed**, **processing**
  or **failed**.

<a name="check_items_struct"></a>
The `check_items` block supports:

* `name` - Indicates the check item name.

* `reason` - Indicates the exception cause.

* `severity` - Indicates the exception severity, which can be **Major (severe)** or **Tips (informational)**.

* `subject` - Indicates the check type. **config** indicates configuration check.

* `job_id` - Indicates the task ID.

* `reason_template` - Indicates the exception reason template.

* `reason_params` - Indicates the exception variables, which is used to dynamically generate exception causes based on
  the exception cause template.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

The ELB member check task resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_member_check_task.test <id>
```
