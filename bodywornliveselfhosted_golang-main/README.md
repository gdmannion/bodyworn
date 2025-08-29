# bodywornlive

## Project Structure

```

bodywornlive/
├── go.mod                          ← Go module dependencies
├── go.sum                          ← Go module checksums
├── config.json                     ← Device IP, credentials
├── main.go                         ← Entry point, contains handlers & startup logic
├── static/                         ← Embedded static frontend files
│   ├── index.html                  ← Main UI page
│   ├── main.js                     ← Entry point script (initializes event feed)
│   ├── events.js                   ← Handles WebSocket event feed logic
│   ├── webrtc_video.js             ← Handles WebRTC and signaling logic
├── subscribe\_events/               ← Event subscription logic
│   └── subscribe\_events.go         ← Handles Axis event subscription
├── digest\_auth/                    ← Digest authentication logic
│   └── digest\_auth.go              ← Digest auth client for Axis


````

## Overview

This project is designed to interface with Axis devices and manage live video streaming using WebRTC. The Go backend handles digest authentication and event subscription. The frontend contains a simple UI for interacting with the system.

## Setup

1. Clone the repository:



2. Initialize Go modules:



3. Update the `config.json` file with the necessary device IP and credentials.

## Running the Application

To start the application, run the following command:

```bash
go run main.go
```

This will launch the server, and you can access the UI via `index.html`.

## Directory Breakdown

* `main.go`: Contains the entry point for the application. It sets up the web server and the WebRTC communication logic.
* `config.json`: Holds the Axis device configuration for S3008
* `static/`: Contains static frontend files.

  * `index.html`: The main UI page that interacts with the backend.
  * `websocket_events.js`: JavaScript file that handles the WebRTC logic and communication with the backend 

* `subscribe_events/`: Handles subscription to Axis device events.
  * `subscribe_events.go`: Manages event subscription and handling for Axis devices.

* `digest_auth/`: Manages digest authentication with the Axis device.

  * `digest_auth.go`: Contains functions to perform digest authentication when interacting with Axis devices.

