---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_pay_per_scan_statistics"
description: |-
  Use this data source to get the HSS image pay per scan statistics within HuaweiCloud.
---

# huaweicloud_hss_image_pay_per_scan_statistics

Use this data source to get the HSS image pay per scan statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_pay_per_scan_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `repository_scan_num` - The number of repository image scans successfully.

* `cicd_scan_num` - The number of CICD image scans successfully.

* `free_quota_num` - The number of free scans.

* `already_given` - Whether free quota has been given.

* `image_free_quota` - The number of free quotas given when logging into the container image interface.

* `high_risk` - The high risk images.

  The [high_risk](#hss_imageType_riskInfo) structure is documented below.

* `has_risk` - The risky images.

  The [has_risk](#hss_imageType_riskInfo) structure is documented below.

* `total` - The total number of images.

  The [total](#hss_imageType_riskInfo) structure is documented below.

* `unscan` - The number of unscanned images.

  The [unscan](#hss_imageType_riskInfo) structure is documented below.

<a name="hss_imageType_riskInfo"></a>
The `high_risk`, `has_risk`, `total`, and `unscan` blocks support:

* `local` - The number of local images.

* `registriy` - The number of repository images.

* `cicd` - The number of CICD images.
