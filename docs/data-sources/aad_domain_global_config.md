---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_domain_global_config"
description: |-
  Use this data source to get the Advanced Anti-DDos domain global config.
---

# huaweicloud_aad_domain_global_config

Use this data source to get the Advanced Anti-DDos domain global config.

## Example Usage

```hcl
data "huaweicloud_aad_domain_global_config" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `support_tls` - The list of supported TLS versions.

* `cipher` - The cipher suites.

  The [cipher](#cipher_struct) structure is documented below.

<a name="cipher_struct"></a>
The `cipher` block supports:

* `name` - The suite name.

* `algo` - The cryptographic algorithm.

* `desc_cn` - The Chinese description.

* `desc_en` - The English description.
