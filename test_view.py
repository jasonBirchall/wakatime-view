import os
import shutil
from unittest.mock import patch
import unittest
from typer.testing import CliRunner
import view

PATH = view.CONFIG_FILE.__str__()


class TestTodayCommand(unittest.TestCase):
    def setUp(self):
        self.runner = CliRunner()
        if os.path.exists(view.CONFIG_FILE):
            shutil.move(PATH, PATH + ".bak")

    def tearDown(self):
        if os.path.exists(PATH):
            os.remove(PATH)
            if os.path.exists(PATH + ".bak"):
                shutil.move(PATH + ".bak", PATH)


class TestSetupCommand(unittest.TestCase):
    def setUp(self):
        self.runner = CliRunner()
        if os.path.exists(view.CONFIG_FILE):
            shutil.move(PATH, PATH + ".bak")

    def tearDown(self):
        if os.path.exists(PATH):
            os.remove(PATH)
            if os.path.exists(PATH + ".bak"):
                shutil.move(PATH + ".bak", PATH)

    def test_create_config_file(self):
        result = self.runner.invoke(view.app, ["setup"], input="123456789")
        assert result.exit_code == 0
        assert "Config file created at:" + PATH

    def test_config_file_exists(self):
        result = self.runner.invoke(view.app, ["setup"], input="123456789")
        assert result.exit_code == 0
        result = self.runner.invoke(view.app, ["setup"], input="123456789")
        assert result.exit_code == 0
        assert "Config file already exists. Exiting..." in result.output


if __name__ == "__main__":
    unittest.main()
