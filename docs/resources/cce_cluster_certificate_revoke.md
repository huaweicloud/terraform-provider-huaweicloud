---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_certificate_revoke"
description: |-
  Use this resource to revoke the certificate of a CCE cluster within HuaweiCloud.
---

# huaweicloud_cce_cluster_certificate_revoke

Use this resource to revoke the certificate of a CCE cluster within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}
variable "user_id" {}

resource "huaweicloud_cce_cluster_certificate_revoke" "test" {
  cluster_id = var.cluster_id
  user_id    = var.user_id
}
```

~> Deleting certificate revoke resource is not supported, it will only be removed from the state.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the node sync resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `user_id` - (Optional, String, NonUpdatable) Specifies the user ID.

* `agency_id` - (Optional, String, NonUpdatable) Specifies the agency ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `cluster_id`.
