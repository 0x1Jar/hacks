# Software Lab Environments

This `lab/` directory contains configurations, primarily Dockerfiles, for setting up various software lab environments. Each subdirectory is intended to house a specific lab setup.

## Purpose

The goal is to provide easily reproducible environments for testing, development, or learning purposes related to different software stacks or vulnerabilities.

## General Usage

Typically, each subdirectory within `lab/` will contain:
*   A `Dockerfile` to build the lab environment.
*   Supporting configuration files (e.g., `docker-compose.yml`, scripts, application source code).
*   A specific `README.md` within that subdirectory explaining its purpose and how to build and run that particular lab.

Navigate into a specific lab's subdirectory (e.g., `cd lab/wordpress`) and follow the instructions in its local `README.md` to set it up.

## Available Labs

*   **`wordpress/`**: Contains a setup for a WordPress environment. (See `lab/wordpress/README.md` for details).

*(More labs may be added here in the future.)*

## Contributing

To add a new lab:
1.  Create a new subdirectory under `lab/`.
2.  Add your `Dockerfile` and any necessary supporting files.
3.  Create a `README.md` inside your new subdirectory explaining what the lab is for and how to use it (build commands, run commands, exposed ports, default credentials if any, etc.).
4.  Update this main `lab/README.md` to include your new lab in the "Available Labs" section.
