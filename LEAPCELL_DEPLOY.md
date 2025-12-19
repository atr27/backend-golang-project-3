# Deploying Go Backend to Leapcell (Native Runtime)

This guide details how to deploy your Go application to Leapcell using the native Go runtime (without Docker).

## Prerequisites

-   **Leapcell Account**: Sign up at [leapcell.io](https://leapcell.io).
-   **GitHub Repository**: Your code must be pushed to a GitHub repository.

## Step 1: Prepare Your Project

Ensure your project structure is ready.

1.  **Check `go.mod`**: Ensure your `go.mod` file is in the `backend` directory (or root if you adjusted it).
2.  **Verify Build**: Locally, you should be able to run `go build -o main ./cmd/api` from the `backend` folder.

## Step 2: Create a Service on Leapcell

1.  **Log in to Leapcell Dashboard**.
2.  Click **"Create Service"**.
3.  **Connect Repository**:
    -   Select your GitHub repository.
    -   Choose the branch you want to deploy (usually `main`).

## Step 3: Configure Deployment

Since we are not using Docker, we will configure the Go runtime directly.

**Configuration Settings:**

*   **Runtime**: `Go` (Select the version matching your `go.mod`, e.g., 1.21).
*   **Root Directory**: `backend` (Important: This tells Leapcell where your code lives).
*   **Build Command**: `go build -o main ./cmd/api`
    *   This command compiles your application located in `cmd/api` and outputs a binary named `main`.
*   **Start Command**: `./main`
    *   This executes the binary created in the build step.
*   **Port**: `8080`

**Environment Variables:**

Add your production environment variables in the Leapcell Dashboard.

*   `PORT`: `8080`
*   `GIN_MODE`: `release`
*   `POSTGRES_URL`: (Your production database connection string)
*   `REDIS_URL`: (Your production Redis URL if used)
*   Any other variables from your `.env`.

## Step 4: Deploy

1.  Click **"Deploy"**.
2.  Leapcell will install dependencies (`go mod download`), run your build command, and start the application.

## Step 5: Verify Deployment

1.  Visit the provided public URL.
2.  Test the health check endpoint to confirm the service is up.
