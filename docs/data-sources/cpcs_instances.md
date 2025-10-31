---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_instances"
description: |-
  Use this data source to get the list of CPCS instances.
---

# huaweicloud_cpcs_instances

Use this data source to get the list of CPCS instances.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
data "huaweicloud_cpcs_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `name` - (Optional, String) Specifies the instance name.

* `is_normal` - (Optional, String) Specifies the instance health status.
  The valid values are as follows:
  + **true**: Indicates normal.
  + **false**: Indicates abnormal.

* `sort_key` - (Optional, String) Specifies the sort attribute.
  The default value is **create_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  The default value is **DESC**. Valid values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `result` - Indicates the applications list.
  The [result](#instances_result_struct) structure is documented below.

<a name="instances_result_struct"></a>
The `result` block supports:

* `instance_id` - The instance ID.

* `resource_id` - The CBC resource ID.

* `instance_name` - The instance name.

* `service_type` - The service type to which the instance belongs.

* `cluster_id` - The cluster ID to which the instance belongs.

* `is_normal` - The instance health status.

* `status` - The instance service status.

* `image_name` - The instance image name.

* `specification` - The instance specification.

* `az` - The vailability zone.

* `expired_time` - The expire time.

* `create_time` - The creation time of the instance, in UNIX timestamp format.
