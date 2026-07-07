---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_datastore_versions"
description: |-
  Use this data source to get the list of database versions within HuaweiCloud.
---

# huaweicloud_gaussdb_datastore_versions

Use this data source to get the list of database versions within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_datastore_versions" "test" {}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the data source.

* `database_versions` - Indicates the list of database versions.

  The [database_versions](#database_versions_struct) structure is documented below.

<a name="database_versions_struct"></a>
The `database_versions` block supports:

* `software_version` - Indicates the three-digit DB engine version.

* `hotfixes` - Indicates the hot patch information corresponding to the three-digit DB engine version.

  The [hotfixes](#hotfixes_struct) structure is documented below.

<a name="hotfixes_struct"></a>
The `hotfixes` block supports:

* `version` - Indicates the hot patch version.

* `create_time` - Indicates the creation time of the hot patch. The value is a UNIX timestamp in milliseconds.

* `deploy_modes` - Indicates the deployment model information of the instance for which the hot patch can be installed.

  The [deploy_modes](#deploy_modes_struct) structure is documented below.

<a name="deploy_modes_struct"></a>
The `deploy_modes` block supports:

* `default_upgrade` - Indicates the Default patch fix policy.
  + **true**: You do not need to select the hot patch version to be installed.
  The current hot patch is installed by default.
  + **false**: You need to select the hot patch version to be installed.

* `update_time` - Indicates the modify time of the hot patch installation policy.
  The value is a UNIX timestamp in milliseconds.

* `mode` - Indicates the deployment model of the instance for which the patch can be installed.
  + **distributed**: The deployment model can be either independent or combined (basic edition).
  + **centralization_standard**: The deployment model can be centralized.
