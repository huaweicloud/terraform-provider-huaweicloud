---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_intrusion_events"
description: |-
  Use this data source to get the list of intrusion events.
---

# huaweicloud_hss_event_intrusion_events

Use this data source to get the list of intrusion events.

## Example Usage

```hcl
variable "category" {}

data "huaweicloud_hss_event_intrusion_events" "test" {
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

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `last_days` - (Optional, Int) Specifies the query time range of days.
  The valid value range from `1` to `30`.
  This parameter is manually exclusive with `begin_time` and `end_time`.

* `host_name` - (Optional, String) Specifies the host name.

* `host_id` - (Optional, String) Specifies the host ID.

* `private_ip` - (Optional, String) Specifies the host private IP address.

* `public_ip` - (Optional, String) Specifies the host EIP.

* `container_name` - (Optional, String) Specifies the container instance name.

* `event_types` - (Optional, List) Specifies the event type.
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

* `begin_time` - (Optional, String) Specifies the query start time.

* `end_time` - (Optional, String) Specifies the query end time.

* `event_class_ids` - (Optional, List) Specifies the event flag.
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

* `att_ck` - (Optional, String) Specifies the ATT and CK attack phase.
  The valid values are as follows:
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

* `auto_block` - (Optional, String) Specifies whether to automatically block alarms.
  The valid values are as follows:
  + **true**
  + **false**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of events.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `event_id` - The event ID.

* `event_class_id` - The event class ID.

* `event_type` - The event type.

* `event_name` - The event name.

* `severity` - The threat level.

* `container_name` - The container instance name.

* `image_name` - The image name.

* `host_name` - The host name.

* `host_id` - The host ID.

* `private_ip` - The host private IP address.

* `public_ip` - The host EIP.

* `os_type` - The OS type.
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `host_status` - The host status.
  The valid values are as follows:
  + **ACTIVE**
  + **SHUTOFF**
  + **BUILDING**
  + **ERROR**

* `agent_status` - The agent status.
  The valid values are as follows:
  + **installed**
  + **not_installed**
  + **online**
  + **offline**
  + **install_failed**
  + **installing**

* `protect_status` - The protection status.
  The valid values are as follows:
  + **closed**
  + **opened**

* `asset_value` - The asset importance.

* `attack_phase` - The event attack phase.

* `attack_tag` - The attack flag.

* `occur_time` - The event occurrence time, in milliseconds.

* `handle_time` - The handling time, in milliseconds.

* `handle_status` - The handling status.

* `handle_method` - The handle method.
  The valid values are as follows:
  + **mark_as_handled**
  + **ignore**
  + **add_to_alarm_whitelist**
  + **add_to_login_whitelist**
  + **isolate_and_kill**

* `handler` - The remarks.

* `operate_accept_list` - The supported processing operation.

* `operate_detail_list` - The operation details list.

  The [operate_detail_list](#event_operate_detail_list_struct) structure is documented below.

* `forensic_info` - The attack information, in JSON format.

* `resource_info` - The resource information.

  The [resource_info](#event_resource_info_struct) structure is documented below.

* `geo_info` - The geographical location, in JSON format.

* `malware_info` - The malware information, in JSON format.

* `network_info` - The network information, in JSON format.

* `app_info` - The application information, in JSON format.

* `system_info` - The system information, in JSON format.

* `extend_info` - The extended event information, in JSON format.

* `recommendation` - The handle suggestion.

* `description` - The alarm description.

* `event_abstract` - The alarm summary.

* `process_info_list` - The processes information list.

  The [process_info_list](#event_process_info_list_struct) structure is documented below.

* `user_info_list` - The users information list.

  The [user_info_list](#event_user_info_list_struct) structure is documented below.

* `file_info_list` - The files information list.

  The [file_info_list](#event_file_info_list_struct) structure is documented below.

* `event_details` - The event brief description.

* `tag_list` - The event tags list

* `event_count` - The event occurrences.

* `operate_type` - The operation type.
  The valid values are as follows:
  + **add**
  + **delete**
  + **change_attribute**
  + **modify**
  + **move**

<a name="event_operate_detail_list_struct"></a>
The `operate_detail_list` block supports:

* `agent_id` - The agent ID.

* `process_pid` - The process ID.

* `is_parent` - Whether a process is a parent process.

* `file_hash` - The file hash.

* `file_path` - The file path.

* `file_attr` - The file attribute.

* `private_ip` - The host private IP address.

* `login_ip` - The login source IP address.

* `login_user_name` - The login user name.

* `keyword` - The alarm event keyword.

* `hash` - The alarm event hash.

<a name="event_resource_info_struct"></a>
The `resource_info` block supports:

* `domain_id` - The account ID.

* `project_id` - The project ID.

* `enterprise_project_id` - The enterprise project ID.

* `region_name` - The region name.

* `vpc_id` - The vpc ID.

* `cloud_id` - The cloud host ID.

* `vm_name` - The VM name.

* `vm_uuid` - The VM UUID.

* `container_id` - The container ID.

* `container_status` - The container status.

* `pod_uid` - The pod UID.

* `pod_name` - The pod name.

* `namespace` - The namespace.

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `image_id` - The image ID.

* `image_name` - The image name.

* `host_attr` - The host attribute.

* `service` - The business service.

* `micro_service` - The microservice.

* `sys_arch` - The system CPU architecture.

* `os_bit` - The OS bit version.

* `os_type` - The OS type.

* `os_name` - The OS name.

* `os_version` - The OS version.

<a name="event_process_info_list_struct"></a>
The `process_info_list` block supports:

* `process_name` - The process name.

* `process_path` - The process file path.

* `process_pid` - The process ID.

* `process_uid` - The process user ID.

* `process_username` - The process user name.

* `process_cmdline` - The process file command line.

* `process_filename` - The process file name.

* `process_start_time` - The process start time.

* `process_gid` - The process group ID.

* `process_egid` - The effective process group ID.

* `process_euid` - The effective process user ID.

* `ancestor_process_path` - The grandparent process path

* `ancestor_process_pid` - The grandparent process ID.

* `ancestor_process_cmdline` - The grandparent process file command line.

* `parent_process_name` - The parent process name.

* `parent_process_path` - The parent process file path.

* `parent_process_pid` - The parent process ID.

* `parent_process_uid` - The parent process user ID.

* `parent_process_cmdline` - The parent process file command line.

* `parent_process_filename` - The parent process file name.

* `parent_process_start_time` - The parent process start time.

* `parent_process_gid` - The parent process group ID.

* `parent_process_egid` - The parent effective process group ID.

* `parent_process_euid` - The parent effective process user ID.

* `child_process_name` - The subprocess name.

* `child_process_path` - The subprocess file path.

* `child_process_pid` - The subprocess ID.

* `child_process_uid` - The subprocess user ID.

* `child_process_cmdline` - The subprocess file command line.

* `child_process_filename` - The subprocess file name.

* `child_process_start_time` - The subprocess start time.

* `child_process_gid` - The subprocess group ID.

* `child_process_egid` - The subprocess effective group ID.

* `child_process_euid` - The subprocess effective user ID.

* `virt_cmd` - The virtualization command.

* `virt_process_name` - The virtualization process name.

* `escape_mode` - The espace method.

* `escape_cmd` - The command executed after the espace.

* `process_hash` - The process startup file hash.

* `process_file_hash` - The process file hash.

* `parent_process_file_hash` - The parent process file hash.

* `block` - Whether the blocking is successful.
  `1` indicates blocking is successed and `0` indicates blocking is failed.

<a name="event_user_info_list_struct"></a>
The `user_info_list` block supports:

* `user_id` - The user ID.

* `user_gid` - The user group ID.

* `user_name` - The user name.

* `user_group_name` - The user group name.

* `user_home_dir` - The user home directory.

* `login_ip` - The user login IP address.

* `service_type` - The service type.
  The valid values are as follows:
  + **system**
  + **mysql**
  + **redis**

* `service_port` - The login service port.

* `login_mode` - The login method.

* `login_last_time` - The user last login time.

* `login_fail_count` - The number of user filed login.

* `pwd_hash` - The password hash.

* `pwd_with_fuzzing` - The anonymized password.

* `pwd_used_days` - The password used days.

* `pwd_min_days` - The minimum password validity period.

* `pwd_max_days` - The maximum password validity period.

* `pwd_warn_left_days` - The advance warning of a password expiration days.

<a name="event_file_info_list_struct"></a>
The `file_info_list` block supports:

* `file_path` - The file path.

* `file_alias` - The file alias.

* `file_size` - The file size.

* `file_mtime` - The file last modified time.

* `file_atime` - The file last access time.

* `file_ctime` - The file status was changed last time.

* `file_hash` - The file hash.

* `file_md5` - The file MD5 value.

* `file_sha256` - The file SHA256 value.

* `file_type` - The file type.

* `file_content` - The file content.

* `file_attr` - The file attributes before and after the change.

* `file_operation` - The file operation type.

* `file_action` - The file action.

* `file_change_attr` - The file change attribute.

* `file_new_path` - The new file path.

* `file_desc` - The file description.

* `file_key_word` - The file keyword.

* `is_dir` - Whether it is a directory.

* `fd_info` - The file handle information.

* `fd_count` - The number of file handles.
