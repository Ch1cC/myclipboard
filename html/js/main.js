import { Popover } from "bootstrap";
import pako from "pako";

let protocol = "ws://";
const port = 9090;
if (window.location.protocol === "https:") protocol = "wss://";
let ws = {};
const go = new Go();
const wasm = fetch(new URL("../statics/wasm.wasm", import.meta.url));
if ("instantiateStreaming" in WebAssembly) {
    WebAssembly.instantiateStreaming(wasm, go.importObject).then(function (
        obj
    ) {
        go.run(obj.instance);
        connect_WebSocket(encryptedFunc());
    });
} else {
    wasm.then((resp) => resp.arrayBuffer()).then((bytes) =>
        WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
            go.run(obj.instance);
            connect_WebSocket(encryptedFunc());
        })
    );
}
const tableRef = document.getElementsByTagName("tbody")[0];
// 获取输入容器和内容元素
const container = document.getElementById("imageContainer");
const spinner = document.getElementById("spinner");
const main = document.getElementById("main");
// 创建 DOMParser 实例
const parser = new DOMParser();
function renderItem(item) {
    tableRef.insertRow().innerHTML = `<th scope='row'>
            ${new Date(item.unix * 1000).toLocaleString("zh-CN")}
            </th>
            <td>
                <div style="white-space:normal;word-wrap:break-word;word-break:break-all;">
                    ${item.msg}
                </div>
            </td>
            <td>
            <span class="d-grid gap-2">
                <a 
                    role="button"  
                    tabindex="0"
                    type="button" 
                    ${
                        extractFirstUrl(item.msg)
                            ? `class="btn btn-primary copy">打开`
                            : `data-bs-container="body" 
                    class="btn btn-dark copy" 
                    data-bs-trigger="focus" 
                    data-bs-toggle="popover" 
                    data-bs-placement="left" 
                    data-bs-content="复制成功">
                    复制`
                    }
                </a>
                ${
                    navigator.share
                        ? `<a 
                role="button"  
                tabindex="0"
                type="button"
                data-bs-container="body" 
                class="btn btn-success share" 
                data-bs-trigger="focus">
                分享
            </a>`
                        : ``
                }
                
            </span></td>`;
}

function isSafari() {
    const ua = navigator.userAgent;
    const vendor = navigator.vendor;

    return /Safari/.test(ua) || /Apple Computer/.test(vendor);
}
async function fetchBlob(src) {
    const data = await fetch(src);
    return await data.blob();
}
function getImgByBase64(base64Img) {
    const doc = parser.parseFromString(base64Img, "text/html");
    // 从 DOM 文档中提取图像元素
    const img = doc.querySelector("img[src]");
    if (img) {
        // 将 Blob 转换为图像 URL
        const blob = base64ToBlob(img.src);
        img.src = URL.createObjectURL(blob);
        return img.outerHTML;
    }
    return "";
}
function extractFirstUrl(text) {
    // 定义匹配 URL 的正则表达式模式
    var pattern = /https?:\/\/\S+/i;

    // 使用正则表达式进行匹配
    var match = text.match(pattern);

    // 如果找到匹配则返回第一个 URL
    if (match && !match[0].includes(window.location.hostname)) {
        return match[0];
    } else {
        return false;
    }
}
// 将 Base64 编码字符串转换为 Blob 对象
function base64ToBlob(base64) {
    const binaryString = window.atob(base64.split(",")[1]);
    const arrayBuffer = new ArrayBuffer(binaryString.length);
    const uint8Array = new Uint8Array(arrayBuffer);
    for (let i = 0; i < binaryString.length; i++) {
        uint8Array[i] = binaryString.charCodeAt(i);
    }
    return new Blob([uint8Array], { type: `image/png` });
}
// 添加点击事件监听器
document.getElementById("submit").addEventListener("click", function () {
    submit();
});
function submit(value) {
    const text = value ? value : container.textContent;
    // 解析 HTML 字符串
    var doc = parser.parseFromString(text, "text/html");
    if (
        !Array.from(doc.body.childNodes).some((node) => node.nodeType === 1) &&
        text
    ) {
        send(text);
    } else {
        //判断是否有img.src
        var img = doc.querySelector("img[src]");
        // 使用 Fetch API 获取远程图片数据
        if (img) {
            fetch(img.src)
                .then((response) => response.blob()) // 将响应数据转换为 Blob 对象
                .then((blob) => {
                    handleImagePaste(blob); // 处理粘贴的图片
                })
                .catch((error) =>
                    console.error("Failed to fetch image:", error)
                );
        }
    }
    container.textContent = "";
}
function handleImagePaste(blob) {
    const img = new Image();
    // 将 Blob 转换为图像 URL
    img.src = URL.createObjectURL(blob);
    img.onload = function () {
        const { width, height } = img;
        // 创建canvas画布
        const canvas = document.createElement("canvas");
        const context = canvas.getContext("2d");
        if (!context) {
            console.error("Canvas 2D context is not supported.");
            return;
        }
        canvas.width = width;
        canvas.height = height;
        context.drawImage(img, 0, 0, canvas.width, canvas.height);
        // 将 Canvas 数据导出为 Blob
        const base64String = canvas.toDataURL(
            "image/jpeg", // 输出格式
            0.7 // 输出质量（0-1）
        );
        // 将 Base64 字符串用作图像的 src，回显到页面上
        const compressedImg = new Image();
        compressedImg.src = base64String;
        send(compressedImg.outerHTML);
    };
    img.onopen;
}
function handleTextPaste(item) {
    item.getAsString(function (text) {
        submit(text);
    });
}

