---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_configurations"
description: |-
  Use this data source to get the configurations of a CCE cluster within HuaweiCloud.
---

# huaweicloud_cce_cluster_configurations

Use this data source to get the configurations of a CCE cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_cluster_configurations" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE cluster configurations. If omitted, the
  provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID in which to query the configurations.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `configurations` - The map of configurations.
