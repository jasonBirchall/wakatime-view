import os
import shutil
import unittest
from typer.testing import CliRunner
import view


# Test for the "setup" command
class TestSetupCommand(unittest.TestCase):
    def setUp(self):
        self.runner = CliRunner()
        path = view.CONFIG_FILE.__str__()
        if os.path.exists(view.CONFIG_FILE):
            shutil.move(path, path + ".bak")

    def tearDown(self):
        path = view.CONFIG_FILE.__str__()
        if os.path.exists(path):
            os.remove(path)
            if os.path.exists(path + ".bak"):
                shutil.move(path + ".bak", path)

    def test_create_config_file(self):
        result = self.runner.invoke(view.app, ["setup"], input="123456789")
        assert result.exit_code == 0
        print(result.output)
        assert "Config file created. Please edit the file and add your API key." in result.output


if __name__ == "__main__":
    unittest.main()
