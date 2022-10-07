from encodings import utf_8
from fastapi import FastAPI, Response
from fastapi.responses import JSONResponse
from fastapi.encoders import jsonable_encoder
import boto3
import trp
import json
from pydantic import BaseModel
import base64

app = FastAPI()

class UserRequestIn(BaseModel):
    """
    Class for holding the data from user request
    """
    data: bytes


@app.post("/classification")
def classification(user_request_in: UserRequestIn):
    """
    Process the user request with image data.
    Return the JSON with Entities
    """


    client = boto3.client(service_name='textract', region_name='us-east-2')

    response = client.detect_document_text(
        Document={
            'Bytes': base64.b64decode(user_request_in.data),
        }
    )
    doc = trp.Document(response)

    maxLength = 20000

    comprehendResponse = []

    comprehend_medical_client = boto3.client(
        service_name='comprehendmedical', region_name='us-east-2')

    for page in doc.pages:
        pageText = page.text

        for i in range(0, len(pageText), maxLength):
            response = comprehend_medical_client.detect_entities_v2(
                Text=pageText[0+i:maxLength+i])
            comprehendResponse.append(response)
            
    return comprehendResponse[0]
