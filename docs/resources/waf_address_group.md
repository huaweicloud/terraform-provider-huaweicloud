---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_address_group"
description: |-
  Manages a WAF address group resource within HuaweiCloud.
---

# huaweicloud_waf_address_group

Manages a WAF address group resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The address group resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable enterprise_project_id {}

resource "huaweicloud_waf_address_group" "test" {
  name                  = "example_address_name"
  description           = "example_description"
  ip_addresses          = ["192.168.1.0/24"]
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the address group. The value consists of `1` to `128` characters.
  Only letters, digits, hyphens (-), underscores (_), colons (:) and periods (.) are allowed.
  The name of each enterprise project by one account must be unique.

* `ip_addresses` - (Required, List) Specifies the IP addresses or IP address ranges.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF address group.
  Changing this parameter will create a new resource.
  For enterprise users, if omitted, default enterprise project will be used.

* `description` - (Optional, String) Specifies the description of the address group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `rules` - The list of rules that use the IP address group.
  The [rules](#AddressGroup_rules) structure is documented below.

<a name="AddressGroup_rules"></a>
The `rules` block supports:

* `rule_id` - The ID of rule.

* `rule_name` - The name of rule.

* `policy_id` - The ID of policy.

* `policy_name` - The name of policy.

## Import

There are two ways to import WAF address group state.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_address_group.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_address_group.test <id>/<enterprise_project_id>
```
