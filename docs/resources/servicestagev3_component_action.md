---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_component_action"
description: |-
  Use this resource to operate ServiceStage component within HuaweiCloud.
---

# huaweicloud_servicestagev3_component_action

Use this resource to operate ServiceStage component within HuaweiCloud.

-> This resource is only a one-time action resource for operating component. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "application_id" {}
variable "component_id" {}

resource "huaweicloud_cae_component_action" "test" {
  application_id = var.application_id
  component_id   = var.component_id
  action         = "sync_workload"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the component to be operated is located.  
  If omitted, the provider-level region will be used.
  Changing this will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the application ID to which the componnet belongs.  
  Changing this will create a new resource.

* `component_id` - (Required, String, ForceNew) Specifies the ID of the component to be operated.  
  Changing this will create a new resource.

* `action` - (Required, String, ForceNew) Specifies the action type of the component execution.  
  The valid values are as follows:
  + **start**
  + **stop**
  + **restart**
  + **scale**
  + **rollback**
  + **rollback_current**
  + **continue_deploy**
  + **check_gray_release**
  + **modify_gray_rule**
  + **sync_workload**

  Changing this will create a new resource.

* `parameters` - (Optional, String) Specifies the metadata of this action request.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
