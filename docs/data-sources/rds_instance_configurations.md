---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_instance_configurations"
description: |-
  Use this data source to retrieve the configuration parameters of a specific RDS instance in HuaweiCloud.
---

# huaweicloud_rds_instance_configurations

Use this data source to retrieve the configuration parameters of a specific RDS instance in HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_instance_configurations" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datastore_version_name` - Indicates the database engine version.

* `datastore_name` - Indicates the database engine type.

* `created` - Indicates the creation time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.
  T is the separator between the calendar and the hourly notation of time. Z indicates the time
  zone offset. For example, in the Beijing time zone, the time zone offset is shown as +0800.

* `updated` - Indicates the last update time of the configuration. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.
  T is the separator between the calendar and the hourly notation of time. Z indicates the time
  zone offset. For example, in the Beijing time zone, the time zone offset is shown as +0800.

* `configuration_parameters` - Indicates the list of configuration parameters for the RDS instance.

  The [configuration_parameters](#configuration_parameters_struct) structure is documented below.

<a name="configuration_parameters_struct"></a>
The `configuration_parameters` block supports:

* `name` - Indicates the name of the configuration parameter.

* `value` - Indicates the current value of the configuration parameter.

* `restart_required` - Indicates whether a restart is required for the parameter to take effect. Value can be:
  + **true**: Indicates that the instance must be restarted for the change to apply.
  + **false**: Indicates that the change takes effect immediately without restart.

* `readonly` - Indicates whether the parameter is read-only. Value can be:
  + **true**: Indicates the permission is read-only.
  + **false**: Indicates the permission is read/write.

* `value_range` - Indicates the allowed value range for the parameter.

* `type` - Indicates the data type of the configuration parameter.

* `description` - Indicates the description of the configuration parameter.
