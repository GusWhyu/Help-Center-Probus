@echo off
title Help Center Probus - All in One Runner
cls

echo ========================================================
echo   Membersihkan port yang masih aktif (8088, 5173)...
echo ========================================================
for /f "tokens=5" %%a in ('netstat -aon ^| findstr ":8088" ^| findstr "LISTENING"') do (
    taskkill /F /PID %%a >nul 2>&1
)
for /f "tokens=5" %%a in ('netstat -aon ^| findstr ":5173" ^| findstr "LISTENING"') do (
    taskkill /F /PID %%a >nul 2>&1
)

echo ========================================================
echo   Memulai Help Center Probus (Backend Go + Frontend Vite)
echo ========================================================
echo.
echo   - Backend  : http://localhost:8088/api/v1
echo   - Frontend : http://localhost:5173
echo.
echo   Tekan Ctrl + C untuk mematikan kedua service sekaligus.
echo ========================================================
echo.

npx.cmd --yes concurrently -n "BACKEND,FRONTEND" -c "blue.bold,green.bold" --kill-others "go run ./cmd/server" "npm.cmd --prefix frontend run dev"
