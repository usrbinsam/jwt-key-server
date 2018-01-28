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

from datetime import datetime
from enum import IntEnum
from typing import Any  # NOQA: F401

from flask_sqlalchemy import SQLAlchemy

db = SQLAlchemy()  # type: Any


class Application(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String, nullable=False, unique=True)
    support_message = db.Column(db.String)


class Key(db.Model):
    """
    Database representation of a software key provided by MKS.

    Id: identifier for a kkey
    token: the license token fed to the program
    remaining: remaining activations for a key. -1 if unlimited
    enabled: if the license is able to
    """
    id = db.Column(db.Integer, primary_key=True)
    token = db.Column(db.String, unique=True)
    remaining = db.Column(db.Integer)
    enabled = db.Column(db.Boolean, default=True)
    memo = db.Column(db.String)
    cutdate = db.Column(db.DateTime)
    app_id = db.Column(db.Integer, db.ForeignKey("application.id"), nullable=False)
    app = db.relationship("Application", uselist=False, backref="keys")

    def __init__(self, token: str, remaining: int, app_id: int,
                 enabled: bool=True, memo: str="") -> None:
        self.token = token
        self.remaining = remaining
        self.enabled = enabled
        self.memo = memo
        self.app_id = app_id

    def __str__(self):
        return f"<Key({self.token})>"


class Event(IntEnum):
    Info = 0
    Warn = 1
    Error = 2
    AppActivation = 3
    FailedActivation = 4
    KeyModified = 5
    KeyCreated = 6
    KeyAccess = 7
    AppCreated = 8
    AppModified = 9


class AuditLog(db.Model):
    """
    Database representation of an audit log.
    """
    id = db.Column(db.Integer, primary_key=True)
    key_id = db.Column(db.Integer, db.ForeignKey("key.id"), nullable=False)
    key = db.relationship("Key", uselist=False, backref="logs")
    app_id = db.Column(db.Integer, db.ForeignKey("application.id"), nullable=False)
    app = db.relationship("Application", backref="logs")
    message = db.Column(db.String)
    event_type = db.Column(db.Integer)
    timestamp = db.Column(db.DateTime)

    def __init__(self, key_id: int, app_id: int, message: str, event_type: Event) -> None:
        self.key_id = key_id
        self.app_id = app_id
        self.message = message
        self.event_type = int(event_type)
        self.timestamp = datetime.now()

    @classmethod
    def from_key(cls, key: Key, message: str, event_type: Event):
        audit = cls(key.id, key.app.id, message, event_type)
        db.session.add(audit)
        db.session.commit()
