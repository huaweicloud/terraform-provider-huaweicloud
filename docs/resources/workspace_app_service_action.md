---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_service_action"
description: |-
  Use this resource to active or deactive the Workspace APP service within HuaweiCloud.
---

# huaweicloud_workspace_app_service_action

Use this resource to active or deactive the Workspace APP service within HuaweiCloud.

-> 1. Before using this resource, ensure that the Workspace service has been registered.
   <br>2. This resource is a one-time action resource for activating or deactivating the Workspace APP service. Deleting
   this resource will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_workspace_app_service_action" "test" {
  service_status = "active"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Workspace APP service is located.  
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `service_status` - (Required, String, NonUpdatable) Specifies the status of the Workspace APP service.  
  The valid values are as follows:
  + **active**: Activate the service.
  + **inactive**: Deactivate the service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
