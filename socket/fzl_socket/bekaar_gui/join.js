io = require("socket.io")
if (window["WebSocket"]) {
    var roomid = $('#roomid').val();
    conn = new WebSocket("ws://localhost:8080/ws/join/" + roomid);

    console.log("join new scoket")
    conn.onclose = function(evt) {
        appendLog($("<div><b>Connection closed.<\/b><\/div>"))
    }
    conn.onmessage = function(evt) {
        appendLog($("<div/>").text(evt.data))
    }
} else {
    appendLog($("<div><b>Your browser does not support WebSockets.<\/b><\/div>"))
}