---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_alarm_templates"
description: |-
  Use this data source to get the list of CES alarm templates.
---

# huaweicloud_ces_alarm_templates

Use this data source to get the list of CES alarm templates.

## Example Usage

```hcl
data huaweicloud_ces_alarm_templates "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of an alarm template.

* `type` - (Optional, String) Specifies the alarm template type.
  The valid values are as follows:
  + **system**: default metric template.
  + **custom**: custom metric template.
  + **system_event**: default event template.
  + **custom_event**: custom event template.
  + **system_custom_event**: all event templates.

* `namespace` - (Optional, String) Specifies the namespace of a service.

* `dimension_name` - (Optional, String) Specifies the resource dimension.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarm_templates` - The alarm template list.

  The [alarm_templates](#alarm_templates_struct) structure is documented below.

<a name="alarm_templates_struct"></a>
The `alarm_templates` block supports:

* `template_id` - The alarm template ID.

* `name` - The alarm template name.

* `type` - The alarm template type.

* `created_at` - The creation time of the alarm template.

* `description` - The alarm template description.
