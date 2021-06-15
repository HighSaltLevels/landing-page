from flask import Flask, request
import werkzeug

APP = Flask(__name__)
SADIE1 = "assets/sadie1.jpg"
SADIE2 = "assets/sadie2.jpg"
CSS = "html/styles.css"

WHITELIST = ["/" + endpoint + "?" for endpoint in {SADIE1, SADIE2, CSS, "register" ""}]


@APP.route("/")
def index():
    with open("html/index.html") as stream:
        return stream.read()


@APP.route("/" + SADIE1)
def get_sadie1():
    with open(SADIE1, "rb") as stream:
        return stream.read()


@APP.route("/" + SADIE2)
def get_sadie2():
    with open(SADIE2, "rb") as stream:
        return stream.read()


@APP.route("/register")
def get_registration_page():
    with open("html/register.html") as stream:
        return stream.read()


@APP.before_request
def prevent_unexpected_endpoint_hits():
    """I keep getting hits from scam servers with embedded urls in the url
    Let's prevent that by whitelisting expected endpoints
    """
    if request.full_path not in WHITELIST:
        raise werkzeug.exceptions.NotFound()


APP.run(host="0.0.0.0", port=80)
