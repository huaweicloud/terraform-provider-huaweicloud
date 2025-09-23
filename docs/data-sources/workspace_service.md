---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_service"
description: |-
  Use this data source to get the configuration of the Workspace service within HuaweiCloud.
---

# huaweicloud_workspace_service

Use this data source to get the configuration of the Workspace service within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_service" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The service ID.
  If the service is closed, the value will be a random UUID.

* `vpc_id` - The VPC ID to which the service belongs.

* `network_ids` - The network ID list of subnets that the service have.

* `access_mode` - The access mode of Workspace service.
  + **INTERNET**: Access through Internet.
  + **DEDICATED**: Access through Direct Connect.
  + **DEDICATED**: Access through Internet or Direct Connect.

* `auth_type` - The authentication type of Workspace service.
  + **LITE_AS**: Local authentication.
  + **LOCAL_AD**: Local AD.

* `ad_domain` - The configuration of AD domain.  
  The [ad_domain](#workspace_data_service_ad_domain) structure is documented below.

* `enterprise_id` - The enterprise ID.

* `internet_access_port` - The internet access port.

* `internet_access_address` - The internet access address.

* `dedicated_subnets` - The subnet segments of the dedicated access.

* `management_subnet_cidr` - The subnet segment of the management component.

* `infrastructure_security_group` - The management component security group automatically created under the specified
  VPC after service is registered.  
  The [infrastructure_security_group](#workspace_data_service_security_group) structure is documented below.

* `desktop_security_group` - The desktop security group automatically created under the specified VPC after the service
  is registered.  
  The [desktop_security_group](#workspace_data_service_security_group) structure is documented below.

* `status` - The current status of the Workspace service.
  + **SUBSCRIBED**: The service has been subscribed.
  + **SUBSCRIPTION_FAILED**: The service cannot be subscribed.
  + **DEREGISTERING**: The service is being unsubscribed.
  + **DEREGISTRATION_FAILED**: The service cannot be unsubscribed.
  + **CLOSED**: The service has been unsubscribed and is not subscribed.

* `otp_config_info` - The configuration of auxiliary authentication.  
  The [otp_config_info](#workspace_data_service_otp_config_info) structure is documented below.

* `is_locked` - Whether the service is locked.
  + **0**: unlocked.
  + **1**: locked.

* `lock_time` - The time when the service is locked.

* `lock_reason` - The lock reason of the service.

<a name="workspace_data_service_ad_domain"></a>
The `ad_domain` block supports:

* `name` - The domain name.

* `admin_account` - The domain administrator account.

* `active_domain_ip` - The IP address of primary domain controller.

* `active_domain_name` - The name of the primary domain controller.

* `active_dns_ip` - The primary DNS IP address.

* `standby_domain_ip` - The IP address of standby domain controller.

* `standby_domain_name` - The name of the standby domain controller.

* `standby_dns_ip` - The standby DNS IP address.

* `delete_computer_object` - Whether to delete the corresponding computer object on AD while deleting the desktop.

<a name="workspace_data_service_security_group"></a>
The `infrastructure_security_group` and `desktop_security_group` block support:

* `id` - The security group ID.

* `name` - The security group name.

<a name="workspace_data_service_otp_config_info"></a>
The `otp_config_info` block supports:

* `enable` - Whether to enable auxiliary authentication.

* `receive_mode` - The verification code receiving mode.
  + **VMFA**: Virtual MFA device
  + **HMFA**: Hardware MFA device

* `auth_url` - The auxiliary authentication server address.

* `app_id` - The auxiliary authentication server access account.

* `app_secret` - The authentication service access password.

* `auth_server_access_mode` - The authentication service access mode.

* `cert_content` - The certificate content, in PEM format.

* `rule_type` - The type of the object to which authentication applies.
  + **ACCESS_MODE**: Access type.

* `rule` - The object to which authentication applies.
  + **INTERNET**: Internet access.
  + **PRIVATE**: private line access.
