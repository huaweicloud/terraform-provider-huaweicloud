---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_lts_configs"
description: |-
  Use this data source to query the LTS configurations of TaurusDB HTAP StarRocks instances within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_lts_configs

Use this data source to query the LTS configurations of TaurusDB HTAP StarRocks instances within HuaweiCloud.

## Example Usage

### Query All Instances LTS Configurations

```hcl
data "huaweicloud_taurusdb_htap_starrocks_lts_configs" "test" {}
```

### Query by Instance ID

```hcl
variable "instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_lts_configs" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the LTS configurations.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the HTAP StarRocks instance.

* `instance_name` - (Optional, String) Specifies the name of the HTAP StarRocks instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_lts_configs` - The list of instance LTS configurations.
  The [instance_lts_configs](#instance_lts_configs) structure is documented below.

<a name="instance_lts_configs"></a>
The `instance_lts_configs` block supports:

* `instance` - The instance basic information.
  The [instance](#instance_lts_configs_instance) structure is documented below.

* `lts_configs` - The LTS configuration list.
  The [lts_configs](#instance_lts_configs_lts_configs) structure is documented below.

<a name="instance_lts_configs_instance"></a>
The `instance` block supports:

* `id` - The instance ID.

* `name` - The instance name.

* `mode` - The instance type. The valid values are **Cluster** and **Single**.

* `engine_name` - The engine name.

* `engine_version` - The engine version.

* `status` - The instance status.

* `enterprise_project_id` - The enterprise project ID.

* `enterprise_project_name` - The enterprise project name.

<a name="instance_lts_configs_lts_configs"></a>
The `lts_configs` block supports:

* `log_type` - The log type. The valid values are **error_log** and **slow_log**.

* `lts_group_id` - The LTS log group ID.

* `lts_stream_id` - The LTS log stream ID.

* `enabled` - Whether the LTS configuration is enabled.
