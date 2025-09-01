---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_gateway_jobs"
description: |-
  Use this data source to get the list of VPN gateway jobs.
---

# huaweicloud_vpn_gateway_jobs

Use this data source to get the list of VPN gateway jobs.

## Example Usage

```hcl
data "huaweicloud_vpn_gateway_jobs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `resource_id` - (Optional, String) Specifies the instance ID of a VPN gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the job list.
  The [jobs](#attrblock--jobs) structure is documented below.

<a name="attrblock--jobs"></a>
The `jobs` block supports:

* `id` - Indicates the job ID.

* `job_type` - Indicates the upgrade operation.

* `resource_id` - Indicates the instance ID of a VPN gateway.

* `status` - Indicates the job status.

* `sub_jobs` - Indicates the sub-job info.
  The [sub_jobs](#attrblock--jobs--sub_jobs) structure is documented below.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the update time.

<a name="attrblock--jobs--sub_jobs"></a>
The `sub_jobs` block supports:

* `id` - Indicates the job ID.

* `job_type` - Indicates the job type.

* `status` - Indicates the job status.

* `error_message` - Indicates error information.

* `created_at` - Indicates the creation time.

* `finished_at` - Indicates the end time.
