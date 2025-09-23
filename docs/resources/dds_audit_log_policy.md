---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_audit_log_policy"
description: ""
---

# huaweicloud_dds_audit_log_policy

Manages a DDS audit log policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "keep_days" {}

resource "huaweicloud_dds_audit_log_policy" "test"{
  instance_id = var.instance_id
  keep_days   = var.keep_days
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DDS instance.

  Changing this parameter will create a new resource.

* `keep_days` - (Required, Int) Specifies the number of days for storing audit logs.  
  The value ranges from `7` to `732`.

* `audit_scope` - (Optional, String) Specifies the audit scope.
  If this parameter is left blank or set to **all**, all audit log policies are enabled.
  You can enter the database or collection name. Use commas (,) to separate multiple databases
  or collections. If the name contains a comma (,), add a dollar sign ($) before the comma
  to distinguish it from the separators. Enter a maximum of 1024 characters. The value
  cannot contain spaces or the following special characters "[]{}():? The dollar sign ($)
  can be used only in escape mode.

* `audit_types` - (Optional, List) Specifies the audit type. Value options:
  + **auth**
  + **insert**
  + **delete**
  + **update**
  + **query**
  + **command**

* `reserve_auditlogs` - (Optional, String) Specifies whether the historical audit logs are
  retained when SQL audit is disabled.
    + **true** (default value): indicates that historical audit logs are retained
      when SQL audit is disabled.
    + **false**: indicates that existing historical audit logs are deleted when
      SQL audit is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The DDS audit log policy can be imported using the instance ID, e.g.:

```bash
$ terraform import huaweicloud_dds_audit_log.test <instance_id>
```
