# mini-key-server

This web application provides a restful API for your desktop and other applications licensing needs.

### Key View

![key view](etc/KeyView.png)
![key detail](etc/KeyDetail.png)
![add key](etc/AddKey.png)

### Application View

![app view](etc/AppView.png)
![app detail](etc/AppDetail.png)

### API Example

![cURL](etc/cURLExample.png)

## Requirements

Aside from the python module requirements listed in [requirements.txt](requirements.txt), the following is required:
* Python 3.6 or later.
* PostgreSQL (or other SQLAlchemy supported backend)


## Installation

This software should be used from a [viritualenv](https://virtualenv.pypa.io/en/stable/)
environment.

```sh
virtualenv venv
source venv/bin/activate
pip3 install -U -r requirements.txt
```

Then edit the config:

```sh
mv keyserv/config.example.py keyserv/config.py
```

Make sure you set `SECRET_KEY` to a randomly generated value, then change `SQLALCHEMY_DATABASE_URI`
to the URI for the database you create below.

## Database Setup

The following commands will create a suitable database for the keyserver to use.

```sh
su - postgres
createuser keyserver
createdb -O keyserver keyserver
```

## User Setup

This creates a user and password on the command line. Currently there's no user creation available
in the user interface.

```sh
export FLASK_APP=keyserv.py
flask create-user username password
```

## Key Creation & Usage

1. Create an Application at the `/add/app` URL.
2. Create a Key at the `/add/key` URL. Activations set to `-1` means unlimited activations

### API Endpoints

#### `/api/check` GET

Used to check if a key is valid. Your application should exit if the response code is not `201`.
A response of `404` means the key does not exist. This endpoint only accepts the GET method.

404 response:
```json
{"result": "failure", "error": "invalid key"}
```

201 OK response:
```json
{"result": "ok"}
```

Arguments:
- `token` - The token of the key to check for
- `machine` - The NetBIOS or domain name of the machine
- `user` - The name of the currently logged in user

#### `/api/activate` POST

Used to activate the application. If successful, the number of remaining activations will decrement
by one. After activation, your application should store the token in an obscure location and use the
`/api/check` endpoint each time it starts up. This endpoint only supports the POST method.

404 Invalid Key response:
```json
{"result": "failure", "error": "invalid activation token"}
```

410 Out of Activations response:
```json
{"result": "failure", "error": "key is out of activations"}
```

201 Activation Successful response:
```json
{"result": "ok", "remainingActivations": 1}
```
The number of remaining activations will be returned in the JSON payload. `-1` indicates unlimited
activations.

Arguments:
- `token` - The token of the key to check for
- `machine` - The NetBIOS or domain name of the machine
- `user` - The name of the currently logged in user

Example:

```sh
curl localhost:5001/api/activate -X POST -d token=2SZRHXZBNB3GUCHM375FTB8DJ -d machine=ICEBREAKER -d user=sam
{
    "result": "ok",
    "remainingActivations": "9"
}
```
