---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_ssl_certificates"
description: |-
  Use this data source to get SSL certificate list associated with the specified instance within HuaweiCloud.
---

# huaweicloud_apig_instance_ssl_certificates

Use this data source to get SSL certificate list associated with the specified instance within HuaweiCloud.

## Example Usage

### Querying all SSL certificates associated with the specified instance

```hcl
variable "instance_id" {}

data "huaweicloud_apig_instance_associated_ssl_certificates" "test" {
  instance_id = var.instance_id
}
```

### Querying SSL certificates associated with the instance using specified common domain name

```hcl
variable "instance_id" {}
variable "common_domain_name" {}

data "huaweicloud_apig_instance_associated_ssl_certificates" "test" {
  instance_id = var.instance_id
  common_name = var.common_domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the SSL certificates belong.

* `name` - (Optional, String) Specifies the name of the SSL certificate.

* `common_name` - (Optional, String) Specifies the domain name of the SSL certificate.

* `signature_algorithm` - (Optional, String) Specifies the signature algorithm of the SSL certificate.

* `type` - (Optional, String) Specifies the visibility range of the SSL certificate.

* `algorithm_type` - (Optional, String) Specifies the algorithm type of the SSL certificate(RSA, ECC, SM2).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - All associated SSL certificates that match the filter parameters.  
  The [certificates](#instance_associated_ssl_certificates) structure is documented below.

<a name="instance_associated_ssl_certificates"></a>
The `certificates` block supports:

* `instance_id` - The ID of the dedicated instance to which the SSL certificates belong.

* `project_id` - The ID of the region project to which the SSL certificates belong.

* `name` - The name of the SSL certificate.

* `type` - The type of the SSL certificate.

* `san` - The san extended domain of the SSL certificate.

* `signature_algorithm` - The signature algorithm of the SSL certificate.

* `algorithm_type` - The algorithm type of the SSL certificate(RSA, ECC, SM2).

* `is_has_trusted_root_ca` - The certificate has trusted root certificate authority or not.

* `not_after` - The expiration date of the SSL certificate.

* `create_time` - The create time of the SSL certificate.

* `update_time` - The update time of the SSL certificate.
