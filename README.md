# desmume-launcher

##About DeSmuME Launcher

Welcome to DeSmuME Launcher, a sleek and user-friendly emulator launcher designed for enthusiasts of Nintendo DS gaming! Crafted with passion and powered by the Wails framework in Golang, this application provides a seamless experience for running your favorite DS games on your computer.

## Building

To build this project in debug mode, use `wails build`. For production, use `wails build -production`.
To generate a platform native package, add the `-package` flag.

## Live Development

To run in live development mode, run `wails dev` in the project directory. In another terminal, go into the `frontend`
directory and run `npm run dev`. The frontend dev server will run on http://localhost:34115. Connect to this
in your browser and connect to your application.
