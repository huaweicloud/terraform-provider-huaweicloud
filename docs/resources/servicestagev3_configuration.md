---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_configuration"
description: |-
  Manages a configuration file resource within HuaweiCloud.
---

# huaweicloud_servicestagev3_configuration

Manages a configuration file resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}
variable "configuration_name" {}

resource "huaweicloud_servicestagev3_configuration" "test" {
  config_group_id = var.group_id
  name            = var.configuration_name
  type            = "properties"
  content         = <<EOF
spring.application.name=service
cloud.servicecomb.service.name=$${spring.application.name}
cloud.servicecomb.service.version=$${CAS_INSTANCE_VERSION}
cloud.servicecomb.service.application=$${CAS_APPLICATION_NAME}
cloud.servicecomb.discovery.address=$${PAAS_CSE_SC_ENDPOINT}
cloud.servicecomb.discovery.healthCheckInterval=10
cloud.servicecomb.discovery.pollInterval=15000
cloud.servicecomb.discovery.waitTimeForShutDownInMillis=15000
cloud.servicecomb.config.serverAddr=$${PAAS_CSE_CC_ENDPOINT}
cloud.servicecomb.config.serverType=kie
cloud.servicecomb.config.fileSource=governance.yaml,application.yaml
EOF
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `config_group_id` - (Required, String, NonUpdatable) Specifies the ID of the configuration group to which
  the configuration file belongs.

* `name` - (Required, String, NonUpdatable) Specifies the name of the configuration file.  
  The valid length is limited from `2` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.  
  The name must start with a letter and end with a letter or a digit.

* `content` - (Required, String) Specifies the content of the configuration file.  
  For details about the system variables that can be written into a configuration file, please refer to the [document](https://support.huaweicloud.com/intl/en-us/usermanual-servicestage/servicestage_03_0316.html#servicestage_03_0316__en-us_topic_0000001904007428_table62206277581).

* `type` - (Required, String) Specifies the type of the configuration file.  
  The valid values are as follows:
  + **yaml**
  + **properties**

* `sensitive` - (Optional, Bool) Specifies whether to enable data encryption.  
  Defaults to **false**.  
  Only container-deployed components support data encryption. If data encryption is enabled, the configuration file is
  mounted using secret. If data encryption is disabled, the configuration file is mounted using ConfigMap.
  
* `description` - (Optional, String) Specifies the description of the configuration file.  
  The maximum length is `128` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also configuration file ID.

* `components` - The list of the components associated with the configuration file.  
  The [components](#configuration_attr_components) structure is documented below.

* `version` - The version of the configuration file.

* `creator` - The creator of the configuration file.

* `created_at` - The creation time of the configuration file, in RFC3339 format.

* `updated_at` - The latest update time of the configuration file, in RFC3339 format.

<a name="configuration_attr_components"></a>
The `components` block supports:

* `environment_id` - The ID of the environment.

* `application_id` - The ID of the application.

* `component_id` - The ID of the component.

* `component_name` - The name of the component.

## Import

The resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_servicestagev3_configuration.test <id>
```
