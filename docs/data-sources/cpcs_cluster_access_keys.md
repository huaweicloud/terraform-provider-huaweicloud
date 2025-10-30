---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_cluster_access_keys"
description: |-
  Use this data source to get the access keys of a CPCS cluster.
---

# huaweicloud_cpcs_cluster_access_keys

Use this data source to get the access keys of a CPCS cluster.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cpcs_cluster_access_keys" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CPCS cluster.

* `app_name` - (Optional, String) Specifies the name of the application to filter access keys.

* `sort_key` - (Optional, String) Specifies the key for sorting the access keys. Defaults to **create_time**.

* `sort_dir` - (Optional, String) Specifies the direction for sorting. Valid values are **ASC** and **DESC**.
  Defaults to **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `access_keys` - The list of access keys.
  The [access_keys](#access_keys_struct) structure is documented below.

<a name="access_keys_struct"></a>
The `access_keys` block supports:

* `access_key_id` - The ID of the access key.

* `status` - The status of the access key.

* `app_name` - The name of the application that the access key belongs to.

* `access_key` - The access key value.

* `key_name` - The name of the access key.

* `create_time` - The creation time of the access key, in milliseconds.
