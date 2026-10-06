resource "aws_lambda_function" "finops_scanner" {
  filename         = "${path.module}/lambda_function.zip"
  function_name    = "finops-waste-scanner"
  role             = aws_iam_role.lambda_role.arn
  handler          = "lambda_function.lambda_handler"
  source_code_hash = filebase64sha256("${path.module}/lambda_function.zip")
  runtime          = "python3.12"
  timeout          = 60

  environment {
    variables = {
      SNS_TOPIC_ARN = aws_sns_topic.finops_alerts.arn
    }
  }
}
