---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_lts_configs"
description: |-
  Use this data source to get the list of log configuration (LTS) configs.
---

# huaweicloud_rds_lts_configs

Use this data source to get the list of log configuration (LTS) configs.

## Example Usage

```hcl
variable "engine" {}

data "huaweicloud_rds_lts_configs" "test" {
  engine = var.engine
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `engine` - (Required, String) Specifies the RDS engine type. Value options:
  **mysql**, **postgresql**, **sqlserver**.

* `enterprise_project_id` - (Optional, String) Specifies the project ID.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `instance_name` - (Optional, String) Specifies the instance name.

* `sort` - (Optional, String) Specifies the sort criteria for the returned instances.

* `instance_status` - (Optional, String) Specifies the instance status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_lts_configs` - Indicates the list of LTS configuration objects.

  The [instance_lts_configs](#instance_lts_configs_struct) structure is documented below.

<a name="instance_lts_configs_struct"></a>
The `instance_lts_configs` block supports:

* `lts_configs` - Indicates the list of log types and their LTS upload configuration.

  The [lts_configs](#lts_configs_struct) structure is documented below.

* `instance` - Indicates the detail of instance.

  The [instance](#instance_struct) structure is documented below.

<a name="lts_configs_struct"></a>
The `lts_configs` block supports:

* `log_type` - Indicates the type of log.

* `lts_group_id` - Indicates the log group ID associated with LTS.

* `lts_stream_id` - Indicates the log stream ID associated with LTS.

* `enabled` - Indicates whether logs of this type are uploaded to LTS.

<a name="instance_struct"></a>
The `instance` block supports:

* `id` - Indicates the ID of the RDS instance.

* `name` - Indicates the name of the RDS instance.

* `engine_name` - Indicates the engine name.

* `engine_version` - Indicates the version of the database engine.

* `engine_category` - Indicates the category of the engine.

* `status` - Indicates the current status of the instance.

* `enterprise_project_id` - Indicates the enterprise project ID associated with the instance.

* `actions` - Indicates the ongoing actions on the instance.
