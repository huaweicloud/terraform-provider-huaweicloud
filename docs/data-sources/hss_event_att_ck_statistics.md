---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_att_ck_statistics"
description: |-
  Use this data source to query the statistics of ATT and CK phases.
---

# huaweicloud_hss_event_att_ck_statistics

Use this data source to query the statistics of ATT and CK phases.

## Example Usage

```hcl
variable "category" {}

data "huaweicloud_hss_event_att_ck_statistics" "test" {
  category = var.category
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the event category.
  The valid values are as follows:
  + **host**
  + **container**
  + **serverless**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `last_days` - (Optional, Int) Specifies the query time range of days.
  The valid value range from `1` to `30`.

* `host_name` - (Optional, String) Specifies the host name.

* `host_id` - (Optional, String) Specifies the host ID.

* `private_ip` - (Optional, String) Specifies the host private IP address.

* `public_ip` - (Optional, String) Specifies the host EIP.

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

* `handle_status` - (Optional, String) Specifies the handling status.
  The valid values are as follows:
  + **unhandled**
  + **handled**

* `severity` - (Optional, String) Specifies the threat level.
  The valid values are as follows:
  + **Security**
  + **Low**
  + **Medium**  
  + **High**
  + **Critical**

* `severity_list` - (Optional, List) Specifies the threat level.
  The valid values are as follows:
  + **Security**
  + **Low**
  + **Medium**  
  + **High**
  + **Critical**

* `attack_tag` - (Optional, String) Specifies the attack flag.
  The valid values are as follows:
  + **attack_success**: Indicates attack success.
  + **attack_attempt**: Indicates attack attempt.
  + **attack_blocked**: Indicated blocked attack.
  + **abnormal_behavior**: Indicates abnormal behavior.
  + **collapsible_host**: Indicates compromised host.
  + **system_vulnerability**: Indicates system vulnerability.

* `asset_value` - (Optional, String) Specifies the asset importance.
  The valid values are as follows:
  + **important**
  + **common**
  + **test**

* `tag_list` - (Optional, List) Specifies the event tags list.

* `event_name` - (Optional, String) Specifies the alarm name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The attack phase details.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `att_ck` - The attack phase name.

* `num` - The quantity.
