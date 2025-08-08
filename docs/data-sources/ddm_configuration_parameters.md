---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_configuration_parameters"
description: |-
  Use this data source to obtain detailed information about a specified DDM parameter template.
---

# huaweicloud_ddm_configuration_parameters

Use this data source to obtain detailed information about a specified DDM parameter template.

## Example Usage

```hcl
variable "config_id" {}

data "huaweicloud_ddm_configuration_parameters" "test" {
  config_id = var.config_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `config_id` - (Required, String) Specifies the ID of parameter template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - Indicates the name of the parameter template.

* `description` - Indicates the description of the parameter template.

* `datastore_name` - Indicates the database type.

* `created` - Indicates the creation time, in the format **yyyy-MM-ddTHH:mm:ssZ**.

* `updated` - Indicates the update time, in the format **yyyy-MM-ddTHH:mm:ssZ**.

* `configuration_parameters` - Indicates the list of DDM configurations.
  The [ConfigurationParameter](#DdmConfigurationParameters_ConfigurationParameter) structure is documented below.

<a name="DdmConfigurationParameters_ConfigurationParameter"></a>
The `ConfigurationParameters` block supports:

* `name` - Indicates the name of the parameter.

* `value` - Indicates the parameter value.

* `restart_required` - Indicates whether a restart is required for this parameter. The value can be **false** (indicates
  that a restart is not required) or **true** (indicates that a restart is required)

* `readonly` - Indicates whether the parameter is read-only. The value can be **false** (indicates that the parameter
  is not read-only) or **true** (indicates that the parameter is read-only).

* `value_range` - Indicates the valid value range of parameter.

* `type` - Indicates the parameter type.

* `description` - Indicates the parameter description.
