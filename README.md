# xBrowser

## About

xBrowser is an Application of Wails(Golang + Vue3).

## Live Development

To run in live development mode, run `wails dev` in the project directory. In another terminal, go into the `frontend`
directory and run `npm run dev`. The frontend dev server will run on http://localhost:34115. Connect to this in your
browser and connect to your application.

## Building

To build a redistributable, production mode package, use `wails build`.

## Building prepare
go install github.com/wailsapp/wails/v2/cmd/wails@latest

wails build -platform windows/amd64 2>&1
wails build -platform windows/amd64 -webview2 embed 2>&1