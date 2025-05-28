---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_migration_project"
description: |-
  Manages an SMS migration project resource within HuaweiCloud.
---

# huaweicloud_sms_migration_project

Manages an SMS migration project resource within HuaweiCloud.

-> Migration project can only be destroyed when it is not the default project.

## Example Usage

```hcl
variable "project_name" {}

resource "huaweicloud_sms_migration_project" "test" {
  name          = var.project_name
  region        = "cn-north-9"
  use_public_ip = true
  exist_server  = true
  type          = "MIGRATE_BLOCK"
  syncing       = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the migration project name.

* `region` - (Required, String) Specifies the region name.

* `use_public_ip` - (Required, Bool) Specifies whether to use a public IP address for migration.

* `exist_server` - (Required, Bool) Specifies whether the server already exists.

* `type` - (Required, String) Specifies the migration project type.
  Values can be as follows:
  + **MIGRATE_BLOCK**: Block-level migration.
  + **MIGRATE_FILE**: File-level migration.

* `syncing` - (Required, Bool) Specifies whether to continue syncing after the first copy or sync.

* `description` - (Optional, String) Specifies the migration project description.

* `is_default` - (Optional, Bool) Specifies whether to use the default template. Defaults to **false**.

  -> Only support to update from **false** to **true**.

* `start_target_server` - (Optional, Bool) Specifies whether to start the destination virtual machine after migration.
  Defaults to **false**.

* `speed_limit` - (Optional, Int) Specifies the migration rate limit in Mbps. Defaults to **0**.

* `enterprise_project` - (Optional, String) Specifies the name of the enterprise project. Defaults to **default**.

* `start_network_check` - (Optional, Bool) Specifies whether to enable network quality detection. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

SMS migration projects can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_sms_migration_project.demo <id>
```
