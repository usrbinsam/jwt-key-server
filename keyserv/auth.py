# MIT License

# Copyright(c) 2018 Samuel Hoffman

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files(the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and / or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

import secrets

import argon2
from flask_login import LoginManager, UserMixin

from keyserv.models import db

login_manager = LoginManager()
login_manager.session_protection = "strong"


def add_user(username: str, password: bytes, level=500):
    passwd = argon2.hash_password(password, secrets.token_bytes(None))
    user = Users(username, passwd, level)
    db.session.add(user)
    db.session.commit()


class Users(db.Model, UserMixin):
    id = db.Column(db.Integer(), primary_key=True)
    username = db.Column(db.String(), unique=True, nullable=False)
    passwd = db.Column(db.LargeBinary(), nullable=False)
    level = db.Column(db.Integer())

    def __init__(self, username=None, passwd=None, level=0):
        self.username = username
        self.passwd = passwd
        self.level = level

    def get_id(self):
        return self.id

    def check_password(self, passwd):
        try:
            return argon2.verify_password(self.passwd, bytes(passwd, "UTF-8"))
        except argon2.exceptions.VerifyMismatchError:
            return False


@login_manager.user_loader
def user_loader(user_id) -> Users:
    return Users.query.get(user_id)
