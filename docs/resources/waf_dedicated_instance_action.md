---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_dedicated_instance_action"
description: |-
  Manages a WAF dedicated instance action resource within HuaweiCloud.
---

# huaweicloud_waf_dedicated_instance_action

Manages a WAF dedicated instance action resource within HuaweiCloud.

-> This resource is only a one-time action resource using to operate WAF dedicated instance. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

### Modify security group

```hcl
variable instance_id {}
variable security_group_id {}

resource "huaweicloud_waf_dedicated_instance_action" "test" {
  instance_id = var.instance_id
  action      = "security_groups"
  params      = [var.security_group_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this setting will create a new instance.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of WAF dedicated instance.

* `action` - (Required, String, NonUpdatable) Specifies the operation name.

* `params` - (Optional, List, NonUpdatable) Specifies the specific request body, which is a list of strings.

-> The valid values for parameters `action` and `params` can be found in the Example Usage section of the document.

## Attribute Reference

The following attributes are exported:

* `id` - The resource ID (also the WAF dedicated instance ID).

* `instancename` - The name of the WAF dedicated instance.

* `server_id` - The ID of the ECS hosting the dedicated engine.

* `zone` - The availability zone code.

* `arch` - The CPU architecture code.

* `cpu_flavor` - The ECS specification code.

* `vpc_id` - The VPC ID where the dedicated engine is located.

* `subnet_id` - The subnet ID of the VPC where the dedicated engine is located.

* `service_ip` - The service plane IP address of the WAF dedicated instance.

* `service_ipv6` - The service plane IPv6 address of the WAF dedicated instance.

* `float_ip` - The management plane IP of the dedicated engine.

* `security_group_ids` - The security groups bound to the dedicated engine ECS.

* `status` - The billing status of the instance.  
  The valid values are as follows:
  + **0**: Normal billing.
  + **1**: Frozen (resources and data will be retained, but tenants can no longer use cloud services normally).
  + **2**: Terminated (resources and data will be cleared).

* `run_status` - The running status of the instance.  
  The valid values are as follows:
  + **0**: Creating.
  + **1**: Running.
  + **2**: Deleting.
  + **3**: Deleted.
  + **4**: Creation failed.
  + **5**: Frozen.
  + **6**: Abnormal.
  + **7**: Updating.
  + **8**: Update failed.

* `access_status` - The access status of the instance. `0`: inaccessible, `1`: accessible.

* `upgradable` - Whether the dedicated engine can be upgraded. `0`: Cannot be upgraded, `1`: Can be upgraded.

* `cloud_service_type` - The cloud service code.

* `resource_type` - The cloud service resource type.

* `resource_spec_code` - The cloud service resource code.

* `specification` - The ECS specification of the dedicated engine, such as **8vCPUs | 16GB**.

* `hosts` - The domains protected by the dedicated engine.

  The [hosts](#dedicated_domain_instance_hosts) structure is documented below.

* `volume_type` - The storage type.

* `cluster_id` - The storage resource pool ID.

* `pool_id` - The ID of the WAF group where the dedicated engine is located (only applicable to special dedicated mode).

* `charge_mode` - The billing mode.  
  The valid values are as follows:
  + **0**: Package cycle.
  + **1**: On-demand.

<a name="dedicated_domain_instance_hosts"></a>
The `hosts` block supports:

* `id` - The ID of the protected domain.

* `hostname` - The protected domain name.
