---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_statistic_configuration"
description: |-
  Manages the CDN statistic configuration resource within HuaweiCloud.
---

# huaweicloud_cdn_statistic_configuration

Manages the CDN statistic configuration resource within HuaweiCloud.

-> This resource is a one-time action resource for setting CDN statistic configuration. Deleting this resource will
   not restore the corresponding configuration, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_cdn_statistic_configuration" "test" {
  resource_type = "domain"
  resource_name = "www.example.com,www.example2.com"
  config_type   = 0

  config_info {
    url {
      enable = true
    }
    ua {
      enable = true
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String, NonUpdatable) Specifies the resource type.  
  The valid values are as follows:
  + **domain**: The `resource_name` is a domain name.
  + **account**: The `resource_name` is an account.

* `resource_name` - (Required, String, NonUpdatable) Specifies the resource name, which can be an account or domain
  name.  
  Multiple domain names are separated by commas.  
  When `resource_type` is set to **account**, the value can be **all** to represent all accounts.

* `config_info` - (Required, List, NonUpdatable) Specifies the statistics configuration information.  
  Top indicators only support ua, refer, url, and origin url.  
  The [config_info](#cdn_statistic_configuration_config_info) structure is documented below.

* `config_type` - (Optional, Int, NonUpdatable) Specifies the configuration category.  
  The valid values are as follows:
  + **0**: Hotspot statistics category.

<a name="cdn_statistic_configuration_config_info"></a>
The `config_info` block supports:

* `url` - (Optional, List) Specifies the top URL statistics configuration.  
  The [url](#cdn_statistic_configuration_url) structure is documented below.

* `ua` - (Optional, List) Specifies the top UA statistics configuration.  
  The [ua](#cdn_statistic_configuration_ua) structure is documented below.

<a name="cdn_statistic_configuration_url"></a>
The `url` block supports:

* `enable` - (Required, Bool) Specifies whether to enable the top URL statistics configuration.

<a name="cdn_statistic_configuration_ua"></a>
The `ua` block supports:

* `enable` - (Required, Bool) Specifies whether to enable the top UA statistics configuration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
