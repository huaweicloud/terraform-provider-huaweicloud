---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_iac_files"
description: |-
  Use this data source to get the list of HSS container IAC files within HuaweiCloud.
---

# huaweicloud_hss_container_iac_files

Use this data source to get the list of HSS container IAC files within HuaweiCloud.

## Example Usage

```hcl
variable "scan_type" {}

data "huaweicloud_hss_container_iac_files" "filtered" {
  scan_type = var.scan_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `scan_type` - (Required, String) Specifies the scan type.  
  The valid values are as follows:
  + **manual_scan**: Manual scan.
  + **cicd_scan**: Cicd scan.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `file_id` - (Optional, String) Specifies the file ID.

* `file_name` - (Optional, String) Specifies the file name.

* `file_type` - (Optional, String) Specifies the file type.  
  The valid values are as follows:
  + **dockerfile**
  + **k8s_yaml**

* `risky` - (Optional, Bool) Specifies whether there are risks.  
  The valid values are as follows:
  + **true**: There are risks.
  + **false**: There are no risks.

* `cicd_id` - (Optional, String) Specifies the CI/CD ID.

* `cicd_name` - (Optional, String) Specifies the CI/CD name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of IAC file information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `file_id` - The file ID.

* `file_name` - The file name.

* `file_type` - The file type.

* `risky` - Whether there are risks.

* `risk_num` - The number of risk items.

* `first_scan_time` - The first scan time.

* `last_scan_time` - The last scan time.

* `upload_file_time` - The file upload time.

* `cicd_id` - The CI/CD ID.

* `cicd_name` - The CI/CD name.

* `scan_type` - The scan type.
