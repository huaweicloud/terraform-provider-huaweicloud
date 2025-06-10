---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_l7rule"
description: ""
---

# huaweicloud_elb_l7rule

Manages an ELB L7 Rule resource within HuaweiCloud.

## Example Usage

### Create by value

```hcl
variable l7policy_id {}

resource "huaweicloud_elb_l7rule" "l7rule_1" {
  l7policy_id  = var.l7policy_id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
```

### Create by conditions and type is HOST_NAME

```hcl
variable l7policy_id {}

resource "huaweicloud_elb_l7rule" "l7rule_1" {
  l7policy_id  = var.l7policy_id
  type         = "HOST_NAME"
  compare_type = "EQUAL_TO"

  conditions {
    value = "test.com"
  }
}
```

### Create by conditions and type is HEADER

```hcl
variable l7policy_id {}

resource "huaweicloud_elb_l7rule" "l7rule_1" {
  l7policy_id  = var.l7policy_id
  type         = "HEADER"
  compare_type = "EQUAL_TO"

  conditions {
    key   = "testKey"
    value = "testValue"
  }
}
```

### Create by conditions and type is SOURCE_IP

```hcl
variable l7policy_id {}

resource "huaweicloud_elb_l7rule" "l7rule_1" {
  l7policy_id  = var.l7policy_id
  type         = "SOURCE_IP"
  compare_type = "EQUAL_TO"

  conditions {
    value = "192.168.0.2/32"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the L7 Rule resource. If omitted, the
  provider-level region will be used. Changing this creates a new L7 Rule.

* `l7policy_id` - (Required, String, ForceNew) Specifies the ID of the L7 Policy. Changing this creates a new L7 Rule.

* `type` - (Required, String, ForceNew) Specifies the L7 Rule type. Value options:
  + **HOST_NAME**: A domain name will be used for matching.
  + **PATH**: A URL will be used for matching.
  + **METHOD**: An HTTP request method will be used for matching.
  + **HEADER**: The request header will be used for matching.
  + **QUERY_STRING**: A query string will be used for matching.
  + **SOURCE_IP**: The source IP address will be used for matching.
  + **COOKIE**: The cookie will be used for matching.
  
  Changing this creates a new L7 Rule.

  -> **NOTE:** If `type` is set to **HOST_NAME**, **PATH**, **METHOD**, or **SOURCE_IP**, only one forwarding rule can
  be created for each type.

* `compare_type` - (Required, String) Specifies how requests are matched with the forwarding rule. Value options:
  + **EQUAL_TO**: Exact match.
  + **REGEX**: Regular expression match.
  + **STARTS_WITH**: Prefix match.
  
  Instructions for use:
  + If `type` is set to **HOST_NAME**, the value can only be **EQUAL_TO**, and asterisks (*) can be used as wildcard
    characters.
  + If `type` is set to **PATH**, the value can be **REGEX**, **STARTS_WITH**, or **EQUAL_TO**.
  + If `type` is set to **METHOD** or **SOURCE_IP**, the value can only be **EQUAL_TO**.
  + If `type` is set to **HEADER** or **QUERY_STRING**, the value can only be **EQUAL_TO**, asterisks (*) and question
    marks (?) can be used as wildcard characters.

* `value` - (Optional, String) Specifies the value of the match content. This parameter is valid only when `conditions`
  are left blank.
  + If `type` is set to **HOST_NAME**, the value can contain letters, digits, hyphens (-), periods (.), and
    asterisks (\*) and must start with a letter or digit. If you want to use a wildcard domain name, enter an
    asterisk (\*) as the leftmost label of the domain name.
  + If `type` is set to **PATH** and `compare_type` to **STARTS_WITH** or **EQUAL_TO**, the value can contain only
    letters, digits, and special characters _~';@^-%#&$.*+?,=!:|/()[]{}.
  + If `type` is set to **METHOD**, **SOURCE_IP**, **HEADER**, or **QUERY_STRING**, this parameter will not take effect,
    and `conditions` will be used to specify the key and value.

* `conditions` - (Optional, List) Specifies the matching conditions of the forwarding rule. This parameter is available
  only when `advanced_forwarding_enabled` of the listener is set to **true**. If it is specified, parameter `value` will
  not take effect, and the value will contain all conditions configured for the forwarding rule. The keys in the list
  must be the same, whereas each value must be unique.
  The [condition](#conditions) structure is documented below.

<a name="conditions"></a>
The `condition` block supports:

* `key` - (Optional, String) Specifies the key of match item.
  + If `type` is set to **HOST_NAME**, **PATH**, **METHOD**, or **SOURCE_IP**, this parameter is left blank.
  + If `type` is set to **HEADER**, it indicates the name of the HTTP header parameter. It can contain 1 to 40
    characters, including letters, digits, hyphens (-), and underscores (_).
  + If `type` is set to **QUERY_STRING**, it indicates the name of the query parameter. It is case-sensitive and can
    contain 1 to 128 characters. Spaces, square brackets ([]), curly brackets ({}), angle brackets (<>), backslashes (),
    double quotation marks (" "), pound signs (#), ampersands (&), vertical bars (|), percent signs (%), and tildes (~)
    are not supported.

  -> **NOTE:** All keys in the conditions list in the same rule must be the same.

* `value` - (Required, String) Specifies the value of the match item.
  + If `type` is set to **HOST_NAME**, it indicates the domain name, which can contain 1 to 128 characters, including
    letters, digits, hyphens (-), periods (.), and asterisks (), and must start with a letter, digit, or asterisk ().
    If you want to use a wildcard domain name, enter an asterisk (*) as the leftmost label of the domain name.
  + If `type` is set to **PATH**, it indicates the request path, which can contain 1 to 128 characters. If
    `compare_type` is set to **STARTS_WITH** or **EQUAL_TO** for the forwarding rule, the value must start with a
    slash (/) and can contain only letters, digits, and special characters _~';@^-%#&$.*+?,=!:|/()[]{}.
  + If `type` is set to **HEADER**, it indicates the value of the HTTP header parameter. The value can contain 1 to 128
    characters. Asterisks (*) and question marks (?)are allowed, but spaces and double quotation marks are not allowed.
    An asterisk can match zero or more characters, and a question mark can match 1 character.
  + If `type` is set to **QUERY_STRING**, it indicates the value of the query parameter. The value is case-sensitive
    and can contain 1 to 128 characters. Spaces, square brackets ([]), curly brackets ({}), angle brackets (<>),
    backslashes (), double quotation marks (""), pound signs (#), ampersands (&), vertical bars (|), percent signs (%),
    and tildes (~) are not supported. Asterisks (*)and question marks (?) are allowed. An asterisk can match zero or
    more characters, and a question mark can match 1 character.
  + If `type` is set to **METHOD**, it indicates the HTTP method. The value can be **GET**, **PUT**, **POST**,
    **DELETE**, **PATCH**, **HEAD**, or **OPTIONS**.
  + If `type` is set to **SOURCE_IP**, it indicates the source IP address of the request. The value is an **IPv4** or
    **IPv6** CIDR block, for example, 192.168.0.2/32 or 2049::49/64.

  -> **NOTE:** All values in the conditions list in the same rule must be unique.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the L7 Rule.

* `created_at` - The create time of the L7 Rule.

* `updated_at` - The update time of the L7 Rule.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB L7 rule can be imported using the `l7policy_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_elb_rule.rule_1 <l7policy_id>/<id>
```
