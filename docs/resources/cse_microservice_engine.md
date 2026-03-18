---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_engine"
description: ""
---

# huaweicloud_cse_microservice_engine

Manages a microservice engine (2.0+) resource within HuaweiCloud.

## Example Usage

### Create an engine for the default type CSE

```hcl
variable "engine_name" {}
variable "network_id" {}
variable "floating_ip_id" {}
variable "availability_zones" {
  type = list(string)
}
variable "manager_password" {}

resource "huaweicloud_cse_microservice_engine" "test" {
  name               = var.engine_name
  flavor             = "cse.s1.small2"
  network_id         = var.network_id
  eip_id             = var.floating_ip_id
  availability_zones = var.availability_zones
  auth_type          = "RBAC"
  admin_pass         = var.manager_password
}
```

### Create an engine for the type Nacos

```hcl
variable "engine_name" {}
variable "network_id" {}

resource "huaweicloud_cse_microservice_engine" "test" {
  name       = var.engine_name
  flavor     = "cse.nacos2.c1.large.10"
  network_id = var.network_id
  auth_type  = "NONE"
  version    = "Nacos2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the microservice engine is located.
  If omitted, the provider-level region will be used. Changing this will create a new engine.

* `name` - (Required, String, NonUpdatable) Specifies the name of the microservice engine.
 The name can contain `3` to `24` characters, only letters, digits and hyphens (-) are allowed.
  The name must start with a letter and cannot end with a hyphen (-).

* `flavor` - (Required, String, NonUpdatable) Specifies the flavor of the microservice engine.

* `network_id` - (Required, String, NonUpdatable) Specifies the network ID of the subnet to which the microservice
  engine belongs.

* `auth_type` - (Required, String, NonUpdatable) Specifies the authentication type for the microservice engine.  
  The valid values are as follows:
  + **RBAC**: Enable security authentication.  
    Security authentication applies to the scenario where multiple users use the same engine.
    After security authentication is enabled, all users who use the engine can log in using the account and password.
    You can assign the account and role in the System Management.
  + **NONE**: Disable security authentication.  
    After security authentication is disabled, all users who use the engine can use the engine without using the account
    and password, and have the same operation permissions on all services.

  -> **NONE** is required for the nacos engine.

* `availability_zones` - (Optional, List, NonUpdatable) Specifies the availability zones where the microservice engine
  is deployed.  
  Required if the `version` is **CSE2**.

* `version` - (Optional, String, NonUpdatable) Specifies the version of the microservice engine.  
  The valid values are as follows:
  + **CSE2**
  + **Nacos2**

  Defaults to **CSE2**.

* `admin_pass` - (Optional, String, NonUpdatable) Specifies the administrator account password of the microservice
  engine.  
  The corresponding account name is **root**.
  Required if `auth_type` is **RBAC**.  
  The password format must meet the following conditions:
  + Must be `8` to `32` characters long.
  + A password must contain at least one digit, one uppercase letter, one lowercase letter, and one special character
    (-~!@#%^*_=+?$&()|<>{}[]).
  + Cannot be the account name or account name spelled backwards.
  + The password can only start with a letter.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID to which the
  microservice engine belongs.  
  If omitted and the version is **Nacos2**, the default enterprise project will be used.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the microservice engine.
  The description can contain a maximum of `255` characters.

* `eip_id` - (Optional, String, NonUpdatable) Specifies the EIP ID bound to the microservice engine.

* `extend_params` - (Optional, Map, NonUpdatable) Specifies the extended parameters of the microservice engine.

-> After the engine is created, the system will automatically add a series of additional parameters to it.  
  The specific parameters are subject to the state of the microservice engine.  
  This parameter will be affected by these parameters and will appear when `terraform plan` or `terraform apply`.  
  If it is inconsistent with the script configuration, it can be ignored by `ignore_changes` in non-change scenarios.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `service_limit` - The maximum number of the microservice resources.

* `instance_limit` - The maximum number of the microservice instance resources.

* `service_registry_addresses` - The service registry addresses of the microservice engine.  
  The [service_registry_addresses](#cse_microservice_engine_service_registry_addresses) structure is documented below.

* `config_center_addresses` - The config center addresses of the microservice engine.  
  The [config_center_addresses](#cse_microservice_engine_config_center_addresses) structure is documented below.

<a name="cse_microservice_engine_service_registry_addresses"></a>
The `service_registry_addresses` block supports:

* `private` - The private address of the service registry.

* `public` - The public address of the service registry.  
  This address is only set when EIP is bound.

<a name="cse_microservice_engine_config_center_addresses"></a>
The `config_center_addresses` block supports:

* `private` - The private address of the config center.

* `public` - The public address of the config center.  
  This address is only set when EIP is bound.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 20 minutes.

## Import

Engines can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_cse_microservice_engine.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes are `admin_pass` and `extend_params`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cse_microservice_engine" "test" {
  ...

  lifecycle {
    ignore_changes = [
      admin_pass,
      extend_params,
    ]
  }
}
```

For the engine created with the `enterprise_project_id`, its enterprise project ID needs to be specified additionally
when importing, the format is `<id>/<enterprise_project_id>`, e.g.

```bash
$ terraform import huaweicloud_cse_microservice_engine.test <id>/<enterprise_project_id>
```
