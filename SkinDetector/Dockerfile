FROM python:latest
RUN apt-get update
RUN apt-get install -y python python-dev python-pip python-virtualenv
RUN apt-get install -y python-opencv
RUN apt-get install -y python-matplotlib
ADD requirements.txt /app/requirements.txt
RUN pip install -r /app/requirements.txt
WORKDIR /app
ADD . /app
CMD export FLASK_APP=test.py && flask run
EXPOSE 8080