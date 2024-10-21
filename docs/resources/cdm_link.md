---
subcategory: "Cloud Data Migration (CDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_link"
description: ""
---

# huaweicloud_cdm_link

Manages a link resource within HuaweiCloud. A link enables the CDM cluster to read data from and write data to
 a data source.

## Example Usage

### Link to OBS

```hcl
variable "obs_name" {}
variable "obs_link_name" {}
variable "cdm_cluster_id" {}
variable "access_key" {}
variable "secret_key" {}

resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = var.obs_name
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_cdm_link" "obsLink" {
  name       = var.obs_link_name
  connector  = "obs-connector"
  cluster_id = var.cdm_cluster_id

  config = {
    "storageType" = "OBS"
    "server"      = trimprefix(huaweicloud_obs_bucket.bucket.bucket_domain_name, "${huaweicloud_obs_bucket.bucket.bucket}.")
    "port"        = "443"
    "properties"  = jsonencode(
      {
        connectionTimeout = "10000",
        socketTimeout     = "20000"
      }
    )
  }
  
  access_key   = var.access_key
  secret_key   = var.secret_key
}
```

### Link to MySql

```hcl
variable "mysql_link_name" {}
variable "cdm_cluster_id" {}
variable "mysql_host" {}
variable "db_name" {}
variable "db_password" {}

resource "huaweicloud_cdm_link" "mysqlLink" {
  name       = var.mysql_link_name
  connector  = "generic-jdbc-connector"
  cluster_id = var.cdm_cluster_id

  config = {
    "databaseType" = "MYSQL"
    "host"         = var.mysql_host
    "port"         = "3306"
    "database"     = var.db_name
    "username"     = var.username
  }
  password = var.db_password
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the link.

* `cluster_id` - (Required, String, ForceNew) Specifies the id of CDM cluster which this link belongs to.
 Changing this parameter will create a new resource.

* `connector` - (Required, String, ForceNew) Specifies the connector that is classified based on the type of the data
 source to be connected. Changing this parameter will create a new resource. The options are as follows:

  - **generic-jdbc-connector**: link to a relational database.
  - **obs-connector**: link to OBS.
  - **hdfs-connector**: link to HDFS.
  - **hbase-connector**: link to HBase and link to CloudTable.
  - **hive-connector**: link to Hive.
  - **ftp-connector/sftp-connector**: link to an FTP or SFTP server.
  - **mongodb-connector**: link to MongoDB.
  - **redis-connector**: link to Redis.
  - **kafka-connector**: link to Kafka.
  - **dis-connector**: link to DIS.
  - **elasticsearch-connector**: link to Elasticsearch/Cloud Search Service.
  - **dli-connector**: link to DLI.
  - **opentsdb-connector**: link to CloudTable OpenTSDB.
  - **dms-kafka-connector**: link to DMS Kafka.

* `config` - (Required, Map) Specifies the link configuration parameters. Each type of the data source to be connected
 has different configuration parameters, please refer to the document link below.

  - **Link to a Relational Database**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0030.html)
  - **Link to OBS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0031.html)
  - **Link to HDFS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0032.html)
  - **Link to HBase**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0033.html)
  - **Link to CloudTable**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0085.html)
  - **Link to Hive**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0034.html)
  - **Link to an FTP or SFTP Server**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0036.html)
  - **Link to MongoDB**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0072.html)
  - **Link to Redis**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0037.html)
  - **Link to Kafka**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0035.html)
  - **Link to DIS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0075.html)
  - **Link to Elasticsearch/Cloud Search Service**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0076.html)
  - **Link to DLI**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0070.html)
  - **Link to CloudTable OpenTSDB**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0089.html)
  - **Link to DMS Kafka**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0095.html)

-> Please remove the `linkconfig.` in the parameter key listed in the document. Configuration parameters such as
 `password`, `ak`, `sk`, `accessKey` and `securityKey` do not need to be specified in `config`, they are specified in
  the following parameters.

* `password` - (Optional, String) Specifies the password for accessing the data sources.

* `access_key` - (Optional, String) Specifies access key for accessing the data sources.
  
* `secret_key` - (Optional, String) Specifies security key for accessing the data sources.

* `enabled` - (Optional, Bool) Specifies whether to enable the link. The default value is `true`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of **cluster_id/link_name**. It is composed of the ID of CDM cluster and the name
 of job, separated by a slash.

## Import

The link can be imported by `id`, It is composed of the ID of CDM cluster and the name of job, separated by a slash.
 For example,

```bash
terraform import huaweicloud_cdm_link.test b11b407c-e604-4e8d-8bc4-92398320b847/linkName
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`, `secret_key` and `config`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cdm_link" "test" {
    ...

  lifecycle {
    ignore_changes = [
      password, secret_key, config,
    ]
  }
}
```
