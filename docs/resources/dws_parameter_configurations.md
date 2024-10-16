---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_parameter_configurations"
description: |-
  Use this resource to modify the parameter configurations of DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_parameter_configurations

Use this resource to modify the parameter configurations of DWS cluster within HuaweiCloud.

-> 1. Only one `huaweicloud_dws_parameter_configurations` resource can be created for a DWS cluster.
   <br>2. If the modified parameters require restarting the DWS cluster, the modified values ​​will be displayed in the
   corresponding tfstate file after the resource is created, and the cluster status will be pending restart, and some
   operation and maintenance operations will be disabled. The modified values ​​will not take effect until the cluster
   is restarted, and the cluster will be restored to an available state.
   <br>3. This resource is only used to modify the DWS cluster parameter configurations. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "parameter_name" {}
variable "parameter_value" {}

resource "huaweicloud_dws_parameter_configurations" "test" {
  cluster_id = var.dws_cluster_id

  configurations {
    name  = var.parameter_name
    type  = "cn"
    value = var.parameter_value
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID.
  Changing this creates a new resource.

* `configurations` - (Required, List) Specifies the list of the DWS cluster parameter configurations.
  The [configurations](#parameter_configurations) structure is documented below.

<a name="parameter_configurations"></a>
The `configurations` block supports:

* `name` - (Required, String) Specifies the name of the parameter.

* `type` - (Required, String) Specifies the type of the parameter.  
  The valid values are as follows:
  + **cn**
  + **dn**

* `value` - (Required, String) Specifies the value of the parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
