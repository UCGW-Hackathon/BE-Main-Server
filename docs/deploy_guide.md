# Deploying Go Backend to Hugging Face Spaces

This guide outlines the steps to deploy your Go application (`whatsapp-backend`) to Hugging Face Spaces using Docker.

## Prerequisites

- A [Hugging Face](https://huggingface.co/) account.
- Your project code locally.

## Step 1: Create a New Space

1.  Log in to your Hugging Face account.
2.  Click on your profile picture in the top right and select **New Space**.
3.  **Owner**: Select your username or organization.
4.  **Space Name**: Enter a name for your space (e.g., `whatsapp-backend`).
5.  **License**: Choose a license (e.g., MIT) or leave it empty correctly.
6.  **SDK**: Select **Docker**.
7.  **Blank**: Choose "Blank" to start with an empty repository.
8.  Click **Create Space**.

## Step 2: Prepare Your Code (Already Done by Agent)

Your `Dockerfile` has been configured to:

- Use port `7860` (required by Hugging Face Spaces).
- Bind to address `0.0.0.0`.

Ensure your `Dockerfile` contains:

```dockerfile
ENV HOST_ADDRESS=0.0.0.0
ENV HOST_PORT=7860
EXPOSE 7860
```

## Step 3: Deployment Methods

You can deploy using one of two methods:

### Method A: Upload via Web Interface (Easiest)

1.  In your newly created Space, go to the **Files** tab.
2.  Click **Add file** -> **Upload files**.
3.  Drag and drop all your project files (including `.dockerignore` if you have one, `Dockerfile`, `go.mod`, `go.sum`, `main.go`, and all source directories like `config`, `controllers`, etc.) into the upload area.
    - _Note: Do not upload `.git` folder._
4.  Commit the changes with a message like "Initial commit for deployment".
5.  Hugging Face will automatically start building your Docker image. You can see the progress in the **App** tab.

### Method B: Push via Git (Recommended for updates)

1.  In your Space, click on the **Clone repository** button (or just copy the HTTPS URL).
2.  Open your terminal in your project directory.
3.  Initialize git if you haven't (or use your existing git setup):
    ```bash
    git remote add space https://huggingface.co/spaces/YOUR_USERNAME/SPACE_NAME
    ```
4.  Push your code:
    ```bash
    git push space main
    ```
    _(Note: You might need to use a Hugging Face Access Token as your password if prompted)._

### Method C: Automated Sync (GitHub Actions) - **RECOMMENDED**

This method automatically pushes your code to the Hugging Face Space whenever you push to the `test` branch of your GitHub repository.

1.  **Generate a Hugging Face Token**:

    - Go to your [Hugging Face Settings > Tokens](https://huggingface.co/settings/tokens).
    - Create a new token with **write** permissions.
    - Copy the token.

2.  **Add Secret to GitHub**:

    - Go to your GitHub repository.
    - Navigate to **Settings** > **Secrets and variables** > **Actions**.
    - Click **New repository secret**.
    - Name: `HF_TOKEN`
    - Value: Paste your Hugging Face token.
    - Click **Add secret**.

3.  **Update Workflow File**:

    - Open `.github/workflows/sync-to-hub.yml` in your project.
    - Find the `git push` command at the bottom.
    - Replace `YOUR_USERNAME` and `SPACE_NAME` with your actual Hugging Face username and Space name.
    - _Example_: `git push https://ryu:Start@huggingface.co/spaces/ryu/whatsapp-backend test:main -f`

4.  **Trigger Deployment**:
    - Simply push any change to your `test` branch:
      ```bash
      git checkout test
      git add .
      git commit -m "New feature"
      git push origin test
      ```
    - The GitHub Action will run and sync your code to the Spaces.

## Step 4: Configure Secrets (Environment Variables)

Your application uses environment variables (like Database credentials). You **must** set these in Hugging Face settings, not in the `Dockerfile` (for security).

1.  Go to your Space's **Settings** tab.
2.  Scroll down to **Variables and secrets**.
3.  Click **New secret** for sensitive data (like passwords) or **New variable** for public config.
4.  Add the following based on your `env_config.go` and `.env` needs:

    - `DB_HOST`
    - `DB_PORT`
    - `DB_USER`
    - `DB_PASSWORD`
    - `DB_NAME`
    - `SALT` (Optional, defaults to "Def4u|7")
    - `TZ` (Optional, e.g., "Asia/Jakarta")

    _Note: `HOST_ADDRESS` and `HOST_PORT` are already set in the Dockerfile, so you don't need to add them here unless you want to override them._

## Step 5: Verify Deployment

1.  Go to the **App** tab.
2.  Wait for the status to change from **Building** to **Running**.
3.  If successful, you will see your application output or a blank screen (since it's a backend API).
4.  You can verify it's running by checking the **Logs** tab for server startup messages.
