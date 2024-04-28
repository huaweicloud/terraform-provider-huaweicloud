---
subcategory: "Graph Engine Service (GES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ges_backup"
description: ""
---

# huaweicloud_ges_backup

Manages a GES backup resource within HuaweiCloud.  

## Example Usage

```hcl
variable "graph_id" {}
  
resource "huaweicloud_ges_backup" "test" {
  graph_id = var.graph_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `graph_id` - (Required, String, ForceNew) The ID of the graph that created the backup.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `name` - The backup name.  

* `backup_method` - Backup method. The value can be **auto** or **manual**.  

* `status` - Backup status.  
 The value can be one of the following:
   + **backing_up**: indicates that a graph is being backed up.
   + **success**: indicates that a graph is successfully backed up.
   + **failed**: indicates that a graph fails to be backed up.

* `start_time` - Start time of a backup job.

* `end_time` - End time of a backup job.

* `size` - Backup file size (MB).

* `duration` - Backup duration (seconds).

* `encrypted` - Whether the data is encrypted.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

## Import

The GES backup can be imported using
`graph_id`, `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_ges_backup.test <graph_id>/<id>
```
