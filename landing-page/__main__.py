from flask import Flask

APP = Flask(__name__)
SADIE1 = "assets/sadie1.jpg"
SADIE2 = "assets/sadie2.jpg"
CSS = "html/styles.css"

@APP.route("/")
def index():
    with open("html/index.html") as stream:
        return stream.read()

@APP.route(f"/{SADIE1}")
def get_sadie1():
    with open(SADIE1, "rb") as stream:
        return stream.read()

@APP.route(f"/{SADIE2}")
def get_sadie2():
    with open(SADIE2, "rb") as stream:
        return stream.read()

@APP.route(f"/{CSS}")
def get_css():
    with open(CSS) as stream:
        return stream.read()

APP.run(host="0.0.0.0", port=80)
