---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_components"
description: |-
  Use this data source to query the list of components under specified application within HuaweiCloud.
---

# huaweicloud_servicestagev3_components

Use this data source to query the list of components under specified application within HuaweiCloud.

## Example Usage

### Query all components under specified region

```hcl
data "huaweicloud_servicestagev3_components" "test" {}
```

### Query all components under specified application via its ID

```hcl
variable "application_id" {}

data "huaweicloud_servicestagev3_components" "test" {
  application_id = var.application_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the components are located.  
  If omitted, the provider-level region will be used.

* `application_id` - (Optional, String) Specifies the ID of the application to which the components belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `components` - All components details.  
  The [components](#servicestage_v3_components) structure is documented below.

<a name="servicestage_v3_components"></a>
The `components` block supports:

* `id` - The component ID.

* `environment_id` - The environment ID where the component is deployed.

* `name` - The name of the component.

* `runtime_stack` - The configuration of the runtime stack.  
  The [runtime_stack](#servicestage_v3_components_runtime_stack) structure is documented below.

* `source` - The source configuration of the component, in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0077.html#servicestage_06_0077__en-us_topic_0220056058_ref28944532).

* `version` - The version of the component.

* `refer_resources` - The configuration of the reference resources.  
  The [refer_resources](#servicestage_v3_components_refer_resources) structure is documented below.

* `external_accesses` - The configuration of the external accesses.  
  The [external_accesses](#servicestage_v3_components_external_accesses) structure is documented below.

* `tags` - The key/value pairs to associate with the component.

* `status` - The status of the component.
  + **FAILED**
  + **RUNNING**
  + **DOWN**
  + **RESERVED**
  + **STOPPED**
  + **PENDING**
  + **UNKNOWN**
  + **PARTIALLY_FAILED**

<a name="servicestage_v3_components_runtime_stack"></a>
The `runtime_stack` block supports:

* `name` - The stack name.

* `type` - The stack type.
  + **Java**
  + **Tomcat**
  + **Nodejs**
  + **Php**
  + **Docker**
  + **Python**

* `deploy_mode` - The deploy mode of the stack.
  + **container**
  + **virtualmachine**

* `version` - The stack version.

<a name="servicestage_v3_components_refer_resources"></a>
The `refer_resources` block supports:

* `id` - The resource ID.

* `type` - The resource type.
  + **vpc**
  + **eip**
  + **elb**
  + **cce**
  + **ecs**
  + **as**
  + **cse**
  + **dcs**
  + **rds**

* `parameters` - The resource parameters, in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0076.html#servicestage_06_0076__table838321632514).

<a name="servicestage_v3_components_external_accesses"></a>
The `external_accesses` block supports:

* `protocol` - The protocol of the external access.
  + **http**
  + **https**

* `address` - The address of the external access.

* `forward_port` - The forward port of the external access.
