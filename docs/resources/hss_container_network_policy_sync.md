---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_network_policy_sync"
description: |-
  Manages an HSS container network policy sync resource within HuaweiCloud.
---

# huaweicloud_hss_container_network_policy_sync

Manages an HSS container network policy sync resource within HuaweiCloud.

-> This resource is only a one-time action resource for synchronizing container network policies. Deleting this resource
   will not affect the synchronization status, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "cluster_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_hss_container_network_policy_sync" "test" {
  cluster_id            = var.cluster_id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the HSS container network policy sync resource belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to synchronize network policies for.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  If omitted, the default enterprise project is used.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
