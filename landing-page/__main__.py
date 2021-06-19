from flask import Flask, request
import werkzeug

APP = Flask(__name__)
SADIE1 = "assets/sadie1.jpg"
SADIE2 = "assets/sadie2.jpg"
CSS = "html/styles.css"


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


@APP.route("/favicon.ico")
def get_favicon():
    with open("assets/greeson.png", "rb") as stream:
        return stream.read()


APP.run(host="0.0.0.0", port=80)
