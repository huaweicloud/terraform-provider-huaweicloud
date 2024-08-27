---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_engine"
description: ""
---

# huaweicloud_cse_microservice_engine

Manages a dedicated microservice engine (2.0+) resource within HuaweiCloud.

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

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the dedicated microservice engine.
  If omitted, the provider-level region will be used. Changing this will create a new engine.

* `name` - (Required, String, ForceNew) Specifies the name of the dedicated microservice engine.
 The name can contain `3` to `24` characters, only letters, digits and hyphens (-) are allowed.
  The name must start with a letter and cannot end with a hyphen (-).
  Changing this will create a new engine.

* `flavor` - (Required, String, ForceNew) Specifies the flavor of the dedicated microservice engine.
  Changing this will create a new engine.

* `availability_zones` - (Optional, List, ForceNew) Specifies the list of availability zones.  
  Required if the `version` is **CSE2**.  
  Changing this will create a new engine.

* `network_id` - (Required, String, ForceNew) Specifies the network ID of the subnet to which the dedicated microservice
  engine belongs. Changing this will create a new engine.

* `auth_type` - (Required, String, ForceNew) Specifies the authentication method for the dedicated microservice engine.
  Changing this will create a new engine.
  + **RBAC**: Enable security authentication.
    Security authentication applies to the scenario where multiple users use the same engine.
    After security authentication is enabled, all users who use the engine can log in using the account and password.
    You can assign the account and role in the System Management.
  + **NONE**: Disable security authentication.
    After security authentication is disabled, all users who use the engine can use the engine without using the account
    and password, and have the same operation permissions on all services.

  -> **NONE** is required for the nacos engine.

* `version` - (Optional, String, ForceNew) Specifies the version of the dedicated microservice engine.  
  The valid values are as follows:
  + **CSE2**
  + **Nacos2**

  Defaults to **CSE2**. Changing this will create a new engine.

* `admin_pass` - (Optional, String, ForceNew) Specifies the account password. The corresponding account name is **root**.
  Required if `auth_type` is **RBAC**. Changing this will create a new engine.
  The password format must meet the following conditions:
  + Must be `8` to `32` characters long.
  + A password must contain at least one digit, one uppercase letter, one lowercase letter, and one special character
    (-~!@#%^*_=+?$&()|<>{}[]).
  + Cannot be the account name or account name spelled backwards.
  + The password can only start with a letter.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.  
  If omitted and the version is **Nacos2**, the default enterprise project will be used.  
  Changing this will create a new engine.

* `description` - (Optional, String, ForceNew) Specifies the description of the dedicated microservice engine.
  The description can contain a maximum of `255` characters.
  Changing this will create a new engine.

* `eip_id` - (Optional, String, ForceNew) Specifies the EIP ID to which the dedicated microservice engine assocated.
  Changing this will create a new engine.

* `extend_params` - (Optional, Map, ForceNew) Specifies the additional parameters for the dedicated microservice engine.
  Changing this will create a new engine.

-> After the engine is created, the system will automatically add a series of additional parameters to it.
  The specific parameters are subject to the state of the dedicated microservice engine.
  This parameter will be affected by these parameters and will appear when `terraform plan` or `terraform apply`.
  If it is inconsistent with the script configuration, it can be ignored by `ignore_changes` in non-change scenarios.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the dedicated
  microservice engine belongs.  
  Changing this will create a new engine.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `service_limit` - The maximum number of the microservice resources.

* `instance_limit` - The maximum number of the microservice instance resources.

* `service_registry_addresses` - The connection address of service center.
  The [object](#engine_center_addresses) structure is documented below.

* `config_center_addresses` - The address of config center.
  The [object](#engine_center_addresses) structure is documented below.

<a name="engine_center_addresses"></a>
The `service_registry_addresses` and `config_center_addresses` block supports:

* `private` - The internal access address.

* `public` - The public access address. This address is only set when EIP is bound.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 20 minutes.

## Import

Engines can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_cse_microservice_engine.test eddc5d42-f9d5-4f8e-984b-d6f3e088561c
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
$ terraform import huaweicloud_cse_microservice_engine.test eddc5d42-f9d5-4f8e-984b-d6f3e088561c/ef101e1a-990c-42cd-bb99-a4474e41e461
```
