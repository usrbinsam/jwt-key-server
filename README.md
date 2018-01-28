# mini-key-server

This web application provides a restful API for your desktop and other applications licensing needs.

### Key View

![key view](etc/KeyView.png)
![key detail](etc/KeyDetail.png)

### Application View

![app view](etc/AppView.png)
![app detail](etc/AppDetail.png)

### API

![cURL](etc/cURLExample.png)

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

