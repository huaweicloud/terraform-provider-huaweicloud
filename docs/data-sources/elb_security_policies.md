---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_security_policies"
description: |-
  Use this data source to get the list of ELB security policies.
---

# huaweicloud_elb_security_policies

Use this data source to get the list of ELB security policies.

## Example Usage

```hcl
variable "security_policies_name" {}

data "huaweicloud_elb_security_policies" "test" {
  name = var.security_policies_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the security policy.

* `type` - (Optional, String) Specifies the type of the security policy. Value options: **system**, **custom**.

* `security_policy_id` - (Optional, String) Specifies the ID of the security policy.

* `description` - (Optional, String) Specifies the description of the security policy.

* `protocol` - (Optional, String) Specifies the TLS protocol supported by the security policy.

* `cipher` - (Optional, String) Specifies the cipher suite supported by the security policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `security_policies` - Lists the security policies.
  The [security_policies](#Elb_security_policies) structure is documented below.

<a name="Elb_security_policies"></a>
The `security_policies` block supports:

* `id` - The ID of the security policy.

* `type` - The type of the security policy.

* `name` - The name of the security policy.

* `description` - The description of the security policy.

* `listeners` - The IDs of listeners with which the security policy is associated.
  The [listeners](#Elb_security_policy_listeners) structure is documented below.

* `protocols` - The TLS protocols supported by the security policy.

* `ciphers` - The cipher suites supported by the security policy.

* `created_at` - The time when the custom security policy was created.

* `updated_at` - The time when the custom security policy was updated.

<a name="Elb_security_policy_listeners"></a>
The `listeners` block supports:

* `id` - The listener ID.
