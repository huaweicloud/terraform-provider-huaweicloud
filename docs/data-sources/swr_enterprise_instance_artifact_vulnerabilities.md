---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_artifact_vulnerabilities"
description: |-
  Use this data source to get the list of SWR enterprise instance artifact vulnerabilities.
---

# huaweicloud_swr_enterprise_instance_artifact_vulnerabilities

Use this data source to get the list of SWR enterprise instance artifact vulnerabilities.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "repository_name" {}
variable "reference" {}

data "huaweicloud_swr_enterprise_instance_artifact_vulnerabilities" "test" {
  instance_id     = var.instance_id
  namespace_name  = var.namespace_name
  repository_name = var.repository_name
  reference       = var.reference
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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `reports` - Indicates the reports.
  The [reports](#attrblock--reports) structure is documented below.

<a name="attrblock--reports"></a>
The `reports` block supports:

* `content` - Indicates the report content.
  The [content](#attrblock--reports--content) structure is documented below.

* `type` - Indicates the report type.

<a name="attrblock--reports--content"></a>
The `content` block supports:

* `generated_at` - Indicates the generation time.

* `scanner` - Indicates the scanner infos.
  The [scanner](#attrblock--reports--content--scanner) structure is documented below.

* `severity` - Indicates the report severity.

* `vulnerabilities` - Indicates the vulnerabilities.
  The [vulnerabilities](#attrblock--reports--content--vulnerabilities) structure is documented below.

<a name="attrblock--reports--content--scanner"></a>
The `scanner` block supports:

* `name` - Indicates the scanner name.

* `vendor` - Indicates the vendor of the scanner.

* `version` - Indicates the version of the scanner.

<a name="attrblock--reports--content--vulnerabilities"></a>
The `vulnerabilities` block supports:

* `id` - Indicates the vulnerability ID.

* `severity` - Indicates the severity of the vulnerability.

* `artifact_digests` - Indicates the digests containing this vulnerability.

* `description` - Indicates the description of the vulnerability.

* `links` - Indicates the links of the vulnerability.

* `package` - Indicates the package name which has vulnerability.

* `cwe_ids` - Indicates the CWE IDs of the vulnerability.

* `preferred_cvss` - Indicates the vulnerability scoring and attack analysis based on CVSS3 and CVSS2.
  The [preferred_cvss](#attrblock--reports--content--vulnerabilities--preferred_cvss) structure is documented below.

* `version` - Indicates the package version which has vulnerability.

* `fix_version` - Indicates the fixed package version of the vulnerability.

<a name="attrblock--reports--content--vulnerabilities--preferred_cvss"></a>
The `version` block supports:

* `score_v2` - Indicates the vCVSS-2 score of the vulnerability.

* `score_v3` - Indicates the CVSS-3 score of the vulnerability.

* `vector_v2` - Indicates the CVSS-2 attack vector of the vulnerability.

* `vector_v3` - Indicates the CVSS-3 attack vector of the vulnerability.