function decompressDataAndRender(list) {
    main.style.display = "block";
    spinner.style.display = "none";
    tableRef.innerHTML = "";
    list.forEach((element) => {
        // 使用 atob 将 base64字符串解码为二进制数据
        const binaryString = atob(element.msg);
        const plain = new Uint8Array(
            [...binaryString].map((char) => char.charCodeAt(0))
        );

        try {
            // 解压gzip数据
            const uncompressedData = pako.ungzip(plain, {
                to: "string",
            });
            // 输出解压后的数据;
            const img = getImgByBase64(uncompressedData);
            if (img) {
                element.msg = img;
            } else {
                element.msg = uncompressedData;
            }
            renderItem(element);
        } catch (error) {
            // console.log(error);
            //未压缩的预设值
            const decoder = new TextDecoder();
            const text = decoder.decode(plain);
            element.msg = text;
            renderItem(element);
        }
    });
    const popoverTriggerList = document.querySelectorAll(
        '[data-bs-toggle="popover"]'
    );
    const popoverList = [...popoverTriggerList].map(
        (popoverTriggerEl) => new Popover(popoverTriggerEl)
    );
}
function share(event) {
    const text =
        event.target.parentNode.parentNode.parentNode.children[1].textContent;
    //如果text节点没文本.代表是图片
    if (!text.trim().length) {
        fetchBlob(
            event.target.parentNode.parentNode.parentNode.children[1]
                .children[0].children[0].src
        )
            .then(function (blob) {
                // 创建 File 对象
                var file = new File([blob], "image.png", {
                    type: "image/png",
                });

                // 调用 Web Share API 进行分享
                navigator
                    .share({
                        files: [file],
                    })
                    .then(function () {
                        // console.log("Image shared successfully");
                    })
                    .catch(function (error) {
                        console.error("Error sharing image:", error);
                    });
            })
            .catch(function (error) {
                console.error("Error getting Blob:", error);
            });
    } else {
        navigator.share({
            text: text.trim(),
        });
    }
}
//showImage
function showImage(event) {
    const img = event.target.src;
    //直接打开
    window.open(img);
}
function removeAllEventListener() {
    const copys = document.getElementsByClassName("copy");
    for (let index = 0; index < copys.length; index++) {
        const element = copys[index];
        element.removeEventListener("click", copy);
    }
    const shares = document.getElementsByClassName("share");
    for (let index = 0; index < shares.length; index++) {
        const element = shares[index];
        element.removeEventListener("click", share);
    }
    //获取所有img tag
    const imgs = document.getElementsByTagName("img");
    for (let index = 0; index < imgs.length; index++) {
        const element = imgs[index];
        element.removeEventListener("click", showImage);
    }
}
function addAllEventListener() {
    const copys = document.getElementsByClassName("copy");
    for (let index = 0; index < copys.length; index++) {
        const element = copys[index];
        // 添加点击事件监听器
        element.addEventListener("click", copy);
    }
    const shares = document.getElementsByClassName("share");
    for (let index = 0; index < shares.length; index++) {
        const element = shares[index];
        // 添加点击事件监听器
        element.addEventListener("click", share);
    }
    //获取所有img tag
    const imgs = document.getElementsByTagName("img");
    for (let index = 0; index < imgs.length; index++) {
        //添加点击事件
        const element = imgs[index];
        element.addEventListener("click", showImage);
    }
}
function copy(event) {
    const text =
        event.target.parentNode.parentNode.parentNode.children[1].textContent;
    //如果text节点没文本.代表是图片
    if (!text.trim().length) {
        const clipboardItem = new ClipboardItem({
            [`image/png`]: fetchBlob(
                event.target.parentNode.parentNode.parentNode.children[1]
                    .children[0].children[0].src
            ),
        });
        navigator.clipboard.write([clipboardItem]).catch(function (error) {
            console.log(error);
        });
        //
    } else {
        const haveUrl = extractFirstUrl(text.trim());
        if (haveUrl) {
            // 使用 window.open() 打开 URL 在新窗口中
            window.open(haveUrl, "_blank");
        } else {
            navigator.clipboard.writeText(text.trim());
        }
    }
}
function connect_WebSocket(encryptedData) {
    ws = new WebSocket(
        protocol +
            window.location.hostname +
            ":" +
            port +
            window.location.pathname +
            `ws?token=${encryptedData}`
    );
    ws.onmessage = function (e) {
        var reader = new FileReader();
        reader.onload = function () {
            const uncompressedData = pako.ungzip(reader.result, {
                to: "string",
            });
            const list = JSON.parse(uncompressedData);
            removeAllEventListener();
            decompressDataAndRender(list);
            addAllEventListener();
        };
        reader.readAsArrayBuffer(e.data);
    };
    ws.onopen = function (e) {
        // console.log("开启了");
    };
    ws.onerror = function (e) {
        // console.log(e)
        console.log("错误了");
    };
    ws.onclose = function (e) {
        // console.log(e)
        console.log("关闭了");
        main.style.display = "none";
        spinner.style.display = "flex";
        connect_WebSocket(encryptedFunc());
    };
}

function send(params) {
    ws.send(pako.gzip(params, { level: 6 }));
}

// 监听容器的粘贴事件
container.addEventListener("paste", function (e) {
    // 取消默认粘贴行为
    e.preventDefault();

    // 获取粘贴的数据
    const items = (e.clipboardData || e.originalEvent.clipboardData).items;
    for (let i = 0; i < items.length; i++) {
        const item = items[i];
        // 判断是否为图片
        if (item.type.indexOf("image") !== -1) {
            handleImagePaste(item.getAsFile());
        } else {
            handleTextPaste(item);
        }
        break;
    }
});
