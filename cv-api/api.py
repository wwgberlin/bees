from SkinDetector import skin_detector
from flask import Flask, request, redirect, url_for, Response

import cv2
import face_detector
import json
import urllib.request
import numpy as np

app = Flask(__name__)

def url_to_image(image_url):
    resp = urllib.request.urlopen(image_url)
    image = np.asarray(bytearray(resp.read()), dtype="uint8")
    return cv2.imdecode(image, cv2.IMREAD_COLOR)

@app.route('/face/')
def detect_face():
    img_path = request.args.get('image')
    image = url_to_image(img_path)
    cropped = face_detector.process(image)
    return Response(cv2.imencode('.jpg', cropped)[1].tostring(), mimetype="image/jpeg")

@app.route('/api/face/')
def detect_face_api():
    img_path = request.args.get('image')
    image = url_to_image(img_path)
    cropped = face_detector.process(image)
    return Response(json.dumps(cropped.tolist()), mimetype="application/json")

@app.route('/skin/')
def detect_skin():
    img_path = request.args.get('image')
    image = url_to_image(img_path)
    mask = skin_detector.process(image)
    mask_rgb = cv2.bitwise_and(image, image, mask=mask)
    return Response(cv2.imencode('.jpg', mask_rgb)[1].tostring(), mimetype='image/jpeg')

@app.route('/api/skin/')
def detect_skin_api():
    img_path = request.args.get('image')
    image = url_to_image(img_path)
    mask = skin_detector.process(image)
    mask_rgb = cv2.bitwise_and(image, image, mask=mask)
    return Response(json.dumps(mask_rgb.tolist()), mimetype='application/json')

app.run(host='0.0.0.0', port=8080)