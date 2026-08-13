# Variable definitions for authentication
variable "region_name" {
  description = "The region where resources will be created"
  type        = string
}

variable "access_key" {
  description = "The access key of the IAM user"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "The secret key of the IAM user"
  type        = string
  sensitive   = true
}

variable "parameter_template_configuration" {
  description = "The configuration of the parameter template"
  type        = map(object({
    template_name        = string
    template_description = string
    parameter_values     = map(string)
  }))

  default = {}
}

variable "datastore_type" {
  description = "The database type. Valid values are cassandra, mongodb, influxdb, redis, dynamodb, hbase"
  type        = string
  default     = "cassandra"
}

variable "datastore_version" {
  description = "The database version. Cassandra: 3.11, Mongo: 4.0, Influx: 1.8, Redis: 5.0"
  type        = string
  default     = "3.11"
}

variable "datastore_mode" {
  description = "The database instance mode. Valid values are ReplicaSet, InfluxdbSingle, EnhancedCluster, CloudNativeCluster"
  type        = string
  default     = "CloudNativeCluster"
}
