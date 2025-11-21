---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_external_config_rule_compliances"
description: |-
  Use this data source to get the external config rule compliance in Resource Governance Center.
---

# huaweicloud_rgc_external_config_rule_compliances

Use this data source to get the external config rule compliance in Resource Governance Center.

## Example Usage

```hcl
variable managed_account_id {}

data "huaweicloud_rgc_external_config_rule_compliances" "test" {
  managed_account_id = var.managed_account_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `managed_account_id` - (Required, String) The ID of the managed account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `config_rule_compliances` - Information about the external config rule compliance status.

The [config_rule_compliances](#config_rule_compliances) structure is documented below.

<a name="config_rule_compliances"></a>
The `config_rule_compliances` block supports:

* `rule_name` - The name of the config rule.

* `status` - The compliance status of the config rule.

* `region` - The region where the config rule is located.

* `control_id` - The ID of the control associated with the config rule.
