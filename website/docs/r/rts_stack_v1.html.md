---
layout: "huaweicloud"
page_title: "huaweicloud: huaweicloud_rts_stack_v1"
sidebar_current: "docs-huaweicloud-resource-rts-stack-v1"
description: |-
  Provides an Huawei cloud RTS Stack resource.
---

# huaweicloud_rts_stack_v1_

Provides an Huawei Cloud Stack resource.

## Example Usage

 ```hcl
 variable "name" { }
 variable "network_id" { }
 variable "instance_type" { }
variable "image_id" { }
  
resource "huaweicloud_rts_stack_v1" "stack" {
  name = "${var.name}"
    disable_rollback = true
    timeout_mins=60
    parameters = {
        "network_id" = "${var.network_id}"
        "instance_type" = "${var.instance_type}"
        "image_id" = "${var.image_id}"
      }
    template_body = <<STACK
    {
      "heat_template_version": "2016-04-08",
      "description": "Simple template to deploy",
      "parameters": {
          "image_id": {
              "type": "string",
              "description": "Image to be used for compute instance",
              "label": "Image ID"
          },
          "network_id": {
              "type": "string",
              "description": "The Network to be used",
              "label": "Network UUID"
          },
          "instance_type": {
              "type": "string",
              "description": "Type of instance (Flavor) to be used",
              "label": "Instance Type"
          }
      },
      "resources": {
          "my_instance": {
              "type": "OS::Nova::Server",
              "properties": {
                  "image": {
                      "get_param": "image_id"
                  },
                  "flavor": {
                      "get_param": "instance_type"
                  },
                  "networks": [{
                      "network": {
                          "get_param": "network_id"
                      }
                  }]
              }
          }
      },
      "outputs":  {
        "InstanceIP":{
          "description": "Instance IP",
          "value": {  "get_attr": ["my_instance", "first_address"]  }
        }
      }
  }
  STACK
}
 ```
## Argument Reference
The following arguments are supported:


* `stack_name` - (Required) Specifies the stack name. The value must meet the regular expression rule (^[a-zA-Z][a-zA-Z0-9_.-]{0,254}$). Changing this will create a new stack.

* `stack_id` - (Required) Specifies the stack UUID.

* `template` - (Optional) Specifies the template. The template content must use the json syntax.

* `environment` - (Optional) Specifies the environment information about the stack.

* `files` - (Optional) Specifies files used in the environment.

* `parameters` - (Optional) Specifies parameter information of the stack.

* `timeout_mins` - (Optional) Specifies the timeout duration.

* `template_url` - (Optional) Specifies the template URL.

* `disable_rollback` - (Optional) Specifies whether to perform a rollback if the creation fails.

## Attributes Reference
The following attributes are exported:

* `outputs` - A map of outputs from the stack.

* `capabilities` - List of stack capabilities for stack.

* `notification_topics` - List of notification topics for stack.

* `status` - Specifies the stack status.


## Import

RTS Stacks can be imported using the `name`, e.g.

```
$ terraform import huaweicloud_rts_stack_v1.stack rts-stack
```


<a id="timeouts"></a>
## Timeouts

`huaweicloud_rts_stack_v1` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for Creating Stacks
- `update` - (Default `30 minutes`) Used for Stack modifications
- `delete` - (Default `30 minutes`) Used for destroying stacks.
