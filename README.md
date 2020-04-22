# WeatherApp

```shell script
API_KEY=xxx docker-compose up
```

The above command will start the server using the provided `API_KEY`

Now call http://localhost:8080/?c=london&c=warsaw&cc=uk&cc=pl&c=Kraków&c=Wrocław&c=Łomża&c=Chełm&c=Brwinów&c=Hobbiton to get the weather conditions for London, Warsaw, etc.

As you might notice, there is also a non-existent city in the query, but I chose to not return any errors when bulk data is received. Server returns only what is found.

## Cache

I configured cache to be valid for a minute per location (there are many optimisations that could be made, but that is not the issue here).

Try repeating the above API call multiple times and see what the logs show. You will not see valid and existing GETs for a minute.

## Testing

### Unit Tests
Tests are being run on build

If you want to run them again, just do
```shell script
docker build .
```
`--no-cache` should not be needed as the build is non-deterministic (could be tweaked to be, but that is not the point)

### Integration Tests
I have compiled the tests into a separate binary, so they can be run in a suitable environment.

```shell script
API_KEY=xxx docker-compose run --rm --entrypoint="./weatherapp.test" app -test.v
```

This runs tests which actually GET the data from a real server, so a valid `API_KEY` is required. 

## API Payload

I chose to not concern myself with the output data format.

## Docker and Compose

I chose to use `docker-compose` for setting up a running service.

### Build

The build is a straightforward one, but I chose to use "build stages" just form fun.

## Notes

Generally there are many shortcuts taken here. I wanted to demonstrate various techniques at my disposal.
I will be glad to discuss the intricacies of this solution in more detail.

### Package Organisation

I chose to not use `pkg` or `cmd` package directories on purpose. My goal was to keep a flat structure here. I really like Clean Architecture. 
Of course, as the project grows it might be necessary to move files around, but that is something I let to be changed as the project progresses.

### Repository Name
> Name of the repository must be as follows: base64 of (name + last name + “recruitment task”)

I took this requirement literally, but seeing my last name concatenated without any spaces with the "recruitment task" phrase makes my eyes bleed. 