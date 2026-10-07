import unittest
from unittest.mock import MagicMock, patch
import os
import sys

# Add src to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '../src')))

class TestFinOpsScanner(unittest.TestCase):

    @patch('boto3.client')
    def test_lambda_handler_execution(self, mock_boto_client):
        """Test Lambda handler execution and notification flow."""
        mock_ec2 = MagicMock()
        mock_ec2.describe_volumes.return_value = {
            'Volumes': [
                {'VolumeId': 'vol-12345', 'State': 'available', 'Size': 20}
            ]
        }
        mock_boto_client.return_value = mock_ec2

        from lambda_function import lambda_handler
        response = lambda_handler({}, None)

        self.assertEqual(response['statusCode'], 200)
        self.assertIn('Unused EBS Volumes', response['body'])

if __name__ == '__main__':
    unittest.main()
