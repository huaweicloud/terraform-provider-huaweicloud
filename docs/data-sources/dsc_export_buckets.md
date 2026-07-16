---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_export_buckets"
description: |-
  Use this data source to get the list of DSC export buckets within HuaweiCloud.
---

# huaweicloud_dsc_export_buckets

Use this data source to get the list of DSC export buckets within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_export_buckets" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `buckets` - The export bucket information list.

  The [buckets](#buckets_struct) structure is documented below.

<a name="buckets_struct"></a>
The `buckets` block supports:

* `asset_name` - The asset name of the bucket.

* `bind_task` - The number of tasks bound to the bucket.

* `bucket_location` - The location of the bucket.

* `bucket_name` - The name of the bucket.

* `bucket_policy` - The access control policy of the bucket.

* `create_time` - The creation timestamp of the bucket.

* `enable_audit_status` - Whether the audit function is enabled for the bucket.

* `id` - The unique ID of the bucket.

* `is_deleted` - Whether the bucket is logically deleted.
