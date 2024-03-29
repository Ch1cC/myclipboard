let protocol = "ws://";
let ws = {};
const go = new Go();
const WASM_URL = "static/wasm.wasm";
if ("instantiateStreaming" in WebAssembly) {
    WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(
        function (obj) {
            go.run(obj.instance);
            webSocket(encryptedFunc());
        }
    );
} else {
    fetch(WASM_URL)
        .then((resp) => resp.arrayBuffer())
        .then((bytes) =>
            WebAssembly.instantiate(bytes, go.importObject).then(function (
                obj
            ) {
                go.run(obj.instance);
                webSocket(encryptedFunc());
            })
        );
}
if (window.location.protocol === "https:") protocol = "wss://";
function webSocket(encryptedData) {
    const webSocket = new WebSocket(
        protocol +
            window.location.host +
            window.location.pathname +
            `ws?token=${encryptedData}`
    );
    webSocket.onmessage = function (e) {
        render(e.data);
    };
    webSocket.onopen = function (e) {
        ws = webSocket;
        console.log("开启了");
    };
    webSocket.onerror = function (e) {
        // console.log(e)
        console.log("错误了");
    };
    webSocket.onclose = function (e) {
        // console.log(e)
        console.log("关闭了");
        window.location.reload();
    };
}
function send(params) {
    ws.send(params);
}
