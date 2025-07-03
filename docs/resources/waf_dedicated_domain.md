---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_dedicated_domain"
description: |-
  Manages a dedicated mode domain resource within HuaweiCloud.
---

# huaweicloud_waf_dedicated_domain

Manages a dedicated mode domain resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The dedicated mode domain name resource can be used in Dedicated Mode.

-> You can use this resource to create a cloud-mode ELB domain name, but there are usage restrictions:
<br/>1. To use cloud load balancer WAF, you need to submit a service ticket to enable it for you first.
<br/>2. If you want to use the load balancer access mode, make sure you are using standard, professional, or enterprise
cloud WAF.
<br/>3. You have purchased a dedicated load balancer with Specifications set to Application load balancing (HTTP/HTTPS).
<br/>For more usage restrictions, please refer to the documentation:
[Cloud Mode - Load Balancer Access](https://support.huaweicloud.com/intl/en-us/usermanual-waf/waf_01_0287.html)

## Example Usage

### Dedicated Mode domain

```hcl
variable "certificated_id" {}
variable "vpc_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_waf_dedicated_domain" "test" {
  domain                = "www.example.com"
  certificate_id        = var.certificated_id
  enterprise_project_id = var.enterprise_project_id
  protect_status        = 1
  website_name          = "websiteName"
  description           = "test description"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "192.168.1.100"
    port            = 8080
    type            = "ipv4"
    vpc_id          = var.vpc_id
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

  forward_header_map = {
    "key1" = "$time_local"
    "key2" = "$tenant_id"
  }
  
  connection_protection {
    error_threshold                       = 1000
    error_percentage                      = 90
    initial_downtime                      = 200
    multiplier_for_consecutive_breakdowns = 5
    pending_url_request_threshold         = 7000
    duration                              = 10000
    status                                = true
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
}
```

### Cloud ELB Mode domain

```hcl
variable "enterprise_project_id" {}
variable "loadbalancer_id" {}

resource "huaweicloud_waf_dedicated_domain" "test" {
  loadbalancer_id       = var.loadbalancer_id
  enterprise_project_id = var.enterprise_project_id
  domain                = "119.8.5.14"
  mode                  = "elb-shared"
  description           = "test description"
  protocol_port         = 0
  website_name          = "test-web-tag"

  lifecycle {
    ignore_changes = [ proxy ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the dedicated mode domain resource.
  If omitted, the provider-level region will be used. Changing this setting will push a new domain.

* `domain` - (Required, String, NonUpdatable) Specifies the protected domain name or IP address (port allowed). For example,
  `www.example.com` or `*.example.com` or `www.example.com:89`.

* `server` - (Optional, List) Specifies the server configuration list of the domain.
  A maximum of `80` can be configured. The [server](#DedicatedDomain_server) structure is documented below.

  -> Field `server` is required only when creating dedicated domain. This field does not need to be configured when
  creating a cloud mode ELB domain name.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID of WAF dedicated domain.
  For enterprise users, if omitted, default enterprise project will be used.

* `certificate_id` - (Optional, String) Specifies the certificate ID. This parameter is mandatory when `client_protocol`
  is set to **HTTPS**.

* `policy_id` - (Optional, String) Specifies the policy ID associated with the domain. If not specified, a new policy
  will be created automatically.

* `proxy` - (Optional, Bool) Specifies whether a proxy is configured. Defaults to **false**.

  -> **NOTE:** WAF forwards only HTTP/S traffic. So WAF cannot serve your non-HTTP/S traffic, such as UDP, SMTP, FTP,
  and basically all other non-HTTP/S traffic. If a proxy such as public network ELB (or Nginx) has been used, set
  proxy `true` to ensure that the WAF security policy takes effect for the real source IP address.

* `keep_policy` - (Optional, Bool) Specifies whether to retain the policy when deleting a domain name.
  Defaults to **true**.

* `protect_status` - (Optional, Int) Specifies the protection status of domain, `0`: suspended, `1`: enabled.
  Defaults to `0`.

* `tls` - (Optional, String) Specifies the minimum required TLS version. The valid values are: **TLS v1.0**,
  **TLS v1.1** and **TLS v1.2**.

* `cipher` - (Optional, String) Specifies the cipher suite of domain. The valid values are: **cipher_1**, **cipher_2**,
  **cipher_3**, **cipher_4**, **cipher_5**, **cipher_6**, and **cipher_default**.

* `pci_3ds` - (Optional, Bool) Specifies the status of the PCI 3DS compliance certification check. The options
  include **true** and **false**. This parameter must be used together with `tls` and `cipher`.

  -> **NOTE:** Tls must be set to **TLS v1.2**, and cipher must be set to **cipher_2**. The PCI 3DS compliance certification
  check cannot be disabled after being enabled.

* `pci_dss` - (Optional, Bool) Specifies the status of the PCI DSS compliance certification check. The options
  include **true** and **false**. This parameter must be used together with `tls` and `cipher`.

  -> **NOTE:** Tls must be set to **TLS v1.2**, and cipher must be set to **cipher_2**.

* `website_name` - (Optional, String) Specifies the website name. This website name must start with a letter and only
  letters, digits, underscores (_), hyphens (-), colons (:) and periods (.) are allowed. The value contains `1` to `128`
  characters. The website name must be unique within this account.

* `custom_page` - (Optional, List) Specifies the custom page. Only supports one custom alarm page.
  The [custom_page](#DedicatedDomain_custom_page) structure is documented below.

* `redirect_url` - (Optional, String) Specifies the URL of the redirected page. The root domain name of the redirection
  address must be the name of the currently protected domain (including a wildcard domain name).
  The available **${http_host}** can be used to indicate the currently protected domain name and port.
  For example: **${http_host}/error.html**.

-> The fields `redirect_url` and `custom_page` are mutually exclusive and cannot be specified simultaneously.

* `description` - (Optional, String) Specifies the description of the WAF dedicated domain.

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

* `connection_protection` - (Optional, List) Specifies the connection protection configuration to let WAF protect your
  origin servers from being crashed when WAF detects a large number of `502`/`504` error codes or pending requests.
  Only supports one protection configuration.
  The [connection_protection](#DedicatedDomain_connection_protection) structure is documented below.

* `timeout_settings` - (Optional, List) Specifies the timeout setting. Only supports one timeout setting.
  The [timeout_settings](#DedicatedDomain_timeout_settings) structure is documented below.

* `traffic_mark` - (Optional, List) Specifies the traffic identifier. WAF uses the configurations to identify the
  malicious client IP address (proxy mode) in the header, session in the cookie, and user attribute in the parameter,
  and then triggers the corresponding known attack source rules to block attack sources.
  Only supports one traffic identifier.
  The [traffic_mark](#DedicatedDomain_traffic_mark) structure is documented below.

* `mode` - (Optional, String) Specifies the mode. If you use the cloud ELB access mode, enter enter **elb-shared**.
  Otherwise, leave this parameter blank. This field is required only when creating cloud ELB mode domain.

* `loadbalancer_id` - (Optional, String, NonUpdatable) Specifies the load balancer ID.
  Available ELB IDs can be found through the datasource `huaweicloud_elb_loadbalancers`.
  This field is required only when creating cloud ELB mode domain.

* `listener_id` - (Optional, String, NonUpdatable) Specifies the listener ID.
  Available listener IDs can be found through the datasource `huaweicloud_elb_listeners`.
  All listeners configured and to be configured for the load balancer will be protected by WAF.
  In cloud ELB access mode, you are advised to set this parameter.

* `protocol_port` - (Optional, Int, NonUpdatable) Specifies the protocol port. Valid value ranges from `0` to `65,535`.
  This field is required only when creating cloud ELB mode domain.

* `pool_ids` - (Optional, List) Specifies the dedicated engine group the domain name was added to.
  This field is valid only when creating cloud ELB mode domain.

-> Fields `mode`, `loadbalancer_id`, `listener_id`, `protocol_port`, and `pool_ids` are only available when creating
cloud ELB mode domain. Due to API restrictions, the configuration of many fields will be invalid when creating a cloud
mode ELB domain name. This means that the cloud mode elb domain name does not support the configuration of these fields.
Please use `lifecycle` to ignore these fields. For example, you need to ignore the `proxy` or configure
the `proxy` to **true** before creating cloud ELB mode domain. Before using this resource, please refer to the example.

<a name="DedicatedDomain_server"></a>
The `server` block supports:

* `client_protocol` - (Required, String) Specifies the protocol type of the client. The options include **HTTP**
  and **HTTPS**.

* `server_protocol` - (Required, String) Specifies the protocol used by WAF to forward client requests to the
  server. The valid values are **HTTP** and **HTTPS**.

* `vpc_id` - (Required, String) Specifies the ID of the VPC used by the server.

* `type` - (Required, String) Specifies the server network type, IPv4 or IPv6.
  Valid values are **ipv4** and **ipv6**.

* `address` - (Required, String) Specifies the IP address or domain name of the web server accessed by the
  client. For example, `192.168.1.1` or `www.example.com`.

* `port` - (Required, Int) Specifies the port number used by the web server. The value ranges from `0` to `65,535`.

<a name="DedicatedDomain_custom_page"></a>
The `custom_page` block supports:

* `http_return_code` - (Required, String) Specifies the HTTP return code.
  The value can be a positive integer in the range of `200`-`599` except `408`, `444` and `499`.

* `block_page_type` - (Required, String) Specifies the content type of the custom alarm page.
  The value can be **text/html**, **text/xml** or **application/json**.

* `page_content` - (Required, String) Specifies the page content. The page content based on the selected page type.
  The available **${waf_event_id}** in the page content indicates an event ID, and only one **${waf_event_id}** variable
  can be available.

<a name="DedicatedDomain_connection_protection"></a>
The `connection_protection` block supports:

* `error_threshold` - (Optional, Int) Specifies the `502`/`504` error threshold for every 30 seconds. Valid value ranges
  from `0` to `2,147,483,647`.

* `error_percentage` - (Optional, Float) Specifies the `502`/`504` error percentage. A breakdown protection is triggered
  when the `502`/`504` error threshold and percentage threshold have been reached. Valid value ranges from `0` to `99`.

* `initial_downtime` - (Optional, Int) Specifies the breakdown duration (s) when the breakdown is triggered for the first
  time. Valid value ranges from `0` to `2,147,483,647`.

* `multiplier_for_consecutive_breakdowns` - (Optional, Int) Specifies the maximum multiplier for consecutive breakdowns
  that occur within an hour. Valid value ranges from `0` to `2,147,483,647`.
  For example: Assume that you set the initial downtime to `180s` and the maximum multiplier to `3`. If the breakdown
  protection is triggered for the second time, the website downtime is 360s (180s X 2).
  If the breakdown protection is triggered for the third or fourth time, the website downtime is 540s (180s x 3).
  The breakdowns are calculated every one hour.

* `pending_url_request_threshold` - (Optional, Int) Specifies the pending URL request threshold. Connection protection
  is triggered when the number of read URL requests reaches the threshold you configure. Valid value ranges from `0` to
  `2,147,483,647`.

* `duration` - (Optional, Int) Specifies the protection duration (s) for connection protection. During this period, WAF
  stops forwarding website requests. Valid value ranges from `0` to `2,147,483,647`.

* `status` - (Optional, Bool) Specifies whether to enable connection protection. Defaults to **false**.

<a name="DedicatedDomain_timeout_settings"></a>
The `timeout_settings` block supports:

* `connection_timeout` - (Optional, Int) Specifies the timeout for WAF to connect to the origin server. The unit is second.
  Valid value ranges from `0` to `180`.

* `read_timeout` - (Optional, Int) Specifies the timeout for WAF to receive responses from the origin server.
  The unit is second. Valid value ranges from `0` to `3,600`.

* `write_timeout` - (Optional, Int) Specifies the timeout for WAF to send requests to the origin server. The unit is second.
  Valid value ranges from `0` to `3,600`.

<a name="DedicatedDomain_traffic_mark"></a>
The `traffic_mark` block supports:

* `ip_tags` - (Optional, List) Specifies the IP tags. HTTP request header field of the original client IP address.
  This field is used to store the real IP address of the client. After the configuration, WAF preferentially reads the
  configured field to obtain the real IP address of the client. If multiple fields are configured, WAF reads the IP
  address list in order. Note:
  + If you want to use a TCP connection IP address as the client IP address, set IP Tag to `$remote_addr`.
  + If WAF does not obtain the real IP address of a client from fields you configure, WAF reads the `cdn-src-ip`,
    `x-real-ip`, `x-forwarded-for` and `$remote_addr` fields in sequence to read the client IP address.

* `session_tag` - (Optional, String) Specifies the session tag. This tag is used by known attack source rules to block
  malicious attacks based on cookie attributes. This parameter must be configured in known attack source rules to block
  requests based on cookie attributes.

* `user_tag` - (Optional, String) Specifies the user tag. This tag is used by known attack source rules to block malicious
  attacks based on params attributes. This parameter must be configured to block requests based on the params attributes.

## Attribute Reference

The following attributes are exported:

* `id` - ID of the domain.

* `certificate_name` - The name of the certificate used by the domain name.

* `access_status` - Whether a domain name is connected to WAF. Valid values are:
  + `0` - The domain name is not connected to WAF,
  + `1` - The domain name is connected to WAF.

* `protocol` - The protocol type of the client. The options are `HTTP` and `HTTPS`.

* `compliance_certification` - The compliance certifications of the domain, values are:
  + `pci_dss` - The status of domain PCI DSS, `true`: enabled, `false`: disabled.
  + `pci_3ds` - The status of domain PCI 3DS, `true`: enabled, `false`: disabled.

* `alarm_page` - The alarm page of domain. Valid values are:
  + `template_name` - The template of alarm page, values are: `default`, `custom` and `redirection`.
  + `redirect_url` - The redirection URL when `template_name` is set to `redirection`.

* `traffic_identifier` - The traffic identifier of domain. Valid values are:
  + `ip_tag` - The IP tag of traffic identifier.
  + `session_tag` - The session tag of traffic identifier.
  + `user_tag` - The user tag of traffic identifier.

* `timestamp` - The time the domain name was added to WAF.

* `extend` - The extended field, which is used to save some configuration information about the protected domain name.
  The value is key-value pairs.

* `flag` - The identifier, which is used on the console.
  The [flag](#DedicatedDomain_flag) structure is documented below.

* `block_page` - The alarm page configuration.
  The [block_page](#DedicatedDomain_block_page) structure is documented below.

<a name="DedicatedDomain_flag"></a>
The `flag` block supports:

* `pci_3ds` - Whether the website passes the PCI 3DS certification check. Valid values are **true** and **false**.

* `pci_dss` - Whether the website passed the PCI DSS certification check. Valid values are **true** and **false**.

* `cname` - The CNAME record being used. Valid values are **old** and **new**.

* `is_dual_az` - Whether WAF support Multi-AZ DR. Valid values are **true** and **false**.

* `ipv6` - Whether IPv6 protection is supported. Valid values are **true** and **false**.

<a name="DedicatedDomain_block_page"></a>
The `block_page` block supports:

* `template` - Template name. Enter **default** for the default page. Enter **custom** for the customized alarm page.
  Enter **redirect** for the redirection page.

* `custom_page` - The custom alarm page.
  The [custom_page](#DedicatedDomain_block_page_custom_page) structure is documented below.

* `redirect_url` - The URL of the redirected page.

<a name="DedicatedDomain_block_page_custom_page"></a>
The `custom_page` block supports:

* `status_code` - The status code.

* `content_type` - The content type of the custom alarm page. The value can be **text/html**, **text/xml** or **application/json**.

* `content` - The page content based on the selected page type.

## Import

There are two ways to import WAF dedicated domain state.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_dedicated_domain.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_dedicated_domain.test <id>/<enterprise_project_id>
```
