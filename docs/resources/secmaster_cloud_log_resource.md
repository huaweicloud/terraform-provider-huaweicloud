---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_cloud_log_resource"
description: |-
  Manages a resource to create cloud log resource within HuaweiCloud.
---

# huaweicloud_secmaster_cloud_log_resource

Manages a resource to create cloud log resource within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "domain_id" {}
variable "resource_list" {
  type = list(object({
    enable    = string  
    region_id = string
  }))
}

resource "huaweicloud_secmaster_cloud_log_resource" "test" {
  workspace_id = var.workspace_id
  domain_id    = var.domain_id

    dynamic "resources" {
    for_each = var.resource_list
  
    content {
      enable    = resource_list.value.enable
      region_id = resource_list.value.region_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `domain_id` - (Required, String, NonUpdatable) Specifies the account ID.

* `resources` - (Required, List, NonUpdatable) Specifies the resource data.
  The [resources](#cloud_log_resource_data) structure is documented below.

<a name="cloud_log_resource_data"></a>
The `resources` block supports:

* `enable` - (Required, String, NonUpdatable) Specifies Whether enabled.
  The valid values are as follows:
  + **active**
  + **inactive**

* `region_id` - (Required, String, NonUpdatable) Specifies the region ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
