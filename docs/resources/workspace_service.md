---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_service"
description: ""
---

# huaweicloud_workspace_service

Use this resource to register or unregister the Workspace service in HuaweiCloud.

-> **NOTE:** Only one resource can be created in a region.

## Example Usage

### Register the Workspace service and use local authentication

```hcl
variable "vpc_id" {}
variable "network_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = var.vpc_id
  network_ids = var.network_ids

  otp_config_info {
    enable       = true
    receive_mode = "VMFA"
    rule_type    = "ACCESS_MODE"
    rule         = "PRIVATE"
  }
}
```

### Register the Workspace service and connect to the AD domain

```hcl
variable "vpc_id" {}
variable "network_ids" {
  type = list(string)
}

variable "ad_domain_name" {}
variable "ad_server_admin_account" {}
variable "ad_server_admin_password" {}
variable "ad_master_domain_ip" {}
variable "ad_server_name" {}
variable "ad_master_dns_ip" {}

resource "huaweicloud_workspace_service" "test" {
  auth_type   = "LOCAL_AD"
  access_mode = "INTERNET"
  vpc_id      = var.vpc_id
  network_ids = var.network_ids

  ad_domain {
    name               = var.ad_domain_name
    admin_account      = var.ad_server_admin_account
    password           = var.ad_server_admin_password
    active_domain_ip   = var.ad_master_domain_ip
    active_domain_name = format("%s.%s", var.ad_server_name, var.ad_domain_name)
    active_dns_ip      = var.ad_master_dns_ip
  }

  otp_config_info {
    enable       = true
    receive_mode = "VMFA"
    rule_type    = "ACCESS_MODE"
    rule         = "PRIVATE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to register the Workspace service.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the service belongs.
  Changing this will create a new resource.

  -> The resources required by Workspace will be created in the selected VPC subnet. After the configuration is saved,
     the VPC cannot be modified.

* `network_ids` - (Required, List) The network ID list of subnets that the service have.
  The subnets corresponding to this parameter must be included in the VPC resource corresponding to `vpc_id`.
  These subnet segments cannot conflict with `172.16.0.0/12`.

  -> The subnet of first registry must be selected. The DNS server address of the selected subnet will be automatically
     changed. Do not manually change it. You are advised to select a dedicated Workspace subnet and ensure that the DHCP
     function of the subnet is enabled.

* `access_mode` - (Required, String) Specifies the access mode of Workspace service.
  The valid values are as follows:
  + **INTERNET**: internet access.
  + **DEDICATED**: dedicated line access.
  + **BOTH**: both internet access and dedicated access are supported.

* `auth_type` - (Optional, String, ForceNew) Specifies the authentication type of Workspace service.
  The valid values are as follows:
  + **LITE_AS**: Local authentication.
  + **LOCAL_AD**: Connect to AD domain.

  Defaults to **LITE_AS**. Changing this will create a new resource.

* `ad_domain` - (Optional, List) Specifies the configuration of AD domain.
  Required if `auth_type` is **LOCAL_AD**. Make sure that the selected VPC network and the network to which AD
  belongs can be connected. The [object](#service_domain) structure is documented below.

  -> If AD domain is enabled, you need to connect the cloud desktop and Windows AD network.  
     If Windows AD is deployed in the intranet of the customer data center, these
     [ports](#secgroup_rules_for_ad_domain_connection) need to be opened in the firewall.

* `enterprise_id` - (Optional, String) Specifies the enterprise ID.
  The enterprise ID is the unique identification in the Workspace service.
  If omitted, the system will automatically generate an enterprise ID.
  The ID can contain `1` to `32` characters, only letters, digits, hyphens (-) and underscores (_) are allowed.

* `internet_access_port` - (Optional, Int) Specifies the internet access port.
  The valid value is range from `1,025` to `65,535`.
  
  -> If you want to modify the internet access port, please open a service ticket to enable this function.  
  
* `dedicated_subnets` - (Optional, List) The subnet segments of the dedicated access.

* `management_subnet_cidr` - (Optional, String, ForceNew) The subnet segment of the management component.

* `lock_enabled` - (Optional, Bool) Specifies whether to allow the provider to automatically unlock locked service
  when it is running. The default value is **false**.

* `otp_config_info` - (Optional, List) Specifies the configuration of auxiliary authentication.
  The [object](#config_info) structure is documented below.

<a name="service_domain"></a>
The `ad_domain` block supports:

* `name` - (Required, String) Specifies the domain name.
  The domain name must be an existing domain name on the AD server, and the length cannot exceed `55`.

* `admin_account` - (Required, String) Specifies the domain administrator account.
  It must be an existing domain administrator account on the AD server.

* `password` - (Required, String) Specifies the account password of domain administrator.

* `active_domain_ip` - (Required, String) Specifies the IP address of primary domain controller.

* `active_domain_name` - (Required, String) Specifies the name of primary domain controller.

* `standby_domain_ip` - (Optional, String) Specifies the IP address of the standby domain controller.

* `standby_domain_name` - (Optional, String) Specifies the name of the standby domain controller.

* `active_dns_ip` - (Optional, String) Specifies the primary DNS IP address.

* `standby_dns_ip` - (Optional, String) Specifies the standby DNS IP address.

* `delete_computer_object` - (Optional, Bool) Specifies whether to delete the corresponding computer object on AD
  while deleting the desktop.

<a name="config_info"></a>
The `otp_config_info` block supports:

* `enable` - (Required, Bool) Specifies whether to enable auxiliary authentication.

* `receive_mode` - (Required, String) Specifies the verification code receiving mode.
  + **VMFA**: Indicates virtual MFA device.
  + **HMFA**: Indicates hardware MFA device.
  
* `auth_url` - (Optional, String) Specifies the auxiliary authentication server address.

* `app_id` - (Optional, String) Specifies the auxiliary authentication server access account.

* `app_secret` - (Optional, String) Specifies the authentication service access password.

* `auth_server_access_mode` - (Optional, String) Specifies the authentication service access mode.
  + **INTERNET**: Indicates internet access.
  + **DEDICATED**: Indicates dedicated access.
  + **SYSTEM_DEFAULT**: Indicates system default.

* `cert_content` - (Optional, String) Specifies the PEM format certificate content.

* `rule_type` - (Optional, String) Specifies authentication application object type.
  + **ACCESS_MODE**: Indicates access type.

* `rule` - (Optional, String) Specifies authentication application object.
  + **INTERNET**: Indicates Internet access. Optional only when rule_type is **ACCESS_MODE**.
  + **PRIVATE**: Indicates dedicated line access. Optional only when rule_type is **ACCESS_MODE**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `internet_access_address` - The internet access address.
  This attribute is returned only when the access_mode is **INTERNET** or **BOTH**.

* `infrastructure_security_group` - The management component security group automatically created under the specified
  VPC after the service is registered. The [object](#service_security_group) structure is documented below.

* `desktop_security_group` - The desktop security group automatically created under the specified VPC after the service
  is registered. The [object](#service_security_group) structure is documented below.

* `status` - The current status of the Workspace service.

* `is_locked` - Whether the Workspace service is locked. The valid values are as follows:
  + **0**: Indicates not locked.
  + **1**: Indicates locked.

* `lock_time` - The time of the Workspace service is locked.

* `lock_reason` - The reason of the Workspace service is locked.

<a name="service_security_group"></a>
The `infrastructure_security_group` and `desktop_security_group` block supports:

* `id` - Security group ID.

* `name` - Security group name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

Service can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_service.test <id>
```

