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

function decompressData(list) {
    list.forEach((element) => {
        // 使用 atob 将 base64 字符串解码为二进制数据
        const binaryString = atob(element.msg);
        const plain = new Uint8Array(
            [...binaryString].map((char) => char.charCodeAt(0))
        );
        try {
            // 解压gzip数据
            const uncompressedData = pako.inflate(plain, {
                to: "string",
            });
            // 输出解压后的数据
            element.msg = uncompressedData;
        } catch (error) {
            // 使用 TextDecoder 将二进制数据解码为文本
            const decoder = new TextDecoder();
            const text = decoder.decode(plain);
            element.msg = text;
        }
    });
}

function webSocket(encryptedData) {
    const webSocket = new WebSocket(
        protocol +
            window.location.host +
            window.location.pathname +
            `ws?token=${encryptedData}`
    );
    webSocket.onmessage = function (e) {
        var reader = new FileReader();
        reader.onload = function () {
            var text = reader.result;
            const uncompressedData = pako.inflate(text, {
                to: "string",
            });
            const list = JSON.parse(uncompressedData);
            decompressData(list);
            render(list);
        };
        reader.readAsArrayBuffer(e.data);
    };
    webSocket.onopen = function (e) {
        ws = webSocket;
        // console.log("开启了");
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
    ws.send(pako.gzip(params, { to: "binary", level: 9 }));
}
