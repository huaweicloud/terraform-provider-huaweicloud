---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_parameter_templates"
description: |-
  Use this data source to query the parameter templates of GeminiDB within HuaweiCloud.
---

# huaweicloud_geminidb_parameter_templates

Use this data source to query the parameter templates of GeminiDB within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_geminidb_parameter_templates" "test" {}
```

### Filter by datastore name

```hcl
data "huaweicloud_geminidb_parameter_templates" "test" {
  datastore_name = "cassandra"
}
```

### Filter by mode

```hcl
data "huaweicloud_geminidb_parameter_templates" "test" {
  mode = "Cluster"
}
```

### Filter by user-defined

```hcl
data "huaweicloud_geminidb_parameter_templates" "test" {
  user_defined = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the parameter templates.  
  If omitted, the provider-level region will be used.

* `datastore_name` - (Optional, String) Specifies the database name.  
  The valid values are as follows:
  + **cassandra**
  + **redis**
  + **influxdb**
  + **mongodb**

* `mode` - (Optional, String) Specifies the database instance type.  
  The valid values are as follows:
  + **CloudNativeCluster**
  + **Cluster**
  + **InfluxdbSingle**
  + **ReplicaSet**
  + **All**

* `name` - (Optional, String) Specifies the parameter template name.

* `user_defined` - (Optional, String) Specifies whether to query user-defined parameter templates.
  The valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The list of parameter templates.  
  The [configurations](#configurations) structure is documented below.

<a name="configurations"></a>
The `configurations` block supports:

* `id` - The parameter template ID.

* `name` - The parameter template name.

* `description` - The parameter template description.

* `datastore_version_name` - The database version name.

* `datastore_name` - The database name.

* `created` - The creation time, in the format of **yyyy-MM-ddTHH:mm:ssZ**.

* `updated` - The update time, in the format of **yyyy-MM-ddTHH:mm:ssZ**.

* `mode` - The database instance type.

* `user_defined` - Whether the parameter template is user-defined.
