---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_alarm_white_lists"
description: |-
  Use this data source to get the list of HSS alarm white lists within HuaweiCloud.
---
# huaweicloud_hss_event_alarm_white_lists

Use this data source to get the list of HSS alarm white lists within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_event_alarm_white_lists" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `hash` - (Optional, String) Specifies the event white list SHA256.

* `event_type` - (Optional, Int) Specifies the event type.  
  The valid values are as follows:
  + **1001**: General malware.
  + **1002**: Virus.
  + **1003**: Worm.
  + **1004**: Trojan.
  + **1005**: Botnet.
  + **1006**: Backdoor.
  + **1010**: Rootkit.
  + **1011**: Ransomware.
  + **1012**: Hacker tool.
  + **1015**: Webshell.
  + **1016**: Mining.
  + **1017**: Reverse Shell.
  + **2001**: General vulnerability exploitation.
  + **2012**: Remote code execution.
  + **2047**: Redis vulnerability exploitation.
  + **2048**: Hadoop vulnerability exploitation.
  + **2049**: MySQL vulnerability exploitation.
  + **3002**: File privilege escalation.
  + **3003**: Process privilege escalation.
  + **3004**: Key file change.
  + **3005**: File/directory change.
  + **3007**: Process abnormal behavior.
  + **3015**: High-risk command execution.
  + **3018**: Abnormal Shell.
  + **3027**: Crontab suspicious task.
  + **3029**: System security protection disabled.
  + **3030**: Backup deletion.
  + **3031**: Abnormal registry operation.
  + **3036**: Container image blocking.
  + **4002**: Brute force cracking.
  + **4004**: Abnormal login.
  + **4006**: Illegal system account.
  + **4014**: User account addition.
  + **4020**: User password theft.
  + **6002**: Port scanning.
  + **6003**: Host scanning.
  + **13001**: Kubernetes event deletion.
  + **13002**: Pod abnormal behavior.
  + **13003**: Enumerate user information.
  + **13004**: Bind cluster user role.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of white lists.

* `remain_num` - The number of white lists that can continue to be added.

* `limit_num` - The maximum number of white lists.

* `event_type_list` - The event types that support filtering.

* `data_list` - The list of white lists details.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `enterprise_project_name` - The enterprise project name.

* `hash` - The event white list SHA256.

* `description` - The description.

* `event_type` - The event type.

* `white_field` - The white field.  
  The valid values for this field, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-hss2.0/ListAlarmWhiteList.html).

* `field_value` - The field value.

* `judge_type` - The judge type.  
  The valid values for this field, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-hss2.0/ListAlarmWhiteList.html).

* `update_time` - The update time in milliseconds.
