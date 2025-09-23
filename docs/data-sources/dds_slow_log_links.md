---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_slow_log_links"
description: |-
  Use this data source to get the slow log links of DDS instance.
---

# huaweicloud_dds_slow_log_links

Use this data source to get the slow log links of DDS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dds_slow_log_links" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the instance.

* `file_name_list` - (Optional, List) Specifies the names of the files.

* `node_id_list` - (Optional, List) Specifies the node IDs to which the files belong.
  Nodes that can be queried:
  + mongos, shard, and config nodes in a cluster.
  + All nodes in a replica set or single node instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `links` - Indicates the list of the slow logs.
  The [links](#attrblock--links) structure is documented below.

<a name="attrblock--links"></a>
The `links` block supports:

* `file_link` - Indicates the file link.

* `file_name` - Indicates the file name.

* `file_size` - Indicates the file size.

* `node_name` - Indicates the node name.

* `status` - Indicates the link status.

* `updated_at` - Indicates the update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `read` - Default is 10 minutes.
