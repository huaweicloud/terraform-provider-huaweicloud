---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_protection_rule"
description: |-
  Manages a CFW protection rule resource within HuaweiCloud.
---

# huaweicloud_cfw_protection_rule

!> **WARNING:** It has been deprecated, use `huaweicloud_cfw_acl_rule` instead.

Manages a CFW protection rule resource within HuaweiCloud.

## Example Usage

### Create a basic rule

```hcl
variable "name" {}
variable "description" {}
variable "object_id" {}

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = var.name
  object_id           = var.object_id
  description         = var.description
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type    = 0
    address = "192.168.0.1"
  }

  destination {
    type    = 0
    address = "192.168.0.2"
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
```

### Create a rule with the source address using the region list

```hcl
variable "name" {}
variable "description" {}
variable "object_id" {}

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = var.name
  object_id           = var.object_id
  description         = var.description
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type = 3

    region_list {
      region_id      = "GR"
      description_cn = "希腊"
      description_en = "Greece"
      region_type    = 0
    }

    region_list {
      region_id      = "ZJ"
      description_cn = "浙江"
      description_en = "ZHEJIANG"
      region_type    = 1
    }
  }

  destination {
    type    = 0
    address = "192.168.0.2"
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
```

### Create a rule with the custom service

```hcl
resource "huaweicloud_cfw_protection_rule" "test" {
  name                = var.name
  object_id           = var.object_id
  description         = var.description
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type    = 0
    address = "1.1.1.2"
  }

  service {
    type = 2

    custom_service {
      protocol    = 6
      source_port = 80
      dest_port   = 80
    }

    custom_service {
      protocol    = 6
      source_port = 8080
      dest_port   = 8080
    }
  }

  sequence {
    top = 1
  }
}
```

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The rule name.

* `object_id` - (Required, String, ForceNew) The protected object ID.

  Changing this parameter will create a new resource.

* `type` - (Required, Int) The rule type.
  The value can be `0` (Internet rule), `1` (VPC rule), or `2` (NAT rule).

* `action_type` - (Required, Int) The action type.
  The value can be `0` (allow) or `1` (deny).

* `address_type` - (Required, Int) The address type.
  The value can be `0` (IPv4) or `1` (IPv6).

