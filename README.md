# SKFS Cleaner

**SKFS Cleaner** is a Go-based tool designed to synchronize device records between a local PostgreSQL database (using **OpenLNS/ChirpStack**) and Helium's **Session Key Filters (SKFS)**. It identifies discrepancies, such as missing or outdated devices, and ensures both systems remain in sync.

## Project Goal

The goal of **SKFS Cleaner** is to automate the synchronization of SKFS routes with the local device database by:
- **Adding missing devices** to SKFS.
- **Removing outdated devices** from SKFS.

## Features

- **Automatic SKFS Syncing**: Synchronizes SKFS with your local OpenLNS device records.
- **Device Management**: Detects and handles missing or outdated devices.
- **Dry Run Mode**: Safely test updates without committing them.
- **Batch Processing**: Efficiently processes large device sets.
- **JSON Logging**: Outputs device updates in JSON format for audit.

## Prerequisites

- **Helium Config CLI** must be installed and configured to communicate with the **Helium gRPC** API. You can install it from the official Helium documentation.

## Usage

1. Configure environment variables in a `.env` file:
    ```ini
    DRY_RUN=true
    DB_HOST=localhost
    DB_PORT=5432
    DB_USERNAME=chirpstack
    DB_PASSWORD=chirpstack
    ROUTE_ID=your-route-id
    MAX_COPIES=5
    ```

2. Build and run the project:
    ```bash
    go build .
    ./skfs-cleaner
    ```

3. **Dry Run**:
    ```bash
    DRY_RUN=true ./skfs-cleaner
    ```

4. **Live Mode**:
    ```bash
    DRY_RUN=false ./skfs-cleaner
    ```