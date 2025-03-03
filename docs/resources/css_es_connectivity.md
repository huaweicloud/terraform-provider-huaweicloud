---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_es_connectivity"
description: |-
  Manages CSS ElasticSearch cluster connectivity test resource within HuaweiCloud.
---

# huaweicloud_css_snapshot_restore

Manages CSS ElasticSearch cluster connectivity test resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource is only removed from the state.

## Example Usage

```hcl
variable "source_cluster_id" {}
variable "target_cluster_id" {}

resource "huaweicloud_css_es_connectivity" "test" {
  source_cluster_id = var.target_cluster_id
  target_cluster_id = var.source_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `source_cluster_id` - (Required, String, NonUpdatable) Specifies the source cluster ID.

* `target_cluster_id` - (Required, String, NonUpdatable) Specifies the target cluster ID.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID.
