---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_alarm_white_list_delete"
description: |-
  Manages an HSS event alarm white list deletion resource within HuaweiCloud.
---

# huaweicloud_hss_event_alarm_white_list_delete

Manages an HSS event alarm white list deletion resource within HuaweiCloud.

-> This resource is a one-time action resource using to delete HSS alarm white list. Deleting this resource will not
  clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

## Delete All alarm white lists

```hcl
variable "enterprise_project_id" {}
variable "restore_alarm" {}
variable "event_type" {}

resource "huaweicloud_hss_event_alarm_white_list_delete" "test" {
  enterprise_project_id = var.enterprise_project_id
  restore_alarm         = var.restore_alarm
  delete_all            = true
  event_type            = var.event_type
}
```

## Delete an alarm white list

```hcl
variable "enterprise_project_id" {}
variable "event_type" {}
variable "hash" {}
variable "description" {}

resource "huaweicloud_hss_event_alarm_white_list_delete" "test" {
  enterprise_project_id = var.enterprise_project_id

  data_list {
    event_type  = var.event_type
    hash        = var.hash
    description = var.description
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `restore_alarm` - (Optional, Bool, NonUpdatable) Specifies whether to restore related alarms.  
  The valid values are as follows:
  + **true**
  + **false**

  Defaults to **false**.

* `delete_all` - (Optional, Bool, NonUpdatable) Specifies whether to delete all white lists.  
  The valid values are as follows:
  + **true**
  + **false**

* `event_type` - (Optional, Int, NonUpdatable) Specifies the event type.  
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

* `data_list` - (Optional, List, NonUpdatable) Specifies the details of deleting the alarm white list.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `event_type` - (Required, Int, NonUpdatable) Specifies the event type.  
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

* `hash` - (Required, String, NonUpdatable) Specifies the event white list SHA256.

* `description` - (Required, String, NonUpdatable) Specifies the description.

* `delete_white_rule` - (Optional, Bool, NonUpdatable) Specifies whether to delete the alarm white list rule.
  This field is only used when the deleted white list is of rule type.

* `white_field` - (Optional, String, NonUpdatable) Specifies the white field.  
  The valid values for this field, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-hss2.0/ListAlarmWhiteList.html).

* `judge_type` - (Optional, String, NonUpdatable) Specifies the judge type.  
  The valid values for this field, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-hss2.0/ListAlarmWhiteList.html).

* `field_value` - (Optional, String, NonUpdatable) Specifies the white field value.

* `file_hash` - (Optional, String, NonUpdatable) Specifies the file hash.

* `file_path` - (Optional, String, NonUpdatable) Specifies the file path.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
