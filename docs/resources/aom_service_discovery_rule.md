---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_service_discovery_rule"
description: ""
---

# huaweicloud_aom_service_discovery_rule

Manages an AOM service discovery rule resource within HuaweiCloud.

## Example Usage

### Basic example

```hcl
resource "huaweicloud_aom_service_discovery_rule" "discovery_rule" {
  name                   = "test-rule"
  priority               = 9999
  detect_log_enabled     = "true"
  discovery_rule_enabled = "true"
  is_default_rule        = "false"
  log_file_suffix        = ["log"]
  service_type           = "Java"

  discovery_rules {
    check_content = ["java"]
    check_mode    = "contain"
    check_type    = "cmdLine"
  }

  log_path_rules {
    name_type = "cmdLineHash"
    args      = ["java"]
    value     = ["/tmp/log"]
  }

  name_rules {
    service_name_rule {
      name_type = "str"
      args      = ["java"]
    }
    application_name_rule {
      name_type = "str"
      args      = ["java"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the service discovery rule resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the rule name, which contains `4` to `63` characters. It must start
  with a lowercase letter but cannot end with a hyphen (-). Only digits, lowercase letters, and hyphens are allowed.

* `service_type` - (Required, String) Specifies the service type, which is used only for rule classification and UI display.
  You can enter any field. For example, enter Java or Python by technology stack, or enter collector or database by function.

* `discovery_rules` - (Required, List) Specifies the discovery rule. If the array contains multiple conditions, only the
  processes that meet all the conditions will be matched. If the value of `check_type` is **cmdLine**, set the value of
  `check_mode` to **contain**. `check_content` is in the format of ["xxx"], indicating that the process must contain
  the xxx parameter. If the value of `check_type` is **env**, set the value of `check_mode` to **contain**.
  `check_content` is in the format of ["k1","v1"], indicating that the process must contain the environment variable
  whose name is k1 and value is v1. If the value of `check_type` is **scope**, set the value of `check_mode`
  to **equals**. `check_content` is in the format of ["hostId1","hostId2"], indicating that the rule takes effect only
  on specified nodes. If no nodes are specified, the rule applies to all nodes of the project.
  + `check_type` - (Required, String) Specifies the match type. The values can be **cmdLine**, **env** and **scope**.
  + `check_mode` - (Required, String) Specifies the match condition. The values can be **contain** and **equals**.
  + `check_content` - (Required, List) Specifies the matched value. This is a list of strings.

* `name_rules` - (Required, List) Specifies the naming rules for discovered services and applications.
  The [object](#name_rules_object) structure is documented below.

* `log_file_suffix` - (Required, List) Specifies the log file suffix. This is a list of strings.
  The values can be: **log**, **trace**, and **out**.

* `detect_log_enabled` - (Optional, Bool) Specifies whether to enable log collection. The default value is true.

* `priority` - (Optional, Int) Specifies the rule priority. Value range: 1 to 9999. The default value is 9999.

* `discovery_rule_enabled` - (Optional, Bool) Specifies whether the rule is enabled. The default value is true.

* `is_default_rule` - (Optional, Bool) Specifies whether the rule is the default one. The default value is false.

* `log_path_rules` - (Optional, List) Specifies the log path configuration rule. If cmdLineHash is a fixed string,
  logs in the specified log path or log file are collected. Otherwise, only the files whose names end with
  .log or .trace are collected. If the value of `name_type` is **cmdLineHash**, args is in the format of ["00001"] and
  value is in the format of ["/xxx/xx.log"], indicating that the log path is /xxx/xx.log when the startup command is 00001.
  + `name_type` - (Required, String) Specifies the value type, which can be **cmdLineHash**.
  + `args` - (Required, List) Specifies the command. This is a list of strings.
  + `value` - (Required, List) Specifies the log path. This is a list of strings.

* `description` - (Optional, String) Specifies the rule description.

<a name="name_rules_object"></a>
The `name_rules` block supports:

* `service_name_rule` - (Required, List) Specifies the service name rule. If there are multiple objects in the array,
  the character strings extracted from these objects constitute the service name. If the value of `name_type` is
  **cmdLine**, `args` is in the format of ["start", "end"], indicating that the characters between start and end
  in the command are extracted. If the value of `name_type` is **env**, `args` is in the format of ["aa"],
  indicating that the environment variable named aa is extracted. If the value of `name_type` is **str**, `args` is in the
  format of ["fix"], indicating that the service name is suffixed with fix. If the value of `name_type` is
  **cmdLineHash**, `args` is in the format of ["0001"] and `value` is in the format of ["ser"], indicating that the
  service name is ser when the startup command is 0001. The [object](#basic_name_rule_object) structure is
  documented below.

* `application_name_rule` - (Required, List) Specifies the application name rule. If the value of `name_type` is
  **cmdLine**, `args` is in the format of ["start", "end"], indicating that the characters between start and end in
  the command are extracted. If the value of `name_type` is **env**, `args` is in the format of ["aa"], indicating that
  the environment variable named aa is extracted. If the value of `name_type` is **str**, `args` is in the format of
  ["fix"], indicating that the application name is suffixed with fix. If the value of `name_type` is **cmdLineHash**,
  `args` is in the format of ["0001"] and `value` is in the format of ["ser"], indicating that the application name is
  ser when the startup command is 0001. The [object](#basic_name_rule_object) structure is documented below.

<a name="basic_name_rule_object"></a>
The `service_name_rule` block and `application_name_rule` block support:

* `name_type` - (Required, String) Specifies the value type. The value can be **cmdLineHash**, **cmdLine**, **env**
and **str**.

* `args` - (Required, List) Specifies the input value.

* `value` - (Optional, List) Specifies the application name, which is mandatory only when the value of `name_type` is
  **cmdLineHash**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID of the service discovery rule. The value is the rule name.

* `rule_id` - The rule ID in uuid format.

* `created_at` - The rule create time.

## Import

AOM service discovery rules can be imported using the `name`, e.g.

```bash
$ terraform import huaweicloud_aom_service_discovery_rule.alarm_rule <name>
```
