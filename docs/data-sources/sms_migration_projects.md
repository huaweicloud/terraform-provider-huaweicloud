---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_migration_projects"
description: |-
  Use this data source to get the list of SMS migration projects.
---

# huaweicloud_sms_migration_projects

Use this data source to get the list of SMS migration projects.

## Example Usage

```hcl
data "huaweicloud_sms_migration_projects" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `migprojects` - Indicates the details of the queried migration projects.

  The [migprojects](#migprojects_struct) structure is documented below.

<a name="migprojects_struct"></a>
The `migprojects` block supports:

* `id` - Indicates the migration project ID.

* `name` - Indicates the migration project name.

* `use_public_ip` - Indicates whether to use a public IP address for migration.

* `is_default` - Indicates whether the migration project is the default project.

* `start_target_server` - Indicates whether to start the target server after the migration.

* `region` - Indicates the region name.

* `speed_limit` - Indicates the migration rate limit configured in the project.

* `exist_server` - Indicates whether there are servers in the migration project.

* `description` - Indicates the migration project description.

* `type` - Indicates the type of the migration project.

* `enterprise_project` - Indicates the name of the enterprise project to which the migration project belongs.

* `syncing` - Indicates whether to perform a continuous synchronization after the full replication is complete.

* `start_network_check` - Indicates whether to enable network performance measurement.
