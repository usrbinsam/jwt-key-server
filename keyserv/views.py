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

import psycopg2
from flask import (Blueprint, abort, current_app, flash, redirect,
                   render_template, request, url_for)
from flask_login import current_user, login_required, login_user, logout_user

from keyserv.auth import Users
from keyserv.forms import AppForm, KeyForm, LoginForm
from keyserv.keymanager import cut_key_unsafe
from keyserv.models import Application, AuditLog, Event, Key, db

frontend = Blueprint("frontend", __name__)


@frontend.route("/", methods=["GET", "POST"])
def index():

    form = LoginForm(request.form)

    if request.method == "POST" and form.validate():
        current_app.logger.debug("login form was submitted")
        user = Users.query.filter_by(username=form.username.data).first()
        if user and user.check_password(form.password.data):
            if login_user(user):
                current_app.logger.debug(f"login for {user}")
        else:
            flash("Invalid username or password.", "error")
        return redirect(url_for("frontend.index"))

    return render_template("index.html", form=form, current_user=current_user)


@frontend.route("/logout")
def logout():
    logout_user()
    return redirect(url_for("frontend.index"))


@frontend.route("/keys")
@login_required
def keys():
    return render_template("keys.html", keys=Key.query.all())


@frontend.route("/applications")
@login_required
def apps():
    return render_template("applications.html", apps=Application.query.all())


@frontend.route("/logs")
@login_required
def logs():
    return render_template("logs.html", logs=AuditLog.query.all())


@frontend.route("/modify/key/<int:key_id>", methods=["GET", "POST"])
@login_required
def modify_key(key_id: int):

    key = Key.query.get(key_id)
    if not key:
        abort(404)

    form = KeyForm(request.form)
    form.application.choices = [(app.id, app.name) for app in Application.query.all()]

    if request.method == "POST" and form.validate_on_submit():
        changes = []

        if key.remaining != form.activations.data:
            changes.append(f"activations changed from {key.remaining} to {form.activations.data}")
            key.remaining = form.activations.data
        if key.memo != form.memo.data:
            changes.append(f"memo changed from {key.memo!r} to {form.memo.data!r}")
            key.memo = form.memo.data
        if key.app_id != form.application.data:
            changes.append(f"app changed from {key.app} to {form.application.data}")
            key.application = form.application.data
        if key.enabled != form.active.data:
            changes.append(f"active changed from {key.enabled} to {form.active.data}")
            key.enabled = form.active.data

        AuditLog.from_key(key, f"edited by {current_user.username} ({request.remote_addr}):"
                          f" {', '.join(changes)}", Event.KeyModified)

        try:
            db.session.commit()
            flash("Changes successful!")
            return redirect(url_for("frontend.detail_key", key_id=key.id))
        except psycopg2.Error as error:
            flash(f"Failed to update key: {error}")

    form.application.data = key.app_id
    form.active.data = key.enabled
    form.memo.data = key.memo
    form.activations.data = key.remaining

    return render_template("add_modify.html", header=f"Modify Key {key.id}", form=form)


@frontend.route("/add/key", methods=["GET", "POST"])
@login_required
def add_key():
    form = KeyForm(request.form)
    form.application.choices = [(app.id, app.name) for app in Application.query.all()]

    if request.method == "POST" and form.validate_on_submit():
        try:
            token = cut_key_unsafe(form.activations.data, form.application.data,
                                   form.active.data, form.memo.data)
            flash(f"Key added! Token: {token}", "success")
        except psycopg2.Error as error:
            flash(f"Unable to add key: {error}", "error")

    return render_template("add_modify.html", header="Add Key", form=form)


@frontend.route("/add/app", methods=["GET", "POST"])
@login_required
def add_app():
    form = AppForm(request.form)

    if request.method == "POST" and form.validate_on_submit():
        app = Application()
        app.name = form.name.data

        db.session.add(app)
        try:
            db.session.commit()
            flash("Success!")
        except psycopg2.Error as error:
            flash(f"Failed to add application: {error}")

    return render_template("add_modify.html", form=form, header="Add Application")


@frontend.route("/modify/app/<int:app_id>", methods=["GET", "POST"])
@login_required
def modify_app(app_id: int):
    app = Application.query.get(app_id)

    if not app:
        abort(404)

    form = AppForm(request.form)
    if request.method == "POST" and form.validate_on_submit():

        app.name = form.name.data
        try:
            db.session.commit()
            flash("Success.")
        except Exception as error:
            flash(f"Failed to modify application: {error}", "error")

    form.name.data = app.name

    return render_template("add_modify.html", form=form)


@frontend.route("/detail/key/<int:key_id>")
@login_required
def detail_key(key_id: int):

    key = Key.query.get(key_id)

    if not key:
        abort(404)

    return render_template("detail_key.html", key=key)


@frontend.route("/detail/app/<int:app_id>")
@login_required
def detail_app(app_id: int):

    app = Application.query.get(app_id)

    if not app:
        abort(404)

    return render_template("detail_app.html", app=app)


@frontend.route("/keys/app/<int:app_id>")
@login_required
def keys_for_app(app_id):

    app = Application.query.get(app_id)

    if not app:
        abort(404)

    return render_template("keys.html", keys=app.keys)
