---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_certificate_associate_domains"
description: |-
  Use this resource to associate a certificate with multiple CDN domains within HuaweiCloud.
---

# huaweicloud_cdn_certificate_associate_domains

Use this resource to associate a certificate with multiple CDN domains within HuaweiCloud.

-> This resource is only a one-time action resource for bind certificate on the list of domains. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable target_domain_names {}
variable certificate_name {}

resource "huaweicloud_cdn_certificate_associate_domains" "test" {
  domain_names = var.target_domain_names
  https_switch = 1
  cert_name    = var.certificate_name
  certificate  = file("path/to/certificate.pem")
  private_key  = file("path/to/private.key")
}
```

### With Force Redirect Configuration

```hcl
variable target_domain_names {}
variable certificate_name {}

resource "huaweicloud_cdn_certificate_associate_domains" "test" {
  domain_names = var.target_domain_names
  https_switch = 1
  cert_name    = var.certificate_name
  certificate  = file("path/to/certificate.pem")
  private_key  = file("path/to/private.key")
  
  force_redirect_config {
    switch        = 1
    redirect_type = "https"
  }
  
  http2 = 1
}
```

## Argument Reference

The following arguments are supported:

* `domain_names` - (Required, String, NonUpdatable) The list of domain names to associate with the certificate.
  When there are multiple domains, separate them with comma (,).

* `https_switch` - (Required, Int, NonUpdatable) The HTTPS certificate configuration switch.
  The valid values are as follows:
  + **0**: Disable HTTPS
  + **1**: Enable HTTPS

* `access_origin_way` - (Optional, Int, NonUpdatable) The origin protocol configuration.
  The valid values are as follows:
  + **1**: same as user
  + **2**: HTTP protocol
  + **3**: HTTPS protocol

  Defaults to **2**.

* `force_redirect_https` - (Optional, Int, NonUpdatable) Whether to enable HTTPS force redirect to HTTPS to force clients
  to use HTTPS to access CDN PoPs. The valid values are as follows:
  + **0**: Disabled
  + **1**: Enabled

  Defaults to **0**.

* `force_redirect_config` - (Optional, List, NonUpdatable) Whether to force clients to use HTTPS when accessing CDN PoPs.
  The [force_redirect_config](#cdn_force_redirect_config)structure is documented below.

* `http2` - (Optional, Int, NonUpdatable) Whether to enable HTTP/2 to allow clients to use HTTP/2 when accessing CDN PoPs.
  The valid values are as follows:
  + **0**: Disabled
  + **1**: Enabled

  Defaults to **0**.

* `cert_name` - (Optional, String, NonUpdatable) The certificate name.

* `certificate` - (Optional, String, NonUpdatable) The SSL certificate content. The certificate chain cannot exceed 20 KB.
  + Only the PEM format is supported.
  + This parameter is optional if a certificate is not required.
  + This parameter is mandatory when a certificate is configured for the first time.
  + A complete certificate chain is required.

* `private_key` - (Optional, String, NonUpdatable) The private key of the SSL certificate.
  + Only the PEM format is supported.
  + This parameter is optional if a certificate is not required.
  + This parameter is mandatory when a certificate is configured for the first time.

* `certificate_type` - (Optional, Int, NonUpdatable) The certificate type. The valid values are as follows:
  + **0**: your certificate
  + **2**: SSL Certificate Manager (SCM) certificate

  Defaults to **0**.

<a name="cdn_force_redirect_config"></a>
The `force_redirect_config` block supports:

* `switch` - (Required, Int, NonUpdatable) Whether to enable force redirect to force clients to use HTTPS or HTTP to
  access CDN PoPs. The valid values are as follows:
  + **0**: Disabled
  + **1**: Enabled

* `redirect_type` - (Optional, String, NonUpdatable) The protocol to which requests are forcibly redirected.
  The valid values are as follows:
  + **http**: force redirect to HTTP
  + **https**: force redirect to HTTPS

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
