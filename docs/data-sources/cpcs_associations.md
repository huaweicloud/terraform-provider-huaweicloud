---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_associations"
description: |-
  Use this data source to get list of binding relationships between the password service clusters and the applications.
---

# huaweicloud_cpcs_associations

Use this data source to get list of binding relationships between the password service clusters and the applications.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
data "huaweicloud_cpcs_associations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `app_id` - (Optional, String) Specifies the application ID.

* `sort_key` - (Optional, String) Specifies the sort attribute.
  The default value is **create_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  The default value is **DESC**. Valid values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `result` - Indicates the applications list.
  The [result](#associations_result_struct) structure is documented below.

<a name="associations_result_struct"></a>
The `result` block supports:

* `id` - The association ID.

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `app_id` - The application ID.

* `app_name` - The application name.

* `vpc_name` - The VPC name to which the application belongs.

* `subnet_name` - The subnet name to which the application belongs.

* `cluster_server_type` - The service type to which the cluster belongs.

* `vpcep_address` - The access address.

* `update_time` - The latest update time of the association, in UNIX timestamp format.

* `create_time` - The creation time of the association, in UNIX timestamp format.
