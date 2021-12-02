import oprf
import oprfs
import bcl
import base64


def tests():
    k = oprfs.key()
    print(f"Key : {k}")

    # base64_key = base64.standard_b64decode(k)
    # print(f"Base64 Key : {base64_key}")

    ###
    # r = oprfs.handler(k, {})
    oprf_mask = oprf.mask()
    print(f"oprf mask : {oprf_mask}, {type(oprf_mask)}")

    mask_encrypted = bcl.symmetric.encrypt(k, oprf_mask)
    print(f"Mask encrypted : {mask_encrypted}, {type(mask_encrypted)}")

    base64_mask_encrypted = base64.standard_b64encode(mask_encrypted).decode("utf-8")
    print(f"Base64 Mask : {base64_mask_encrypted}")
    ###

    new_mask_encrypted = bcl.cipher(
        base64.standard_b64decode(base64_mask_encrypted.encode("utf-8"))
    )  # bcl.cipher to avoir the error 'can only decrypt a ciphertext'
    print(f"New mask encrypted : {new_mask_encrypted}, {type(new_mask_encrypted)}")

    mask_decrypted = bcl.symmetric.decrypt(
        k, new_mask_encrypted
    )  # OK with mask_encrypted
    print(f"Mask decrypted : {mask_decrypted}")
    assert oprf_mask == mask_decrypted

    # r = oprfs.handler(k, {})
    # assert r['status'] == 'success'

    # r = oprfs.handler(k, '{}')
    # assert r['status'] == 'success'

    # m = oprf.mask.from_base64(r['mask'][0])
    # d = oprf.data.hash('abc')
    # assert oprfs.mask(k, m, d) == oprf.mask(bcl.symmetric.decrypt(k, m))(d)

    # r = oprfs.handler(k, {'mask': [m.to_base64()], 'data': [d.to_base64()]})
    # assert r['status'] == 'success'

    # (m_str, d_str) = (str(m.to_base64()), str(d.to_base64()))
    # r = oprfs.handler(k, '{"mask": ["' + m_str + '"], "data": ["' + d_str + '"]}')
    # assert r['status'] == 'success'

    # assert oprf.data.from_base64(r['data'][0]) == oprf.mask(bcl.symmetric.decrypt(k, m))(d)

    # r = oprfs.handler(k, {'mask': [m.to_base64()]})
    # assert r['status'] == 'failure'


def main():
    tests()
    # key = oprfs.key()

    # r = oprfs.handler(key, {})
    # print(f"First request result : {r}")

    # print(
    #     {
    #         "status": "success",
    #         "mask": [base64.standard_b64encode(bcl.symmetric.encrypt(key, oprf.mask())).decode("utf-8")],
    #     }
    # )

    # # Checks

    # mask_encrypted = r["mask"][0]
    # print(f"Mask encrypted : {mask_encrypted}")

    # mask = oprf.mask.from_base64(mask_encrypted)
    # print(f"Encrypted mask : {mask}")
    # data = oprf.data.hash("abc")
    # print(f"Data : {data}")
    # assert oprfs.mask(key, mask, data) == oprf.mask(bcl.symmetric.decrypt(key, mask))(data)

    # rr = oprfs.handler(key, {"mask": [mask_encrypted], "data": [data]})
    # print(f"Second request result : {rr}")

    # data_masked = oprf.data.from_base64(rr["data"][0])
    # print(f"Data masked : {data_masked}")


if __name__ == "__main__":
    main()
