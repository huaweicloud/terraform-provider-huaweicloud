---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_other_resource_uniagent_sync"
description: |-
  Manages a COC other resource uniagent sync resource within HuaweiCloud.
---

# huaweicloud_coc_other_resource_uniagent_sync

Manages a COC other resource uniagent sync resource within HuaweiCloud.

~> Deleting other resource uniagent sync resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "region_id" {}
variable "resource_id" {}

resource "huaweicloud_coc_other_resource_uniagent_sync" "test" {
  resource_infos {
    region_id   = var.region_id
    resource_id = var.resource_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_infos` - (Optional, List, NonUpdatable) Specifies the resource information corresponding to the resources
 that need to be synchronized.

  The [resource_infos](#resource_infos_struct) structure is documented below.

* `vendor` - (Optional, String, NonUpdatable) Specifies the cloud vendor.
  The value is empty when synchronizing Huawei Cloud resources. The value is **ALI** when synchronizing Alibaba Cloud
  resources.

<a name="resource_infos_struct"></a>
The `resource_infos` block supports:

* `region_id` - (Optional, String) Specifies the region ID to which the resource belongs.

* `resource_id` - (Optional, String) Specifies the resource ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
