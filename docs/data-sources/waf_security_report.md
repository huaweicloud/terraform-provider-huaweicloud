---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_security_report"
description: |-
  Use this data source to get the security report of WAF within HuaweiCloud.
---

# huaweicloud_waf_security_report

Use this data source to get the security report of WAF within HuaweiCloud.

## Example Usage

```hcl
variable "report_id" {}
variable "subscription_id" {}

data "huaweicloud_waf_security_report" "test" {
  report_id       = var.report_id
  subscription_id = var.subscription_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `report_id` - (Required, String) Specifies the report ID.

* `subscription_id` - (Required, String) Specifies the subscription ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sending_period` - The sending time period indicates the preset sending time for the report.

* `report_name` - The report name.

* `topic_urn` - The topic urn is the unique identifier of the SMN topic to which the associated report is sent.

* `subscription_type` - The subscription type, indicating how the security report is subscribed.

* `stat_period` - The statistical period indicates the time range within which the current security report statistics
  are collected.

  The [stat_period](#stat_period_struct) structure is documented below.

* `report_category` - The report category.

* `report_content_info` - The subscribe to the report, which includes various statistical data details for the security
  report.

  The [report_content_info](#report_content_info_struct) structure is documented below.

* `create_time` - The creation time of the report.

<a name="stat_period_struct"></a>
The `stat_period` block supports:

* `begin_time` - The timestamp (in milliseconds) at the start of the statistical period.

* `end_time` - The end timestamp (in milliseconds) of the statistical period.

<a name="report_content_info_struct"></a>
The `report_content_info` block supports:

* `request_statistics_info_list` - The request count statistics, including request counts for various dimensions and
  timelines.

  The [request_statistics_info_list](#report_content_info_request_statistics_info_list_struct) structure is documented
  below.

* `bandwidth_statistics_info` - The bandwidth statistics, including timeline statistics for average and peak bandwidth
  across various dimensions.

  The [bandwidth_statistics_info](#report_content_info_bandwidth_statistics_info_struct) structure is documented below.

* `response_code_statistics_info` - The response code statistics, including timeline statistics for each response code
  from the WAF and upstream servers.

  The [response_code_statistics_info](#report_content_info_response_code_statistics_info_struct) structure is documented
  below.

* `attack_type_distribution_info_list` - The attack type distribution statistics list contains the number of attacks for
  each attack type.

  The [attack_type_distribution_info_list](#report_content_info_attack_type_distribution_info_list_struct) structure is
  documented below.

* `top_attacked_urls_info_list` - The top attacked URL list contains URLs that have been attacked the most times, sorted
  by attack count.

  The [top_attacked_urls_info_list](#report_content_info_top_attacked_urls_info_list_struct) structure is documented below.

* `qps_statistics_info` - QPS statistics information, containing the average QPS and peak QPS for each dimension of the
  timeline.

  The [qps_statistics_info](#report_content_info_qps_statistics_info_struct) structure is documented below.

* `top_attacked_domains_info_list` - The top attacked domain list contains domains that have been attacked the most times,
  sorted by attack count.

  The [top_attacked_domains_info_list](#report_content_info_top_attacked_domains_info_list_struct) structure is documented
  below.

* `top_attack_source_ips_info_list` - The top attack source IP list contains IP addresses that have been attacked the most
  times, sorted by attack count.

  The [top_attack_source_ips_info_list](#report_content_info_top_attack_source_ips_info_list_struct) structure is documented
  below.

* `top_attack_source_locations_info_list` - The top attack source location list contains locations that have been
  attacked the most times, sorted by attack count.

  The [top_attack_source_locations_info_list](#report_content_info_top_attack_source_locations_info_list_struct) structure
  is documented below.

* `top_abnormal_urls_info` - The top abnormal URL list contains URLs that have returned 502, 500, 404, etc. errors.

  The [top_abnormal_urls_info](#report_content_info_top_abnormal_urls_info_struct) structure is documented below.

* `overview_statistics_list_info` - The overview statistics list contains summary statistics for each dimension and the
  top domain details.

  The [overview_statistics_list_info](#report_content_info_overview_statistics_list_info_struct) structure is documented
  below.

<a name="report_content_info_request_statistics_info_list_struct"></a>
The `request_statistics_info_list` block supports:

* `key` - The statistical dimension identifier.

* `timeline` - The timeline data, statistics arranged in chronological order.

  The [timeline](#request_statistics_info_list_timeline_struct) structure is documented below.

<a name="request_statistics_info_list_timeline_struct"></a>
The `timeline` block supports:

* `time` - The timestamp (millisecond level) identifies the time point corresponding to the statistical data.

* `num` - The number of statistical dimensions at this point in time.

<a name="report_content_info_bandwidth_statistics_info_struct"></a>
The `bandwidth_statistics_info` block supports:

* `average_info_list` - The average bandwidth statistics list contains average bandwidth data for various dimensions
  over time.

  The [average_info_list](#bandwidth_statistics_info_average_info_list_struct) structure is documented below.

* `peak_info_list` - The peak bandwidth statistics list contains peak bandwidth data for various dimensions and timelines.

  The [peak_info_list](#bandwidth_statistics_info_peak_info_list_struct) structure is documented below.

<a name="bandwidth_statistics_info_average_info_list_struct"></a>
The `average_info_list` block supports:

* `key` - The statistical dimension identifier.

* `timeline` - The timeline data, average bandwidth values ​​arranged in chronological order.

  The [timeline](#average_info_list_timeline_struct) structure is documented below.

<a name="average_info_list_timeline_struct"></a>
The `timeline` block supports:

* `num` - The average bandwidth value for the statistical dimension at this point in time.

* `time` - The timestamp (in milliseconds) identifies the point in time corresponding to statistical data.

<a name="bandwidth_statistics_info_peak_info_list_struct"></a>
The `peak_info_list` block supports:

* `key` - The statistical dimension identifier.

* `timeline` - The timeline data, peak bandwidth values ​​arranged in chronological order.

  The [timeline](#peak_info_list_timeline_struct) structure is documented below.

<a name="peak_info_list_timeline_struct"></a>
The `timeline` block supports:

* `time` - The timestamp (in milliseconds) identifies the point in time corresponding to statistical data.

* `num` - The peak bandwidth value corresponding to this time point in the statistical dimension.

<a name="report_content_info_response_code_statistics_info_struct"></a>
The `response_code_statistics_info` block supports:

* `response_source_waf_info_list` - The WAF response code statistics list, including the number of WAF responses for
  each response code over time.

  The [response_source_waf_info_list](#response_code_statistics_info_response_source_waf_info_list_struct) structure
  is documented below.

* `response_source_upstream_info_list` - The upstream response code statistics list, including the number of upstream
  responses for each response code over time.

  The [response_source_upstream_info_list](#response_code_statistics_info_response_source_upstream_info_list_struct) structure
  is documented below.

<a name="response_code_statistics_info_response_source_waf_info_list_struct"></a>
The `response_source_waf_info_list` block supports:

* `key` - The response code identifier.

* `timeline` - The timeline data, WAF response code counts arranged in chronological order.

  The [timeline](#response_source_waf_info_list_timeline_struct) structure is documented below.

<a name="response_source_waf_info_list_timeline_struct"></a>
The `timeline` block supports:

* `time` - The timestamp (in milliseconds) identifies the point in time corresponding to statistical data.

* `num` - The number of WAF responses for the response code at this point in time.

<a name="response_code_statistics_info_response_source_upstream_info_list_struct"></a>
The `response_source_upstream_info_list` block supports:

* `key` - The response code identifier.

* `timeline` - The timeline data, upstream response code counts arranged in chronological order.

  The [timeline](#response_source_upstream_info_list_timeline_struct) structure is documented below.

<a name="response_source_upstream_info_list_timeline_struct"></a>
The `timeline` block supports:

* `time` - The timestamp (in milliseconds) identifies the point in time corresponding to statistical data.

* `num` - The number of upstream responses for the response code at this point in time.

<a name="report_content_info_attack_type_distribution_info_list_struct"></a>
The `attack_type_distribution_info_list` block supports:

* `key` - The attack type identifier.

* `num` - The total number of attacks for this attack type.

<a name="report_content_info_top_attacked_urls_info_list_struct"></a>
The `top_attacked_urls_info_list` block supports:

* `key` - The attacked URL path.

* `num` - The total number of attacks for this URL.

* `host` - The domain name identifier to which this URL belongs (e.g., *:80 indicates all domains on port 80).

<a name="report_content_info_qps_statistics_info_struct"></a>
The `qps_statistics_info` block supports:

* `average_info_list` - The average QPS statistics list contains average QPS data for various dimensions over time.

  The [average_info_list](#qps_statistics_info_average_info_list_struct) structure is documented below.

* `peak_info_list` - The peak QPS statistics list contains peak QPS data for various dimensions and timelines.

  The [peak_info_list](#qps_statistics_info_peak_info_list_struct) structure is documented below.

<a name="qps_statistics_info_average_info_list_struct"></a>
The `average_info_list` block supports:

* `key` - The statistical dimension identifier.

* `timeline` - The timeline data, average QPS values arranged in chronological order.

  The [timeline](#average_info_list_timeline_struct) structure is documented below.

<a name="average_info_list_timeline_struct"></a>
The `timeline` block supports:

* `time` - The timestamp (in milliseconds) identifies the point in time corresponding to statistical data.

* `num` - The average QPS value for this statistical dimension at this point in time.

<a name="qps_statistics_info_peak_info_list_struct"></a>
The `peak_info_list` block supports:

* `key` - The statistical dimension identifier.

* `timeline` - The timeline data, peak QPS values arranged in chronological order.

  The [timeline](#peak_info_list_timeline_struct) structure is documented below.

<a name="peak_info_list_timeline_struct"></a>
The `timeline` block supports:

* `time` - The timestamp (in milliseconds) identifies the point in time corresponding to statistical data.

* `num` - The peak QPS value for this statistical dimension at this point in time.

<a name="report_content_info_top_attacked_domains_info_list_struct"></a>
The `top_attacked_domains_info_list` block supports:

* `key` - The domain name identifier, containing the domain name and port (e.g., *:80 indicates all domains on port 80).

* `num` - The total number of attacks for this domain.

* `web_tag` - The Web tag of the domain, used to identify the business type of the domain.

<a name="report_content_info_top_attack_source_ips_info_list_struct"></a>
The `top_attack_source_ips_info_list` block supports:

* `key` - The attack source IP address.

* `num` - The total number of attacks initiated by this IP.

<a name="report_content_info_top_attack_source_locations_info_list_struct"></a>
The `top_attack_source_locations_info_list` block supports:

* `key` - The attack source location identifier.

* `num` - The total number of attacks initiated from this location.

<a name="report_content_info_top_abnormal_urls_info_struct"></a>
The `top_abnormal_urls_info` block supports:

* `abnormal_502_info_list` - The list of URLs that returned 502 errors, sorted by error count.

  The [abnormal_502_info_list](#top_abnormal_urls_info_abnormal_502_info_list_struct) structure is documented below.

* `abnormal_500_info_list` - The list of URLs that returned 500 errors, sorted by error count.

  The [abnormal_500_info_list](#top_abnormal_urls_info_abnormal_500_info_list_struct) structure is documented below.

* `abnormal_404_info_list` - The list of URLs that returned 404 errors, sorted by error count.

  The [abnormal_404_info_list](#top_abnormal_urls_info_abnormal_404_info_list_struct) structure is documented below.

<a name="top_abnormal_urls_info_abnormal_502_info_list_struct"></a>
The `abnormal_502_info_list` block supports:

* `key` - The URL path that returned 502 errors.

* `num` - The total number of 502 errors for this URL.

* `host` - The domain name to which this URL belongs.

<a name="top_abnormal_urls_info_abnormal_500_info_list_struct"></a>
The `abnormal_500_info_list` block supports:

* `host` - The domain name to which this URL belongs.

* `key` - The URL path that returned 500 errors.

* `num` - The total number of 500 errors for this URL.

<a name="top_abnormal_urls_info_abnormal_404_info_list_struct"></a>
The `abnormal_404_info_list` block supports:

* `key` - The URL path that returned 404 errors.

* `num` - The total number of 404 errors for this URL.

* `host` - The domain name to which this URL belongs.

<a name="report_content_info_overview_statistics_list_info_struct"></a>
The `overview_statistics_list_info` block supports:

* `key` - The statistical dimension identifier.

* `num` - The total number of attacks for this domain.

* `top_domains` - TOP domain list, sorted by attack count.

  The [top_domains](#overview_statistics_list_info_top_domains_struct) structure is documented below.

<a name="overview_statistics_list_info_top_domains_struct"></a>
The `top_domains` block supports:

* `num` - The number of attacks for this domain.

* `host` - The domain name identifier, containing the domain name and associated identifier.
