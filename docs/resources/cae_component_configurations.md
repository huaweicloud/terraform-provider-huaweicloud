---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_component_configurations"
description: ""
---

# huaweicloud_cae_component_configurations

Using this resource to manage develop configurations for a component within HuaweiCloud.

-> A component can only have one resource.

## Example Usage

```hcl
variable "environment_id" {}
variable "application_id" {}
variable "component_id" {}

resource "huaweicloud_cae_component_configurations" "test" {
  environment_id = var.environment_id
  application_id = var.application_id
  component_id   = var.component_id

  items {
    type = "lifecycle"
    data = jsonencode({
      "spec": {
        "postStart": {
          "exec": {
            "command": [
              "/bin/bash",
              "-c",
              "sleep",
              "10",
              "done",
            ]
          }
        }
      }
    })
  }
  items {
    type = "env"
    data = jsonencode({
      "spec": {
        "envs": {
            "key": "value",
            "foo": "bar"
        }
      }
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the develop environment where the applications
  and components are located.  
  Changing this parameter will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the ID of the application where the components are
  located.  
  Changing this parameter will create a new resource.

* `component_id` - (Required, String, ForceNew) Specifies the ID of the component to which the configurations belong.  
  Changing this parameter will create a new resource.

* `items` - (Required, List, ForceNew) Specifies the list of configurations for component.  
  The [items](#component_configuration_items) structure is documented below.  
  Changing this parameter will create a new resource.

<a name="component_configuration_items"></a>
The `items` block supports:

* `type` - (Required, String, ForceNew) Specifies the type of the configuration.  
  Please following [reference documentation](https://support.huaweicloud.com/api-cae/CreateComponentConfiguration.html#CreateComponentConfiguration__request_ConfigurationItem).

* `data` - (Required, String, ForceNew) Specifies the configuration detail, in JSON format.  
  Please following [reference documentation](https://support.huaweicloud.com/api-cae/CreateComponentConfiguration.html#CreateComponentConfiguration__request_ConfigurationData).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also component ID), in UUID format.

## Import

The resource can be imported using `environment_id`, `application_id` and `component_id`, e.g.

```bash
$ terraform import huaweicloud_cae_component_configurations.test <environment_id>/<application_id>/<component_id>
```
