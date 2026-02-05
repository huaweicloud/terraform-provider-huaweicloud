---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_vulnerabilities"
description: |-
  Use this data source to get the list of HSS image vulnerabilities of tenants within HuaweiCloud.
---

# huaweicloud_hss_image_vulnerabilities

Use this data source to get the list of HSS image vulnerabilities of tenants within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_vulnerabilities" "test" {}
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

* `repair_necessity` - (Optional, String) Specifies repair urgency.  
  The valid values are as follows:
  + **immediate_repair**: Need to repair as soon as possible.
  + **delay_repair**: Can be postponed for repair.
  + **not_needed_repair**: Not yet fixed.

* `vul_id` - (Optional, String) Specifies the vulnerability ID (supports fuzzy query).

* `type` - (Optional, String) Specifies the vulnerability type.  
  The valid values are as follows:
  + **linux_vul**: Linux vulnerability.
  + **app_vul**: Application vulnerability.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of image vulnerabilities.

* `data_list` - The records corresponding to vulnerabilities.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `vul_name` - The vulnerability name.

* `vul_id` - The vulnerability ID.

* `repair_necessity` - The repair urgency.

* `description` - The vulnerability description.

* `solution` - The solution.

* `url` - The URL link.

* `history_number` - The number of affected mirrors in history.

* `undeal_number` - The unprocessed image format.

* `data_list` - The CVE list.

  The [data_list](#sub_data_list_struct) structure is documented below.

<a name="sub_data_list_struct"></a>
The `data_list` block supports:

* `cve_id` - The CVE ID.

* `cvss_score` - The CVSS score.

* `publish_time` - The CVE announcement time, time unit: milliseconds (ms).

* `description` - The CVE description.
