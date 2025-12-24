---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_create_ignore_rules"
description: |-
  Manages a resource to batch create WAF ignore rules (whitelist rules) within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_create_ignore_rules

Manages a resource to batch create WAF ignore rules (whitelist rules) within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch creating ignore rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_ids" {
  type = list(string)
}
variable "domain" {
  type = list(string)
}

resource "huaweicloud_waf_batch_create_ignore_rules" "test" {
  rule                  = "xss"
  description           = "test description"
  policy_ids            = var.policy_ids
  enterprise_project_id = var.enterprise_project_id
  domain                = var.domain

  conditions {
    category        = "url"
    contents        = ["/admin"]
    logic_operation = "equal"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain` - (Required, List, NonUpdatable) Specifies the protected domain name or website.
  When the array length is `0`, the rule applies to all domain names or protected websites.
  When the access mode of the protected domain name is cloud mode-ELB access, this parameter must be filled in the
  format `<domain name>:[/topic/body/section/table/tgroup/tbody/row/entry/p/br {""}) (br]` (e.g., `www.example.com:id`).
  If all listeners under the load balancer (ELB) bound to the domain name are connected to WAF protection, the ID entered
  should be the load balancer (ELB) ID; otherwise, the ID entered should be the specified listener ID.
  You can find the ELB instance ID bound to the domain name on the domain details page of the WAF console, and find its
  listener ID under the ELB-side listener tab.

* `conditions` - (Required, List, NonUpdatable) Specifies the matching conditions for the rule.
  The [conditions](#Rule_conditions) structure is documented below.

* `rule` - (Required, String, NonUpdatable) Specifies the rule to be ignored. You can block one or more characters;
  when blocking multiple characters, use a half-width character (;) to separate them.

  + If you want to block a specific built-in rule, the value of this parameter is the rule ID.
  To query the rule ID, go to the WAF console, choose **Policies** and click the target policy name. On the displayed
  page, in the **Basic Web Protection** area, select the **Protection Rules** tab, and view the ID of the specific rule.
  You can also query the rule ID in the event details.

  + If you want to mask a type of basic web protection rules, set this parameter to the name of the type of basic web
  protection rules. Valid values are: **xss**(XSS attacks), **webshell**(Web shells), **vuln**(Other types of attacks),
  **sqli**(SQL injection attack), **robot**(Malicious crawlers), **rfi**(Remote file inclusion),
  **lfi**(Local file inclusion), **cmdi**(Command injection attack).

  + To bypass the basic web protection, set this parameter to **all**.

  + To bypass all WAF protection, set this parameter to **bypass**.

* `policy_ids` - (Required, List, NonUpdatable) Specifies the list of policy IDs to which the rule will be applied.

* `advanced` - (Optional, List, NonUpdatable) Specifies the advanced matching conditions for the rule.
  The [advanced](#Rule_advanced) structure is documented below.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the rule.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

<a name="Rule_advanced"></a>
The `advanced` block supports:

* `index` - (Optional, String, NonUpdatable) Specifies the field type.
  Supported field types include: **params**, **cookie**, **header**, **body**, and **multipart**.

* `contents` - (Optional, List, NonUpdatable) Specifies the content of the condition.

<a name="Rule_conditions"></a>
The `conditions` block supports:

* `category` - (Required, String, NonUpdatable) Specifies the field type.
  The value can be: **ip**, **url**, **params**, **cookie**, or **header**.

* `contents` - (Required, List, NonUpdatable) Specifies the content of the condition.
  The array length is limited to `1`. The content format varies depending on the field type.
  For example, when the `category` is **ip**, the content format must be an IP address or a range of IP addresses;
  when the `category` is **url**, the content format must be a standard URL;
  when the `category` is **params**, **cookie**, or **header**, there are no restrictions on the content format.

* `logic_operation` - (Required, String, NonUpdatable) Specifies the matching logic.
  The matching logic varies depending on the field type.
  When the `category` is **ip**, the matching logic supports **equal** and **not_equal**.
  When the `category` is **url**, **header**, **params**, or **cookie**, the matching logic supports **equal**,
  **not_equal**, **contain**, **not_contain**, **prefix**, **not_prefix**, **suffix**, **not_suffix**, **regular_match**,
  and **regular_not_match**.

* `check_all_indexes_logic` - (Optional, Int, NonUpdatable) Specifies how to check subfields.
  When using custom subfields or `category` is **url** or **ip**, the `check_all_indexes_logic` parameter is not required.
  In other cases, the value can be:
  + `1`: Check all subfields
  + `2`: Check any subfield

* `index` - (Optional, String, NonUpdatable) Specifies the subfield name.
  When the `category` is **ip** and the subfield is the client's IP, the `index` parameter is not required.
  When the subfield type is **X-Forwarded-For**, the value is **x-forwarded-for**;
  when the `category` is **params**, **header**, or **cookie** and the subfield is custom, the value of `index` is the
  custom subfield.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
