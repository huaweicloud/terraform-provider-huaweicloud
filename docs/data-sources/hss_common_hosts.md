---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_common_hosts"
description: |-
  Use this data source to get the common hosts list of HSS within HuaweiCloud.
---

# huaweicloud_hss_common_hosts

Use this data source to get the common hosts list of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_common_hosts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For querying assets under all enterprise projects, pass **all_granted_eps**.
  If omitted, the default enterprise project **0** will be used.

* `host_id` - (Optional, String) Specifies the server ID.

* `host_name` - (Optional, String) Specifies the server name.

* `private_ip` - (Optional, String) Specifies the server private IP address.

* `public_ip` - (Optional, String) Specifies the server elastic IP address.

* `feature_name` - (Optional, String) Specifies the policy name.
  The valid values are as follows:
  + **av_detect_feature**: AV policy.

* `group_id` - (Optional, String) Specifies the unique identifier ID of the server group.

* `asset_value` - (Optional, String) Specifies the asset importance.
  The valid values are as follows:
  + **important**: Important asset.
  + **common**: Common asset.
  + **test**: Test asset.

* `agent_status` - (Optional, String) Specifies the Agent status.
  The valid values are as follows:
  + **installed**: Installed.
  + **not_installed**: Not installed.
  + **online**: Online.
  + **offline**: Offline.
  + **install_failed**: Installation failed.
  + **installing**: Installing.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `cluster_name` - (Optional, String) Specifies the cluster name.

* `version_name_upper` - (Optional, String) Specifies the host version higher than this version.
  The valid values are as follows:
  + **hss.version.basic**: Basic version.
  + **hss.version.advanced**: Advanced version.
  + **hss.version.enterprise**: Enterprise version.
  + **hss.version.premium**: Premium version.
  + **hss.version.wtp**: Web page tamper protection version.
  + **hss.version.container.enterprise**: Container version.

* `version_name_lower` - (Optional, String) Specifies the host version lower than this version.
  The valid values are as follows:
  + **hss.version.basic**: Basic version.
  + **hss.version.advanced**: Advanced version.
  + **hss.version.enterprise**: Enterprise version.
  + **hss.version.premium**: Premium version.
  + **hss.version.wtp**: Web page tamper protection version.
  + **hss.version.container.enterprise**: Container version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The common hosts list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The unique identifier ID of the server.

* `host_name` - The server name.

* `public_ip` - The server elastic IP address.

* `private_ip` - The server private IP address.

* `agent_id` - The unique identifier ID of the antivirus Agent installed on the host,
  used to associate the host with the antivirus service.

* `os_type` - The operating system type.
  The valid values are as follows:
  + **Linux**: Linux.
  + **Windows**: Windows.

* `host_status` - The server status.
  The valid values are as follows:
  + **ACTIVE**: Running.
  + **SHUTOFF**: Shut down.
  + **BUILDING**: Creating.
  + **ERROR**: Error.

* `agent_status` - The Agent status.
  The valid values are as follows:
  + **installed**: Installed.
  + **not_installed**: Not installed.
  + **online**: Online.
  + **offline**: Offline.
  + **install_failed**: Installation failed.
  + **installing**: Installing.

* `os_name` - The operating system name.

* `os_version` - The operating system version.

* `asset_value` - The asset importance.
  The valid values are as follows:
  + **important**: Important asset.
  + **common**: Common asset.
  + **test**: Test asset.

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `group_id` - The server group ID.

* `group_name` - The server group name.
