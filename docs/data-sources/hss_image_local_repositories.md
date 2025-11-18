---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_local_repositories"
description: |-
  Use this data source to get the list of HSS local image repositories within HuaweiCloud.
---

# huaweicloud_hss_image_local_repositories

Use this data source to get the list of HSS local image repositories within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_local_repositories" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `image_name` - (Optional, String) Specifies the image name.

* `image_version` - (Optional, String) Specifies the image version.

* `scan_status` - (Optional, String) Specifies the scan status. Valid values are:
  + **unscan**: Not scanned.
  + **success**: Scan completed.
  + **scanning**: Scanning.
  + **failed**: Scan failed.
  + **waiting_for_scan**: Waiting for scan.

* `local_image_type` - (Optional, String) Specifies the image type. Valid values are:
  + **other_image**: Non-SWR image.
  + **swr_image**: SWR image.

* `image_size` - (Optional, Int) Specifies the image size in bytes.

* `start_latest_update_time` - (Optional, Int) Specifies the start time for the latest update time search, in
  milliseconds.

* `end_latest_update_time` - (Optional, Int) Specifies the end time for the latest update time search, in milliseconds.

* `start_latest_scan_time` - (Optional, Int) Specifies the start time for the latest scan completion time search, in
  milliseconds.

* `end_latest_scan_time` - (Optional, Int) Specifies the end time for the latest scan completion time search, in
  milliseconds.

* `has_vul` - (Optional, Bool) Specifies whether there are software vulnerabilities.

* `host_name` - (Optional, String) Specifies the name of the server associated with the local image.

* `host_id` - (Optional, String) Specifies the ID of the server associated with the local image.

* `host_ip` - (Optional, String) Specifies the IP address (public or private) of the server associated with the local
  image.

* `container_id` - (Optional, String) Specifies the ID of the container associated with the local image.

* `container_name` - (Optional, String) Specifies the name of the container associated with the local image.

* `pod_id` - (Optional, String) Specifies the ID of the Pod associated with the local image.

* `pod_name` - (Optional, String) Specifies the name of the Pod associated with the local image.

* `app_name` - (Optional, String) Specifies the name of the software associated with the local image.

* `has_container` - (Optional, Bool) Specifies whether there is a container. Valid values are:
  + **true**
  + **false**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of local images.

* `data_list` - The list of local image data.

  The [data_list](#image_data_list) structure is documented below.

<a name="image_data_list"></a>
The `data_list` block supports:

* `image_name` - The image name.

* `image_id` - The image ID.

* `image_digest` - The image digest.

* `image_version` - The image version.

* `local_image_type` - The local image type.

* `scan_status` - The scan status.

* `image_size` - The image size in bytes.

* `latest_update_time` - The last update time of the image version, in milliseconds.

* `latest_scan_time` - The last scan time, in milliseconds.

* `vul_num` - The number of vulnerabilities.

* `unsafe_setting_num` - The number of baseline scan failures.

* `malicious_file_num` - The number of malicious files.

* `host_num` - The number of associated hosts.

* `container_num` - The number of associated containers.

* `component_num` - The number of associated components.

* `scan_failed_desc` - The reason for scan failure. For details, please refer to the
  document link [reference](https://support.huaweicloud.com/intl/en-us/api-hss2.0/ListImageLocal.html).

* `severity_level` - The image risk level, display after image scanning is completed. Valid values are:
  + **Security**: Safety.
  + **Low**: Low risk.
  + **Medium**: Medium risk.
  + **High**: High-risk.

* `host_name` - The server name.

* `host_id` - The host ID.

* `agent_id` - The agent ID.

* `non_scan_reason` - The reason why the image cannot be scanned. If this field is empty, it means that the image can
  be scanned.
