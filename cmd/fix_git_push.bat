@echo off
echo ==========================================
echo Fix Git Push Issue
echo ==========================================

echo [1/3] Checking Git Status...
git status
if %ERRORLEVEL% NEQ 0 (
    echo Error: Not a git repository or git is not installed.
    pause
    exit /b
)

echo.
echo [2/3] Setting up 'dev' branch...
REM Try to create the branch. If it fails, it might exist.
git checkout -b dev 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Branch 'dev' may already exist. Switching to it...
    git checkout dev
)

echo.
echo [3/3] Pushing to remote...
git push -u origin dev

if %ERRORLEVEL% EQU 0 (
    echo.
    echo SUCCESS: Successfully pushed to dev branch!
) else (
    echo.
    echo FAILED: Could not push to remote. Please check:
    echo 1. You have internet connection.
    echo 2. You have permission to push to this repository.
    echo 3. The remote URL is correct.
)

echo.
pause
