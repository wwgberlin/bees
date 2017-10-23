import cv2

def process(img):
    face_cascade = cv2.CascadeClassifier('/app/face_detector/haarcascade_profileface.xml')
    if face_cascade.empty():
        raise Exception("failed locating file")

    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)

    faces = face_cascade.detectMultiScale(gray, 1.3, 5)

    for (x,y,w,h) in faces:
        cropped = img[y:y+h, x:x+w]
        # cv2.rectangle(img, (x, y), (x+w, y+h), (0, 255, 0), 2)
        return cropped
    return img


