---
subcategory: "Content Delivery Network (CDN)"
---

# huaweicloud_cdn_domain

Manages a CDN domain resource within HuaweiCloud.

## Example Usage

### Create a CDN domain

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

### Create a CDN domain with cache rules

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

### Create a CDN domain with configs

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

* `name` - (Required, String, ForceNew) Specifies acceleration domain name. Changing this parameter will create a new
  resource. The domain name consists of one or more parts, representing domains at different levels.
  Domain names at all levels can only be composed of letters, digits, and hyphens (-), and the letters are equivalent in
  upper and lower case. Domain names at all levels are connected with (.). The domain name can contain up to `75` characters.

* `type` - (Required, String, ForceNew) Specifies the service type of the domain name. Changing this parameter will
  create a new resource. The valid values are as follows:
  + **web**: Static acceleration. For websites with many images and small files, such as portals and e-commerce websites.
  + **download**: Download acceleration. For large files, such as apps in app stores and game clients.
  + **video**: Streaming media acceleration. For video on demand (VOD) websites and online education websites.
  + **wholeSite**: Whole site acceleration. For websites with both dynamic and static content, such as online exam
    platforms, forums, and blogs.

* `sources` - (Required, List) Specifies an array of one or more objects specifying origin server settings.
  A maximum of `50` origin site configurations can be configured.
  The [sources](#sources_cdn_domain) structure is documented below.

* `service_area` - (Optional, String, ForceNew) Specifies the area covered by the acceleration service.
  Changing this parameter will create a new resource. Valid values are as follows:
  + **mainland_china**: Indicates that the service scope is mainland China.
  + **outside_mainland_china**: Indicates that the service scope is outside mainland China.
  + **global**: Indicates that the service scope is global.

* `configs` - (Optional, List) Specifies the domain configuration items. The [configs](#configs_object) structure is
  documented below.

* `cache_settings` - (Optional, List) Specifies the cache configuration. The [cache_settings](#cache_settings_object) structure
  is documented below.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID. Changing this parameter
  will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the domain.

<a name="sources_cdn_domain"></a>
The `sources` block supports:

* `origin` - (Required, String) Specifies the unique domain name or IP address of the origin server.
  + If `origin_type` is set to **ipaddr**, this field can only be set to IPv4 address.
  + If `origin_type` is set to **domain**, this field can only be set to domain name.
  + If `origin_type` is set to **obs_bucket**, this field can only be set to OBS bucket domain name. The OBS bucket
    domain name must end with `.myhuaweicloud.com` or `.myhuaweicloud.cn`.

* `origin_type` - (Required, String) Specifies the origin server type. The valid values are as follows:
  + **ipaddr**: Origin server IP address.
  + **domain**: Origin server domain name.
  + **obs_bucket**: OBS bucket domain name.

* `active` - (Optional, Int) Specifies whether the origin server is primary or standby. Valid values are as follows:
  + **1**: Primary.
  + **0**: Standby.

  Defaults to **1**.

* `obs_web_hosting_enabled` - (Optional, Bool) Specifies whether to enable static website hosting for the OBS bucket.
  This parameter is valid only when the `origin_type` is set to **obs_bucket**. Defaults to **false**.

* `http_port` - (Optional, Int) Specifies the HTTP port. The port number ranges from `1` to `65535`.
  Defaults to **80**.

* `https_port` - (Optional, Int) Specifies the HTTPS port. The port number ranges from `1` to `65535`.
  Default value: **443**.

-> Fields `http_port` and `https_port` are valid only when `origin_type` is set to **ipaddr** or **domain**.

* `retrieval_host` - (Optional, String) Specifies the retrieval host. Things to note when using this field are as follows:
  + If `origin_type` is set to **ipaddr** or **domain**, the acceleration domain name will be used by default.
  + If `origin_type` is set to **obs_bucket**, the bucket's domain name will be used by default.

<a name="configs_object"></a>
The `configs` block support:

* `origin_protocol` - (Optional, String) Specifies the content retrieval protocol. Possible values:
  + **follow**: Same as user requests.
  + **http**: HTTP, which is the default value.
  + **https**: HTTPS.

* `ipv6_enable` - (Optional, Bool) Specifies whether to enable IPv6.

* `range_based_retrieval_enabled` - (Optional, Bool) Specifies whether to enable range-based retrieval.

  -> The prerequisite for enabling range-based retrieval is that your origin site supports Range requests, that is, the
  HTTP request header contains the Range field. Otherwise, the back-to-origin may fail.

* `https_settings` - (Optional, List) Specifies the certificate configuration. The [https_settings](#https_settings_object)
  structure is documented below.

* `retrieval_request_header` - (Optional, List) Specifies the retrieval request header settings.
  The [retrieval_request_header](#retrieval_request_header_object) structure is documented below.

* `http_response_header` - (Optional, List) Specifies the HTTP response header settings.
  The [http_response_header](#http_response_header_object) structure is documented below.

* `url_signing` - (Optional, List) Specifies the URL signing.
  The [url_signing](#url_signing_object) structure is documented below.

* `force_redirect` - (Optional, List) Specifies the force redirect.
  The [force_redirect](#force_redirect_object) structure is documented below.

* `compress` - (Optional, List) Specifies the smart compression. The [compress](#compress_object) structure
  is documented below.

* `cache_url_parameter_filter` - (Optional, List) Specifies the settings for caching URL parameters.
  The [cache_url_parameter_filter](#cache_url_parameter_filter_object) structure is documented below.

* `ip_frequency_limit` - (Optional, List) Specifies the IP access frequency limit.
  The [ip_frequency_limit](#ip_frequency_limit_object) structure is documented below.

  -> Restricting the IP access frequency can effectively defend against CC attacks, but it may affect normal access.
  Please set access thresholds carefully.

* `websocket` - (Optional, List) Specifies the websocket settings. This field can only be configured if `type` is
  set to **wholeSite**. The [websocket](#websocket_object) structure is documented below.

  -> Websocket and HTTP/2 are incompatible and cannot be both enabled. Websocket will not take effect when
  origin cache control is enabled in the cache configuration.

<a name="https_settings_object"></a>
The `https_settings` block support:

* `https_enabled` - (Optional, Bool) Specifies whether to enable HTTPS. Defaults to **false**.

* `certificate_name` - (Optional, String) Specifies the certificate name. The value contains `3` to `32` characters.
  This parameter is mandatory when a certificate is configured.

* `certificate_body` - (Optional, String) Specifies the content of the certificate used by the HTTPS protocol.
  This parameter is mandatory when a certificate is configured. The value is in PEM format.

* `private_key` - (Optional, String) Specifies the private key used by the HTTPS protocol. This parameter is mandatory
  when a certificate is configured. The value is in PEM format.

* `certificate_source` - (Optional, Int) Specifies the certificate type. Possible values are:
  + **1**: Huawei-managed certificate.
  + **0**: Your own certificate.
  
  Defaults to **0**.

* `http2_enabled` - (Optional, Bool) Specifies whether HTTP/2 is used. Defaults to **false**.
  When `https_enabled` is set to **false**, this parameter does not take effect.

* `tls_version` - (Optional, String) Specifies the transport Layer Security (TLS). Currently, **TLSv1.0**,
  **TLSv1.1**, **TLSv1.2**, and **TLSv1.3** are supported. By default, **TLS 1.1**, **TLS 1.2**, and **TLS 1.3** are
  enabled. You can enable a single version or consecutive versions. To enable multiple versions, use commas (,) to
  separate versions, for example, **TLSv1.1,TLSv1.2**.

<a name="retrieval_request_header_object"></a>
The `retrieval_request_header` block support:

* `name` - (Required, String) Specifies the name of a retrieval request header. The value contains `1` to `64` characters,
  including digits, letters, and hyphens (-). The value must start with a letter.

* `action` - (Required, String) Specifies the operation type of the retrieval request header. Valid values are **delete**
  and **set**. If the header does not exist in the original retrieval request, add the header before setting its value.

* `value` - (Optional, String) Specifies the value of the retrieval request header. The value contains `1` to `1000`
  characters, including letters, digits, and special characters `.-_*#!&+|^~'"/:;,=@?<>`. Variables, for example,
  `$client_ip` and `$remote_port`, are not supported.

<a name="http_response_header_object"></a>
The `http_response_header` block support:

* `name` - (Required, String) Specifies the HTTP response header. Valid values are **Content-Disposition**, **Content-Language**,
  **Access-Control-Allow-Origin**, **Access-Control-Allow-Methods**, **Access-Control-Max-Age**, **Access-Control-Expose-Headers**,
  **Access-Control-Allow-Headers** or custom headers. A header contains `1` to `100` characters, including letters, digits,
  and hyphens (-), and starts with a letter.

* `action` - (Required, String) Specifies the operation type of the HTTP response header. The value can be **set** or **delete**.

* `value` - (Optional, String) Specifies the value of the HTTP response header. The value contains `1` to `128` characters,
  including letters, digits, and special characters `.-_*#!&+|^~'"/:;,=@?<>`.

<a name="url_signing_object"></a>
The `url_signing` block support:

* `enabled` - (Required, Bool) Specifies whether to enable of A/B/C URL signing method.

* `type` - (Optional, String) Specifies the signing Method, possible values are:
  **type_a**: Method A.
  **type_b**: Method B.
  **type_c1**: Method C1.
  **type_c2**: Method C2.

* `key` - (Optional, String) Specifies the authentication key contains `6` to `32` characters, including letters and digits.

* `time_format` - (Optional, String) Specifies the time format. Possible values are:
  **dec**: Decimal, can be used in Method A, Method B and Method C2.
  **hex**: Hexadecimal, can be used in Method C1 and Method C2.

* `expire_time` - (Optional, Int) Specifies the expiration time. The value ranges from `0` to `31536000`, in seconds.

<a name="force_redirect_object"></a>
The `force_redirect` blocks support:

* `enabled` - (Required, Bool) Specifies whether to enable force redirect.

* `type` - (Optional, String) Specifies the force redirect type.
  Possible values are: **http** (force redirect to HTTP) and **https** (force redirect to HTTPS).

<a name="compress_object"></a>
The `compress` blocks support:

* `enabled` - (Required, Bool) Specifies whether to enable smart compression.

* `type` - (Optional, String) Specifies the smart compression type.
  Possible values are: **gzip** (gzip) and **br** (Brotli).

<a name="cache_url_parameter_filter_object"></a>
The `cache_url_parameter_filter` block support:

* `type` - (Optional, String) Specifies the operation type for caching URL parameters. Valid values are:
  **full_url**: Cache all parameters
  **ignore_url_params**: Ignore all parameters
  **del_params**: Ignore specific URL parameters
  **reserve_params**: Reserve specified URL parameters

* `value` - (Optional, String) Specifies the parameter values. Multiple values are separated by semicolons (;).

<a name="ip_frequency_limit_object"></a>
The `ip_frequency_limit` block support:

* `enabled` - (Required, Bool) Specifies whether to enable IP access frequency.

* `qps` - (Optional, Int) Specifies the access threshold, in times/second. The value ranges from **1** to **100,000**.
  This field is required when enable IP access frequency.

<a name="websocket_object"></a>
The `websocket` block support:

* `enabled` - (Required, Bool) Specifies whether to enable websocket settings.

* `timeout` - (Optional, Int) Specifies the duration for keeping a connection open, in seconds. The value ranges
  from **1** to **300**. This field is required when enable websocket settings.

<a name="cache_settings_object"></a>
The `cache_settings` block support:

* `follow_origin` - (Optional, Bool) Specifies whether to enable origin cache control. Defaults to **false**.

* `rules` - (Optional, List) Specifies the cache rules, which overwrite the previous rule configurations.
  Blank rules are reset to default rules. The [rules](#rules_object) structure is documented below.

<a name="rules_object"></a>
The `rules` block support:

* `rule_type` - (Required, String) Specifies the rule type. Possible value are:
  + **all**: All types of files are matched. It is the default value.
  + **file_extension**: Files are matched based on their suffixes.
  + **catalog**: Files are matched based on their directories.
  + **full_path**: Files are matched based on their full paths.
  + **home_page**: Files are matched based on their homepage.

* `ttl` - (Required, Int) Specifies the cache age. The maximum cache age is 365 days.

* `ttl_type` - (Required, String) Specifies the unit of the cache age. Possible values:
  + **s**: Second
  + **m**: Minute
  + **h**: Hour
  + **d**: Day

* `priority` - (Required, Int) Specifies the priority weight of this rule. The default value is 1.
  A larger value indicates a higher priority. The value ranges from 1 to 100. The weight values must be unique.

* `content` - (Optional, String) Specifies the content that matches `rule_type`.
  + If `rule_type` is set to **all** or **home_page**, keep this parameter empty.
  + If `rule_type` is set to **file_extension**, the value of this parameter is a list of file name
    extensions. A file name extension starts with a period (.). File name extensions are separated by semicolons (;),
    for example, `.jpg;.zip;.exe`. Up to 20 file types are supported.
  + If `rule_type` is set to **catalog**, the value of this parameter is a list of directories. A directory starts with
    a slash (/). Directories are separated by semicolons (;), for example, `/test/folder01;/test/folder02`.
    Up to 20 directories are supported.
  + If `rule_type` is set to **full_path**, the value must start with a slash (/) and cannot end with an asterisk.
    Example: `/test/index.html` or `/test/*.jpg`

## Attribute Reference

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

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The CDN domain resource can be imported using the domain `name`, e.g.

```bash
$ terraform import huaweicloud_cdn_domain.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_cdn_domain" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      enterprise_project_id,
    ]
  }
}
```
