---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_microservice_engines"
description: |-
  Use this data source to query available microservice engines within HuaweiCloud.
---

# huaweicloud_cse_microservice_engines

Use this data source to query available microservice engines within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cse_microservice_engines" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where dedicated microservice engines are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `engines` - All queried dedicated microservice engines.  
  The [engines](#cse_microservice_engines) structure is documented below.

<a name="cse_microservice_engines"></a>
The `engines` block supports:

* `id` - The ID of the dedicated microservice engine.

* `name` - The name of the dedicated microservice engine.

* `flavor` - The flavor name of the dedicated microservice engine.

* `availability_zones` - The list of availability zones.

* `network_id` - The network ID of the subnet to which the dedicated microservice engine belongs.

* `auth_type` - The authentication method for the dedicated microservice engine
  + **RBAC**
  + **NONE**

* `version` - The version of the dedicated microservice engine.

* `enterprise_project_id` - The enterprise project ID to which the dedicated microservice engine belongs.

* `description` - The description of the dedicated microservice engine.

* `eip_id` - The EIP ID to which the dedicated microservice engine assocated.

* `extend_params` - The additional parameters for the dedicated microservice engine.

* `service_limit` - The maximum number of the microservice resources.

* `instance_limit` - The maximum number of the microservice instance resources.

* `service_registry_addresses` - The connection addresses of service center.  
  The [service_registry_addresses](#cse_microservice_engine_addresses) structure is documented below.

* `config_center_addresses` - The addresses of config center.  
  The [config_center_addresses](#cse_microservice_engine_addresses) structure is documented below.

* `status` - The status of the dedicated microservice engine.
  + **Creating**
  + **Available**
  + **Unavailable**
  + **Deleting**
  + **Deleted**
  + **Upgrading**
  + **Modifying**
  + **CreateFailed**
  + **DeleteFailed**
  + **UpgradeFailed**
  + **ModifyFailed**
  + **Freezed**

* `created_at` - The creation time of the dedicated microservice engine, in RFC3339 format.

<a name="cse_microservice_engine_addresses"></a>
The `service_registry_addresses` and `config_center_addresses` blocks support:

* `private` - The internal access address.

* `public` - The public access address.
