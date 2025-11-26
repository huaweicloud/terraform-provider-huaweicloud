---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_artifact_delete"
description: |-
  Manages a SWR enterprise instance artifact delete resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_instance_artifact_delete

Manages a SWR enterprise instance artifact delete resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "repository_name" {}
variable "reference" {}

resource "huaweicloud_swr_enterprise_instance_artifact_delete" "test" {
  instance_id     = var.instance_id
  namespace_name  = var.namespace_name
  repository_name = var.repository_name
  reference       = var.reference
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String, NonUpdatable) Specifies the namespace name.

* `repository_name` - (Required, String, NonUpdatable) Specifies the repository name.

* `reference` - (Required, String, NonUpdatable) Specifies the artifact digest.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
