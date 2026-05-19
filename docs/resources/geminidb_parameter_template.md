---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_parameter_template"
description: |-
  Manages a GeminiDB parameter template resource within HuaweiCloud.
---

# huaweicloud_geminidb_parameter_template

Manages a GeminiDB parameter template resource within HuaweiCloud.

## Example Usage

```hcl
var "instance_id" {}

resource "huaweicloud_geminidb_parameter_template" "test" {
  name        = "test-configuration"
  description = "test configuration with custom values"

  datastore {
    type    = "cassandra"
    version = "3.11"
  }

  values = {
    max_connections  = "100"
    concurrent_reads = "64"
  }
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `name` - (Required, String) Specifies the parameter template name.

* `description` - (Optional, String) Specifies the parameter template description.

* `values` - (Optional, Map) Specifies the parameter values.

* `datastore` - (Optional, List, NonUpdatable) Specifies the database object.
  If `instance_id` is not specified, this parameter is required.
  The [datastore](#geminidb_configuration_datastore) structure is documented below.

* `instance_id` - (Optional, String, NonUpdatable) Specifies the instance ID.

<a name="geminidb_configuration_datastore"></a>
The `datastore` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the database type.
  The valid values are as follows:
  + **cassandra**: GeminiDB Cassandra.
  + **mongodb**: GeminiDB Mongo.
  + **influxdb**: GeminiDB Influx.
  + **redis**: GeminiDB Redis.
  + **dynamodb**: GeminiDB compatible with DynamoDB.
  + **hbase**: GeminiDB HBase.

* `version` - (Required, String, NonUpdatable) Specifies the database version.
  The valid values are as follows:
  + GeminiDB Cassandra: **3.11**
  + GeminiDB Mongo: **4.0**
  + GeminiDB Influx: **1.8**
  + GeminiDB Redis: **5.0**

* `mode` - (Optional, String, NonUpdatable) Specifies the database instance mode.
  This parameter is required when creating parameter templates for:
  GeminiDB Mongo ReplicaSet, GeminiDB Influx single node, GeminiDB Influx enhanced cluster,
  GeminiDB Influx CloudNativeCluster, GeminiDB Cassandra CloudNativeCluster,
  or GeminiDB Redis CloudNativeCluster.
  The valid values are as follows:
  + **ReplicaSet**: GeminiDB Mongo ReplicaSet.
  + **InfluxdbSingle**: GeminiDB Influx classic deployment mode single node.
  + **EnhancedCluster**: GeminiDB Influx classic deployment mode enhanced cluster.
  + **CloudNativeCluster**: GeminiDB Influx/Cassandra/Redis cloud-native deployment mode.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `datastore_version_name` - The database version name.

* `datastore_name` - The database name.

* `mode` - The database instance mode.

* `created` - The creation time. The format is **yyyy-MM-ddTHH:mm:ssZ**.

* `updated` - The update time. The format is **yyyy-MM-ddTHH:mm:ssZ**.

* `configuration_parameters` - The list of parameter configurations.
  The [configuration_parameters](#geminidb_configuration_parameters) structure is documented below.

<a name="geminidb_configuration_parameters"></a>
The `configuration_parameters` block supports:

* `name` - The parameter name.

* `value` - The parameter value.

* `restart_required` - The parameter whether a restart is required after modifying.

* `readonly` - The parameter whether is read-only.

* `value_range` - The parameter value range.

* `type` - The parameter type. Valid values: **string**, **integer**, **boolean**, **list**, **float**.

* `description` - The parameter description.

## Import

The GeminiDB parameter template can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_geminidb_parameter_template.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `instance_id` and `datastore`.
It is generally recommended running `terraform plan` after importing a parameter template.
You can then decide if changes should be applied to the parameter template, or the resource definition should be
updated to align with the parameter template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_geminidb_parameter_template" "test" {
  ...

  lifecycle {
    ignore_changes = [
      instance_id, datastore,
    ]
  }
}
