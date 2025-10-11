---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_app_events"
description: |-
  Use this data source to query the HSS app events within HuaweiCloud.
---

# huaweicloud_hss_app_events

Use this data source to query the HSS app events within HuaweiCloud.

## Example Usage

```hcl
variable "begin_time" {
  type = int
}

variable "end_time" {
  type = int
}

data "huaweicloud_hss_app_events" "test" {
  begin_time = var.begin_time
  end_time   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `begin_time` - (Required, Int) Specifies the start time of the custom query.

* `end_time` - (Required, Int) Specifies the end time of the custom query.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `last_days` - (Optional, Int) Specifies the number of days to be queried.

* `host_name` - (Optional, String) Specifies the server name.

* `host_ip` - (Optional, String) Specifies the server IP address.

* `handle_status` - (Optional, String) Specifies the event handling status. Valid values are **handled** and **unhandled**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - data list

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `event_id` - The event ID.

* `event_class_id` - The event type. Valid values:
  + **container_1001**: container namespace
  + **container_1002**: container open port
  + **container_1003**: container security option
  + **container_1004**: container mount directory
  + **containerescape_0001**: high-risk system call
  + **containerescape_0002**: shocker attack
  + **containerescape_0003**: Dirty Cow attack
  + **containerescape_0004**: container file escape
  + **dockerfile_001**: modification of user-defined protected container file
  + **dockerfile_002**: modification of executable files in the container file system
  + **dockerproc_001**: abnormal container process
  + **fileprotect_0001**: file privilege escalation
  + **fileprotect_0002**: critical file change
  + **fileprotect_0003**: critical file path change
  + **fileprotect_0004**: file/directory change
  + **av_1002**: virus
  + **av_1003**: worm
  + **av_1004**: Trojan
  + **av_1005**: botnet
  + **av_1006**: backdoor
  + **av_1007**: spyware
  + **av_1008**: adware
  + **av_1009**: phishing
  + **av_1010**: rootkit
  + **av_1011**: ransomware
  + **av_1012**: hacker tool
  + **av_1013**: grayware
  + **av_1015**: web shell
  + **av_1016**: mining software
  + **login_0001**: brute-force attack attempt
  + **login_0002**: successful brute-force attack
  + **login_1001**: successful login
  + **login_1002**: remote login
  + **login_1003**: weak password
  + **malware_0001**: shell change event
  + **malware_0002**: reverse shell event
  + **malware_1001**: malicious program
  + **procdet_0001**: abnormal process behavior
  + **procdet_0002**: process privilege escalation
  + **procreport_0001**: risky command
  + **user_1001**: account change
  + **user_1002**: risky account
  + **vmescape_0001**: VM sensitive command execution
  + **vmescape_0002**: access from virtualization process to sensitive file
  + **vmescape_0003**: abnormal VM port access
  + **webshell_0001**: web shell
  + **network_1001**: mining
  + **network_1002**: servers exploited to launch DDoS attacks
  + **network_1003**: malicious scan
  + **network_1004**: attack in sensitive areas
  + **ransomware_0001**: ransomware attack
  + **ransomware_0002**: ransomware attack
  + **ransomware_0003**: ransomware attack
  + **fileless_0001**: process injection
  + **fileless_0002**: dynamic library injection
  + **fileless_0003**: critical configuration change
  + **fileless_0004**: environment variable change
  + **fileless_0005**: memory file process
  + **fileless_0006**: VDSO hijacking
  + **crontab_1001**: suspicious crontab task
  + **vul_exploit_0001**: Redis vulnerability exploit
  + **vul_exploit_0002**: Hadoop vulnerability exploit
  + **vul_exploit_0003**: MySQL vulnerability exploit
  + **rootkit_0001**: suspicious rootkit file
  + **rootkit_0002**: suspicious kernel module
  + **RASP_0004**: web shell upload
  + **RASP_0018**: fileless web shell
  + **blockexec_001**: known ransomware attack
  + **hips_0001**: Windows Defender disabled
  + **hips_0002**: suspicious hacker tool
  + **hips_0003**: suspicious ransomware encryption behavior
  + **hips_0004**: hidden account creation
  + **hips_0005**: user password and credential reading
  + **hips_0006**: suspicious SAM file export
  + **hips_0007**: suspicious shadow copy deletion
  + **hips_0008**: backup file deletion
  + **hips_0009**: registry operation probably performed by ransomware
  + **hips_0010**: suspicious abnormal process
  + **hips_0011**: suspicious scan
  + **hips_0012**: suspicious ransomware script execution
  + **hips_0013**: suspicious mining command execution
  + **hips_0014**: suspicious Windows security center disabling
  + **hips_0015**: suspicious firewall disabling
  + **hips_0016**: suspicious disabling of system automatic recovery
  + **hips_0017**: executable file creation in Office
  + **hips_0018**: abnormal file creation with macros in Office
  + **hips_0019**: suspicious registry operation
  + **hips_0020**: Confluence remote code execution
  + **hips_0021**: MSDT remote code execution
  + **portscan_0001**: common port scan
  + **portscan_0002**: secret port scan
  + **k8s_1001**: Kubernetes event deletion
  + **k8s_1002**: privileged pod creation
  + **k8s_1003**: interactive shell used in pod
  + **k8s_1004**: pod created with sensitive directory
  + **k8s_1005**: pod created with server network
  + **k8s_1006**: pod created with host PID space
  + **k8s_1007**: authentication failure when common pods access API server
  + **k8s_1008**: API server access from common pod using cURL
  + **k8s_1009**: exec in system management space
  + **k8s_1010**: pod created in management space
  + **k8s_1011**: static pod creation
  + **k8s_1012**: DaemonSet creation
  + **k8s_1013**: scheduled cluster task creation
  + **k8s_1014**: operation on secrets
  + **k8s_1015**: allowed operation enumeration
  + **k8s_1016**: high privilege RoleBinding or ClusterRoleBinding
  + **k8s_1017**: ServiceAccount creation
  + **k8s_1018**: Cronjob creation
  + **k8s_1019**: interactive shell used for exec in pods
  + **k8s_1020**: unauthorized access to API server
  + **k8s_1021**: access to API server with curl
  + **k8s_1022**: Ingress vulnerability
  + **k8s_1023**: man-in-the-middle (MITM) attack
  + **k8s_1024**: worm, mining, or Trojan
  + **k8s_1025**: K8s event deletion
  + **k8s_1026**: SelfSubjectRulesReview
  + **imgblock_0001**: image blocking based on whitelist
  + **imgblock_0002**: image blocking based on blacklist
  + **imgblock_0003**: image tag blocking based on whitelist
  + **imgblock_0004**: image tag blocking based on blacklist
  + **imgblock_0005**: container creation blocked based on whitelist
  + **imgblock_0006**: container creation blocked based on blacklist
  + **imgblock_0007**: container mount proc
  + **imgblock_0008**: container seccomp unconfined
  + **imgblock_0009**: container privilege blocking
  + **imgblock_0010**: container capabilities blocking

* `event_type` - The event type. Valid values are:
  + `1001`: common malware
  + `1002`: virus
  + `1003`: worm
  + `1004`: Trojan
  + `1005`: botnet
  + `1006`: backdoor
  + `1010`: rootkit
  + `1011`: ransomware
  + `1012`: hacker tool
  + `1015`: web shell
  + `1016`: mining
  + `1017`: reverse shell
  + `2001`: common vulnerability exploit
  + `2012`: remote code execution
  + `2047`: Redis vulnerability exploit
  + `2048`: Hadoop vulnerability exploit
  + `2049`: MySQL vulnerability exploit
  + `3002`: file privilege escalation
  + `3003`: process privilege escalation
  + `3004`: critical file change
  + `3005`: file/directory change
  + `3007`: abnormal process behavior
  + `3015`: high-risk command execution
  + `3018`: abnormal shell
  + `3027`: suspicious crontab task
  + `3029`: system protection disabled
  + `3030`: backup deletion
  + `3031`: suspicious registry operations
  + `3036`: container image blocking
  + `4002`: brute-force attack
  + `4004`: abnormal login
  + `4006`: invalid accounts
  + `4014`: account added
  + `4020`: password theft
  + `6002`: port scan
  + `6003`: server scan
  + `13001`: Kubernetes event deletion
  + `13002`: abnormal pod behavior
  + `13003`: user information enumeration
  + `13004`: cluster role binding

* `event_name` - The event name.

* `severity` - The risk level. Valid values:
  + **Security**
  + **Low**
  + **Medium**
  + **High**
  + **Critical**

* `host_name` - The server name.

* `host_id` - The host ID.

* `private_ip` - The server private IP address.

* `public_ip` - The EIP.

* `attack_phase` - The attack phase. Valid values:
  + **reconnaissance**
  + **weaponization**
  + **delivery**
  + **exploit**
  + **installation**
  + **command_and_control**
  + **actions**

* `attack_tag` - The attack tag. Valid values:
  + **attack_success**: successful attack
  + **attack_attempt**: attack attempt
  + **attack_blocked**: attack blocked
  + **abnormal_behavior**: abnormal behavior
  + **collapsible_host**: server compromised
  + **system_vulnerability**: system vulnerability

* `occur_time` - The occurrence time, accurate to milliseconds.

* `handle_time` - The handling time, in milliseconds. This parameter is available only for handled alarms.

* `handle_status` - The handling status. Valid values are **unhandled** and **handled**.

* `handle_method` - The handling method. Valid values:
  + **mark_as_trust**
  + **mark_as_suspicious**
  + **isolate_and_kill**

* `operate_accept_list` - The supported processing operation.

* `operate_detail_list` - The operation details list (Not displayed on the page).
  The [operate_detail_list](#data_list_operate_detail_list_struct) structure is documented below.

* `resource_info` - The resource information (not displayed currently).
  The [resource_info](#data_list_resource_info_struct) structure is documented below.

* `recommendation` - The suggestion.

* `process_info` - The offset, process information list.
  The [process_info](#data_list_process_info_struct) structure is documented below.

* `policy_id` - The policy ID.

* `policy_name` - The policy name.

* `os_type` - The OS type. Valid values are **Linux** and **Windows**.

* `asset_value` - The asset importance. Its value can be **important**, **common**, or **test**.

* `host_status` - The server status. Valid values are:
  + **ACTIVE**: running
  + **SHUTOFF**: shut down
  + **BUILDING**: creating
  + **ERROR**: faulty

* `agent_status` - The agent status. Valid values are:
  + **installed**
  + **not_installed**
  + **online**
  + **offline**
  + **install_failed**
  + **installing**
  + **not_online**

* `protect_status` - The Protection status. Valid values are:
  + **closed**
  + **opened**

<a name="data_list_operate_detail_list_struct"></a>
The `operate_detail_list` block supports:

* `agent_id` - The agent ID.

* `process_pid` - The process ID.

* `file_hash` - The file hash.

* `file_path` - The file path.

* `file_attr` - The file attribute.

* `private_ip` - The server private IP address.

* `login_ip` - The login source IP address.

* `login_user_name` - The login username.

<a name="data_list_process_info_struct"></a>
The `process_info` block supports:

* `process_name` - The process name.

* `process_path` - The process path.

* `process_pid` - The process ID.

* `process_uid` - The process name.

* `process_username` - The process username.

* `process_cmdline` - The process command line.

* `process_filename` - The process file name.

* `process_start_time` - The process start time.

* `process_gid` - The process group ID.

* `process_egid` - The effective process group ID.

* `process_euid` - The effective process user ID.

* `parent_process_name` - The parent process name.

* `parent_process_path` - The parent process file path.

* `parent_process_pid` - The parent process ID.

* `parent_process_uid` - The user ID associated with the parent process.

* `parent_process_cmdline` - The parent process file command line.

* `parent_process_filename` - The parent process file name.

* `parent_process_start_time` - The parent process start time.

* `parent_process_gid` - The parent process group ID.

* `parent_process_egid` - The effective parent process group ID.

* `parent_process_euid` - The effective parent process user ID.

* `child_process_name` - The subprocess name.

* `child_process_path` - The subprocess file path.

* `child_process_pid` - The subprocess ID.

* `child_process_uid` - The user ID associated with the subprocess.

* `child_process_cmdline` - The subprocess file command line.

* `child_process_filename` - The subprocess file name.

* `child_process_start_time` - The subprocess start time.

* `child_process_gid` - The subprocess group ID.

* `child_process_egid` - The effective subprocess group ID.

* `child_process_euid` - The effective subprocess user ID.

* `virt_cmd` - The virtualization command.

* `virt_process_name` - The virtualization process name.

* `escape_mode` - The escape method.

* `escape_cmd` - The command executed after the escape.

* `process_hash` - The process startup file hash.

* `mode` - The file attribute.

* `rule` - The rule.

* `score` - The score.

* `process_file_hash` - The process file hash.

* `parent_process_file_hash` - The hash of the parent process file.

* `ancestor_process_pid` - The grandparent process ID.

* `ancestor_process_cmdline` - The grandparent process command line.

* `ancestor_process_path` - The grandparent process path.

* `operate_type` - The operation type.

* `session_id` - The session ID.

<a name="data_list_resource_info_struct"></a>
The `resource_info` block supports:

* `domain_id` - The tenant account ID.

* `project_id` - The project ID.

* `enterprise_project_id` - The ID of the enterprise project that the server belongs to.

* `region_name` - The region ID

* `vpc_id` - The VPC ID.

* `cloud_id` - The server ID.

* `vm_name` - The VM name.

* `vm_uuid` - The VM UUID.

* `container_id` - The container ID.

* `image_id` - The image ID.

* `image_name` - The image name. This parameter is available only for container alarms.

* `host_attr` - The server attribute.

* `service` - The business service.

* `micro_service` - The microservice.

* `sys_arch` - The system CPU architecture.

* `os_bit` - The OS bit version.

* `os_type` - The OS type. Valid values are **Linux** and **Windows**.

* `os_name` - The OS name.

* `host_name` - The server name.

* `host_ip` - The server IP address.

* `public_ip` - The EIP.

* `host_id` - The Host ID.

* `pod_uid` - The pod uid.

* `pod_name` - The pod name.

* `namespace` - The namespace.

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `asset_value` - The asset importance. Valid values are:
  + **important**:
  + **common**:
  + **test**:

* `container_status` - The container status.

* `os_version` - The OS version.

* `agent_version` - The agent version.
