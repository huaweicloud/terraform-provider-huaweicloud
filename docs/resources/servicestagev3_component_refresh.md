---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_component_refresh"
description: |-
  Use this resource to refresh ServiceStage component within HuaweiCloud.
---

# huaweicloud_servicestagev3_component_refresh

Use this resource to refresh ServiceStage component within HuaweiCloud.

-> This resource is only a one-time action resource for doing component refresh. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "application_id" {}
variable "component_id" {}

resource "huaweicloud_servicestagev3_component_refresh" "test" {
  application_id = var.application_id
  component_id   = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the component to be refreshed is located.  
  If omitted, the provider-level region will be used.
  Changing this will create a new resource.

* `application_id` - (Required, String, NonUpdatable) Specifies the application ID to which the componnet belongs.  

* `component_id` - (Required, String, NonUpdatable) Specifies the ID of the component to be refreshed.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
