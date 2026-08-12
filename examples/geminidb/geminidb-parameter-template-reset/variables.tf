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

# Variable definitions for the parameter template
variable "template_name" {
  description = "The name of the GeminiDB parameter template"
  type        = string
}

variable "template_description" {
  description = "The description of the GeminiDB parameter template"
  type        = string
  default     = "test configuration for reset"
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

variable "parameter_values" {
  description = "The parameter values map for the GeminiDB parameter template"
  type        = map(string)
  default     = {
    request_timeout_in_ms = "20000"
  }
}
