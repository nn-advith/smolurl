### SmolURL

<hr/>

A simple url-compress applicatoin that generates a short-url from a long one. Uses Couchbase as database and other databases can be included with intuitive code changes. Refer `docs/` in `dev` branch for detailed design doc.

Example:

- URL generation:

```
curl -v -X POST "http://localhost:4000/hash/" -d '{"url":"https://nnadvith.netlify.app"}'
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:4000...
* Connected to localhost (127.0.0.1) port 4000 (#0)
> POST /hash/ HTTP/1.1
> Host: localhost:4000
> User-Agent: curl/7.81.0
> Accept: */*
> Content-Length: 38
> Content-Type: application/x-www-form-urlencoded
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sat, 10 May 2025 06:52:27 GMT
< Content-Length: 35
< Content-Type: text/plain; charset=utf-8
<
http://localhost:4000/TNjaBPT_NTU=
* Connection #0 to host localhost left intact
```

- Visit smolurl:

```
curl -v http://localhost:4000/TNjaBPT_NTU=
*   Trying 127.0.0.1:4000...
* Connected to localhost (127.0.0.1) port 4000 (#0)
> GET /TNjaBPT_NTU= HTTP/1.1
> Host: localhost:4000
> User-Agent: curl/7.81.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 301 Moved Permanently
< Content-Type: text/html; charset=utf-8
< Location: https://nnadvith.netlify.app
< Date: Sat, 10 May 2025 06:54:31 GMT
< Content-Length: 63
<
<a href="https://nnadvith.netlify.app">Moved Permanently</a>.

* Connection #0 to host localhost left intact
```
