---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_custom_rule"
description: |-
  Manages a custom rule resource within HuaweiCloud HSS.
---

# huaweicloud_hss_custom_rule

Manages a custom rule resource within HuaweiCloud HSS.

## Example Usage

```hcl
variable "agent_ids" {
  type = list(string)
}

variable "rule_values" {
  type = list(string)
}

resource "huaweicloud_hss_custom_rule" "test" {
  rule_name   = "test-name"
  is_all_host = true
  agent_ids   = var.agent_ids

  custom_rule_value_info {
    auto_block  = 1
    hash_type   = "sha1"
    rule_type   = "black_hash"
    rule_values = var.rule_values
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `rule_name` - (Required, String) Specifies the custom rule name. The name does not support editing at the moment.

* `custom_rule_value_info` - (Required, List) Specifies the custom rule value information.  
  The [custom_rule_value_info](#custom_rule_value_info_struct) structure is documented below.

* `is_all_host` - (Optional, Bool) Specifies whether to apply the rule to all hosts.  
  Defaults to **false**.

* `agent_ids` - (Optional, List) Specifies the list of agent IDs to which the rule applies.  

* `rule_status` - (Optional, Int) Specifies the rule status.  
  Defaults to `1`. Valid values are:
  + `0`: disabled
  + `1`: enabled

<a name="custom_rule_value_info_struct"></a>
The `custom_rule_value_info` block supports:

* `rule_type` - (Required, String) Specifies the rule type. Valid value is **black_hash**.

* `hash_type` - (Required, String) Specifies the hash type. Valid values are: **sha256**, **md5**, and **sha1**.

* `auto_block` - (Required, Int) Specifies whether to automatically block.  
  Valid values are `0` (no auto block) and `1` (auto block).

* `rule_values` - (Required, List) Specifies the list of rule values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as the custom rule ID).

* `create_time` - The creation time of the custom rule.

* `update_time` - The last update time of the custom rule.

* `host_num` - The number of hosts associated with the custom rule.

* `agent_ids_attr` - The list of agent IDs that are associated with the custom rule.

## Import

Custom rule can be imported using the custom rule ID, e.g.

```bash
$ terraform import huaweicloud_hss_custom_rule.test <rule_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `agent_ids`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_custom_rule" "test" { 
  # ...

  lifecycle {
    ignore_changes = [
      agent_ids,
    ]
  }
}
```
