---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_message_templates"
description: |-
  Use this data source to get the list of AOM message templates.
---

# huaweicloud_aom_message_templates

Use this data source to get the list of AOM message templates.

## Example Usage

```hcl
data "huaweicloud_aom_message_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the templates belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `message_templates` - Indicates the message templates.
  The [message_templates](#attrblock--message_templates) structure is documented below.

<a name="attrblock--message_templates"></a>
The `message_templates` block supports:

* `name` - Indicates the template name.

* `source` - Indicates the template type.
  + If it is empty, means it is a metric or event template.
  + If it is **LTS**, means it is a log template.

* `locale` - Indicates the meesage template language.

* `description` - Indicates the meesage template description.

* `enterprise_project_id` - Indicates the enterprise project ID to which the template belongs.

* `templates` - Indicates the templates.
  The [templates](#attrblock--message_templates--templates) structure is documented below.

* `created_at` - Indicates the message template create time.

* `updated_at` - Indicates the message template update time.

<a name="attrblock--message_templates--templates"></a>
The `templates` block supports:

* `content` - Indicates the content of the template.

* `sub_type` - Indicates the subscription type of the template.

* `topic` - Indicates the topic of the template.

* `version` - Indicates the version of the template.
