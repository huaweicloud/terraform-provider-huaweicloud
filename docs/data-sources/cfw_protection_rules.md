---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_protection_rules"
description: |-
  Use this data source to get the list of CFW protection rules.
---

# huaweicloud_cfw_protection_rules

Use this data source to get the list of CFW protection rules.

## Example Usage

```hcl
variable "object_id" {}
variable "name" {}

data "huaweicloud_cfw_protection_rules" "test" {
  object_id = var.object_id
  name      = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

* `name` - (Optional, String) Specifies the rule name.

* `type` - (Optional, String) Specifies the rule type.
  The value can be **0** (Internet rule), **1** (VPC rule), or **2** (NAT rule).

* `direction` - (Optional, String) Specifies the rule direction.
  The options are as follows:
  + **0**: inbound;
  + **1**: outbound.

* `rule_id` - (Optional, String) Specifies the rule ID.

* `status` - (Optional, String) Specifies the rule status.
  The options are as follows:
  + **0**: disabled;
  + **1**: enabled.

* `action_type` - (Optional, String) Specifies the rule action type.
  The options are as follows:
  + **0**: allow;
  + **1**: deny.

* `source` - (Optional, String) Specifies the source address.

* `destination` - (Optional, String) Specifies the destination address.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the protection rule.
  Tags should have only one key/value pair.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The protection rule list.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `rule_id` - The rule ID.

* `name` - The rule name.

* `type` - The rule type.

* `status` - The rule status.

* `description` - The rule description.

* `action_type` - The rule action type.

* `destination` - The destination configuration.

  The [destination](#records_destination_struct) structure is documented below.

* `service` - The service.

  The [service](#records_service_struct) structure is documented below.

* `address_type` - The address type.

* `long_connect_time_hour` - The persistent connection duration (hour).

* `long_connect_time_minute` - The persistent connection duration (minute).

* `long_connect_time_second` - The persistent connection duration (second).

* `long_connect_enable` - Whether to support persistent connections.

* `long_connect_time` - The persistent connection duration.

* `source` - The source configuration.

  The [source](#records_source_struct) structure is documented below.

* `tags` - The tag of a rule.

* `direction` - The direction of a rule.

* `created_date` - The created time of a rule.

* `last_open_time` - The last open time.

<a name="records_destination_struct"></a>
The `destination` block supports:

* `address_set_id` - The ID of the associated IP address group.

* `region_list` - The region list of a rule.

  The [region_list](#region_list_struct) structure is documented below.

* `domain_set_name` - The name of domain group.

* `domain_set_id` - The ID of the domain group.

* `ip_address` - The IP address list.

* `address_group` - The address group.

* `type` - The destination type.

* `address_type` - The destination address type.

* `address` - The destination IP address.

* `address_set_name` - The IP address group name.

* `domain_address_name` - The name of the domain name address.

* `address_set_type` - The destination address set type.

<a name="records_service_struct"></a>
The `service` block supports:

* `dest_port` - The destination port of the service.

* `service_set_id` - The service group ID.

* `service_set_name` - The service group name.

* `type` - The service input type.

* `source_port` - The source port.

* `service_group` - The service group list.

* `protocol` - The protocol type.

* `custom_service` - The custom service.

  The [custom_service](#service_custom_service_struct) structure is documented below.

* `service_set_type` - The service set type.

* `protocols` - The protocols.

<a name="service_custom_service_struct"></a>
The `custom_service` block supports:

* `dest_port` - The destination port.

* `description` - The custom service description.

* `name` - The custom service name.

* `protocol` - The protocol type of the custom service.

* `source_port` - The source port of the custom service.

<a name="records_source_struct"></a>
The `source` block supports:

* `address_set_id` - The ID of the associated IP address group.

* `region_list` - The region list of a rule.

  The [region_list](#region_list_struct) structure is documented below.

* `ip_address` - The IP address list.

* `address_group` - The address group.

* `type` - The source type.

* `address_type` - The address type.

* `address` - The source IP address.

* `address_set_name` - The IP address group name.

* `domain_address_name` - The name of the domain address.

* `address_set_type` - The address set type.

<a name="region_list_struct"></a>
The `region_list` block supports:

* `region_id` - The region ID.

* `description_cn` - The Chinese description of a region.

* `description_en` - The English description of a region.

* `region_type` - The region type.
