---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_agent_versions"
description: |-
  Use this data source to get the list of HSS agent versions within HuaweiCloud.
---

# huaweicloud_hss_agent_versions

Use this data source to get the list of HSS agent versions within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_agent_versions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `os_type` - (Optional, String) Specifies the operating system type.
  Valid values are **Linux** and **Windows**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of agent versions.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `os_type` - The operating system type.

* `latest_version` - The latest version number.

* `version_list` - The version list.

  The [version_list](#version_list_struct) structure is documented below.

<a name="version_list_struct"></a>
The `version_list` block supports:

* `release_version` - The release version.

* `release_note` - The release note.

* `update_time` - The update time, in milliseconds.
