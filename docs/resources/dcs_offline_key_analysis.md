---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_offline_key_analysis"
description: |-
  Manages a DCS offline key analysis resource within HuaweiCloud.
---

# huaweicloud_dcs_offline_key_analysis

Manages a DCS offline key analysis resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_dcs_offline_key_analysis" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

### With Backup ID

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "backup_id" {}

resource "huaweicloud_dcs_offline_key_analysis" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  backup_id   = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the offline key analysis.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance node.

* `backup_id` - (Optional, String, NonUpdatable) Specifies the ID of the DCS instance backup.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the task ID of the offline key analysis.

* `group_name` - The name of the shard.

* `node_ip` - The IPv4 address of the node.

* `node_ipv6` - The IPv6 address of the node.

* `node_type` - The type of the node. The valid values are **master** and **slave**.

* `analysis_type` - The type of the analysis. The valid values are **new_backup** and **old_backup**.

* `started_at` - The start time of the analysis task.

* `finished_at` - The end time of the analysis task.

* `total_bytes` - The total size of keys, in bytes.

* `total_num` - The total number of keys.

* `largest_key_prefixes` - The list of prefix keys.
  The [largest_key_prefixes](#largest_key_prefixes_struct) structure is documented below.

* `largest_keys` - The big key list.
  The [largest_keys](#largest_keys_struct) structure is documented below.

* `type_bytes` - The total size of keys, in bytes.
  The [type_bytes](#type_bytes_struct) structure is documented below.

* `type_num` - The total number keys.
  The [type_num](#type_num_struct) structure is documented below.

<a name="largest_key_prefixes_struct"></a>
The `largest_key_prefixes` block supports:

* `key_prefix` - The prefix of the key.

* `type` - The type of the key. The valid values are **string**, **list**, **set**, **zset**, and **hash**.

* `bytes` - The size of the key, in bytes.

* `num` - The number of keys with the same prefix.

<a name="largest_keys_struct"></a>
The `largest_keys` block supports:

* `key` - The name of the key.

* `type` - The type of the key. The valid values are **string**, **list**, **set**, **zset**, and **hash**.

* `bytes` - The total size of keys, in bytes.

* `num_of_elem` - The element quantity (for the string type, in bytes) or size (for non-string type).

<a name="type_bytes_struct"></a>
The `type_bytes` block supports:

* `type` - The type of the key. The valid values are **string**, **list**, **set**, **zset**, and **hash**.

* `bytes` - The total size of each key type, in bytes.

<a name="type_num_struct"></a>
The `type_num` block supports:

* `type` - The type of the key. The valid values are **string**, **list**, **set**, **zset**, and **hash**.

* `num` - The total number of keys.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

DCS offline key analysis can be imported using the `instance_id` and the `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dcs_offline_key_analysis.test <instance_id>/<id>
```
