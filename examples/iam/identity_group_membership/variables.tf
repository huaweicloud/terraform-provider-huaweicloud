variable "user_A_name" {
  type    = string
  default = "user_A"
}

variable "user_A_password" {
  type    = string
  default = "hw123456"
}

variable "iden_group" {
    description = "List of identity group."
    type = list(object({
        name        = string
        description = string
    }))
    default = [
        {name = "iden_group_A", description = "This is a identity group."},
        {name = "iden_group_B", description = "This is a identity group."}
    ]
}