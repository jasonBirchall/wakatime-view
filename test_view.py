import pytest
from unittest.mock import Mock, patch
from view import today, setup  # importing from your script named view.py

# Test for the "today" command
def test_today():
    with patch('view.get_api_key', return_value="test_api_key"), \
         patch('requests.get') as mock_get:
        
        # Mock the API response
        mock_get.return_value.status_code = 200
        mock_get.return_value.json.return_value = {'data': {'grand_total': {'text': '5 hrs'}}}
        
        # Call the "today" command
        today()
        
        # Assert that the API was called with the correct endpoint and headers
        mock_get.assert_called_once_with(
            'https://wakatime.com/api/v1/users/current/status_bar/today', 
            headers={'Authorization': 'Basic dGVzdF9hcGlfa2V5'}
        )

# Test for the "setup" command
def test_setup():
    with patch('typer.prompt', return_value="test_api_key"), \
         patch('builtins.open', new_callable=Mock) as mock_open, \
         patch('toml.dump') as mock_dump, \
         patch('view.CONFIG_FILE.is_file', return_value=False):
        
        # Call the "setup" command
        setup()
        
        # Assert that the config file was written with the correct content
        mock_open.assert_called_once_with(view.CONFIG_FILE, 'w')  # using view.CONFIG_FILE directly
        mock_dump.assert_called_once_with({'wakatime': {'api_key': 'test_api_key'}}, mock_open())
