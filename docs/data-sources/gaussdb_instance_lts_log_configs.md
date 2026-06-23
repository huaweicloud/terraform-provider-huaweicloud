---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_lts_log_configs"
description: |-
  Use this data source to query the LTS log configurations of GaussDB instances within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_lts_log_configs

Use this data source to query the LTS log configurations of GaussDB instances within HuaweiCloud.

## Example Usage

### Query All Instance LTS Log Configurations

```hcl
data "huaweicloud_gaussdb_instance_lts_log_configs" "test" {
}
```

### Query by Instance ID

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_instance_lts_log_configs" "test" {
  instance_id = instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the LTS log configurations.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the GaussDB instance.

* `instance_mode` - (Optional, String) Specifies the instance mode. The valid values are as follows:
  + **Ha**: Centralized type.
  + **Independent**: Independent deployment type.
  + **Combined**: Combined deployment type.

* `instance_name` - (Optional, String) Specifies the name of the GaussDB instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_lts_configs` - The list of instance LTS log configurations.
  The [instance_lts_configs](#instance_lts_configs) structure is documented below.

<a name="instance_lts_configs"></a>
The `instance_lts_configs` block supports:

* `instance` - The instance basic information.
  The [instance](#instance_lts_configs_instance) structure is documented below.

* `lts_configs` - The LTS log configuration list.
  The [lts_configs](#instance_lts_configs_lts_configs) structure is documented below.

<a name="instance_lts_configs_instance"></a>
The `instance` block supports:

* `id` - The instance ID.

* `name` - The instance name.

* `mode` - The instance mode.

* `status` - The instance status.

* `datastore` - The database information.
  The [datastore](#instance_datastore) structure is documented below.

* `frozen_flag` - Whether the instance is frozen.

* `actions` - Ongoing actions of the instance.

* `enterprise_project_id` - The enterprise project ID.

<a name="instance_datastore"></a>
The `datastore` block supports:

* `type` - The database engine type.

* `version` - The database engine version.

<a name="instance_lts_configs_lts_configs"></a>
The `lts_configs` block supports:

* `log_type` - The log type. The value is **audit_log**.

* `lts_group_id` - The LTS log group ID.

* `lts_stream_id` - The LTS log stream ID.

* `enabled` - Whether the LTS log is enabled.
