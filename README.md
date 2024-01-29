# mock-es
A mock elasticsearch server for testing.  This provides the absolute minimum responses for something like filebeat to send bulk requests to and to provide a response.  The data is thrown away.  By default successful responses are sent for all data.

It is also possible to return errors.  It is possible to configure so that each POST to the bulk endpoint can result in an error.  It is also possible to configure it so that each create action in the bulk request has the possibility of an error.  Any other action always return StatusOK.

| Flag      | Meaning                                                                           |
|-----------|-----------------------------------------------------------------------------------|
| -a string | address to listen on ip:port (default ":9200")                                    |
| -d uint   | percent chance StatusConflict is returned for create action                       |
| -l uint   | percent chance StatusEntityTooLarge is returned for POST method on _bulk endpoint |
| -n uint   | percent chance StatusNotAcceptable is returned for create action                  |
| -t uint   | percent chance StatusTooManyRequests is returned for create action                |

`-d`, `-n`, and `-t` cannot total more than 100 percent, if the total is less than 100, the remaining percentage chance is for StatusOK.

eg:

`-d 20 -n 10 -t 5` would result in a 20 percent chance StatusConflict, 10 percent chance StatusNotAcceptable, 5 percent chance StatusTooManyRequests and 65 percent chance that StatusOK would be returned for each create action in the bulk request

`-l` cannot be larger than 100, if the total is less than 100, the remaining percentage chance is for StatusOK

eg:

`-l 25`  would result in a 25 percent chance that each POST to the bulk endpoint would result in StatusEntityTooLarge, the remaining 75 percent would be StatusOK.

