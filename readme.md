# Demo Programs

An api proxy and frontend to demo communication to the actual sms global api. No makefile or other build script is included in the mono repo so follow instructions in the 'to run' sections.

# Api

Written in go. 

## To run

Having the `go` binary on your PATH is the only requirement. `cd` into the 'api' folder and run a `go build`. Code code block below provides examples of all possible options that can be parsed to the application.

```sh
go build -o /tmp/api .
API_KEY=your_key API_SECRET=your_secret /tmp/api --human --debug --port 9090
```

## What does it do

The go api simply forwards requests it receives on its /api/message endpoint to the /v2/sms endpoint of the actual sms global api.

It writes the sms global response directly into its own response without doing any parsing

The programs uses only a few external dependencies mostly to provide a few nice extras around middleware grouping, logging and environment variable parsing. The primary function of the application uses go std lib http handlers without any external router library.

# Frontend

React based frontend app using  snowpack as the bundler.

## To run

Having `node` + `npm` installed on your path are the only requirements. First you must create a `./fronted/.env` file that includes `SNOWPACK_PUBLIC_API_URL=localhost:9090` as an environment variable. If you are running the api on a different machine or change the port that the api listens on then you will need to adjust this value.

One the .env file is setup you can run `npm install` and `npm run start` in the `frontend` folder.

## Notes

I used a few libraries I haven't used much before as an opportunity to try them out. A few elements of the application are a little rushed. The css styling with tailwind especially. While the application does function correctly it could use a bit of polish.

## Libraries in use

- React
- React Hook Form
- Tailwindcss
- Headless UI
- swr
- snowpack
- computed-types