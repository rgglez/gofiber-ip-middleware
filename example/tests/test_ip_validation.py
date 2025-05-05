import os
import pytest
import requests

ENDPOINT = "http://localhost:3000"

def test_hello_world_endpoint():
    # Define the URL of the endpoint
    url = ENDPOINT

    # Send the GET request
    response = requests.get(url)

    # Assert the status code
    assert response.status_code == 200, f"Expected status code 200, but got {response.status_code}"

    # Assert the response text
    assert response.text == "Hello world", f"Expected response text 'Hello world', but got {response.text}"
