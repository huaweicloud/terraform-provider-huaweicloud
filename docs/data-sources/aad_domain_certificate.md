---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_domain_certificate"
description: |-
  Use this data source to get the Advanced Anti-DDos domain certificate within HuaweiCloud.
---

# huaweicloud_aad_domain_certificate

Use this data source to get the Advanced Anti-DDos domain certificate within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_domain_certificate" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_id` - (Required, String) Specifies the AAD domain ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domain_name` - The domain name.

* `cert_info` - The certificate information.

  The [cert_info](#cert_info_struct) structure is documented below.

<a name="cert_info_struct"></a>
The `cert_info` block supports:

* `expire_time` - The certificate expiration time.

* `expire_status` - The certificate expiration status.

* `cert_name` - The certificate name.

* `id` - The certificate ID.

* `apply_domain` - The domain name that the certificate applies to.
