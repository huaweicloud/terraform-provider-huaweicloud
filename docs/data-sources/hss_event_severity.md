---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_severity"
description: |-
  Use this data source to get the list of HSS event severity levels within HuaweiCloud.
---

# huaweicloud_hss_event_severity

Use this data source to get the list of HSS event severity levels within HuaweiCloud.

## Example Usage

```hcl
variable category {}

data "huaweicloud_hss_event_severity" "test" {
  category = var.category
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the event category.  
  Valid values are:
  + **host**: Host security event.
  + **container**: Container security event.
  + **serverless**: Serverless scenario security event.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `begin_time` - (Optional, Int) Specifies the start time, 13-digit timestamp. Must be less than or equal to `end_time`.
  If `end_time` is not passed, the current time will be queried by default.

* `end_time` - (Optional, Int) Specifies the end time, 13-digit timestamp. Must be greater than or equal to
  `begin_time`. If `begin_time` is not passed, the query will start from timestamp 0 by default.

* `last_days` - (Optional, Int) Specifies the number of days for query time range.
  Mutually exclusive with custom query time `begin_time` and `end_time`.

* `host_name` - (Optional, String) Specifies the server name.

* `host_id` - (Optional, String) Specifies the server ID.

* `private_ip` - (Optional, String) Specifies the server private IP.

* `public_ip` - (Optional, String) Specifies the server public IP.

* `container_name` - (Optional, String) Specifies the container instance name.

* `event_type` - (Optional, Int) Specifies the event type.  
  The valid values are as follows:
  + **1001**: Common malware.
  + **1002**: Virus.
  + **1003**: Worm.
  + **1004**: Trojan.
  + **1005**: Botnet.
  + **1006**: Backdoor.
  + **1010**: Rootkit.
  + **1011**: Ransomware.
  + **1012**: Hacker tool.
  + **1015**: Web shell.
  + **1016**: Mining.
  + **1017**: Reverse shell.
  + **2001**: Common vulnerability exploit.
  + **2012**: Remote code execution.
  + **2047**: Redis vulnerability exploit.
  + **2048**: Hadoop vulnerability exploit.
  + **2049**: MySQL vulnerability exploit.
  + **3002**: File privilege escalation.
  + **3003**: Process privilege escalation.
  + **3004**: Critical file change.
  + **3005**: File/directory change.
  + **3007**: Abnormal process behavior.
  + **3015**: High-risk command execution.
  + **3018**: Abnormal shell.
  + **3026**: Crontab privilege escalation.
  + **3027**: Suspicious crontab task.
  + **3029**: System protection disabled.
  + **3030**: Backup deletion.
  + **3031**: Suspicious registry operations.
  + **3036**: Container image blocking.
  + **4002**: Brute-force attack.
  + **4004**: Abnormal login.
  + **4006**: Invalid accounts.
  + **4014**: Account added.
  + **4020**: Password theft.
  + **6002**: Port scan.
  + **6003**: Server scan.
  + **13001**: Kubernetes event deletion.
  + **13002**: Abnormal pod behavior.
  + **13003**: Enumerating user information.
  + **13004**: Cluster role binding.

* `handle_status` - (Optional, String) Specifies the handle status.  
  Valid values are:
  + **unhandled**: Not handled.
  + **handled**: Handled.

* `severity` - (Optional, String) Specifies the threat level.  
  Valid values are:
  + **Security**: Security.
  + **Low**: Low risk.
  + **Medium**: Medium risk.
  + **High**: High risk.
  + **Critical**: Critical.

* `severity_list` - (Optional, List) Specifies the threat level list. Valid values are the same as `severity`.

* `attack_tag` - (Optional, String) Specifies the attack identifier.  
  Valid values are:
  + **attack_success**: Attack succeeded.
  + **attack_attempt**: Attack attempted.
  + **attack_blocked**: Attack blocked.
  + **abnormal_behavior**: Abnormal behavior.
  + **collapsible_host**: Host compromised.
  + **system_vulnerability**: System vulnerability.

* `asset_value` - (Optional, String) Specifies the asset importance.  
  Valid values are:
  + **important**: Important asset.
  + **common**: Common asset.
  + **test**: Test asset.

* `tag_list` - (Optional, List) Specifies the event tag list, for example: `["热点事件"]`.

* `att_ck` - (Optional, String) Specifies the ATT&CK attack level.  
  Valid values are:
  + **Reconnaissance**
  + **Initial Access**
  + **Execution**
  + **Persistence**
  + **Privilege Escalation**
  + **Defense Evasion**
  + **Credential Access**
  + **Command and Control**
  + **Impact**

* `event_name` - (Optional, String) Specifies the alarm name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number.

* `low_num` - The number of low risk.

* `medium_num` - The number of medium  risk.

* `high_num` - The number of high risk.

* `critical_num` - The number of critical risk.
