---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_scan_rule"
description: |-
  Manages a DSC scan rule resource within HuaweiCloud.
---

# huaweicloud_dsc_scan_rule

Manages a DSC scan rule resource within HuaweiCloud.

## Example Usage

### Create a scan rule with rule matching

```hcl
variable "rule_name" {}
variable "contents" {
  type = list(object({
    effective_mode = string
    location       = string
    rule_content   = string
  }))
}
variable "templates" {
  type = list(object({
    template_id       = string
    classification_id = string
    security_level_id = string
    is_used           = optional(string, "true")
  }))
}

resource "huaweicloud_dsc_scan_rule" "test" {
  rule_name      = var.rule_name
  rule_type      = "REGEX"
  category       = "BUILT_SELF"
  logic_operator = "AND"
  match_rate     = 1
  min_match      = 1

  dynamic "content" {
    for_each = var.contents

    content {
      effective_mode = content.value.effective_mode
      location       = content.value.location
      rule_content   = content.value.rule_content
    }
  }

  dynamic "templates" {
    for_each = var.templates

    content {
      template_id       = templates.value.template_id
      classification_id = templates.value.classification_id
      security_level_id = templates.value.security_level_id
      is_used           = templates.value.is_used
    }
  }
}
```

### Create a scan rule with keyword matching

```hcl
variable "rule_name" {}
variable "contents" {
  type = list(object({
    rule_content = string
  }))
}
variable "templates" {
  type = list(object({
    template_id       = string
    classification_id = string
    security_level_id = string
    is_used           = optional(string, "true")
  }))
}

resource "huaweicloud_dsc_scan_rule" "test" {
  rule_name      = var.rule_name
  rule_type      = "KEYWORD"
  category       = "BUILT_SELF"
  logic_operator = "AND"
  match_rate     = 1
  min_match      = 1

  dynamic "content" {
    for_each = var.contents

    content {
      effective_mode = "KEYWORD"
      location       = "CONTENT"
      rule_content   = content.value.rule_content
    }
  }

  dynamic "templates" {
    for_each = var.templates

    content {
      template_id       = templates.value.template_id
      classification_id = templates.value.classification_id
      security_level_id = templates.value.security_level_id
      is_used           = templates.value.is_used
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the scan rule is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `rule_name` - (Required, String) Specifies the name of the scan rule.  
  The maximum length is `255` characters, only letters, Chinese characters, digits, parentheses (`()`), hyphens (`-`),
  and underscores (`_`) are allowed.

* `rule_type` - (Required, String) Specifies the type of the scan rule.  
  The valid values are as follows:
  + **KEYWORD**
  + **REGEX**

* `category` - (Required, String, NonUpdatable) Specifies the category of the scan rule.  
  The valid values are as follows:
  + **BUILT_SELF**: User-defined rule.

* `logic_operator` - (Required, String) Specifies the logic operator of the scan rule.  
  The valid values are as follows:
  + **AND**
  + **OR**

* `match_rate` - (Required, Int) Specifies the match rate of the scan rule.  
  The valid value ranges from `1` to `100`.

* `min_match` - (Required, Int) Specifies the minimum match count of the scan rule.  
  The valid value ranges from `1` to `100`.

* `content` - (Required, List) Specifies the content list of the scan rule.  
  The [content](#dsc_scan_rule_content) structure is documented below.  
  When the `rule_type` is **KEYWORD**, the `content` block only supports one item.

* `rule_desc` - (Optional, String) Specifies the description of the scan rule.  
  The description maximum length is `255` characters, only letters, Chinese characters, digits, parentheses (`()`),
  hyphens (`-`), and underscores (`_`) are allowed.

* `templates` - (Optional, List) Specifies the template list associated with the scan rule.  
  The [templates](#dsc_scan_rule_templates) structure is documented below.

<a name="dsc_scan_rule_content"></a>
The `content` block supports:

* `effective_mode` - (Required, String) Specifies the effective mode of the rule content.  
  The valid values are as follows:
  + **IN**
  + **NOT_IN**
  + **REGEX**
  + **KEYWORD**

  When the `rule_type` is **KEYWORD**, this parameter must be **KEYWORD**.

* `location` - (Required, String) Specifies the application location of the rule content.  
  The valid values are as follows:
  + **NAME**: Column name.
  + **REMARK**: Column remark.
  + **CONTENT**: Column content.
  + **TABLE_NAME**: Table name.
  + **TABLE_REMARK**: Table remark.

  When the `rule_type` is **KEYWORD**, this parameter must be **CONTENT**.

* `rule_content` - (Required, String) Specifies the content of the rule.  
  When the `rule_type` is **KEYWORD**, multiple contents are separated by commas (`,`).

<a name="dsc_scan_rule_templates"></a>
The `templates` block supports:

* `template_id` - (Required, String) Specifies the template ID associated with the rule.

* `classification_id` - (Required, String) Specifies the classification ID associated with the rule.

* `security_level_id` - (Required, String) Specifies the security level ID associated with the rule.

* `is_used` - (Optional, String) Specifies whether the rule is enabled in the template.  
  The valid values are as follows:
  + **true**
  + **false**

  Defaults to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `templates` - The template list associated with the scan rule.  
  The [templates](#dsc_scan_rule_templates_attr) structure is documented below.

<a name="dsc_scan_rule_templates_attr"></a>
The `templates` block supports:

* `id` - The ID of the rule template.

* `template_name` - The name of the template.

* `classification_name` - The classification name associated with the rule.

* `security_level_color` - The color number corresponding to the security level.

* `security_level_name` - The name of the security level.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dsc_scan_rule.test <id>
```
