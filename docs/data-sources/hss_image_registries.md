---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_registries"
description: |-
  Use this data source to get the image registries of HSS within HuaweiCloud.
---

# huaweicloud_hss_image_registries

Use this data source to get the image registries of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_registries" "test" {}
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

* `registry_name` - (Optional, String) Specifies the repository name.

* `registry_id` - (Optional, String) Specifies the registry ID.

* `registry_type` - (Optional, String) Specifies the image repository type. If this parameter is not specified, all types
  are returned by default. To query multiple types, separate them with commas (,). Valid values are:
  + **Harbor**: Harbor
  + **Jfrog**: Jfrog
  + **SwrPrivate**: SWR private repository
  + **SwrShared**: SWR shared repository
  + **SwrEnterprise**: SWR enterprise repository
  + **Other**: Other repository

* `registry_addr` - (Optional, String) Specifies the image repository address.

* `status` - (Optional, String) Specifies the access status. Valid values are:
  + **success**: The access is successful.
  + **fail**: The access is abnormal.
  + **accessing**: Accessing.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The image repository list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `registry_name` - The registry name.

* `protocol` - The protocol type. Valid values are:
  + **http**: HTTP protocol.
  + **https**: HTTPS protocol.

* `registry_username` - The username for logging in to the image repository.

* `latest_sync_time` - The last synchronization time, in ms.

* `id` - The image repository ID.

* `namespace` - The image repository project, which is used to specify the image repository directory that the scan
  component is to be uploaded to. This value is returned when `get_scan_image_channel` is set to **Other**.

* `fail_reason` - The reason. Valid values are:
  + **CREATE_JOB_FAILED**: The cluster access status is abnormal. Check the cluster access status.
  + **REQUEST_API_ERROR**: The network is disconnected, and the image repository fails to be accessed. Check whether the
    container cluster can access the image repository or retry on the Third-party Image Repository page.
  + **SERVER_ERROR**: Internal system error. Try again later.

* `registry_type` - The image repository type. Valid values are:
  + **Harbor**: Harbor repository
  + **Jfrog**: JFrog repository
  + **SwrPrivate**: SWR private repository
  + **SwrShared**: SWR shared repository
  + **SwrEnterprise**: SWR enterprise repository
  + **Other**: Other repository

* `registry_addr` - The image repository address.

* `get_scan_image_channel` - The method of obtaining the scan component. Valid values are:
  + **Swr**: Access SWR to obtain the scan component.
  + **Other**: Manually upload the scan component to the jumper cluster.

* `images_num` - The number of images.

* `api_version` - The image repository API version. Valid values are:
  + **V1**: V1 version
  + **V2**: V2 version

* `connect_cluster_id` - The jumper cluster ID.

* `status` - The access status. Valid values are:
  + **success**: The access is successful.
  + **fail**: The access is abnormal.
  + **accessing**: Accessing.
