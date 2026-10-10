# 1. Create GitHub OIDC Provider
resource "aws_iam_openid_connect_provider" "github" {
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]
}

# 2. Define IAM Role Trust Policy
data "aws_iam_policy_document" "github_actions_assume_role" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }

    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:achh-914/finops-waste-scanner:*"]
    }
  }
}

# 3. Create IAM Role for GitHub Actions
resource "aws_iam_role" "github_actions_finops" {
  name               = "GitHubActions-FinOpsScanner-Role"
  assume_role_policy = data.aws_iam_policy_document.github_actions_assume_role.json
}

# 4. Attach Permissions Policy
resource "aws_iam_role_policy_attachment" "finops_read_only" {
  role       = aws_iam_role.github_actions_finops.name
  policy_arn = "arn:aws:iam::aws:policy/SecurityAudit"
}

# 5. Output Role ARN
output "github_actions_role_arn" {
  value       = aws_iam_role.github_actions_finops.arn
  description = "The ARN of the IAM Role for GitHub Actions OIDC"
}
