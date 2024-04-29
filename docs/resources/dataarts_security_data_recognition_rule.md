---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_data_recognition_rule"
description: ""
---

# huaweicloud_dataarts_security_data_recognition_rule

Manages a data recognition rule resource of DataArts Security within HuaweiCloud.

## Example Usage

### CUSTOM RULE CHECK

```hcl
variable "workspace_id" {}
variable "secrecy_level_id" {}
variable "category_id" {}

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id       = var.workspace_id
  rule_type          = "CUSTOM"
  name               = "ruleName"
  secrecy_level_id   = var.secrecy_level_id
  category_id        = var.category_id
  method             = "REGULAR"
  content_expression = "^male$|^female&"
  column_expression  = "phoneNumber|email"
  comment_expression = ".*comment*."
  description        = "rule_description_custom_update1"
}
```

### BUILTIN RULE CHECK

```hcl
variable "workspace_id" {}
variable "name" {}
variable "secrecy_level_id" {}
variable "category_id" {}
variable "builtin_rule_id" {}

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id      = var.workspace_id
  rule_type         = "BUILTIN"
  name              = var.name
  secrecy_level_id  = var.secrecy_level_id
  category_id       = var.category_id
  builtin_rule_id   = var.builtin_rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of DataArts Studio workspace.
  Changing this parameter will create a new resource.

* `rule_type` - (Required, String, ForceNew) Spedifies the type of data recognition rule.
  The valid values are **CUSTOM** and **BUILTIN**.
  Changing this parameter will create a new data recognition rule.

* `name` - (Required, String, ForceNew) Specifies the rule name.
  If the value of `rule_type` is **BUILTIN**,
  then the value of `name` should correspond to the value of `builtin_rule_id`.
  Changing this parameter will create a new data recognition rule.

* `builtin_rule_id` - (Optional, String, ForceNew) Specifies the ID of built-in data recognition rule.
  The system has built-in 100+ sensitive data identification and desensitization rules,
  which can detect personal sensitive information (ID card, bank card, etc.),
  corporate sensitive information (business license number, tax registration certificate number, etc.),
  key sensitive information (PEM certificate, HEY private key, etc.),
  device sensitive information (IP address, MAC address, IPV6 address, etc.),
  location sensitive information (province, city, GPS location, address, etc.)
  and general sensitive information (date) Identify and desensitize sensitive information.
  The field `builtin_rule_id` is required when `rule_type` is **BUILTIN**.
  Changing this parameter will create a new data recognition rule.

* `secrecy_level_id` - (Required, String) Spedifies the ID of data secrecy level.
  Before operating on data, you need to define a confidentiality level for the data
  and describe the confidentiality level accordingly.
  The larger the number, the higher the level of confidentiality.

* `method` - (Optional, String) Specifies the method of custom data recognition rule.
  The valid values are **NONE** and **REGULAR**.
  If the value of `rule_type` is **CUSTOM**, the field `method` must be set.
  If the value of `method` is **REGULAR**, `content_expression`, `column_expression`,
  `comment_expression`, one of the three fields must be set.

* `content_expression` - (Optional, String) - Specifies the regular expression used to match the data content.
  For example **^male$|^female&**.

* `column_expression` - (Optional, String) Specifies the regular expression used to match the data column.
  This expression is used for both exact matching and fuzzy matching of field names.
  It supports multiple field matching currently.
  For example **phoneNumber|email**.

* `comment_expression` - (Optional, String) Specifies the regular expression used to match the data comment.
  For example **.*comment\*.**, It represents a fuzzy match on data comment.

* `category_id` - (Optional, String) Sprcifies the ID of data secrecy level category.
  Define data classifications for data of different values to better manage and measure
  your own data in groups, so that all types of groups have a parallel,
  equal and mutually exclusive relationship, making the data clearer.

* `description` - (Optional, String) Specifies the description of data recognition rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `created_by` - The creator.

* `updated_by` - The editor.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `secrecy_level` - The name of data secrecy level.
  It corresponds to `secrecy_level_id` one-to-one.

* `secrecy_level_num` - The level of data secrecy.The larger the secrecy level number,
  the higher the secrecy level. Currently, a maximum of 10 levels of confidentiality can be created.
  It corresponds to `secrecy_level_id` one-to-one.

* `enable` - Whether the current data recognition rule is available.

## Import

The DataArts Security data recognition rule can be imported using `<workspace_id>/<id>`, e.g.

```sh
terraform import huaweicloud_dataarts_security_data_recognition_rule.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `secrecy_level_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      secrecy_level_id,
    ]
  }
}
```
