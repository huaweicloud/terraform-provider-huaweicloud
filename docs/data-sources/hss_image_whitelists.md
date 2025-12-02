---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_whitelists"
description: |-
  Use this data source to get the list of image whitelists.
---

# huaweicloud_hss_image_whitelists

Use this data source to get the list of image whitelists.

## Example Usage

```hcl
variable "image_type" {}
variable "whitelist_type" {}

data "huaweicloud_hss_image_whitelists" "test" {
  global_image_type = var.image_type
  type              = var.whitelist_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `global_image_type` - (Required, String) Specifies the image type.
  The valid values are as follows:
  + **local**
  + **registry**

* `type` - (Required, String) Specifies the whitelist type.
  The valid values are as follows:
  + **vulnerability**: Vulnerability whitelist.

* `vul_name` - (Optional, String) Specifies the vulnerability name.

* `vul_type` - (Optional, String) Specifies the vulnerability type.
  The valid values are as follows:
  + **linux_vul**: Linux vulnerability.
  + **app_vul**: Application vulnerability.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of image whitelists.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `id` - The whitelist ID.

* `vul_id` - The vulnerability ID.

* `vul_name` - The vulnerability name.

* `vul_type` - The vulnerability type.

* `cves` - The vulnerability CVE list.

  The [cve_list](#data_list_cves_struct) structure is documented below.

* `rule_type` - The whitelist rule type.
  The valid values are as follows:
  + **all_images**
  + **specific_image_types**
  + **specific_images**

* `query_info` - The image type information. Only has value when `rule_type` set to **specific_image_types**.

  The [query_info](#data_list_query_info_struct) structure is documented below.

* `image_info` - The images list information. Only has value when `rule_type` set to **specific_images**.

  The [image_info](#data_list_image_info_struct) structure is documented below.

* `description` - The whitelist description.

<a name="data_list_cves_struct"></a>
The `cve_list` block supports:

* `cve_id` - The CVE ID.

* `cvss` - The CVSS score.

<a name="data_list_query_info_struct"></a>
The `query_info` block supports:

* `image_type` - The image type.
  The valid values are as follows:
  + **private_image**: SWR private image repository.
  + **shared_image**: Shared image repository of SWR.
  + **instance_image**: SWR enterprise repository.
  + **harbor**: Harbor repository.
  + **jfrog**: JFrog repository.

<a name="data_list_image_info_struct"></a>
The `image_info` block supports:

* `id` - The repository image iD.

* `image_id` - The local image ID.

* `image_name` - The image name.
