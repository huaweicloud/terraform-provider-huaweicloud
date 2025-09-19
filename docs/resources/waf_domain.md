---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_domain"
description: |-
  Manages a WAF domain resource within HuaweiCloud.
---

# huaweicloud_waf_domain

Manages a WAF domain resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The domain name resource can be used in Cloud Mode.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "certificate_id" {}
variable "certificate_name" {}

resource "huaweicloud_waf_domain" "test" {
  domain                = "www.example.com"
  certificate_id        = var.certificate_id
  certificate_name      = var.certificate_name
  proxy                 = true
  enterprise_project_id = var.enterprise_project_id
  description           = "test description"
  website_name          = "websiteName"
  protect_status        = 1
  
  forward_header_map = {
    "key1" = "$time_local"
    "key2" = "$tenant_id"
  }

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

  traffic_mark {
    ip_tags     = ["ip_tag"]
    session_tag = "session_tag"
    user_tag    = "user_tag"
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.13"
    port            = "8080"
    type            = "ipv4"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF domain resource.
  If omitted, the provider-level region will be used. Changing this setting will push a new certificate.

* `domain` - (Required, String, ForceNew) Specifies the domain name to be protected. For example, `www.example.com` or
  `*.example.com`. Changing this creates a new domain.

* `server` - (Required, List) Specifies an array of origin web servers.
  The [server](#Domain_server) structure is documented below.

* `certificate_id` - (Optional, String) Specifies the certificate ID. This parameter is mandatory when `client_protocol`
  is set to **HTTPS**.

* `certificate_name` - (Optional, String) Specifies the certificate name. This parameter is mandatory
  when `client_protocol` is set to **HTTPS**.

* `policy_id` - (Optional, String) Specifies the policy ID associated with the domain. If not specified, a new
  policy will be created automatically.

* `keep_policy` - (Optional, Bool) Specifies whether to retain the policy when deleting a domain name.
  Defaults to **true**.
  
* `proxy` - (Optional, Bool) Specifies whether a proxy is configured.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the domain. Valid values are **prePaid**
  and **postPaid**, defaults to **prePaid**. Changing this creates a new instance.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF domain.
  For enterprise users, if omitted, default enterprise project will be used.
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
  Defaults to **false**.
  Things to note when using this field are as follows:
  + There must be at least one server configuration with client protocol set to **HTTPS**, or this configuration is unable
    to work.
  + This field cannot not work if the client supports **TLS 1.3**.
  + This field can work only when the client supports **TLS 1.2** or earlier versions.
  + If you want to use HTTP/2 forwarding, use a dedicated WAF instance.

* `ipv6_enable` - (Optional, Bool) Specifies whether IPv6 protection is enabled.
  Enable IPv6 protection if the domain name is accessible using an IPv6 address.
  After you enable it, WAF assigns an IPv6 address to the domain name.
  This field must be set to **true** when `server` contains a value of type **ipv6**.
  Defaults to false.

* `website_name` - (Optional, String) Specifies the website name.
  This website name must start with a letter and only letters, digits, underscores (_),
  hyphens (-), colons (:) and periods (.) are allowed.
  The value contains `1` to `128` characters.
  The website name must be unique within this account.

* `description` - (Optional, String) Specifies the description of the WAF domain.

* `protect_status` - (Optional, Int) The protection status of domain. Valid values are:
  + `0`: The WAF protection is suspended. WAF only forwards requests destined for the domain name and does not detect attacks.
  + `1`: The WAF protection is enabled. WAF detects attacks based on the policy you configure.
  + `-1`: The WAF protection is bypassed. Requests of the domain name are directly sent to the backend server and do
  not pass through WAF.

  Default value is `0`.

* `access_status` - (Optional, Int) Specifies whether a domain name is connected to WAF.  
  `0`: The domain name is not connected to WAF, `1`: The domain name is connected to WAF, `2`: Skip access.

* `pci_3ds` - (Optional, Bool) Specifies the status of the PCI 3DS compliance certification check.
  This parameter must be used together with `tls` and `cipher`.

  -> **NOTE:** `tls` must be set to **TLS v1.2**, and `cipher` must be set to **cipher_2**. The PCI 3DS compliance certification
  check cannot be disabled after being enabled.
  The field `pci_3ds` is meaningful only if `certificate_id` is specified.

* `pci_dss` - (Optional, Bool) Specifies the status of the PCI DSS compliance certification check.
  This parameter must be used together with `tls` and `cipher`.

  -> **NOTE:** `tls` must be set to **TLS v1.2**, and `cipher` must be set to **cipher_2**.
  The field `pci_dss` is meaningful only if `certificate_id` is specified.

* `cipher` - (Optional, String) Specifies the cipher suite of domain.
  The options include **cipher_1**, **cipher_2**,**cipher_3**, **cipher_4**, **cipher_default**.

* `tls` - (Optional, String) Specifies the minimum required TLS version. The options include **TLS v1.0**, **TLS v1.1**,
  **TLS v1.2**.

* `lb_algorithm` - (Optional, String) Specifies the load balancing algorithms used to
  distribute requests across origin servers.
  Only the professional edition (original enterprise edition) and platinum edition
  (original ultimate edition) support configuring the load balancing algorithm.
  The options of value are as follows:
  + **ip_hash** : Requests from the same IP address are routed to the same backend server.
  + **round_robin** : Requests are distributed across backend servers in turn based on the
  weight you assign to each server.
  + **session_hash** : Direct requests with the same session ID to the same origin server.
  Before using this configuration, please make sure to configure the traffic identifier for
  attack punishment after adding the domain name, otherwise the session hash configuration will not take effect.

* `forward_header_map` - (Optional, Map) Specifies the field forwarding configuration. WAF inserts the added fields into
  the header and forwards the header to the origin server. The key cannot be the same as the native Nginx field.
  The options of value are as follows:
  + **$time_local**
  + **$request_id**
  + **$connection_requests**
  + **$tenant_id**
  + **$project_id**
  + **$remote_addr**
  + **$remote_port**
  + **$scheme**
  + **$request_method**
  + **$http_host**
  + **$origin_uri**
  + **$request_length**
  + **$ssl_server_name**
  + **$ssl_protocol**
  + **$ssl_curves**
  + **$ssl_session_reused**

* `traffic_mark` - (Optional, List) Specifies the traffic identifier.
  WAF uses the configurations to identify the malicious client IP address (proxy mode) in the header,
  session in the cookie, and user attribute in the parameter,
  and then triggers the corresponding known attack source rules to block attack sources.
  Only supports one traffic identifier.
  The [traffic_mark](#Domain_traffic_mark) structure is documented below.

* `timeout_settings` - (Optional, List) Specifies the timeout setting. Only supports one timeout setting.
  The [timeout_settings](#Domain_timeout_settings) structure is documented below.

<a name="Domain_server"></a>
The `server` block supports:

* `client_protocol` - (Required, String) Specifies the protocol type of the client. The options include **HTTP** and **HTTPS**.

* `server_protocol` - (Required, String) Specifies the protocol used by WAF to forward client requests to the server.
  The options include **HTTP** and **HTTPS**.

* `address` - (Required, String) Specifies the IP address or domain name of the web server that the client accesses.

* `port` - (Required, Int) Specifies the port number used by the web server. The value ranges from `0` to `65,535`,
  for example, `8,080`.

* `type` - (Required, String) Specifies the server network type. Valid values are: **ipv4** and **ipv6**.
  + When this field is set to **ipv4**, `address` must be set to an IPv4 address.
  + When this field is set to **ipv6**, `address` must be set to an IPv6 address.

* `weight` - (Optional, Int) Specifies the load balancing algorithm will assign requests to the origin
  site according to this weight.
  Defaults to `1`.

<a name="Domain_custom_page"></a>
The `custom_page` block supports:

* `http_return_code` - (Required, String) Specifies the HTTP return code.
  The value can be a positive integer in the range of `200` to `599` except `408`, `444` and `499`.

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

<a name="Domain_traffic_mark"></a>
The `traffic_mark` block supports:

* `ip_tags` - (Optional, List) Specifies the IP tags. HTTP request header field of the original client IP address.
  This field is used to store the real IP address of the client. After the configuration, WAF preferentially reads the
  configured field to obtain the real IP address of the client. If multiple fields are configured, WAF reads the IP
  address list in order. Note:
  + If you want to use a TCP connection IP address as the client IP address, set IP Tag to **$remote_addr**.
  + If WAF does not obtain the real IP address of a client from fields you configure, WAF reads the **cdn-src-ip**,
    **x-real-ip**, **x-forwarded-for** and **$remote_addr** fields in sequence to read the client IP address.
  + When the website setting `proxy` is configured as **false**, this field does not take effect,
    and the client IP is only obtained through the `$remote_addr` field.

* `session_tag` - (Optional, String) Specifies the session tag. This tag is used by known attack source rules to block
  malicious attacks based on cookie attributes. This parameter must be configured in known attack source rules to block
  requests based on cookie attributes.

* `user_tag` - (Optional, String) Specifies the user tag. This tag is used by known attack source rules to block malicious
  attacks based on params attributes. This parameter must be configured to block requests based on the params attributes.

## Attribute Reference

The following attributes are exported:

* `id` - ID of the domain.

* `access_code` - The CNAME prefix. The CNAME suffix is `.vip1.huaweicloudwaf.com`.

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