* `sequence` - (Required, List) The sequence configuration.
The [Order Rule](#ProtectionRule_OrderRuleAcl) structure is documented below.

* `service` - (Required, List) The service configuration.
The [Rule Service](#ProtectionRule_RuleService) structure is documented below.

* `source` - (Required, List) The source configuration.
The [Rule Source Address](#ProtectionRule_RuleSourceAddress) structure is documented below.

* `destination` - (Required, List) The destination configuration.
The [Rule Destination Address](#ProtectionRule_RuleDestinationAddress) structure is documented below.

* `status` - (Required, Int) The rule status. The options are as follows:
  + **0**: disabled;
  + **1**: enabled;

* `long_connect_enable` - (Required, Int) Whether to support persistent connections.
  The options are as follows:
  + **0**: supported;
  + **1**: not supported;

* `long_connect_time_hour` - (Optional, Int) The persistent connection duration (hour).

* `long_connect_time_minute` - (Optional, Int) The persistent connection duration (minute).

* `long_connect_time_second` - (Optional, Int) The persistent Connection Duration (second).

* `description` - (Optional, String) The description.

* `direction` - (Optional, Int) The direction. The options are as follows:
  + **0**: inbound;
  + **1**: outbound;

* `rule_hit_count` - (Optional, String) The number of times the protection rule is hit.
  Setting the value to **0** will clear the hit count. Value options: **0**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the protection rule.
  Tags should have only one key/value pair.

<a name="ProtectionRule_OrderRuleAcl"></a>
The `sequence` block supports:

* `dest_rule_id` - (Optional, String) The ID of the rule that the added rule will follow.
  This parameter cannot be left blank if the rule is not pinned on top, and is empty when the added rule is pinned on top.

* `top` - (Optional, Int) Whether to pin on top.
  The options are as follows:
  + **0**: no;
  + **1**: yes;

<a name="ProtectionRule_RuleService"></a>
The `service` block supports:

* `type` - (Required, Int) The service input type.
  The options are as follows:
  + **0**: manual input;
  + **1**: automatic input;
  + **2**: multiple objects;

* `dest_port` - (Optional, String) The destination port.

* `protocol` - (Optional, Int) The protocol type. The options are as follows:
  + **6**: TCP;
  + **17**: UDP;
  + **1**: ICMP;
  + **58**: ICMPv6;
  + **-1**: any protocol;
  
  Regarding the addition type, a null value indicates it is automatically added.

* `service_set_id` - (Optional, String) The service group ID.
  This parameter is left blank for the manual type and cannot be left blank for the automatic type.

* `service_set_name` - (Optional, String) The service group name.
  This parameter is left blank for the manual type and cannot be left blank for the automatic type.

* `source_port` - (Optional, String) The source port.

* `custom_service` - (Optional, List) The custom service list.
  When using this parameter, the `type` must be set to **2** (multiple objects).

  The [custom_service](#ProtectionRule_RuleCustomService) structure is documented below.

* `service_group` - (Optional, List) The list of service group IDs.
  When using this parameter, the `type` must be set to **2** (multiple objects).

<a name="ProtectionRule_RuleSourceAddress"></a>
The `source` block supports:

* `type` - (Required, Int) The source type. The options are as follows:
  + **0**: manual input;
  + **1**: associated IP address group;
  + **2**: domain name;
  + **3**: region;
  + **4**: application domain name group;
  + **5**: multiple objects;
  + **6**: network domain name group;
  + **7**: application domain name;

* `address` - (Optional, String) The IP address.
  The value cannot be empty for the manual type, and cannot be empty for the automatic or domain type.

* `address_set_id` - (Optional, String) The ID of the associated IP address group.
  This parameter cannot be left blank when `type` is set to **1** (associated IP address group).

* `address_set_name` - (Optional, String) The IP address group name.
  This parameter cannot be left blank when `type` is set to **1** (associated IP address group).

* `address_type` - (Optional, Int) The address type. The options are as follows:
  + **0**: IPv4;
  + **1**: IPv6;

* `domain_address_name` - (Optional, String) The name of the domain name address.
  This parameter cannot be left empty for the domain name type, and is empty for the manual or automatic type.

* `region_list` - (Optional, List) The region list.
  The [region_list](#ProtectionRule_RuleRegionList) structure is documented below.

* `ip_address` - (Optional, List) The IP address list.
  When using this parameter, the `type` must be set to **5** (multiple objects).

* `address_group` - (Optional, List) The list of address group IDs.
  When using this parameter, the `type` must be set to **5** (multiple objects).

<a name="ProtectionRule_RuleDestinationAddress"></a>
The `destination` block supports:

* `type` - (Required, Int) The destination type. The options are as follows:
  + **0**: manual input;
  + **1**: associated IP address group;
  + **2**: domain name;
  + **3**: region;
  + **4**: application domain name group;
  + **5**: multiple objects;
  + **6**: network domain name group;
  + **7**: application domain name;

* `address` - (Optional, String) The IP address.
  The value cannot be empty for the manual type, and cannot be empty for the automatic or domain type.

* `address_set_id` - (Optional, String) The ID of the associated IP address group.
  This parameter cannot be left blank when `type` is set to **1** (associated IP address group).

* `address_set_name` - (Optional, String) The IP address group name.
  This parameter cannot be left blank when `type` is set to **1** (associated IP address group).

* `address_type` - (Optional, Int) The address type. The options are as follows:
  + **0**: IPv4;
  + **1**: IPv6;

* `domain_address_name` - (Optional, String) The name of the domain name address.
  This parameter is valid when `type` is set to **2** (domain name) or **7** (application domain name).

* `region_list` - (Optional, List) The region list.
  The [region_list](#ProtectionRule_RuleRegionList) structure is documented below.

* `ip_address` - (Optional, List) The IP address list.
  When using this parameter, the `type` must be set to **5** (multiple objects).

* `domain_set_id` - (Optional, String) The ID of the domain group.
  The value cannot be left blank when `type` is set to **4** (application domain name group) or **6** (network domain
  name group).
  
* `domain_set_name` - (Optional, String) The name of domain group.
  The value cannot be left blank when `type` is set to **4** (application domain name group) or **6** (network domain
  name group).

* `address_group` - (Optional, List) The list of address group IDs.
  When using this parameter, the `type` must be set to **5** (multiple objects).

<a name="ProtectionRule_RuleCustomService"></a>
The `custom_service` block supports:

* `protocol` - (Required, Int) The protocol type. The options are as follows:
  + **6**: TCP;
  + **17**: UDP;
  + **1**: ICMP;
  + **58**: ICMPv6;
  + **-1**: any protocol;

* `source_port` - (Required, String) The source port.

* `dest_port` - (Required, String) The destination port.

<a name="ProtectionRule_RuleRegionList"></a>
The `region_list` block supports:

* `region_id` - (Required, String) The region ID.

* `region_type` - (Required, Int) The region type. The options are as follows:
  + **0**: country;
  + **1**: province;
  + **2**：continent;

* `description_cn` - (Optional, String) The Chinese description of the region.

* `description_en` - (Optional, String) The English description of the region.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The protection rule can be imported using `object_id`, `id`, separated by a slash, e.g.

```sh
$ terraform import huaweicloud_cfw_protection_rule.test <object_id>/<id>
```
