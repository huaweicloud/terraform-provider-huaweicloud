---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_message_template"
description:  |-
  Manages an AOM message template resource within HuaweiCloud.
---

# huaweicloud_aom_message_template

Manages an AOM message template resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "enterprise_project_id" {}

resource "huaweicloud_aom_message_template" "test" {
  name                  = var.name
  locale                = "en-us"
  enterprise_project_id = var.enterprise_project_id
  description           = "test"

  templates {
    sub_type = "email"
    topic    = "$${region_name}[$${event_severity}_$${event_type}_$${clear_type}] have a new alert at $${starts_at}."
    content  = <<EOF
Alarm Name:$${event_name};
Alarm ID:$${id};
Occurred:$${starts_at};
Event Severity:$${event_severity};
Alarm Info:$${alarm_info};
Resource Identifier:$${resources_new};
Suggestion:$${alarm_fix_suggestion_zh};
EOF
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the meesage template name.
  Changing this creates a new resource.

* `locale` - (Required, String) Specifies the meesage template language. Valid values are **en-us** and **zh-cn**.

* `templates` - (Required, List) Specifies the templates.
  The [templates](#block--templates) structure is documented below.

* `description` - (Optional, String) Specifies the meesage template description.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the template belongs.
  Changing this creates a new resource.

* `source` - (Optional, String, ForceNew) Specifies the template type.
  + If it is empty, means it is a metric or event template.
  + If it is **LTS**, means it is a log template.
  Changing this creates a new resource.

<a name="block--templates"></a>
The `templates` block supports:

* `content` - (Required, String) Specifies the content of the template.

* `sub_type` - (Required, String) Specifies the subscription type of the template.
  Valid value are **email**, **sms**, **wechat**, **dingding**, **webhook**, **voice**, **espace**, **feishu**, **welink**.

* `topic` - (Optional, String) Specifies the topic of the template.

* `version` - (Optional, String) Specifies the version of the template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is same as `name`.

* `created_at` - Indicates the meesage template create time.

* `updated_at` - Indicates the meesage template update time.

## Import

The message template can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_aom_message_template.test <name>
```
