let protocol = "ws://";
if (window.location.protocol === "https:") protocol = "wss://";
let ws = {};
function str2ab(str) {
    const buf = new ArrayBuffer(str.length);
    const bufView = new Uint8Array(buf);
    for (let i = 0, strLen = str.length; i < strLen; i++) {
        bufView[i] = str.charCodeAt(i);
    }
    return buf;
}
const publicKeyPem = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA64gDlAFUlVNc/Fm1dN4k
njxrok2Y6C4mmnt0FDrC/jYO+pxzZ6mPkVS/JQuweHmrYVkQ6RJSXKew8I+2ukJc
Ny+N43ZuSPDqCHVECQlQkClTAug139cGBiMaUNnSWj2/d6R8DLXVYgfDuqPWBaCp
PJ9+jjy9WYGheoE/n5MPAhNSMqP4PDqt+auVJcWgVCrizeO/GuUn84Fm4J98Ln9s
9CqWcg/JSTGF1Za55FG9BfThW5bM0L+SpTKXzGco7jQ4QDF+bPFCzbzbUR638AbE
lHCBT+jGRhwwzWBJ8Z9bWY7NaYHJUv6OiQ+5J3OhcNkMa8rf8tIfCovNN3BPmJhl
9wIDAQAB
-----END PUBLIC KEY-----`;
// 加载明文的公钥
async function loadPublicKey(pem) {
    // 获取 PEM 字符串在头部和尾部之间的部分
    const pemHeader = "-----BEGIN PUBLIC KEY-----";
    const pemFooter = "-----END PUBLIC KEY-----";
    const pemContents = pem.substring(
        pemHeader.length,
        pem.length - pemFooter.length
    );
    // 将字符串通过 base64 解码为二进制数据
    const binaryDerString = window.atob(pemContents);
    // 将二进制字符串转换为 ArrayBuffer
    const binaryDer = str2ab(binaryDerString);
    return window.crypto.subtle.importKey(
        "spki",
        binaryDer,
        {
            name: "RSA-OAEP",
            hash: "SHA-256",
        },
        true,
        ["encrypt"]
    );
}
// 3. 使用公钥加密数据
async function encryptWithPublicKey(publicKey, data) {
    return await window.crypto.subtle.encrypt(
        {
            name: "RSA-OAEP",
        },
        publicKey,
        data
    );
}
// 3. 使用公钥加密数据
// 示例
async function encryptedData() {
    const data = new TextEncoder().encode(Math.floor(Date.now() / 1000));
    const publicKey = await loadPublicKey(publicKeyPem);
    // 3. 使用公钥加密数据
    const encryptedData = await encryptWithPublicKey(publicKey, data);
    webSocket(encryptedData);
}
function webSocket(encryptedData) {
    const webSocket = new WebSocket(
        protocol +
            window.location.host +
            window.location.pathname +
            `ws?token=${Array.from(new Uint8Array(encryptedData))
                .map((byte) => byte.toString(16).padStart(2, "0"))
                .join("")}`
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
encryptedData();
function send(params) {
    ws.send(params);
}
