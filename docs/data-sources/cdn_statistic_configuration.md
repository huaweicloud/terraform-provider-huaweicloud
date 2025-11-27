---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_statistic_configuration"
description: |-
  Use this data source to get a list of CDN statistic configurations within HuaweiCloud.
---

# huaweicloud_cdn_statistic_configuration

Use this data source to get a list of CDN statistic configurations within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_statistic_configuration" "test" {
  config_type = 0
}
```

## Argument Reference

The following arguments are supported:

* `config_type` - (Required, Int) Specifies the configuration category.  
  The valid values are as follows:
  + **0**: Hotspot statistics category.
  + **1**: CES reporting category.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The list of statistic configurations.  
  The [configurations](#cdn_statistic_configuration_configurations) structure is documented below.

<a name="cdn_statistic_configuration_configurations"></a>
The `configurations` block supports:

* `config_type` - The configuration category.

* `resource_type` - The resource type.
  + **domain**: The `resource_name` is a domain name.
  + **account**: The `resource_name` is an account.

* `resource_name` - The resource name.  
  Multiple domain names are separated by commas.

* `config_info` - The statistics configuration information.  
  The [config_info](#cdn_statistic_configuration_config_info) structure is documented below.

* `expired_time` - The expiration time of the statistics configuration, in seconds timestamp.

<a name="cdn_statistic_configuration_config_info"></a>
The `config_info` block supports:

* `url` - The top URL statistics configuration.  
  The [url](#cdn_statistic_configuration_url) structure is documented below.

* `ua` - The top UA statistics configuration.  
  The [ua](#cdn_statistic_configuration_ua) structure is documented below.

<a name="cdn_statistic_configuration_url"></a>
The `url` block supports:

* `enable` - Whether the top URL statistics configuration is enabled.

* `limit` - The number of top URL statistics to report.

* `sort_by_code` - Whether to support reporting by status code.

<a name="cdn_statistic_configuration_ua"></a>
The `ua` block supports:

* `enable` - Whether the top UA statistics configuration is enabled.
