import cv2
from SkinDetector import skin_detector
import face_detector
import json

from flask import Flask
from flask import Response
app = Flask(__name__)


@app.route('/face/')
def detect_face():
    img_path = "beyonce.jpg"
    image = cv2.imread(img_path)
    cropped = face_detector.process(image)
    # return Response(cv2.imencode('.jpg', cropped)[1].tostring(), mimetype="image/jpeg")
    return Response(json.dumps(cropped.tolist()), mimetype="application/json")

@app.route('/')
def hello_world():

    img_path = "beyonce.jpg"
    image = cv2.imread(img_path)
    mask = skin_detector.process(image)
    mask_rgb = cv2.bitwise_and(image, image, mask=mask)
    img_str = cv2.imencode('.jpg', mask_rgb)[1].tostring()
    # return Response(json.dumps(mask_rgb.tolist()), mimetype='application/json')
    return Response(img_str, mimetype='image/jpeg')

app.run(host='0.0.0.0', port=8080)

