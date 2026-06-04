@echo off
echo ========================================
echo   启动 Demo 用户登录系统
echo ========================================

echo.
echo [1/3] 启动后端服务 (端口 8080)...
cd /d "%~dp0demo-server"
start "Demo-Server" cmd /k "go run main.go"

timeout /t 2 /nobreak > nul

echo [2/3] 启动前端服务 (端口 3000)...
cd /d "%~dp0demo-web"
start "Demo-Web" cmd /k "npm run dev"

echo [3/3] 完成!
echo.
echo 服务地址:
echo   - 前端: http://localhost:3000
echo   - 后端: http://localhost:8080
echo.
echo 按任意键打开浏览器...
pause > nul

start http://localhost:3000
