---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_solution_template_setting"
description: |-
  Use this data source to get the number of replicas, shards, and nodes corresponding to a specified instance or deployment mode.
---

# huaweicloud_gaussdb_opengauss_solution_template_setting

Use this data source to get the number of replicas, shards, and nodes corresponding to a specified instance or
deployment mode.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_solution_template_setting" "test" {
  solution = "single"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `solution` - (Optional, String) Specifies the solution template name.
  Value options: **triset**, **single**.

* `instance_id` - (Optional, String) Specifies the GaussDB OpenGauss instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `shard_num` - Indicates the number of shards.

* `replica_num` - Indicates the number of replicas.

* `initial_node_num` - Indicates the number of initial nodes.
  If `solution` is set to **triset**, this parameter is returned. Otherwise, **null** is returned.
