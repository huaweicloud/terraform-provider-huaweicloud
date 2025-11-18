---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_baseline_risk_configs"
description: |-
  Use this data source to get the list of HSS image baseline risk configs within HuaweiCloud.
---

# huaweicloud_hss_image_baseline_risk_configs

Use this data source to get the list of HSS image baseline risk configs within HuaweiCloud.

## Example Usage

```hcl
variable "image_type" {}

data "huaweicloud_hss_image_baseline_risk_configs" "test" {
  image_type = var.image_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `image_type` - (Required, String) Specifies the image type. Valid values are:
  + **private_image**: Private image repository.
  + **shared_image**: Shared image repository.
  + **local_image**: Local image.
  + **instance_image**: Enterprise image.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `namespace` - (Optional, String) Specifies the namespace.

* `image_name` - (Optional, String) Specifies the image name.

* `image_version` - (Optional, String) Specifies the image version.

* `image_id` - (Optional, String) Specifies the image ID.

* `check_name` - (Optional, String) Specifies the baseline name.

* `severity` - (Optional, String) Specifies the risk level. Valid values are:
  + **Security**: Security.
  + **Low**: Low risk.
  + **Medium**: Medium risk.
  + **High**: High risk.

* `standard` - (Optional, String) Specifies the standard type. Valid values are:
  + **cn_standard**: Compliance standard.
  + **hw_standard**: Cloud security practice standard.

* `instance_id` - (Optional, String) Specifies the enterprise repository instance ID, SWR shared version does not
  require the use of this parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The configure detection list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `severity` - The risk level.

* `check_name` - The baseline name.

* `check_type` - The baseline type.

* `standard` - The standard type.

* `check_rule_num` - The number of inspection items.

* `failed_rule_num` - The number of risk items.

* `check_type_desc` - The baseline description.
