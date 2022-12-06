FROM python:3.10-bullseye

COPY kubeless/requirements.txt /kubeless/requirements.txt
RUN pip install -r /kubeless/requirements.txt
RUN pip install protobuf==3.20.* --force-reinstall

COPY kubeless/ /

WORKDIR /

USER 1000

CMD ["python", "/kubeless.py"]
