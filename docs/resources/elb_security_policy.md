---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_security_policy"
description: ""
---

# huaweicloud_elb_security_policy

Manages an ELB security policy resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_elb_security_policy" "test" {
  name        = "security_policy_test"
  description = "this is a security policy"
  protocols   = ["TLSv1","TLSv1.1","TLSv1.2","TLSv1.3"]
  ciphers     = ["ECDHE-RSA-AES256-GCM-SHA384","ECDHE-RSA-AES128-GCM-SHA256"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `protocols` - (Required, List) Specifies the TSL protocol list which the security policy select.
  Value options: **TLSv1**, **TLSv1.1**, **TLSv1.2**, and **TLSv1.3**.

* `ciphers` - (Required, List) Specifies the cipher suite list of the security policy.
  The protocol and cipher suite must match. That is to say, there must be at least one cipher suite in
  ciphers that matches the protocol. The following cipher suites are supported:
  **ECDHE-RSA-AES256-GCM-SHA384**, **ECDHE-RSA-AES128-GCM-SHA256**, **ECDHE-ECDSA-AES256-GCM-SHA384**,
  **ECDHE-ECDSA-AES128-GCM-SHA256**, **AES128-GCM-SHA256**, **AES256-GCM-SHA384**, **ECDHE-ECDSA-AES128-SHA256**,
  **ECDHE-RSA-AES128-SHA256**, **AES128-SHA256**, **AES256-SHA256**, **ECDHE-ECDSA-AES256-SHA384**,
  **ECDHE-RSA-AES256-SHA384**, **ECDHE-ECDSA-AES128-SHA**, **ECDHE-RSA-AES128-SHA**, **ECDHE-RSA-AES256-SHA**,
  **ECDHE-ECDSA-AES256-SHA**, **AES128-SHA**, **AES256-SHA**, **CAMELLIA128-SHA**, **DES-CBC3-SHA**,
  **CAMELLIA256-SHA**, **ECDHE-RSA-CHACHA20-POLY1305**, **ECDHE-ECDSA-CHACHA20-POLY1305**, **TLS_AES_128_GCM_SHA256**,
  **TLS_AES_256_GCM_SHA384**, **TLS_CHACHA20_POLY1305_SHA256**, **TLS_AES_128_CCM_SHA256**,
  **TLS_AES_128_CCM_8_SHA256**.

* `name` - (Optional, String) Specifies the ELB security policy name.
  The name contains only Chinese characters, letters, digits, underscores (_), and hyphens (-),
  and cannot exceed 255 characters.

* `description` - (Optional, String) Specifies the description of the ELB security policy.
  The value can contain 0 to 255 characters.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the Enterprise
  router belongs.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time of the security policy.

* `updated_at` - The update time of the security policy.

* `listeners` - The listener which the security policy associated with.
  The [ListenerRef](#SecurityPoliciesV3_ListenerRef) structure is documented below.

<a name="SecurityPoliciesV3_ListenerRef"></a>
The `ListenerRef` block supports:

* `id` - The listener id.

## Import

The elb security policies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_security_policy.test 0ce123456a00f2591fabc00385ff1234
```
