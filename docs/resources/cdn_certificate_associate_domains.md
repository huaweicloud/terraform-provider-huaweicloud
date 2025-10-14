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
  + **1**: Protocol follow
  + **2**: HTTP protocol (default)
  + **3**: HTTPS protocol

* `force_redirect_https` - (Optional, Int, NonUpdatable) Whether to enable HTTPS force redirect.  
  The valid values are as follows:
  + **0**: Disable (default)
  + **1**: Enable

* `force_redirect_config` - (Optional, List, NonUpdatable) The force redirect configuration.  
  The [force_redirect_config](#cdn_force_redirect_config)structure is documented below.

* `http2` - (Optional, Int, NonUpdatable) The HTTP/2 protocol switch.  
  The valid values are as follows:
  + **0**: Disable (default)
  + **1**: Enable

* `cert_name` - (Optional, String, NonUpdatable) The certificate name.

* `certificate` - (Optional, String, NonUpdatable, Sensitive) The SSL certificate content in PEM format.

* `private_key` - (Optional, String, NonUpdatable, Sensitive) The SSL certificate private key content in PEM format.

* `certificate_type` - (Optional, Int, NonUpdatable) The certificate type.  
  The valid values are as follows:
  + **0**: Free certificate (default)
  + **1**: Paid certificate

<a name="cdn_force_redirect_config"></a>
The `force_redirect_config` block supports:

* `switch` - (Required, Int, NonUpdatable) The force redirect switch.  
  The valid values are as follows:
  + **0**: Disable
  + **1**: Enable

* `redirect_type` - (Optional, String, NonUpdatable) The redirect type.  
  The valid values are as follows:
  + **http**: Redirect to HTTP
  + **https**: Redirect to HTTPS

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
