---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_alarm_template_association_alarms"
description: |-
  Use this data source to get the list of alarm rules associated with an alarm template.
---

# huaweicloud_ces_alarm_template_association_alarms

Use this data source to get the list of alarm rules associated with an alarm template.

## Example Usage

```hcl
variable "template_id" {}

data "huaweicloud_ces_alarm_template_association_alarms" "test" {
  template_id = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `template_id` - (Required, String) Specifies the ID of an alarm template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarms` - The alarm rules using the given template.

  The [alarms](#alarms_struct) structure is documented below.

<a name="alarms_struct"></a>
The `alarms` block supports:

* `alarm_id` - The ID of an alarm rule.

* `name` - The name of an alarm rule.

* `description` - The description of an alarm rule.
