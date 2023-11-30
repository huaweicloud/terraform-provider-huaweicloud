---
subcategory: "Web Application Firewall (WAF)"
---

# huaweicloud_waf_domain

Manages a WAF domain resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The domain name resource can be used in Cloud Mode.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_waf_certificate" "certificate_1" {
  name                  = "cert_1"
  enterprise_project_id = var.enterprise_project_id
  
  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIFmQl5dh2QUAeo39TIKtadgAgh4zHx09kSgayS9Wph9LEqq7MA+2042L3J9aOa
DAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUR+SosWwALt6PkP0J9iOIxA6RW8gVsLwq
...
+HhDvD/VeOHytX3RAs2GeTOtxyAV5XpKY5r+PkyUqPJj04t3d0Fopi0gNtLpMF=
-----END CERTIFICATE-----
EOT
  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIJwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAM
ATAwMC4GCCsGAQUFBwIBFiJodHRwOi8vY3BzLnJvb3QteDEubGV0c2VuY3J5cHQu
...
he8Y4IWS6wY7bCkjCWDcRQJMEhg76fsO3txE+FiYruq9RUWhiF1myv4Q6W+CyBFC
1qoJFlcDyqSMo5iHq3HLjs
-----END PRIVATE KEY-----
EOT
}

resource "huaweicloud_waf_domain" "domain_1" {
  domain                = "www.example.com"
  certificate_id        = huaweicloud_waf_certificate.certificate_1.id
  certificate_name      = huaweicloud_waf_certificate.certificate_1.name
  proxy                 = true
  enterprise_project_id = var.enterprise_project_id

  custom_page {
    http_return_code = "404"
    block_page_type  = "application/json"
    page_content     = <<EOF
{
  "event_id": "$${waf_event_id}",
  "error_msg": "error message"
}
EOF
  }

  timeout_settings {
    connection_timeout = 100
    read_timeout       = 1000
    write_timeout      = 1000
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.13"
    port            = "8080"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the WAF domain resource. If omitted, the
  provider-level region will be used. Changing this setting will push a new certificate.

* `domain` - (Required, String, ForceNew) Specifies the domain name to be protected. For example, `www.example.com` or
  `*.example.com`. Changing this creates a new domain.

* `server` - (Required, List) Specifies an array of origin web servers. The object structure is documented below.

* `certificate_id` - (Optional, String) Specifies the certificate ID. This parameter is mandatory when `client_protocol`
  is set to HTTPS.

* `certificate_name` - (Optional, String) Specifies the certificate name. This parameter is mandatory
  when `client_protocol` is set to HTTPS.

* `policy_id` - (Optional, String) Specifies the policy ID associated with the domain. If not specified, a new
  policy will be created automatically.

* `keep_policy` - (Optional, Bool) Specifies whether to retain the policy when deleting a domain name.
  Defaults to true.
  
* `proxy` - (Optional, Bool) Specifies whether a proxy is configured.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the domain. Valid values are *prePaid*
  and *postPaid*, defaults to *prePaid*. Changing this creates a new instance.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF domain.
  Changing this parameter will create a new resource.

* `custom_page` - (Optional, List) Specifies the custom page. Only supports one custom alarm page.
  The [custom_page](#Domain_custom_page) structure is documented below.

* `redirect_url` - (Optional, String) Specifies the URL of the redirected page. The root domain name of the redirection
  address must be the name of the currently protected domain (including a wildcard domain name).
  The available **${http_host}** can be used to indicate the currently protected domain name and port.
  For example: **${http_host}/error.html**.

-> The fields `redirect_url` and `custom_page` are mutually exclusive and cannot be specified simultaneously.

* `http2_enable` - (Optional, Bool) Specifies whether to use the http2 protocol.
  This field is only used for communication between clients and WAF.
  Things to note when using this field are as follows:
  + There must be at least one server configuration with client protocol set to `HTTPS`, or this configuration is unable
    to work.
  + This field cannot not work if the client supports TLS 1.3.
  + This field can work only when the client supports TLS 1.2 or earlier versions.
  + If you want to use HTTP/2 forwarding, use a dedicated WAF instance.

  Defaults to **false**.

* `ipv6_enable` - (Optional, Bool) Specifies whether IPv6 protection is enabled.
  Enable IPv6 protection if the domain name is accessible using an IPv6 address.
  After you enable it, WAF assigns an IPv6 address to the domain name.
  Defaults to **false**.

* `timeout_settings` - (Optional, List) Specifies the timeout setting. Only supports one timeout setting.
  The [timeout_settings](#Domain_timeout_settings) structure is documented below.

The `server` block supports:

* `client_protocol` - (Required, String) Protocol type of the client. The options include `HTTP` and `HTTPS`.

* `server_protocol` - (Required, String) Protocol used by WAF to forward client requests to the server. The options
  include `HTTP` and `HTTPS`.

* `address` - (Required, String) IP address or domain name of the web server that the client accesses. For example,
  `192.168.1.1` or `www.a.com`.

* `port` - (Required, Int) Port number used by the web server. The value ranges from 0 to 65535, for example, 8080.

<a name="Domain_custom_page"></a>
The `custom_page` block supports:

* `http_return_code` - (Required, String) Specifies the HTTP return code.
  The value can be a positive integer in the range of 200-599 except 408, 444 and 499.

* `block_page_type` - (Required, String) Specifies the content type of the custom alarm page.
  The value can be **text/html**, **text/xml** or **application/json**.

* `page_content` - (Required, String) Specifies the page content. The page content based on the selected page type.
  The available **${waf_event_id}** in the page content indicates an event ID, and only one **${waf_event_id}** variable
  can be available.

<a name="Domain_timeout_settings"></a>
The `timeout_settings` block supports:

* `connection_timeout` - (Optional, Int) Specifies the timeout for WAF to connect to the origin server. The unit is second.
  Valid value ranges from `0` to `180`.

* `read_timeout` - (Optional, Int) Specifies the timeout for WAF to receive responses from the origin server.
  The unit is second. Valid value ranges from `0` to `3,600`.

* `write_timeout` - (Optional, Int) Specifies the timeout for WAF to send requests to the origin server. The unit is second.
  Valid value ranges from `0` to `3,600`.

## Attribute Reference

The following attributes are exported:

* `id` - ID of the domain.

* `protect_status` - The WAF mode. -1: bypassed, 0: disabled, 1: enabled.

* `access_status` - Whether a domain name is connected to WAF. 0: The domain name is not connected to WAF, 1: The domain
  name is connected to WAF.

* `protocol` - The protocol type of the client. The options are HTTP, HTTPS, and HTTP&HTTPS.

## Import

There are two ways to import WAF domain state.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_domain.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_domain.test <id>/<enterprise_project_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `keep_policy`, `charging_mode`, `ipv6_enable`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_waf_domain" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      keep_policy,
      charging_mode,
      ipv6_enable,
    ]
  }
}
```
