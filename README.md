# Gator CLI

Gator CLI is a command-line application that interacts with a PostgreSQL database. It allows users to manage various features related to the application, including managing users and adding feeds.

## Requirements

Before you start using the Gator CLI, make sure you have the following installed:

- **Docker**: Used to spin up the PostgreSQL database for development.
- **Go (Golang)**: Required to run and build the Go application.
- **Goose**: A Go migration tool for database schema management.

To install **Goose**, run the following command:

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
    
## Setup

1. **Clone the repository**:

   ```
   git clone https://github.com/Pizzu/gator.git
   cd gator
   ```

2. **Set up environment variables**:

   Create a `.env` file in the root of the project with the following content:

   ```
   DB_PASSWORD=rootPassword
   DB_NAME=gator
   ```

   Make sure to replace `rootPassword` with your own secure password.

3. **Start the PostgreSQL container**:

   To make it easy for all developers to spin up a PostgreSQL instance with Docker, we’ve provided a `docker-compose.yml` file. Run the following command to start the database:

   ```
   docker compose -f docker-compose.yml up -d
   ```

   This will pull the necessary Docker image and start the PostgreSQL container. The database `gator` will be created automatically with the environment variables provided.

4. **Run migrations**:

   Once the PostgreSQL container is up, you can run the database migrations with Goose:

   From the root project type the following commands:

   ```
   cd sql/schema
   ```

   ```
   goose postgres "postgres://postgres:<DB_PASSWORD>@127.0.0.1:5432/<DB_NAME>?sslmode=disable" up
   ```

   This will apply the migrations defined in the `./sql/schema` folder.

## Usage

The Gator CLI is used with the following commands:

1. **Register a user**:

   ```
   go run . register alan
   ```

   This will register a user.

2. **To get the current user**:

   ```
   go run . login alan
   ```

   This will login and switch the current user

3. **To get the current user**:

   ```
   go run . users
   ```

   This will display all registered users and the current user.

4. **To add a new feed**:

   ```
   go run . addFeed "TechCrunch" "https://techcrunch.com/feed/"
   ```

   This will add a new feed with the provided name and URL and the current user will automatically follow that feed.

5. **To get feeds**:

   ```
   go run . feeds
   ```

   The will retrieve all the feeds saved on the db.

6. **To follow a feed**:

   ```
   go run . follow "https://techcrunch.com/feed/"
   ```

   The current user will follow the specified feed (created from another user).

7. **To unfollow a feed**:

   ```
   go run . unfollow "https://techcrunch.com/feed/"
   ```

   The current user will unfollow the specified feed.

8. **To retrieve the feeds followed by the user**:

   ```
   go run . following
   ```

   The will show all the feeds the current user follows.

9. **To aggregate feeds**:

   ```
   go run . agg 1min
   ```

   The agg command is a never-ending loop that fetches feeds and saves posts to the database. The intended use case is to leave the agg command running in the background while you interact with the program in another terminal.
   Specify how often you want to collect and update feeds with the following format: 1min, 30min, etc..

10. **To browse feed posts**:

    ```
    go run . browse 3
    ```

    This will display all the posts that belong to the feeds followed by the current user. Use the second argument to set a LIMIT.

11. **Reset**:

    ```
    go run . reset
    ```

    This will reset the entire DB for a fresh start.

## Development

To make development easier, you can spin up the PostgreSQL database using Docker and work with the Go application locally. Just make sure to have Go and Docker installed.
