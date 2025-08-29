Overview


Script File Responsibilities

main.go – Starts the server, initializes logging, and maps endpoints.

server.go – Defines authentication and routes for storage operations.

metadata.go – Parses and stores metadata headers from requests into .meta sidecar files.

containers_objects.go – Handles object I/O, container creation, and mimics Swift behaviors like ETag and HEAD support.

files.go – Generates connection.json and Capabilities.json for BWM integration using system introspection.

logger.go – Implements a pluggable and thread-safe logging interface for use across modules.

This server provides a fully working mock implementation of the Axis Body Worn Integration API, emulating behavior of the OpenStack Swift object storage model over a local filesystem. It is tailored for use as a Content Destination (CD) for testing and integration with Axis Body Worn Systems (BWS).

The server enables third-party applications to:

Authenticate using token-based headers

Stores all metadata and video clips

Serve capability and connection definitions

Integrate with real BWS clients

Key Features

✅ Simulates Swift-compatible 3-level hierarchy (Account > Container > Object)

✅ Local filesystem backend for easy debugging and prototyping

✅ Auto-generates required connection and capability JSON files

✅ Supports PUT/GET/POST/HEAD for objects and metadata

✅ Logs metadata headers and access patterns for traceability

✅ Fully mimics real recording flow: User, Device, System, Recording containers

Developer Setup

Run the Go server:

go run main.go


Auto-Generated Files - connection.json



Directory Structure

StorageAccount/
├── System/
│   ├── Capabilities.json
│   └── <SystemID>
├── Users/
│   └── <UserUUID>
├── Devices/
│   └── <DeviceID>
├── <RecordingContainer>/
│   ├── clip.mkv




Use Cases - Provide a simple code to create a  content destination 


System Architecture Diagram

              +-------------------+
              | Axis Body Worn SCU|
              +--------+----------+
                       |
                [HTTP + Token Auth]
                       |
         +-------------v--------------+
         |     Go Mock Server (CD)    |
         |----------------------------|
         |  /auth/v1.0                |
         |  /v1.0/StorageAccount/...  |
         |                            |
         |  StorageAccount/           |
         |    ├── System/             |
         |    ├── Users/              |
         |    ├── Devices/            |
         |    └── recordings  |
         +----------------------------+
                       |
          [Local Filesystem Objects]



