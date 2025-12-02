---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_registry_images"
description: |-
  Use this data source to get the image registry images of HSS within HuaweiCloud.
---

# huaweicloud_hss_image_registry_images

Use this data source to get the image registry images of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_registry_images" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter only needs to be configured after the Enterprise Project feature is enabled.
  For enterprise users, if omitted, default enterprise project will be used.
  Value **0** means default enterprise project.
  Value **all_granted_eps** means all enterprise projects to which the user has been granted access.

* `namespace` - (Optional, String) Specifies the organization name.

* `image_name` - (Optional, String) Specifies the image name.

* `image_version` - (Optional, String) Specifies the image tag.

* `registry_name` - (Optional, String) Specifies the repository name.

* `image_type` - (Optional, String) Specifies the image type. Valid values are:
  + **private_image**: SWR private image
  + **shared_image**: SWR shared image
  + **instance_image**: SWR enterprise edition image
  + **harbor**: Harbor repository image
  + **jfrog**: JFrog image

* `sort_key` - (Optional, String) Specifies the sorting field. Valid values are:
  + **latest_scan_time**: latest scan time

* `sort_dir` - (Optional, String) Specifies the sorting order. Valid values are:
  + **asc**: ascending
  + **desc**: descending

* `latest_version` - (Optional, Bool) Specifies whether to focus only on the latest version image.

* `image_size` - (Optional, Int) Specifies the image size.

* `scan_status` - (Optional, String) Specifies the scan status. Valid values are:
  + **unscan**: not scanned
  + **success**: the scan is complete
  + **scanning**: scanning
  + **failed**: the scan failed
  + **waiting_for_scan**: waiting for scan

* `start_latest_update_time` - (Optional, Int) Specifies the creation start date, in ms.

* `end_latest_update_time` - (Optional, Int) Specifies the creation end date, in ms.

* `start_latest_scan_time` - (Optional, Int) Specifies the start time of latest scan completion, in ms.

* `end_latest_scan_time` - (Optional, Int) Specifies the end date of the latest scan completion time, in ms.

* `start_latest_sync_time` - (Optional, Int) Specifies the start time of the latest synchronization completion, in ms.

* `end_latest_sync_time` - (Optional, Int) Specifies the end time of the latest synchronization completion, in ms.

* `has_malicious_file` - (Optional, Bool) Specifies whether there are malicious files.

* `has_unsafe_setting` - (Optional, Bool) Specifies whether there are baseline check risks.

* `has_vul` - (Optional, Bool) Specifies whether there are software vulnerabilities.

* `risky` - (Optional, Bool) Specifies whether there are security risks.

* `instance_id` - (Optional, String) Specifies the enterprise repository instance ID. This parameter can be specified
  for the enterprise edition SWR.

* `instance_name` - (Optional, String) Specifies the enterprise image instance name. This parameter can be specified
  for the enterprise edition SWR.

* `is_multarch` - (Optional, Bool) Specifies whether it is a multi-architecture image.

* `severity_level` - (Optional, String) Specifies the image risk level, which is displayed after the image scan is complete.
  Valid values are:
  + **Security**: secure
  + **Low**: low-risk
  + **Medium**: medium-risk
  + **High**: high-risk

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The repository image list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `scan_failed_desc` - The failure cause of the scan. The scan failure cause codes and their description are as follows:
  + **unknown_error**：Unknown error.
  + **authentication_failed**：Authentication failed.
  + **download_failed**：The image download failed. Contact technical support.
  + **image_over_sized**：The image is too large and cannot be scanned. Reduce the image size.
  + **get_detail_info_not_found**：Image details failed to be obtained. The image is not found in the repository.
    Synchronize the latest images.
  + **image_layer_over_sized**：The image has too many layers and cannot be scanned. Reduce the image size.
  + **schema_v1_not_support**：Schema V1 images cannot be scanned. Upgrade to V2.
  + **access_swr_failed**：Failed to access SWR. Contact technical support.
  + **swr_authentication_error**：The SWR authorization is missing. Refer to the image authorization guide to configure
    permissions.
  + **failed_to_scan_vulnerability**：The vulnerability scan failed.
  + **failed_to_scan_file**：The file scan failed.
  + **failed_to_scan_software**：The software scan failed.
  + **failed_to_check_sensitive_information**：Failed to check sensitive information.
  + **failed_to_check_baseline**：Baseline check failed.
  + **failed_to_check_software_compliance**：The software compliance check failed.
  + **failed_to_query_basic_image_information**：Failed to query the basic image information.
  + **failed_to_check_build_cmd**：Failed to scan the image building instructions.
  + **response_timed_out**：The response timed out.
  + **database_error**：Database error.
  + **failed_to_send_the_scan_request**：Failed to send the scan request.

* `latest_scan_time` - The last scanned, in ms.

