---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_check_rules"
description: |-
  Use this data source to get the list of HSS baseline check rules within HuaweiCloud.
---

# huaweicloud_hss_baseline_check_rules

Use this data source to get the list of HSS baseline check rules within HuaweiCloud.

## Example Usage

```hcl
variable "type" {}
variable "image_type" {} 

data "huaweicloud_hss_baseline_check_rules" "test" {
  type       = var.type
  image_type = var.image_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the type of the baseline check.  
  The valid values are as follows:
  + **image**

* `image_type` - (Required, String) Specifies the type of the image.
  The valid values are as follows:
  + **cicd**: CICD image.
  + **registry**: Warehouse image.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `namespace` - (Optional, String) Specifies the namespace.

* `image_name` - (Optional, String) Specifies the name of the image.

* `image_version` - (Optional, String) Specifies the version of the image.

* `instance_id` - (Optional, String) Specifies the enterprise warehouse instance ID, SWR shared version does not require
  this parameter.

* `image_id` - (Optional, String) Specifies the image ID.

* `scan_result` - (Optional, String) Specifies the scan result.  
  The valid values are as follows:
  + **pass**: The check is passed.
  + **failed**: The check is not passed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of baseline check rules.

* `data_list` - The list of baseline check rules.
  
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `severity` - The risk level. Valid values:
  + **Security**
  + **Low**
  + **Medium**
  + **High**
  + **Critical**

* `check_name` - The baseline name.

* `check_type` - The baseline type.

* `standard` - The standard type. Valid values are:
  + **cn_standard**: Waiting for compliance standards.
  + **hw_standard**: Huawei Standard.
  + **qt_standard**: Qingteng Standard.

* `check_type_desc` - The baseline description.

* `check_rule_name` - The check rule name.

* `check_rule_id` - The check rule ID.

* `scan_result` - The scan result.

* `latest_scan_time` - The latest scan time, in Unix timestamp format (milliseconds).

* `image_num` - The number of affected images, the number of images that have undergone current baseline detection.
