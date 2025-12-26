---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_configuration_details"
description: |-
  Use this data source to get the configuration details within HuaweiCloud.
---

# huaweicloud_cce_cluster_configuration_details

Use this data source to get the configurations of a CCE cluster within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cce_cluster_configuration_details" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE cluster configurations. If omitted, the
  provider-level region will be used.

* `cluster_type` - (Optional, String) Specifies the type of cluster.

* `cluster_version` - (Optional, String) Specifies the version of cluster.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `network_mode` - (Optional, String) Specifies the  network mode.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The data of configurations.
