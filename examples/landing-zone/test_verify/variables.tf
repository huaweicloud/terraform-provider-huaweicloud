# Region
variable "region" {
  type    = list(string)
  default = [
    "ap-southeast-3",
    "ap-southeast-1",
    "ap-southeast-2",
    "cn-north-4",
    "cn-south-1",
    "cn-east-3",
    "cn-north-5",
    "cn-southwest-204-dev",
    "cn-east-204-dev"
  ]
}

# Authentication
variable "access_key" {
  type    = string
  default = ""
}

variable "secret_key" {
  type    = string
  default = ""
}

variable "name" {
  type = string
  default = "test_aggregator_account"
}

variable "source_account_list" {
  type = list(string)
  default = [""]
}

variable "topic_urn" {
  type = string
  default = ""
}
