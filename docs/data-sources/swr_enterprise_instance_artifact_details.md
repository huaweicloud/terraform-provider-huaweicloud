---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_artifact_details"
description: |-
  Use this data source to get the list of SWR enterprise artifact details.
---

# huaweicloud_swr_enterprise_instance_artifact_details

Use this data source to get the list of SWR enterprise artifact details.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "repository_name" {}
variable "reference" {}

data "huaweicloud_swr_enterprise_instance_artifact_details" "test" {
  instance_id        = var.instance_id
  namespace_name     = var.namespace_name
  repository_name    = var.repository_name
  reference          = var.reference
  with_scan_overview = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String) Specifies the namespace name.

* `repository_name` - (Required, String) Specifies the repository name.

* `reference` - (Required, String) Specifies the artifact digest.

* `with_scan_overview` - (Optional, String) Specifies whether to return the scan overview infos.
  Value can be **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `artifact_id` - Indicates the artifact ID.

* `digest` - Indicates the digest.

* `manifest_media_type` - Indicates the manifest media type.

* `media_type` - Indicates the media type.

* `namespace_id` - Indicates the namespace ID.

* `pull_time` - Indicates the last pull time.

* `push_time` - Indicates the last push time.

* `repository_id` - Indicates the repository ID.

* `scan_overview` - Indicates the report scan overview.
  The [scan_overview](#attrblock--scan_overview) structure is documented below.

* `size` - Indicates the artifact size, unit is byte.

* `tags` - Indicates the artifact version tags.
  The [tags](#attrblock--tags) structure is documented below.

* `type` - Indicates the artifact type.

<a name="attrblock--scan_overview"></a>
The `scan_overview` block supports:

* `overview` - Indicates the report content.
  The [overview](#attrblock--scan_overview--overview) structure is documented below.

* `type` - Indicates the report type.

<a name="attrblock--scan_overview--overview"></a>
The `overview` block supports:

* `report_id` - Indicates the report ID.

* `complete_percent` - Indicates the completed percent.

* `duration` - Indicates the duration.

* `end_time` - Indicates the end time of the report.

* `scan_status` - Indicates the scan status.

* `scanner` - Indicates the scanner infos.
  The [scanner](#attrblock--scan_overview--overview--scanner) structure is documented below.

* `severity` - Indicates the severity.

* `start_time` - Indicates the start time of the report.

* `summary` - Indicates the vulnerabilities summary.
  The [summary](#attrblock--scan_overview--overview--summary) structure is documented below.

<a name="attrblock--scan_overview--overview--scanner"></a>
The `scanner` block supports:

* `name` - Indicates the scanner name.

* `vendor` - Indicates the vendor of the scanner.

* `version` - Indicates the version of the scanner.

<a name="attrblock--scan_overview--overview--summary"></a>
The `summary` block supports:

* `fixable` - Indicates the fixable count of the vulnerability.

* `summary` - Indicates the summary of the different level vulnerability.

* `total` - Indicates the vulnerability counts.

<a name="attrblock--tags"></a>
The `tags` block supports:

* `id` - Indicates the tag ID.

* `name` - Indicates the tag name.

* `artifact_id` - Indicates the artifact ID.

* `pull_time` - Indicates the pull time of the tag.

* `push_time` - Indicates the push time of the tag.

* `repository_id` - Indicates the repository ID.
