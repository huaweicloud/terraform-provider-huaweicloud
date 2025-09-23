---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_ips_rules"
description: |-
  Use this data source to get the list of CFW IPS basic protection rules.
---

# huaweicloud_cfw_ips_rules

Use this data source to get the list of CFW IPS basic protection rules.

## Example Usage

```hcl
variable "object_id" {}

data "huaweicloud_cfw_ips_rules" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

* `ips_id` - (Optional, String) Specifies the IPS rule ID.

* `ips_name_like` - (Optional, String) Specifies the IPS rule name.
  This parameter supports fuzzy search.

* `ips_status` - (Optional, String) Specifies the IPS rule status.
  The valid value can be **OBSERVE**, **ENABLE**, or **CLOSE**.

* `is_updated_ips_rule_queried` - (Optional, Bool) Specifies whether to check for new update rules.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The IPS rule list.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `ips_group` - The IPS rule group.

* `ips_id` - The IPS rule ID.

* `ips_name` - The IPS rule name.

* `ips_rules_type` - The IPS rule type.

* `affected_application` - The application affected by the rule.

* `create_time` - The creation time.

* `ips_cve` - The CVE.

* `default_status` - The default status of the IPS rule.

* `ips_level` - The risk level.

* `ips_status` - The current status of the IPS rule.
