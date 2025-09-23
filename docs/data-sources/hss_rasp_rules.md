---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_rasp_rules"
description: |-
  Use this data source to get the list of detection rules.
---

# huaweicloud_hss_rasp_rules

Use this data source to get the list of detection rules.

## Example Usage

```hcl
data "huaweicloud_hss_rasp_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `os_type` - (Optional, String) Specifies the OS type.
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of hosts affected by the vulnerability.
  The [data_list](#hosts_data_list) structure is documented below.

<a name="hosts_data_list"></a>
The `data_list` block supports:

* `chk_feature_id` - The detection feature rule ID.

* `chk_feature_name` - The detection feature rule name.

* `chk_feature_desc` - The detection feature rule description.

* `os_type` - The OS type.

* `feature_configure` - The detection feature rule configuration information.

* `optional_protective_action` - The available protection action.
  The valid values are as follows:
  + `1`: Detection.
  + `2`: Detection and blocking/interception.
  + `3`: All.

* `protective_action` - The default protection action.
  The valid values are as follows:
  + `1`: Detection.
  + `2`: Detection and blocking/interception.

* `editable` - Whether the configuration information can be edited.
  The valid values are as follows:
  + `0`: No.
  + `1`: Yes.
