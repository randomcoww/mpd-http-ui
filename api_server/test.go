//
// websocket test adapted from https://github.com/gorilla/websocket/blob/master/examples/filewatch/main.go
//

package main

import (
  "text/template"
	"net/http"
)


var homeTempl = template.Must(template.New("").Parse(homeHTML))

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var v = struct {
		Host    string
	}{
		r.Host,
	}
	homeTempl.Execute(w, &v)
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <h2>WebSocket Example</h2>
        <pre id="message"></pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("message");
                var conn = new WebSocket("ws://{{.Host}}/ws");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('got message');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
