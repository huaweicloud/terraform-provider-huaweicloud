---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_certificate_rotatecredentials"
description: |-
  Use this resource to rotatecredentials the certificate of a CCE cluster within HuaweiCloud.
---

# huaweicloud_cce_cluster_certificate_rotatecredentials

Use this resource to rotatecredentials the certificate of a CCE cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "component" {}

resource "huaweicloud_cce_cluster_certificate_rotatecredentials" "test" {
  cluster_id = var.cluster_id
  component  = var.component
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `component` - (Required, String, NonUpdatable) Specifies the component.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
