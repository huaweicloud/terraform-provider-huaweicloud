---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_resolver_rules"
description: |-
  Use this data source to get the list of DNS resolver rules.
---

# huaweicloud_dns_resolver_rules

Use this data source to get the list of DNS resolver rules.

## Example Usage

```hcl
data "huaweicloud_dns_resolver_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Optional, String) Specifies the domain name of the endpoint rule to be queried.

* `name` - (Optional, String) Specifies the name of the endpoint rule to be queried.

* `endpoint_id` - (Optional, String) Specifies the endpoint ID.

* `resolver_rule_id` - (Optional, String) Specifies the ID of an endpoint rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resolver_rules` - Indicates the list of endpoint rules.

  The [resolver_rules](#resolver_rules_struct) structure is documented below.

<a name="resolver_rules_struct"></a>
The `resolver_rules` block supports:

* `id` - Indicates the ID of an endpoint rule.

* `name` - Indicates the rule name.

* `endpoint_id` - Indicates the ID of the endpoint to which the current rule belongs.

* `status` - Indicates the resource status.

* `rule_type` - Indicates the rule type.

* `routers` - Indicates the VPC associated with the endpoint rule.

  The [routers](#resolver_rules_routers_struct) structure is documented below.

* `domain_name` - Indicates the domain name.

* `ipaddress_count` - Indicates the number of IP addresses in the endpoint rule.

* `create_time` - Indicates the creation time. Format is **yyyy-MM-dd'T'HH:mm:ss.SSS**.

* `update_time` - Indicates the update time. Format is **yyyy-MM-dd'T'HH:mm:ss.SSS**.

<a name="resolver_rules_routers_struct"></a>
The `routers` block supports:

* `router_id` - Indicates the ID of the associated VPC.

* `router_region` - Indicates the region where the associated VPC is located.

* `status` - Indicates the resource status.
