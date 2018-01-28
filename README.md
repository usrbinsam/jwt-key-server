# mini-key-server

## Requirements

Aside from the python module requirements listed in [requirements.txt](requirements.txt), the following is required:
* Python 3.6 or later.
* PostgreSQL


## Installation

This software should be used from a [viritualenv](https://virtualenv.pypa.io/en/stable/) environment.

```sh
virtualenv venv
source venv/bin/activate
pip3 install -U -r requirements.txt
```

## Database Setup

The following commands will create a suitable database for the keyserver to use.

```sh
su - postgres
createuser keyserver
createdb -O keyserver keyserver
```
