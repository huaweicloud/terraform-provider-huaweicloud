# Variable definitions for authentication
variable "region_name" {
  description = "The region where the DataArts Studio instance is located"
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

# Variable definitions for DataArts Studio instance and workspace
variable "workspace_id" {
  description = "The ID of the workspace"
  type        = string
  default     = ""
  nullable    = false
}

variable "instance_id" {
  description = "The ID of the DataArts Studio instance to which the workspace belongs"
  type        = string
  default     = ""
  nullable    = false
}

variable "workspace_name" {
  description = "The name of the workspace used to filter results"
  type        = string
  default     = ""
  nullable    = false
}

# Variable definitions for data connection
variable "connection_name" {
  description = "The name of the DLI data connection"
  type        = string
  default     = ""
  nullable    = false
}

# Variable definitions for DLI database and table
variable "dli_database_name" {
  description = "The name of the DLI database created for the script"
  type        = string
}

variable "dli_database_description" {
  description = "The description of the DLI database"
  type        = string
  default     = ""
}

variable "dli_table_name" {
  description = "The name of the DLI table created for the script"
  type        = string
}

variable "dli_table_description" {
  description = "The description of the DLI table"
  type        = string
  default     = ""
}

variable "dli_table_columns" {
  description = "The column definitions of the DLI table"
  type        = list(object({
    name = string
    type = string
  }))

  default = [
    {
      name = "name"
      type = "string"
    },
    {
      name = "age"
      type = "int"
    },
  ]
}

# Variable definitions for DataArts Factory script
variable "script_name" {
  description = "The name of the DataArts Factory script"
  type        = string
}

variable "script_type" {
  description = "The type of the DataArts Factory script"
  type        = string
  default     = "DLISQL"
}

variable "script_directory" {
  description = "The directory path where the script is stored in DataArts Factory"
  type        = string
  default     = "/terraform"
}

variable "queue_name" {
  description = "The DLI queue name associated with the script"
  type        = string
  default     = "default"
}

variable "script_description" {
  description = "The description of the DataArts Factory script"
  type        = string
  default     = ""
}

variable "script_configuration" {
  description = "The user-defined configuration parameters of the DataArts Factory script"
  type        = map(string)
  default     = {}
}

variable "script_content" {
  description = "The SQL content of the DataArts Factory script"
  type        = string
  default     = ""
  nullable    = false
}

# Variable definitions for script execution
variable "script_execute_params" {
  description = "The execution parameters passed to the script content, in key-value format"
  type        = map(string)
  default     = {
    "spark.sql.adaptive.enabled"                                   = "true"
    "spark.sql.adaptive.join.enabled"                              = "true"
    "spark.sql.adaptive.join.skewedJoin.enabled"                   = "true"
    "spark.sql.forcePartitionPredicatesOnPartitionedTable.enabled" = "true"
    "spark.sql.mergeSmallFiles.enabled"                            = "true"
  }
}
