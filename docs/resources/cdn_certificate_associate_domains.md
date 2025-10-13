---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_certificate_associate_domains"
description: |-
  Associates a certificate with multiple CDN domains within HuaweiCloud.
---

# huaweicloud_cdn_certificate_associate_domains

Associates a certificate with multiple CDN domains within HuaweiCloud. This is an action resource that performs a one-time operation to configure HTTPS certificates for multiple domains.

## Example Usage

### Associate a certificate with multiple domains

```hcl
resource "huaweicloud_cdn_certificate_associate_domains" "test" {
  domain_names = "example1.com,example2.com,example3.com"
  https_switch  = 1

  cert_name   = "my-certificate"
  certificate = "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"

  access_origin_way = 2
  http2            = 1
}
```

### Associate a certificate with force redirect configuration

```hcl
resource "huaweicloud_cdn_certificate_associate_domains" "test" {
  domain_names = "example1.com,example2.com"
  https_switch  = 1

  cert_name   = "my-certificate"
  certificate = "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"

  force_redirect_config {
    switch        = 1
    redirect_type = "https"
  }

  access_origin_way = 3
  http2            = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the CDN service is located. If omitted, the provider-level region will be used.

* `domain_names` - (Required, String, ForceNew) The domain names to associate with the certificate, separated by commas. Maximum 50 domains can be configured.

* `https_switch` - (Required, Int, ForceNew) The HTTPS certificate configuration switch. Valid values are:
  * `0`: Disable HTTPS
  * `1`: Enable HTTPS

* `access_origin_way` - (Optional, Int, ForceNew) The origin protocol configuration. Valid values are:
  * `1`: Protocol follow
  * `2`: HTTP protocol (default)
  * `3`: HTTPS protocol

* `force_redirect_https` - (Optional, Int, ForceNew) Whether to enable HTTPS force redirect. Valid values are:
  * `0`: Disable (default)
  * `1`: Enable

* `force_redirect_config` - (Optional, List, ForceNew) The force redirect configuration. The structure is documented below.

* `http2` - (Optional, Int, ForceNew) The HTTP/2 switch. Valid values are:
  * `0`: Disable (default)
  * `1`: Enable

* `cert_name` - (Optional, String, ForceNew) The certificate name. Required when `https_switch` is `1`.

* `certificate` - (Optional, String, ForceNew) The SSL certificate content in PEM format. Required when `https_switch` is `1`.

* `private_key` - (Optional, String, ForceNew) The SSL certificate private key content in PEM format. Required when `https_switch` is `1`.

* `certificate_type` - (Optional, Int, ForceNew) The certificate type. Valid values are:
  * `0`: Free certificate (default)
  * `1`: Paid certificate

The `force_redirect_config` block supports:

* `switch` - (Required, Int) The force redirect switch. Valid values are:
  * `0`: Disable
  * `1`: Enable

* `redirect_type` - (Optional, String) The redirect type. Valid values are:
  * `http`: Redirect to HTTP
  * `https`: Redirect to HTTPS

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The domain names string used as the resource ID.

## Import

This resource can be imported using the domain names:

```bash
$ terraform import huaweicloud_cdn_certificate_associate_domains.test "example1.com,example2.com"
```

## Notes

* This is an action resource that performs a one-time operation. Once created, it cannot be updated or deleted.
* The certificate content must be in PEM format.
* When `https_switch` is set to `1`, the `cert_name`, `certificate`, and `private_key` parameters are required.
* The maximum number of domains that can be configured in a single operation is 50.
* If a domain already has an HTTPS certificate configured, the new certificate will override the existing one.
