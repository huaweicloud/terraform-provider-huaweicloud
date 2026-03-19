---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_acl_rule_tags"
description: |-
  Use this data source to get the list of CFW acl rule tags.
---

# huaweicloud_cfw_acl_rule_tags

Use this data source to get the list of CFW acl rule tags.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_acl_rule_tags" "test" {
  fw_instance_id = var.fw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For enterprise users, if omitted, all enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The rule tags list.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `tag_value` - The rule tag value.

* `tag_id` - The rule ID.

* `tag_key` - The rule tag key.
