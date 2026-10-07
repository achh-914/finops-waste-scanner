import unittest
import os
import sys

# Ensure src dir is in path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../src')))

class TestFinOpsScanner(unittest.TestCase):

    def test_sample_verification(self):
        """Basic sanity test for CI execution."""
        self.assertTrue(True)

if __name__ == '__main__':
    unittest.main()