* `scan_failed_code` - The failure cause code of the scan. The scan failure cause codes and their description are as
  follows:
  + **unknown_error**: Unknown error.
  + **authentication_failed**: Authentication failed.
  + **download_failed**: The image download failed. Contact technical support.
  + **image_over_sized**: The image is too large and cannot be scanned. Reduce the image size.
  + **get_detail_info_not_found**: Image details failed to be obtained. The image is not found in the repository.
    Synchronize the latest images.
  + **image_layer_over_sized**: The image has too many layers and cannot be scanned. Reduce the image size.
  + **schema_v1_not_support**: Schema V1 images cannot be scanned. Upgrade to V2.
  + **access_swr_failed**: Failed to access SWR. Contact technical support.
  + **swr_authentication_error**: The SWR authorization is missing. Refer to the image authorization guide to configure
    permissions.
  + **failed_to_scan_vulnerability**: The vulnerability scan failed.
  + **failed_to_scan_file**: The file scan failed.
  + **failed_to_scan_software**: The software scan failed.
  + **failed_to_check_sensitive_information**: Failed to check sensitive information.
  + **failed_to_check_baseline**: Baseline check failed.
  + **failed_to_check_software_compliance**: The software compliance check failed.
  + **failed_to_query_basic_image_information**: Failed to query the basic image information.
  + **failed_to_check_build_cmd**: Failed to scan the image building instructions.
  + **response_timed_out**: The response timed out.
  + **database_error**: Database error.
  + **failed_to_send_the_scan_request**: Failed to send the scan request.

* `registry_type` - The image repository type. Valid values are:
  + **SwrPrivate**: SWR private repository.
  + **SwrShared**: SWR shared repository.
  + **SwrEnterprise**: SWR enterprise repository.
  + **Harbor**: Harbor repository.
  + **Jfrog**: JFrog repository.
  + **Other**: Other repository.

* `image_size` - The image size.

* `vul_num` - The number of vulnerabilities.

* `instance_id` - The enterprise image instance ID.

* `scan_status` - The scan status. Valid values are:
  + **unscan**: not scanned
  + **success**: the scan is complete
  + **scanning**: the scan is in progress
  + **failed**: the scan failed
  + **download_failed**: the download failed
  + **image_oversized**: the image is too large
  + **waiting_for_scan**: waiting for scan

* `latest_sync_time` - The last synchronization time, in ms.

* `scannable` - Whether scan or not.

* `namespace` - The organization name.

* `instance_name` - The enterprise image instance name.

* `id` - The ID.

* `image_version` - The image tag.

* `image_type` - The image type. Valid values are:
  + **private_image**: SWR private image
  + **shared_image**: SWR shared image
  + **instance_image**: SWR enterprise edition image
  + **harbor**: Harbor repository image
  + **jfrog**: JFrog image

* `unsafe_setting_num` - The number of settings that failed the baseline check.

* `shared_status` - The shared image status. Valid values are:
  + **expired**: expired
  + **effective**: valid

* `image_digest` - The image digest.

* `registry_name` - The image repository name.

* `latest_version` - The latest version.

* `malicious_file_num` - The number of malicious files.

* `is_multiple_arch` - Whether it is a multi-architecture image.

* `association_images` - The multi-architecture associated image information.

  The [association_images](#data_list_association_images_struct) structure is documented below.

* `severity_level` - The image risk level, which is displayed after the image scan is complete. Valid values are:
  + **Security**: secure
  + **Low**: low-risk
  + **Medium**: medium-risk
  + **High**: high-risk

* `image_name` - The image name.

* `image_id` - The image ID.

* `registry_id` - The image repository ID.

* `latest_update_time` - The last update time of the image tag, in ms.

* `domain_name` - The owner (shared image parameter).

* `instance_url` - The enterprise image instance URL.

<a name="data_list_association_images_struct"></a>
The `association_images` block supports:

* `image_name` - The image name.

* `image_type` - The image type. Valid values are:
  + **private_image**: SWR private image
  + **shared_image**: SWR shared image
  + **instance_image**: SWR enterprise edition image
  + **harbor**: Harbor repository image
  + **jfrog**: JFrog image

* `namespace` - The organization name.

* `vul_num` - The number of vulnerabilities.

* `scan_status` - The scan status. Valid values are:
  + **unscan**: not scanned
  + **success**: the scan is complete
  + **scanning**: the scan is in progress
  + **failed**: the scan failed
  + **download_failed**: the download failed
  + **image_oversized**: the image is too large
  + **waiting_for_scan**: waiting for scan

* `unsafe_setting_num` - The number of settings that failed the baseline check.

* `malicious_file_num` - The number of malicious files.

* `id` - The ID.

* `image_id` - The image ID.

* `image_version` - The image tag.

* `image_digest` - The image digest.
