---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_handle_history"
description: |-
  Use this data source to get the list of historical alarm records.
---

# huaweicloud_hss_event_handle_history

Use this data source to get the list of historical alarm records.

## Example Usage

```hcl
data "huaweicloud_hss_event_handle_history" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `severity` - (Optional, String) Specifies the threat level.
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

* `event_class_ids` - (Optional, List) Specifies the event ID list.
  The valid values are as follows:
  + **container_1001**: Container namespace.
  + **container_1002**: Container port enabled.
  + **container_1003**: Container security options.
  + **container_1004**: Container mount directory.
  + **containerescape_0001**: High-risk system call.
  + **containerescape_0002**: Shocker attack.
  + **containerescape_0003**: DirtCow attack.
  + **containerescape_0004**: Container file escape.
  + **dockerfile_001**: Modification of user-defined protected container file.
  + **dockerfile_002**: Modification of executable files in the container file system.
  + **dockerproc_001**: Abnormal container process.
  + **fileprotect_0001**: File privilege escalation.
  + **fileprotect_0002**: Key file change.
  + **fileprotect_0003**: Key file path change.
  + **fileprotect_0004**: File/directory change.
  + **av_1002**: Virus.
  + **av_1003**: Worm.
  + **av_1004**: Trojan.
  + **av_1005**: Botnet.
  + **av_1006**: Bbackdoor.
  + **av_1007**: Spyware.
  + **av_1008**: Malicious adware.
  + **av_1009**: Phishing.
  + **av_1010**: Rootkit
  + **av_1011**: Ransomware.
  + **av_1012**: Hacker tool.
  + **av_1013**: Grayware.
  + **av_1015**: Webshell.
  + **av_1016**: Mining software.
  + **login_0001**: Brute-force cracking.
  + **login_0002**: Successful cracking.
  + **login_1001**: Successful login.
  + **login_1002**: Remote login.
  + **login_1003**: Weak password.
  + **malware_0001**: Shell change report.
  + **malware_0002**: Reverse shell report.
  + **malware_1001**: Malicious program.
  + **procdet_0001**: Abnormal process behavior detection.
  + **procdet_0002**: Process privilege escalation.
  + **crontab_0001**: Crontab script privilege escalation.
  + **crontab_0002**: Malicious path privilege escalation.
  + **procreport_0001**: Risky commands.
  + **user_1001**: Account change.
  + **user_1002**: Risky account.
  + **vmescape_0001**: VM sensitive command execution.
  + **vmescape_0002**: Access from virtualization process to sensitive file.
  + **vmescape_0003**: Abnormal VM port access.
  + **webshell_0001**: Website backdoor.
  + **network_1001**: Malicious mining.
  + **network_1002**: DDoS attacks.
  + **network_1003**: Malicious scan.
  + **network_1004**: Attack in sensitive areas.
  + **ransomware_0001**: Ransomware attack.
  + **ransomware_0002**: Ransomware attack.
  + **ransomware_0003**: Ransomware attack.
  + **fileless_0001**: Process injection.
  + **fileless_0002**: Dynamic library injection.
  + **fileless_0003**: Key configuration change.
  + **fileless_0004**: Environment variable change.
  + **fileless_0005**: Memory file process.
  + **fileless_0006**: VDSO hijacking
  + **crontab_1001**: Suspicious crontab task.
  + **vul_exploit_0001**: Redis vulnerability exploit.
  + **vul_exploit_0002**: Hadoop vulnerability exploit.
  + **vul_exploit_0003**: MySQL vulnerability exploit.
  + **rootkit_0001**: Suspicious rootkit file.
  + **rootkit_0002**: Suspicious kernel module.
  + **RASP_0004**: Webshell upload.
  + **RASP_0018**: Fileless webshell.
  + **blockexec_001**: Known ransomware attack.
  + **hips_0001**: Windows Defender disabled.
  + **hips_0002**: Suspicious hacker tool.
  + **hips_0003**: Suspicious ransomware encryption behavior.
  + **hips_0004**: Hidden account creation.
  + **hips_0005**: User password and credential reading.
  + **hips_0006**: Suspicious SAM file export.
  + **hips_0007**: Suspicious shadow copy deletion.
  + **hips_0008**: Backup file deletion.
  + **hips_0009**: Registry of suspicious ransomware.
  + **hips_0010**: Suspicious abnormal process.
  + **hips_0011**: Suspicious scan.
  + **hips_0012**: Suspicious ransomware script running.
  + **hips_0013**: Suspicious mining command execution.
  + **hips_0014**: Suspicious windows security center disabling.
  + **hips_0015**: Suspicious behavior of disabling the firewall service.
  + **hips_0016**: Suspicious system automatic recovery disabling.
  + **hips_0017**: Executable file execution in Office.
  + **hips_0018**: Abnormal file creation with macros in Office.
  + **hips_0019**: Suspicious registry operation.
  + **hips_0020**: Confluence remote code execution.
  + **hips_0021**: MSDT remote code execution.
  + **portscan_0001**: Common port scan.
  + **portscan_0002**: Secret port scan.
  + **k8s_1001**: Kubernetes event deletion.
  + **k8s_1002**: Privileged pod creations.
  + **k8s_1003**: Interactive shell used in pod.
  + **k8s_1004**: Pod created with sensitive directory.
  + **k8s_1005**: Pod created with server network.
  + **k8s_1006**: Pod created with host PID space.
  + **k8s_1007**: Authentication failure when common pods access API server.
  + **k8s_1008**: API server access from common pod using curl.
  + **k8s_1009**: Exec in system management space.
  + **k8s_1010**: Pod created in management space.
  + **k8s_1011**: Static pod creation.
  + **k8s_1012**: DaemonSet creation.
  + **k8s_1013**: Scheduled cluster task creation.
  + **k8s_1014**: Operation on secrets.
  + **k8s_1015**: Allowed operation enumeration.
  + **k8s_1016**: High privilege RoleBinding or ClusterRoleBinding.
  + **k8s_1017**: ServiceAccount creation.
  + **k8s_1018**: Cronjob creation.
  + **k8s_1019**: Interactive shell used for exec in pods.
  + **k8s_1020**: Unauthorized access to API server.
  + **k8s_1021**: Access to API server with curl.
  + **k8s_1022**: Ingress vulnerability.
  + **k8s_1023**: MITM attack.
  + **k8s_1024**: Worm or mining or Trojan.
  + **k8s_1025**: K8s event deletion.
  + **k8s_1026**: SelfSubjectRulesReview.
  + **imgblock_0001**: Image blocking based on whitelist.
  + **imgblock_0002**: Image blocking based on blacklist.
  + **imgblock_0003**: Image tag blocking based on whitelist.
  + **imgblock_0004**: Image tag blocking based on blacklist.
  + **imgblock_0005**: Container creation blocked based on whitelist.
  + **imgblock_0006**: Container creation blocked based on blacklist.
  + **imgblock_0007**: Container mount proc blocking.
  + **imgblock_0008**: Container seccomp unconfined blocking.
  + **imgblock_0009**: Container privilege blocking.
  + **imgblock_0010**: Container capabilities blocking.

* `event_name` - (Optional, String) Specifies the alarm name. Supports fuzzy match.

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

* `host_name` - (Optional, String) Specifies the host name.

* `handle_status` - (Optional, String) Specifies the handling status.
  The valid values are as follows:
  + **unhandled**
  + **handled**

* `host_ip` - (Optional, String) Specifies the host IP address.

* `public_ip` - (Optional, String) Specifies the host EIP.

* `private_ip` - (Optional, String) Specifies the host private IP address.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `sort_key` - (Optional, String) Specifies sorting field.
  The valid value is **handle_time**.

* `sort_dir` - (Optional, String) Specifies sorting order.
  The valid values are as follows:
  + **asc**: Ascending order.
  + **desc**: Descending order.

  If `sort_key` is not empty, the returned results are sorted in ascending or descending order by `sort_key`.
  The default order is descending.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The event handle history list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `event_type` - The event type.

* `host_name` - The host name.

* `event_abstract` - The event summary.

* `attack_tag` - The attack flag.

* `private_ip` - The host private IP address.

* `public_ip` - The host EIP.

* `asset_value` - The asset importance.

* `occur_time` - The event occurrence time.

* `handle_status` - The handling status.

* `notes` - The remarke.

* `event_class_id` - The event category.

* `event_name` - The event name.

* `handle_time` - The handling time.

* `operate_type` - The handling method.
  The valid values are as follows:
  + **mark_as_handled**
  + **ignore**
  + **add_to_alarm_whitelist**
  + **add_to_login_whitelist**
  + **isolate_and_kill**
  + **unhandle**
  + **do_not_ignore**
  + **remove_from_alarm_whitelist**
  + **remove_from_login_whitelist**
  + **do_not_isolate_or_kill**

* `severity` - The threat level.

* `user_name` - The user name.
