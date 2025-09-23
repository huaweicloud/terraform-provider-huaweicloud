---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_hda_batch_upgrade"
description: |-
  Use this resource to batch upgrade HDA versions of APP servers within HuaweiCloud.
---

# huaweicloud_workspace_app_hda_batch_upgrade

Use this resource to batch upgrade HDA versions of APP servers within HuaweiCloud.

-> This resource is only a one-time action resource for batch upgrading HDA versions. Deleting this resource will not
   clear the corresponding upgrade request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Upgrade Multiple Servers

```hcl
variable "server_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_app_hda_batch_upgrade" "test" {
  server_ids = var.server_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the servers to be upgraded are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `server_ids` - (Required, List, NonUpdatable) Specifies the list of server IDs to be upgraded HDA in batches.  
  The server SIDs must exist and be valid. This parameter cannot be updated after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
