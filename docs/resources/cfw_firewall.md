---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_firewall"
description: ""
---

# huaweicloud_cfw_firewall

Manages a CFW firewall resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
resource "huaweicloud_cfw_firewall" "test" {
  name = "test"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
```

### PrePaid firewall

```hcl
resource "huaweicloud_cfw_firewall" "test" {
  name = "test"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
```

### firewall with east-west firewall

```hcl
resource "huaweicloud_cfw_firewall" "test" {
  name = "test"

  east_west_firewall_inspection_cidr = "172.16.1.0/24"
  east_west_firewall_er_id           = huaweicloud_er_instance.test.id
  east_west_firewall_mode            = "er"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
```

### firewall with IPS switch and IPS protection mode

```hcl
resource "huaweicloud_cfw_firewall" "test" {
  name = "test"

  flavor {
    version = "Professional"
  }

  tags = {
    key = "value"
    foo = "bar"
  }

  charging_mode        = "prePaid"
  period_unit          = "month"
  period               = 1
  auto_renew           = false
  ips_switch_status    = 1
  ips_protection_mode  = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the firewall name.

  Changing this parameter will create a new resource.

* `flavor` - (Required, List, ForceNew) Specifies the flavor of the firewall.
  Changing this parameter will create a new resource.
  The [flavor](#Firewall_Flavor) structure is documented below.

* `east_west_firewall_inspection_cidr` - (Optional, String) Specifies the inspection cidr of the east-west firewall.

* `east_west_firewall_mode` - (Optional, String) Specifies the mode of the east-west firewall.
  The value can be: **er**.

* `east_west_firewall_er_id` - (Optional, String) Specifies the ER ID of the east-west firewall.

* `east_west_firewall_status` - (Optional, Int) Specifies the protection statue of the east-west firewall.
  The value can be: `0`(enabled) and `1`(disabled). Defaults to `0`.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the firewall.

  Changing this parameter will create a new resource.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the firewall.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

  Changing this parameter will create a new resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period.
    If `period_unit` is set to **month**, the value ranges from 1 to 9.
    If `period_unit` is set to **year**, the value ranges from 1 to 3.
    This parameter is mandatory if `charging_mode` is set to **prePaid**.

  Changing this parameter will create a new resource.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.

  Changing this parameter will create a new resource.

* `ips_switch_status` - (Optional, Int) Specifies the IPS patch switch status of the firewall.
  The value can be `0`(disabled) and `1`(enabled). Defaults to `0`.

* `ips_protection_mode` - (Optional, Int) Specifies the IPS protection mode of the firewall. Defaults to `0`.

  Valid values are as follows:
  + **0**: Observation Mode.
  + **1**: Strict Mode.
  + **2**: Medium Mode.
  + **3**: Loose Mode.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the firewall.

<a name="Firewall_Flavor"></a>
The `flavor` block supports:

* `version` - (Required, String, ForceNew) Specifies the version of the firewall.
  When the charging_mode is **prePaid**: the value can be **Standard** and **Prefessional**.
  When the charging_mode is **postPaid**: the value can be **Prefessional**.
  Changing this parameter will create a new resource.

* `extend_eip_count` - (Optional, Int, ForceNew) Specifies the extend EIP number of the firewall.
  Only works when the charging_mode is **prePaid**.
  Changing this parameter will create a new resource.

* `extend_bandwidth` - (Optional, Int, ForceNew) Specifies the extend bandwidth of the firewall.
  Only works when the charging_mode is **prePaid**.
  Changing this parameter will create a new resource.

* `extend_vpc_count` - (Optional, Int, ForceNew) Specifies the extend VPC number of the firewall.
  Only works when the charging_mode is **prePaid**.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `engine_type` - The engine type

* `ha_type` - The HA type.

* `flavor` - The firewall flavor.
  The [flavor](#Firewall_Flavor_Attribute) structure is documented below.

* `protect_objects` - The protect objects list.
  The [protect_objects](#Firewall_ProtectObject) structure is documented below.

* `service_type` - The service type.

* `status` - The firewall status.

* `support_ipv6` - Whether IPv6 is supported.

* `east_west_firewall_inspection_vpc_id` - The east-west firewall inspection VPC ID.

* `east_west_firewall_er_attachment_id` - Enterprise Router and Firewall Connection ID.

<a name="Firewall_Flavor_Attribute"></a>
The `flavor` block supports:

* `eip_count` - The EIP number of the firewall.

* `vpc_count` - The VPC number of the firewall.

* `bandwidth` - The bandwidth of the firewall.

* `log_storage` - The log storage of the firewall.

* `default_eip_count` - The default EIP number of the firewall.

* `default_vpc_count` - The default VPC number of the firewall.

* `default_bandwidth` - The default bandwidth of the firewall.

* `default_log_storage` - The default log storage of the firewall.

* `vpc_bandwidth` - The VPC bandwidth of the firewall.

* `used_rule_count` - The used rule count of the firewall.

* `total_rule_count` - The total rule count of the firewall.

<a name="Firewall_ProtectObject"></a>
The `protect_objects` block supports:

* `object_id` - The protected object ID.

* `object_name` - The protected object name.

* `type` - The object type.
  The options are as follows: 0: north-south; 1: east-west.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The firewall can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cfw_firewall.test 6cb1ce47-9990-447e-b071-d167c5393871
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`period_unit`, `period` and `auto_renew`. It is generally
recommended running `terraform plan` after importing an CFW firewall. You can then decide if changes should be applied to
the firewall, or the resource definition should be updated to align with the firewall. Also you can ignore changes as
below.

```hcl
resource "huaweicloud_cfw_firewall" "test" {
    ...

  lifecycle {
    ignore_changes = [
      period_unit, period, auto_renew
    ]
  }
}
```
