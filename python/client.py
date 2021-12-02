import json
import requests
import oprf

# Request an encrypted mask.
print("Requesting new mask...")
# server computes base64.standard_b64encode(bcl.symmetric.encrypt(key, oprf.mask())).decode("utf-8")
# - generate a new mask
# - encrypt the mask with the symmetric key
# - encode it to base64 for transmission
response = requests.post("http://localhost:5000", json={})

mask_encrypted = json.loads(response.text)["mask"][0]
print(f"Encrypted mask: {mask_encrypted}")

# Mask some data.
data = oprf.data.hash("abc").to_base64()
print(f"Data : {data}")

print("Sending the data...")
# server computes :
# m = oprf.mask.from_base64(mask_encrypted)
# d = oprf.data.from_base64(data)
# base64.standard_b64encode(oprf.mask(bcl.symmetric.decrypt(key, bcl.cipher(m)))(d)).decode("utf-8")
# - decode the encrypted mask and data from base64
# - decrypt the encrypted mask
# - call the mask on the data
# - encode it to base64 for transmission
response = requests.post(
    "http://localhost:5000", json={"mask": [mask_encrypted], "data": [data]}
)
print(f"Response : {response.text}")

data_masked = json.loads(response.text)["data"][0]  # oprf.data.from_base64()
print(f"Masked data : {data_masked}")
