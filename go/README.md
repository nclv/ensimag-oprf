```bash
curl -X POST http://localhost:1323/evaluate -H 'Content-Type: application/json' -d 
'{"suite": 3, "mode": 1, "blinded_elements": [[49, 50, 51, 52],[50, 51, 51]]}'  # "blinded_elements": ["1234","233"]

{"suite":3,"mode":1,"blinded_elements":["MTIzNA==","MjMz"]}  # Base64 encoded strings
```