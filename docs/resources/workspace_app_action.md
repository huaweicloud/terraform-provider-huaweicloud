---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_action"
description: |-
  Use this resource to active or deactive the Workspace APP tenant within HuaweiCloud.
---

# huaweicloud_workspace_app_action

Use this resource to active or deactive the Workspace APP tenant within HuaweiCloud.

-> This resource is a one-time action resource for activating or deactivating the Workspace APP tenant.
  Deleting this resource will not clear the corresponding request record, but will only remove the resource information
  from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_workspace_app_action" "test" {
  service_status = "active"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `service_status` - (Required, String, NonUpdatable) Specifies the service status of the Workspace APP tenant.  
  The valid values are as follows:
  + **active**: Activate the tenant service.
  + **inactive**: Deactivate the tenant service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
