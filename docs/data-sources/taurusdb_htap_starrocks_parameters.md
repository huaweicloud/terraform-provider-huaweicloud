---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_parameters"
description: |-
  Use this data source to query the parameters of a TaurusDB HTAP StarRocks instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_parameters

Use this data source to query the parameters of a TaurusDB HTAP StarRocks instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_parameters" "test" {
  instance_id = var.htap_instance_id
  node_type   = "be"
}
```

### Query frontend node parameters

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_parameters" "test" {
  instance_id = var.htap_instance_id
  node_type   = "fe"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP StarRocks parameters.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the HTAP StarRocks instance ID.

* `node_type` - (Required, String) Specifies the node type.
  The valid values are as follows:
  + **be**: backend nodes
  + **fe**: frontend nodes

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The parameter configuration information.
  The [configurations](#htap_starrocks_parameters_configurations_attr) structure is documented below.

* `parameter_values` - The list of parameter objects.
  The [parameter_values](#htap_starrocks_parameters_parameter_values_attr) structure is documented below.

<a name="htap_starrocks_parameters_configurations_attr"></a>
The `configurations` block supports:

* `configuration_id` - The parameter template ID.

* `datastore_version_name` - The DB version name.

* `datastore_name` - The DB engine name in the parameter template.

* `created` - The time when the parameter template was created.

* `updated` - The time when the parameter template was updated.

<a name="htap_starrocks_parameters_parameter_values_attr"></a>
The `parameter_values` block supports:

* `name` - The parameter name.

* `value` - The parameter value.

* `restart_required` - Whether a reboot is required.

* `readonly` - Whether the parameter is read-only.

* `value_range` - The value range of the parameter.

* `type` - The parameter type.

* `description` - The parameter description.
