import json
import requests
import oprf

# Request an encrypted mask.
response = requests.post("http://localhost:5000", json={})
mask_encrypted = json.loads(response.text)["mask"][0]

# Mask some data.
data = oprf.data.hash("abc").to_base64()
response = requests.post(
    "http://localhost:5000", json={"mask": [mask_encrypted], "data": [data]}
)

print(f"Response : {response.text}")
data_masked = oprf.data.from_base64(json.loads(response.text)["data"][0])
