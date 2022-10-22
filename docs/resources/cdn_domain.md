---
subcategory: "Content Delivery Network (CDN)"
---

# huaweicloud_cdn_domain

CDN domain management.

## Example Usage

### Create a cdn domain

```hcl
variable "domain_name" {}
variable "origin_server" {}

resource "huaweicloud_cdn_domain" "domain_1" {
  name = var.domain_name
  type = "web"

  sources {
    origin      = var.origin_server
    origin_type = "ipaddr"
    active      = 1
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}
```

### Create a cdn domain with cache rules

```hcl
variable "domain_name" {}
variable "origin_server" {}

resource "huaweicloud_cdn_domain" "domain_1" {
  name = var.domain_name
  type = "web"

  sources {
    origin      = var.origin_server
    origin_type = "ipaddr"
    active      = 1
  }

  cache_settings {
    rules {
      rule_type = 0
      ttl       = 180
      ttl_type  = 4
      priority  = 2
    }
  }
}
```

### Create a cdn domain with configs

```hcl
variable "domain_name" {}
variable "origin_server" {}

resource "huaweicloud_cdn_domain" "domain_1" {
  name = var.domain_name
  type = "web"

  sources {
    origin      = var.origin_server
    origin_type = "ipaddr"
    active      = 1
  }

  configs {
    origin_protocol = "http"

    https_settings {
      certificate_name = "terraform-test"
      certificate_body = file("your_directory/chain.cer")
      http2_enabled    = true
      https_enabled    = true
      private_key      = file("your_directory/server_private.key")
    }

    cache_url_parameter_filter {
      type = "ignore_url_params"
    }

    retrieval_request_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    http_response_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    url_signing {
      enabled = false
    }

    compress {
      enabled = false
    }

    force_redirect {
      enabled = true
      type    = "http"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) The acceleration domain name. Changing this parameter will create a new
  resource.

* `type` - (Required, String, ForceNew) The service type. The valid values are  'web', 'download', 'video' and
  'wholeSite'.  Changing this parameter will create a new resource.

* `sources` - (Required, List, ForceNew) An array of one or more objects specifies the domain name of the origin server.
  The sources object structure is documented below.

* `service_area` - (Optional, String, ForceNew) The area covered by the acceleration service. Valid values are
  `mainland_china`, `outside_mainland_china`, and `global`. Changing this parameter will create a new resource.

* `configs` - (Optional, List) Specifies the domain configuration items. The [object](#configs_object) structure is
  documented below.

* `cache_settings` - (Optional, List) Specifies the cache configuration. The [object](#cache_settings_object) structure
  is documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id. Changing this parameter will create
  a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the domain.

The `sources` block supports:

* `origin` - (Required, String) The domain name or IP address of the origin server.

* `origin_type` - (Required, String) The origin server type. The valid values are 'ipaddr', 'domain', and 'obs_bucket'.

* `active` - (Optional, Int) Whether an origin server is active or standby (1: active; 0: standby). The default value is
  1.

* `obs_web_hosting_enabled` - (Optional, Bool) Whether to enable static website hosting for the OBS bucket.
  This parameter is mandatory when the `origin_type` is **obs_bucket**.

* `http_port` - (Optional, Int) Specifies the HTTP port. Default value: **80**.

* `https_port` - (Optional, Int) Specifies the HTTPS port. Default value: **443**.

* `retrieval_host` - (Optional, String) Specifies the retrieval host. The default value is the acceleration domain name.

<a name="configs_object"></a>
The `configs` block support:

* `origin_protocol` - (Optional, String) Specifies the content retrieval protocol. Possible values:
  + **follow**: same as user requests.
  + **http**: HTTP, which is the default value.
  + **https**: HTTPS.

* `ipv6_enable` - (Optional, Bool) Specifies whether to enable IPv6.

* `range_based_retrieval_enabled` - (Optional, Bool) Specifies whether to enable range-based retrieval.

* `https_settings` - (Optional, List) Specifies the certificate configuration. The [object](#https_settings_object)
  structure is documented below.

* `retrieval_request_header` - (Optional, List) Specifies the retrieval request header settings.
  The [object](#request_and_response_header_object) structure is documented below.

* `http_response_header` - (Optional, List) Specifies the HTTP response header settings.
  The [object](#request_and_response_header_object) structure is documented below.

* `url_signing` - (Optional, List) Specifies the URL signing.
  The [object](#url_signing_object) structure is documented below.

* `force_redirect` - (Optional, List) Specifies the force redirect.
  The [object](#redirect_and_compress_object) structure is documented below.

* `compress` - (Optional, List) Specifies the smart compression. The [object](#redirect_and_compress_object) structure
  is documented below.

* `cache_url_parameter_filter` - (Optional, List) Specifies the settings for caching URL parameters.
  The [object](#cache_url_parameter_filter_object) structure is documented below.

<a name="https_settings_object"></a>
The `https_settings` block support:

* `https_enabled` - (Optional, Bool) Specifies whether to enable HTTPS.

* `certificate_name` - (Optional, String) Specifies the certificate name. The value contains 3 to 32 characters.
  This parameter is mandatory when a certificate is configured.

* `certificate_body` - (Optional, String) Specifies the content of the certificate used by the HTTPS protocol.
  This parameter is mandatory when a certificate is configured. The value is in PEM format.

* `private_key` - (Optional, String) Specifies the private key used by the HTTPS protocol. This parameter is mandatory
  when a certificate is configured. The value is in PEM format.

* `certificate_source` - (Optional, Int) Specifies the certificate type. Possible values are:
  + **1**: Huawei-managed certificate.
  + **0**: your own certificate.
  
  Default value: **0**.
  This parameter is mandatory when a certificate is configured.

* `http2_enabled` - (Optional, Bool) Specifies whether HTTP/2 is used.

* `tls_version` - (Optional, String) Specifies the transport Layer Security (TLS). Currently, **TLSv1.0**,
  **TLSv1.1**, **TLSv1.2**, and **TLSv1.3** are supported. By default, all versions are enabled. You can enable
  a single version or consecutive versions. To enable multiple versions, use commas (,) to separate versions,
  for example, **TLSv1.1,TLSv1.2**.

<a name="request_and_response_header_object"></a>
The `retrieval_request_header` and `http_response_header` block support:

* `name` - (Required, String) Specifies the request or response header.

* `action` - (Required, String) Specifies the operation type of request or response

* `value` - (Optional, String) Specifies the value of request or response header.

<a name="url_signing_object"></a>
The `url_signing` block support:

* `enabled` - (Required, Bool) Specifies whether to enable of A/B/C URL signing method.

* `type` - (Optional, String) Specifies the signing Method, possible values are:
  **type_a**: Method A.
  **type_b**: Method B.
  **type_c1**: Method C1.
  **type_c2**: Method C2.

* `key` - (Optional, String) Specifies the authentication key contains 6 to 32 characters, including letters and digits.

* `time_format` - (Optional, String) Specifies the time format. Possible values are:
  **dec**: Decimal, can be used in Method A, Method B and Method C2.
  **hex**: Hexadecimal, can be used in Method C1 and Method C2.

* `expire_time` - (Optional, Int) Specifies the expiration time. The value ranges from **0** to **31536000**,
  in seconds.

<a name="redirect_and_compress_object"></a>
The `force_redirect` and `compress` blocks support:

* `enabled` - (Required, Bool) Specifies the whether to enable force redirect or smart compression.

* `type` - (Optional, String) Specifies the force redirect or smart compression type.
  Possible values for force redirect: **http** (force redirect to HTTP) and **https** (force redirect to HTTPS).
  Possible values for smart compression: **gzip** (gzip) and **br** (Brotli).

<a name="cache_url_parameter_filter_object"></a>
The `cache_url_parameter_filter` block support:

* `type` - (Optional, String) Specifies the operation type for caching URL parameters. Posiible values are:
  **full_url**: cache all parameters
  **ignore_url_params**: ignore all parameters
  **del_args**: ignore specific URL parameters
  **reserve_args**: reserve specified URL parameters

* `value` - (Optional, String) Specifies the parameter values. Multiple values are separated by semicolons (;).

<a name="cache_settings_object"></a>
The `cache_settings` block support:

* `follow_origin` - (Optional, Bool) Specifies whether to enable origin cache control.

* `rules` - (Optional, List) Specifies the cache rules, which overwrite the previous rule configurations.
  Blank rules are reset to default rules. The [object](#rules_object) structure is documented below.

<a name="rules_object"></a>
The `rules` block support:

* `rule_type` - (Required, Int) Specifies the rule type. Possible value are:
  **0**: All types of files are matched. It is the default value.
  **1**: Files are matched based on their suffixes.
  **2**: Files are matched based on their directories.
  **3**: Files are matched based on their full paths.

* `ttl` - (Required, Int) Specifies the cache age. The maximum cache age is 365 days.

* `ttl_type` - (Required, Int) Specifies the unit of the cache age. Possible values: **1** (second), **2** (minute),
  **3** (hour), and **4** (day).

* `priority` - (Required, Int) Specifies the priority weight of this rule. The default value is 1.
  A larger value indicates a higher priority. The value ranges from 1 to 100. The weight values must be unique.

* `content` - (Optional, String) Specifies the content that matches `rule_type`. If `rule_type` is set to **0**,
  this parameter is empty. If `rule_type` is set to **1**, the value of this parameter is a list of file name
  extensions. A file name extension starts with a period (.). File name extensions are separated by semicolons (;),
  for example, .jpg;.zip;.exe. If `rule_type` is set to **2**, the value of this parameter is a list of directories.
  A directory starts with a slash (/). Directories are separated by semicolons (;), for example,
  /test/folder01;/test/folder02.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The acceleration domain name ID.

* `cname` - The CNAME of the acceleration domain name.

* `domain_status` - The status of the acceleration domain name. The available values are
  'online', 'offline', 'configuring', 'configure_failed', 'checking', 'check_failed' and 'deleting.'

* `configs/https_settings/https_status` - The status of the https. The available values are 'on' and 'off'.

* `configs/https_settings/http2_status` - The status of the http 2.0. The available values are 'on' and 'off'.

* `configs/url_signing/status` - The status of the url_signing. The available values are 'on' and 'off'.

* `configs/force_redirect/status` - The status of the force redirect. The available values are 'on' and 'off'.

* `configs/compress/status` - The status of the compress. The available values are 'on' and 'off'.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minute.
* `delete` - Default is 20 minute.

## Import

Domains can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_cdn_domain.domain_1 fe2462fac09a4a42a76ecc4a1ef542f1
```
