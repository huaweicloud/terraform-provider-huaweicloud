---
subcategory: "Document Database Service (DDS)"
---

# huaweicloud_dds_api_versions

Use this data source to get the list of DDS API versions.

## Example Usage

```hcl
data "huaweicloud_dds_api_versions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - Indicates the list of DDS API versions.
  The [versions](#Dds_versions) structure is documented below.

<a name="Dds_versions"></a>
The `versions` block supports:

* `id` - Indicates the API version. The values are as follows:
  + **v1**: API v1 version.
  + **v3**: API v3 version.

* `status` - Indicates the version status. The values are as follows:
  + **CURRENT**: recommended version.
  + **DEPRECATED**: deprecated version which may be deleted later.

* `updated` - Indicates the version update time.

* `links` - Indicates the API links. The value is empty when the version is **v1** or **v3**.
  The [links](#Dds_links) structure is documented below.

<a name="Dds_links"></a>
The `links` block supports:

* `href` - Indicates the API url. The value is "".

* `rel` - Indicates the href is a local link. The value is **self**.
