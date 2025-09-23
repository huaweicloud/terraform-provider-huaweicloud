---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_runtime_stack"
description: |-
  Manages a custom runtime stack resource within HuaweiCloud.
---

# huaweicloud_servicestagev3_runtime_stack

Manages a custom runtime stack resource within HuaweiCloud.

## Example Usage

```hcl
variable "runtime_stack_name" {}
variable "java_sdk_obs_storage_url" {}

resource "huaweicloud_servicestagev3_runtime_stack" "test" {
  name        = var.runtime_stack_name
  deploy_mode = "virtualmachine"
  type        = "Java"
  version     = "1.0.1"
  spec        = jsonencode({
    "parameters": {
      "jdk_url": var.java_sdk_obs_storage_url
    }
  })
  description = "Created by terraform script"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the runtime stack.  
  The valid length is limited from `2` to `64`, only letters, digits, slash (/), hyphens (-) and underscores (_) are
  allowed. The name must start with a letter and end with a letter or a digit.

* `deploy_mode` - (Required, String, NonUpdatable) Specifies the deploy mode of the runtime stack.
  The valid value is **virtualmachine**.

* `type` - (Required, String, NonUpdatable) Specifies the type of the runtime stack.
  The valid values are as follows:
  + **Java**
  + **Tomcat**

* `version` - (Required, String) Specifies the version of the runtime stack.
  The valid format is **{number}.{number}.{number}**, e.g. **1.0.1**.

* `spec` - (Optional, String) Specifies the configuration of runtime stack, in JSON format.
  For the structure, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-servicestage/servicestage_06_0229.html#servicestage_06_0229__table51321252171513).

* `description` - (Optional, String) Specifies the description of the runtime stack.  
  The value can contain a maximum of `512` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `status` - The current status of the runtime stack.

* `creator` - The creator name of the runtime stack.

* `created_at` - The creation time of the runtime stack, in RFC3339 format.

* `updated_at` - The latest update time of the runtime stack, in RFC3339 format.

* `component_count` - The number of components associated with the runtime stack.

## Import

Runtime stacks can be imported using their `name` or `id`, e.g.

### Import a runtime stack using its name

```bash
$ terraform import huaweicloud_servicestagev3_runtime_stack.test <name>
```

### Import a runtime stack using its ID

```bash
$ terraform import huaweicloud_servicestagev3_runtime_stack.test <id>
```
