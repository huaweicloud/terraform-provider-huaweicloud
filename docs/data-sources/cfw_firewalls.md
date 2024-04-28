---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_firewalls"
description: ""
---

# huaweicloud_cfw_firewalls

Use this data source to get the list of CFW firewalls.

## Example Usage

```hcl
data "huaweicloud_cfw_firewalls" "test" {
  service_type = 0
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.
  If not specified, the first instance will be returned.

* `service_type` - (Optional, Int) Specifies the service type. The value can be:
  + **0**: North-south firewall;
  + **1**: East-west firewall;

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The firewall instance records.
  The [records object](#firewalls_GetFirewallInstanceResponseRecord) structure is documented below.

<a name="firewalls_GetFirewallInstanceResponseRecord"></a>
The `records` block supports:

* `name` - The firewall name.

* `charge_mode` - The billing mode. The value can be 0 (yearly/monthly) or 1 (pay-per-use).

* `engine_type` - The engine type.

* `feature_toggle` - The map of feature toggle.

* `fw_instance_id` - The firewall ID.

* `ha_type` - The cluster type.

* `is_old_firewall_instance` - Whether the engine is an old engine.

* `service_type` - The service type.

* `support_ipv6` - Whether IPv6 is supported.

* `status` - The firewall status. The options are as follows:
  + **-1**: waiting for payment;
  + **0**: creating;
  + **1**: deleting;
  + **2**: running;
  + **3**: upgrading;
  + **4**: deletion completed;
  + **5**: freezing;
  + **6**: creation failed;
  + **7**: deletion failed;
  + **8**: freezing failed;
  + **9**: storage in progress;
  + **10**: storage failed;
  + **11**: upgrade failed;

* `flavor` - The flavor of the firewall.
  The [Flavor](#firewalls_GetFirewallInstanceResponseRecordFlavor) structure is documented below.

* `protect_objects` - The project list.
  The [Protect Object](#firewalls_GetFirewallInstanceResponseRecordProtectObject) structure is documented below.

* `resources` - The firewall instance resources.
  The [Firewall Instance Resource](#firewalls_GetFirewallInstanceResponseRecordFirewallInstanceResource) structure is
  documented below.

<a name="firewalls_GetFirewallInstanceResponseRecordFlavor"></a>
The `flavor` block supports:

* `bandwidth` - The bandwidth.

* `eip_count` - The number of EIPs.

* `log_storage` - The log storage.

* `version` - The firewall version. The value can be 0 (standard edition), 1 (professional edition),
  2 (platinum edition), or 3 (basic edition).

* `vpc_count` - The number of VPCs.

<a name="firewalls_GetFirewallInstanceResponseRecordProtectObject"></a>
The `protect_objects` block supports:

* `object_id` - The protected object ID.

* `object_name` - The protected object name.

* `type` - The project type. The options are as follows:
  + **0**: north-south;
  + **1**: east-west;

<a name="firewalls_GetFirewallInstanceResponseRecordFirewallInstanceResource"></a>
The `resources` block supports:

* `cloud_service_type` - Service type, which is used by CBC. The value is **hws.service.type.cfw**.

* `resource_id` - Resource ID.

* `resource_size` - Resource quantity.

* `resource_size_measure_id` - Resource unit name.

* `resource_spec_code` - Inventory unit code.

* `resource_type` - Resource type. The options are as follows:
  + **CFW**: hws.resource.type.cfw;
  + **EIP**: hws.resource.type.cfw.exp.eip;
  + **Bandwidth**: hws.resource.type.cfw.exp.bandwidth;
  + **VPC**: hws.resource.type.cfw.exp.vpc;
  + **Log storage**: hws.resource.type.cfw.exp.logaudit;