'NA' or other characters can be used to instead of the `id`.

```bash
$ terraform import huaweicloud_workspace_service.test NA
```

## Appendix

<a name="secgroup_rules_for_ad_domain_connection"></a>
If a firewall is deployed between Windows AD and the Workspace service, you need to open the following ports on the
firewall for the desktops of Workspace service to connect to Windows AD or DNS:

| Protocol | Ports | Usage |
| ---- | ---- | ---- |
| TCP | 135 | RPC protocol (required for LDAP, Distributed File System, and Distributed File Replication) |
| UDP | 137 | NetBIOS name resolution (required by the network login service) |
| UDP | 138 | NetBIOS datagram service (distributed file system, network login and other services need to use this port) |
| TCP | 139 | NetBIOS-SSN Service (Network Basic I/O Interface) |
| TCP | 445 | NetBIOS-SSN Service (Network Basic I/O Interface) |
| UDP | 445 | NetBIOS-SSN Service (Network Basic I/O Interface) |
| TCP | 49152-65535 | RPC dynamic ports (ports that are not hardened and open by AD. If AD is hardened, ports 50152-51151 need to be opened) |
| UDP | 49152-65535 | RPC dynamic ports (ports that are not hardened and open by AD. If AD is hardened, ports 50152-51151 need to be opened) |
| TCP | 88 | Kerberos Key Distribution Center Service |
| UDP | 88 | Kerberos Key Distribution Center Service |
| UDP | 123 | Port used by NTP service |
| TCP | 389 | LDAP server |
| UDP | 389 | LDAP server |
| TCP | 464 | Kerberos authentication protocol |
| UDP | 464 | Kerberos Authentication Protocol |
| UDP | 500 | isakmp |
| TCP | 593 | RPC over HTTP |
| TCP | 636 | LDAP SSL |
| TCP | 53 | DNS server |
| UDP | 53 | DNS server |
