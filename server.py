import oprfs
import flask

app = flask.Flask(__name__)

# Normally, a persistent key should be retrieved from secure storage.
# Here, a new key is created each time so older masks cannot be reused
# once the service is restarted.
key = oprfs.key()


@app.route("/", methods=["POST"])
def endpoint():
    # Call the handler with the key and request, then return the response.
    return flask.jsonify(oprfs.handler(key, flask.request.get_json()))


app.run()
