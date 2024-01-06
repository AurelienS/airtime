# Paraglider Utility

## Project Overview
This project is a Go-based utility designed for processing and analyzing paraglider flight data. It offers capabilities such as flight data parsing, elevation chart generation, and web-based visualization.

## Features
- IGC file parsing and processing.
- Elevation chart generation for visual flight analysis.
- Web server for data presentation and interaction.
- PostgreSQL integration for data storage and retrieval.

## Directory Structure

- `/cmd/server`: Contains the entry point for the web server.
- `/internal`: Houses the application's core logic, including flight data processing and web server handling.
- `/pkg`: Shared libraries used across the application, including model definitions and utility functions.
- `/web`: Front-end files and templates for the web interface.

## Getting Started

### Prerequisites
- Ensure you have Go installed on your system.
- PostgreSQL is required for the database components.

