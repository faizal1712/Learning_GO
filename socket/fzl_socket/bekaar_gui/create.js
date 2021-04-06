// $("#form").submit(function() {
//     if (!conn) {
//         return false;
//     }
//     if (!msg.val()) {
//         return false;
//     }
//     conn.send(msg.val());
//     msg.val("");
//     return false
// });

var io = require("socket.io")


console.log("created new socket")
var conn = io("localhost:8080/", {
    transports: ['websocket']
});
console.log("created new socket")
conn.on('message', function(message) {
    console.log('new message');
    console.log(message);
});

conn.on('connect', function() {
    console.log('socket connected');
    //send something
    socket.emit('send', {
        name: "my name",
        message: "hello"
    }, function(result) {
        console.log('sent successfully');
        console.log(result);
    });
});
// conn.onclose = function(evt) {
//     appendLog($("<div><b>Connection closed.<\/b><\/div>"))
// }
// conn.onmessage = function(evt) {
//     appendLog($("<div/>").text(evt.data))
// }