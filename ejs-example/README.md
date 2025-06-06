# EJS Example Web Application

This project is a sample web application demonstrating user registration, login, and a simple API endpoint using Node.js, Express, EJS templates, and MongoDB. It also showcases password hashing with `bcrypt` and session management with `express-session`.

## Features

*   User registration with password hashing.
*   User login and session management.
*   A main page displaying a list (e.g., "people").
*   A page to view users with a basic search/filter API (`/api/users`).
*   Static file serving.

## Prerequisites

*   **Node.js**: Version 12.x or newer recommended.
*   **npm**: Comes with Node.js.
*   **MongoDB**: A running MongoDB server instance (defaults to `mongodb://localhost:27017`).

## Setup

1.  **Clone/Download the project:**
    Ensure you have all project files in a directory (e.g., `ejs-example`).

2.  **Navigate to the project directory:**
    ```bash
    cd path/to/your/ejs-example
    ```

3.  **Install dependencies:**
    This command reads `package.json` and installs the required Node.js modules (Express, EJS, MongoDB driver, bcrypt, express-session) into the `node_modules` directory.
    ```bash
    npm install
    ```

4.  **Ensure MongoDB is running:**
    The application will try to connect to a MongoDB instance at `mongodb://localhost:27017` and use a database named `rushwebapp`.

## Running the Application

Once setup is complete, start the web server:
```bash
node main.js
```
You should see a message like: `Example app listening on port 3000!`

Open your web browser and navigate to `http://localhost:3000`.

## Project Structure

*   `main.js`: The main application file containing the Express server setup, routes, and logic.
*   `package.json`: Defines project metadata and dependencies.
*   `package-lock.json`: Records exact versions of dependencies.
*   `views/`: Directory containing EJS template files (`.ejs`).
    *   `people.ejs`: Main page.
    *   `register.ejs`: User registration page.
    *   `users.ejs`: Page to display/search users.
    *   `error.ejs`: Generic error page.
*   `static/`: Directory for static assets (CSS, client-side JS, images - currently empty but configured).
*   `.gitignore`: Specifies files/directories to be ignored by Git (e.g., `node_modules`, `dbdata`).

### Helper Scripts

*   `generate-hashes.js`: A utility script to generate bcrypt password hashes. This is not part of the main application but can be run manually (e.g., `node generate-hashes.js`) if you need to pre-generate hashes for testing or manual database seeding.
*   `mongo-await.js`: An example script demonstrating asynchronous MongoDB operations using `async/await`. It connects to the database, creates an index, inserts a user, and fetches documents. This is for demonstration/testing and not part of the main web application flow.

## Endpoints

*   `GET /`: Main page, displays a list of people and login/logout/register links.
*   `GET /register`: Displays the user registration form.
*   `POST /register`: Handles user registration.
*   `POST /login`: Handles user login.
*   `GET /logout`: Logs out the current user.
*   `GET /users`: Displays a page for searching users.
*   `GET /api/users?q=<query>`: API endpoint that returns a JSON array of usernames matching the query (prefix search).
*   `/static/*`: Serves static files from the `static` directory.
