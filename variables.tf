variable "aws_region" {
  description = "AWS region for deployment"
  type        = string
  default     = "us-east-1"
}

variable "email_address" {
  description = "Email address for SNS alert notifications"
  type        = string
}
