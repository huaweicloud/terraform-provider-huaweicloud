---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_operations"
description: |-
  Use this data source to get the list of HSS event operations within HuaweiCloud.
---

# huaweicloud_hss_event_operations

Use this data source to get the list of HSS event operations within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_event_operations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `event_type` - (Optional, String) Specifies the type of the event.  
  The valid values are as follows:
  + **1001**: General Malware.
  + **1002**: Virus.
  + **1003**: Worm.
  + **1004**: Trojan.
  + **1005**: Botnet.
  + **1006**: Backdoor.
  + **1010**: Rootkit.
  + **1011**: Ransomware.
  + **1012**: Hacker Tool.
  + **1015**: Web Shell.
  + **1016**: Cryptomining.
  + **1017**: Reverse Shell.
  + **2001**: General Vulnerability Exploitation.
  + **2012**: Remote Code Execution.
  + **2047**: Redis Vulnerability Exploitation.
  + **2048**: Hadoop Vulnerability Exploitation.
  + **2049**: MySQL Vulnerability Exploitation.
  + **3002**: File Privilege Escalation.
  + **3003**: Process Privilege Escalation.
  + **3004**: Critical File Modification.
  + **3005**: File/Directory Modification.
  + **3007**: Abnormal Process Behavior.
  + **3015**: High-Risk Command Execution.
  + **3018**: Abnormal Shell.
  + **3027**: Suspicious Crontab Task.
  + **3029**: System Security Protection Disabled.
  + **3030**: Backup Deletion.
  + **3031**: Abnormal Registry Operation.
  + **3036**: Container Image Blocking.
  + **4002**: Brute-Force Attack.
  + **4004**: Abnormal Login.
  + **4006**: Unauthorized System Account.
  + **4014**: User Account Added.
  + **4020**: User Password Theft.
  + **6002**: Port Scanning.
  + **6003**: Host Scanning.
  + **13001**: Kubernetes Event Deletion.
  + **13002**: Abnormal Pod Behavior.
  + **13003**: User Information Enumeration.
  + **13004**: Cluster User Role Binding.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `operate_accept_list` - The list of supported processing operations.
