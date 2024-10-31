# Go Webhook Transformation

A Go application that transforms incoming webhook events and sends them to a specified webhook URL. This project is designed to process JSON data and output structured information based on the dynamic fields provided.

## Live Demo

You can access the live version of this application at:
[go-webhook-transform-production.up.railway.app](https://go-webhook-transform-production.up.railway.app)

## Features

- Accepts incoming POST requests with JSON payloads.
- Processes dynamic fields to extract and organize data.
- Sends the transformed data to a specified webhook URL.
- Outputs structured JSON responses.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Endpoints](#endpoints)
6. [Accessing the Service via cURL](#accessing-the-service-via-curl)

## Prerequisites

- Go (version 1.16 or higher)
- Basic understanding of JSON and webhooks
- cURL installed on your machine (for testing)

## Installation

1. Clone the repository to your local machine:
    ```bash
    git clone https://github.com/ahdaan98/go-webhook-transformation.git
    cd go-webhook-transformation
    ```

2. Install the required Go packages (if any):
    ```bash
    go mod tidy
    ```

## Usage

1. Build the application:
    ```bash
    go build -o webhook-transform
    ```

2. Run the application:
    ```bash
    ./webhook-transform
    ```

3. The server will start listening on `http://localhost:8080`.

## Endpoints

### POST /

- **Request Body:** JSON formatted data. Example:
    ```json
    {
        "ev": "contact_form_submitted",
        "et": "form_submit",
        "id": "cl_app_id_001",
        "uid": "cl_app_id_001-uid-001",
        "mid": "cl_app_id_001-uid-001",
        "t": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
        "p": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
        "l": "en-US",
        "sc": "1920 x 1080",
        "atrk1": "form_varient",
        "atrv1": "red_top",
        "atrt1": "string",
        "atrk2": "ref",
        "atrv2": "XPOWJRICW993LKJD",
        "atrt2": "string",
        "uatrk1": "name",
        "uatrv1": "iron man",
        "uatrt1": "string",
        "uatrk2": "email",
        "uatrv2": "ironman@avengers.com",
        "uatrt2": "string",
        "uatrk3": "age",
        "uatrv3": "32",
        "uatrt3": "integer"
    }
    ```

- **Response:** JSON response with the following structure:
    ```json
    {
        "event": "contact_form_submitted",
        "event_type": "form_submit",
        "app_id": "cl_app_id_001",
        "user_id": "cl_app_id_001-uid-001",
        "message_id": "cl_app_id_001-uid-001",
        "page_title": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
        "page_url": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
        "browser_language": "en-US",
        "screen_size": "1920 x 1080",
        "attributes": {
            "form_varient": {
                "value": "red_top",
                "type": "string"
            },
            "ref": {
                "value": "XPOWJRICW993LKJD",
                "type": "string"
            }
        },
        "traits": {
            "name": {
                "value": "iron man",
                "type": "string"
            },
            "email": {
                "value": "ironman@avengers.com",
                "type": "string"
            },
            "age": {
                "value": "32",
                "type": "integer"
            }
        }
    }
    ```


## Accessing the Service via cURL
### Sample Input

You can use the following `curl` command to test the API endpoint:

```bash
curl -X POST http://localhost:8080 \
-H "Content-Type: application/json" \
-d '{
    "ev": "contact_form_submitted",
    "et": "form_submit",
    "id": "cl_app_id_001",
    "uid": "cl_app_id_001-uid-001",
    "mid": "cl_app_id_001-uid-001",
    "t": "Vegefoods - Free Bootstrap 4 Template by Colorlib",
    "p": "http://shielded-eyrie-45679.herokuapp.com/contact-us",
    "l": "en-US",
    "sc": "1920 x 1080",
    "atrk1": "form_varient",
    "atrv1": "red_top",
    "atrt1": "string",
    "atrk2": "ref",
    "atrv2": "XPOWJRICW993LKJD",
    "atrt2": "string",
    "uatrk1": "name",
    "uatrv1": "iron man",
    "uatrt1": "string",
    "uatrk2": "email",
    "uatrv2": "ironman@avengers.com",
    "uatrt2": "string",
    "uatrk3": "age",
    "uatrv3": "32",
    "uatrt3": "integer"
}'
