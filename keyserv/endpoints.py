"""RESTful API views."""
# MIT License

# Copyright(c) 2017 Samuel Hoffman

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
from flask import request
from flask_restful import Resource, Api
from flask_restful import reqparse
from keyserv.keymanager import Origin, key_exists_const, key_get_unsafe, activate_key_unsafe
from keyserv.models import Application

api = Api()


class ActivateKey(Resource):
    """Endpoint used for key activation."""

    def post(self):
        """
        Activate a key

        Activates a live key; will either allow key activation or deny if there are no more key
        activations left. Function will log attempts to activate regardless of success or failure.
        """
        parser = reqparse.RequestParser()
        parser.add_argument("token", required=True)
        parser.add_argument("machine", required=True)
        parser.add_argument("user", required=True)
        parser.add_argument("app_id")

        args = parser.parse_args()

        origin = Origin(request.remote_addr, args.machine, args.user)

        if not key_exists_const(args.token, origin):

            resp = {"result": "failure", "error": "invalid activation token",
                    "support_message": None}
            if args.app_id:
                app = Application.query.get(args.app_id)
                if app and app.support_message:
                    resp["support_message"] = app.support_message

            return resp, 404

        key = key_get_unsafe(args.token, origin)

        if key.remaining == 0:
            resp = {"result": "failure", "error": "key is out of activations",
                    "support_message": key.app.support_message}

            return resp, 410

        activate_key_unsafe(args.token, origin)

        return {"result": "ok", "remainingActivations": str(key.remaining)}, 201


class CheckKey(Resource):
    """Endpoint used for checking if a key is valid."""

    def get(self):
        parser = reqparse.RequestParser()
        parser.add_argument("token", required=True)
        parser.add_argument("machine", required=True)
        parser.add_argument("user", required=True)

        args = parser.parse_args()

        origin = Origin(request.remote_addr, args.machine, args.user)

        if key_exists_const(args.token, origin):
            return {"result": "ok"}, 201

        return {"result": "failure", "error": "invalid key"}, 404


api.add_resource(ActivateKey, "/api/activate")
api.add_resource(CheckKey, "/api/check")
