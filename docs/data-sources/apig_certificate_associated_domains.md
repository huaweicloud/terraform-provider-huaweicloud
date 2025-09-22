---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_certificate_associated_domains"
description: |-
  Use this data source to get domain list associated with the specified SSL certificate within HuaweiCloud.
---

# huaweicloud_apig_certificate_associated_domains

Use this data source to get domain list associated with the specified SSL certificate within HuaweiCloud.

## Example Usage

### Query all domains under the specified SSL certificate

```hcl
variable "certificate_id" {}

data "huaweicloud_apig_certificate_associated_domains" "test" {
  certificate_id = var.certificate_id
}
```

### Query domains under the specified SSL certificate by domain name

```hcl
variable "certificate_id" {}
variable "domain_name" {}

data "huaweicloud_apig_certificate_associated_domains" "test" {
  certificate_id = var.certificate_id
  url_domain     = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the domains are located.  
  If omitted, the provider-level region will be used.

* `certificate_id` - (Required, String) Specifies the ID of the certificate associated with the domains.

* `url_domain` - (Optional, String) Specifies the associated domain name to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domains` - All domains that match the filter parameters.  
  The [domains](#apig_data_certificate_associated_domains) structure is documented below.

<a name="apig_data_certificate_associated_domains"></a>
The `domains` block supports:

* `id` - The ID of the associated domain.

* `url_domain` - The associated domain name.

* `instance_id` - The ID of the dedicated instance to which the domain belongs.

* `status` - The CNAME resolution status of the domain name.
  + **1**: Not resolved.
  + **2**: Resolving.
  + **3**: Resolved.
  + **4**: Resolution failed.

* `min_ssl_version` - The minimum SSL protocol version of the domain.

* `verified_client_certificate_enabled` - Whether client certificate verification is enabled.

* `api_group_id` - The ID of the API group to which the domain belongs.

* `api_group_name` - The name of the API group to which the domain belongs.
