---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_bigkey_analysis"
description: ""
---

# huaweicloud_dcs_bigkey_analysis

Manage a DCS big key analysis resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dcs_bigkey_analysis" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DCS instance.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID.

* `scan_type` - Indicates the mode of the big key analysis. The value can be:
  + **manual**: indicates manual analysis.
  + **auto**: indicates automatic analysis.

* `created_at` - Indicates the creation time of the big key analysis. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format.

* `started_at` - Indicates the time when the big key analysis started. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format.

* `finished_at` - Indicates the time when the big key analysis ended. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format.

* `num` - Indicates the number of the big key.

* `status` - Indicates the analysis status. The value can be:
  + **waiting**: The analysis is waiting to begin.
  + **running**: The big key analysis is in progress.
  + **success**: The big key analysis succeeded.
  + **failed**: The big key analysis failed.

* `keys` - Indicates the record of big key.
  The [keys](#dcs_big_keys) structure is documented below.

<a name="dcs_big_keys"></a>
The `keys` block supports:

* `name` - Indicates the name of the big key.

* `type` - Indicates the type of the big key. The value can be **string**, **list**, **set**, **zset**, **hash**.

* `shard` - Indicates the shard where the big key is located.
  This parameter is supported only when the instance type is cluster. The format is **ip:port**.

* `db` - Indicates the database where the big key is located.

* `size` - Indicates the size of the key value.

* `unit` - Indicates the unit of the big key. The value can be:
  + **count**: The number of keys.
  + **byte**: The size of key.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The big key analysis can be imported using `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dcs_bigkey_analysis.test <instance_id>/<id>
```
