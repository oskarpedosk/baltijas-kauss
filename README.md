# Baltijas Kauss™

Baltijas Kauss is an e-sports league web application that enables users to create teams, view player statistics, draft players, enter game results, view the league table, and more. The players' database is created through web scraping since there is no API available for players' attributes.

## Server
Server runs on Ubuntu 22.04.2 LTS  
[baltijaskauss.ee](https://baltijaskauss.ee/)

### Built with
- Go version 1.20.3: https://go.dev/doc/go1.20
Programming language used for the backend development.
- Chi router: https://github.com/go-chi/chi  
Lightweight and fast HTTP router for Go.
- Alex Edwards SCS: https://github.com/alexedwards/scs  
Secure cookie session management for Go.
- Nosurf: https://github.com/justinas/nosurf  
CSRF (Cross-Site Request Forgery) protection for Go web applications.
- Gorilla WebSocket: https://github.com/gorilla/websocket  
Implementation of WebSocket protocol for Go.
- npm: https://www.npmjs.com/  
Package manager for Node.js and JavaScript.
- Caddy: https://caddyserver.com/  
Open-source web server with automatic HTTPS.
- Node.js: https://nodejs.org/  
Used with Puppeteer for web scraping for data.
- Chromium: https://www.chromium.org/Home  
Open-source web browser project on which Google Chrome is based.
- Puppeteer: https://pptr.dev/  
Node.js library for controlling headless Chrome or Chromium.
