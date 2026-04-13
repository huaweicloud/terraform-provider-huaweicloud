---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_acl_rule_hit_info"
description: |-
  Use this data source to get the CFW ACL rule hit count and last hit time within HuaweiCloud.
---

# huaweicloud_cfw_acl_rule_hit_info

Use this data source to get the CFW ACL rule hit count and last hit time within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "rule_ids" {
  type = list(string)
}

data "huaweicloud_cfw_acl_rule_hit_info" "test" {
  fw_instance_id = var.fw_instance_id
  rule_ids       = var.rule_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `rule_ids` - (Required, List) Specifies the list of ACL rule IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The response to data for the number of hits and the last hit time of the rule.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `rule_id` - The rule ID.

* `rule_hit_count` - The rule hit count.

* `rule_last_hit_time` - The last hit time.
