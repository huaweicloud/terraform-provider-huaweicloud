---
subcategory: Content Delivery Network (CDN)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domains"
description: |-
  Use this data source to get a list of CDN domains within HuaweiCloud.
---

# huaweicloud_cdn_domains

Use this data source to get a list of CDN domains within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_domains" "test" {
  type          = "web"
  domain_status = "online"
  service_area  = "mainland_china"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional, String) Specifies the ID of the domain.

* `name` - (Optional, String) Specifies the name of the domain, using fuzzy matching.  
  The valid length is limited from `1` to `255`.

* `type` - (Optional, String) Specifies the business type of the domain.  
  The valid values are as follows:
  + **web**: Accelerate for the website.
  + **download**: Accelerate for file downloads.
  + **video**: Accelerate for on-demand.
  + **wholeSite**: Accelerate for the entire site.

* `domain_status` - (Optional, String) Specifies the status of the domain.  
  The valid values are as follows:
  + **online**
  + **offline**
  + **configuring**
  + **configuring_failed**
  + **checking**
  + **check_failed**
  + **deleting**

* `service_area` - (Optional, String) Specifies the accelerated coverage area for the domain.  
  The valid values are as follows:
  + **mainland_china**
  + **outside_mainland_china**
  + **global**

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource
  belongs.  
  When the user turns on the enterprise project function, this parameter takes effect,
  indicating that the project to which the resource belongs is queried.
  "all" indicates all projects.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domains` - The list of domains that matched filter parameters.  
  The [domains](#cdn_domains) structure is documented below.

<a name="cdn_domains"></a>
The `domains` block supports:

* `id` - The ID of the domain.

* `name` - The name of the domain.

* `type` - The business type of the domain.

* `domain_status` - The status of the domain.

* `cname` - The CNAME of the domain.

* `sources` - An array of one or more objects specifies the domain of the origin server.  
  The [sources](#cdn_domains_sources) structure is documented below.

* `domain_origin_host` - The back-to-origin HOST configuration of accelerate domain.  
  The [domain_origin_host](#cdn_domains_origin_host) structure is documented below.

* `https_status` - Whether the HTTPS protocol is enabled.  
  + **0**: Disable HTTPS acceleration.
  + **1**: Turn on HTTPS acceleration.

* `created_at` - The creation time, in RFC3339 format.

* `updated_at` - The update time, in RFC3339 format.

* `disabled` - The ban status of the domain.  
  + **0**: The domain is not banned.
  + **1**: The domain is banned.

* `locked` - The lock status of the domain.  
  + **0**: The domain is not locked
  + **1**: The domain is locked.

* `auto_refresh_preheat` - Whether to automatically refresh preheating.  
  + **0**: Auto_refresh_preheat is off.
  + **1**: Auto_refresh_preheat is on.

* `service_area` - The accelerated coverage area for the domain.

* `range_based_retrieval_enabled` - Whether to enable range-based retrieval.  
  The valid value can be **true** or **false**.

* `follow_status` - The status of back-to-source following.  
  The valid value can be **on** or **off**.

* `origin_status` - Whether to pause origin site return to origin.  
  The valid value can be **on** or **off**.

* `banned_reason` - The reason why the domain was banned.

* `locked_reason` - The reason why the domain was locked.

* `enterprise_project_id` - The ID of the enterprise project to which the resource belongs.

* `tags` - The key/value pairs to associate with the domain.

<a name="cdn_domains_sources"></a>
The `sources` block supports:

* `origin` - The domain name or IP address of the origin server.

* `origin_type` - The origin server type.  
  + **ipaddr**
  + **domain**
  + **obs_bucket**

* `active` - Whether an origin server is active or standby.  
  + **1**: The origin source is primary source site.
  + **0**: The origin source is backup source site.

* `obs_web_hosting_enabled` - Whether to enable static website hosting for the OBS bucket.  
  The valid value can be **true** or **false**.

<a name="cdn_domains_origin_host"></a>
The `domain_origin_host` block supports:

* `origin_host_type` - The type of origin host.  
  + **accelerate**: Select the accelerate domain as the back-to-origin host domain.
  + **customize**: Use a custom domain as the back-to-origin host domain.

* `customize_domain` - The name of origin host. Return the host domain set by the primary origin site
  of the accelerate domain. If the accelerate domain has multiple primary origin sites and corresponds
  to multiple back-to-origin hosts, the host domain corresponding to the first primary origin site in
  the origin site configuration will be returned.
