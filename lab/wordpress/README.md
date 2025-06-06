# WordPress Lab Environment

This directory contains the necessary files to build and run a Dockerized WordPress lab environment. It includes configurations for Apache, MySQL, and WordPress itself.

## Files

*   `Dockerfile`: Defines the Docker image for the WordPress environment. It likely sets up a base OS, installs Apache, PHP, MySQL client, and WordPress.
*   `apache.conf`, `vhost.conf`: Apache web server configuration files.
*   `configure-mysql.sh`, `run-mysql.ini`: Scripts and configurations related to setting up and running the MySQL database.
*   `index.php`: Likely the main WordPress entry point or a custom landing page.
*   `init-local.sh`, `run-local.sh`, `run.sh`: Shell scripts to initialize, run, or manage the lab environment.
*   `notes`: May contain specific notes or credentials for this lab setup (check its content).
*   `run-apache.ini`: Configuration for running Apache, possibly via a process manager like supervisord.
*   `wordpress.sql`: An SQL dump, likely used to initialize the WordPress database with specific content or a pre-configured state.
*   `wp-config.php`: The WordPress configuration file, containing database connection details and other settings.

## Setup and Usage

While the exact commands might vary based on the content of the shell scripts, a typical workflow to set up and run this lab would be:

1.  **Navigate to this directory:**
    ```bash
    cd path/to/your/new-hacks/lab/wordpress
    ```

2.  **Build the Docker Image:**
    Inspect the `Dockerfile` and any build-related scripts (e.g., `init-local.sh` if it involves pre-build steps). A common command to build a Docker image from a Dockerfile in the current directory is:
    ```bash
    docker build -t wordpress-lab .
    ```
    *(Replace `wordpress-lab` with your preferred image name/tag.)*

3.  **Run the Docker Container:**
    Inspect the run scripts (`run.sh`, `run-local.sh`). These scripts likely handle starting the necessary services (MySQL, Apache) and may set up port mappings. A common command to run a container might look like:
    ```bash
    docker run -d -p 8080:80 --name my-wordpress-lab wordpress-lab
    ```
    *(This maps port 8080 on your host to port 80 in the container. Adjust ports as needed based on the Dockerfile or run scripts.)*

    Alternatively, if a `docker-compose.yml` file were present (it's not listed but is common for multi-service labs), you would use `docker-compose up`.

4.  **Accessing the Lab:**
    Once the container is running, you should be able to access the WordPress site via your browser, typically at `http://localhost:8080` (if you used the port mapping above). Check the `notes` file or script outputs for any specific URLs or default credentials.

## Important Notes

*   **Review Scripts:** Before running any shell scripts (`.sh`), it's crucial to review their content to understand what actions they will perform.
*   **Database Initialization:** The `wordpress.sql` file suggests the database will be pre-populated. The `configure-mysql.sh` script likely handles this.
*   **Configuration:** `wp-config.php` will contain database credentials. These might be hardcoded, set via environment variables in the Dockerfile/scripts, or generated during initialization.
*   **Local Data:** The main `.gitignore` in the parent `lab/` directory specifies `*/local/*` to be ignored. This implies that some scripts might create or expect local data storage (e.g., for MySQL data persistence) in a `local/` subdirectory within `lab/wordpress/`.

Refer to the specific shell scripts and configuration files for detailed setup steps and any custom configurations for this lab.
