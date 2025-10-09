---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_page_notices"
description: |-
  Use this data source to get the list of HSS page notices within HuaweiCloud.
---

# huaweicloud_hss_page_notices

Use this data source to get the list of HSS page notices within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_page_notices" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `page_location` - (Optional, String) Specifies the page location. Valid values are:
  + **hostMgmt**: Host management - ECS.
  + **hostProtectQuota**: Host management - Protection quota.
  + **containerNodeList**: Container management - Container nodes.
  + **containerProtectQuota**: Container management - Container protection quota.
  + **containerMirror**: Container management - Container image.
  + **container**: Container management - Container.
  + **clusterAgent**: Container management - Cluster agent.
  + **vulView**: Vulnerability management - Vulnerability view.
  + **vulHostView**: Vulnerability management - Host view.
  + **ransomwareProtection**: Ransomware protection.
  + **policyMgmt**: Policy management.
  + **antiVirus**: Antivirus.
  + **hostAlarm**: Security alarm events - Host security alarms.
  + **containerAlarm**: Security alarm events - Container security alarms.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of page notice information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `page_location` - The page location.

* `type` - The notice type. Valid values are:
  + **links**: Hyperlink.
  + **text**: Text.

* `content` - The notice content.

* `title` - The notice title.

* `url` - The hyperlink.

* `level` - The notice level. Valid values are:
  + **error**: Emergency.
  + **warn**: Important.
  + **prompt**: Prompt.
