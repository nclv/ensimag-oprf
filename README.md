# ensimag-oprf

[oprfs](https://github.com/nthparty/oprfs) (BROKEN, use [our fork](https://github.com/nclv/oprfs)) use of the distinct [oprf](https://github.com/nthparty/oprf) library to represent a data instance (which is itself a wrapper for an Ed25519 group element as represented by an instance of the point class in the [oblivious](https://github.com/nthparty/oblivious) library). [bcl](https://github.com/nthparty/bcl) provides symmetric and asymmetric encryption.

## Setup

```bash
# Create a virtual environment
python3 -m virtualenv venv-oprf
# Activate the virtual environment
source ./venv-oprf/bin/activate

# Install the requirements (do not contains oprfs from https://github.com/nthparty/oprfs)
pip install -r requirements.txt
# Install fixed oprfs version (from https://github.com/nclv/oprfs)
python setup.py install  # in the directory with venv-oprf activated
```

### oprfs issues

```bash
python oprfs/oprfs.py -v  # tests fails with decryption
```