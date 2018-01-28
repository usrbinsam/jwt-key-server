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

import click
from flask import Flask
from flask_bootstrap import Bootstrap

from .auth import login_manager, add_user
from .endpoints import api
from .models import db, Event
from .views import frontend


def format_event(value):
    return Event(value)


def format_datetime(value):
    if value is None:
        return ""

    try:
        return value.strftime("%Y-%m-%d %H:%M:%S")
    except ValueError:
        return ""


def create_app(config):
    app = Flask(__name__)

    app.config.from_object(__name__)
    app.config.from_object("keyserv.config.{}".format(config))
    app.jinja_env.filters["event"] = format_event
    app.jinja_env.filters["datetime"] = format_datetime

    Bootstrap(app)
    api.init_app(app)
    db.init_app(app)
    login_manager.init_app(app)

    app.register_blueprint(frontend)

    @app.cli.command("initdb")
    def initdb_command():
        db.create_all()
        print("database initialized")

    @app.cli.command("create-user")
    @click.argument("username")
    @click.argument("password")
    def create_user_command(username: str, password: str):
        add_user(username, password.encode())

    return app
