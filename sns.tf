resource "aws_sns_topic" "finops_alerts" {
  name = "finops-alerts-topic"
}

resource "aws_sns_topic_subscription" "email_sub" {
  topic_arn = aws_sns_topic.finops_alerts.arn
  protocol  = "email"
  endpoint  = var.email_address
}
