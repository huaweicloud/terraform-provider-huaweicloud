---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_parameters"
description: |-
  Use this data source to get the list of the parameters under specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_parameters

Use this data source to get the list of the parameters under specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_cluster_parameters" "test" {
  cluster_id = var.dws_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `parameters` - The list of the parameters under specified DWS cluster.

  The [parameters](#cluster_parameters_struct) structure is documented below.

<a name="cluster_parameters_struct"></a>
The `parameters` block supports:

* `name` - The name of the parameter.

* `values` - The list of the parameter values.

  The [values](#parameters_values_struct) structure is documented below.

* `unit` - The unit of the parameter.

* `type` - The type of the parameter value.
  + **boolean**
  + **string**
  + **integer**
  + **float**
  + **list**

* `readonly` - Whether the parameter is read-only.

* `value_range` - The range of the parameter value.

* `restart_required` - Whether the DWS cluster needs to be restarted after modifying the parameter value.

* `description` - The description of the parameter.

<a name="parameters_values_struct"></a>
The `values` block supports:

* `type` - The type of the parameter.
  + **cn**
  + **dn**

* `value` - The value of the parameter.

* `default_value` - The default value of the parameter.
