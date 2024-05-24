let protocol = "ws://";

if (window.location.protocol === "https:") protocol = "wss://";
let ws = {};
const go = new Go();
const wasm = fetch("static/wasm.wasm");
if ("instantiateStreaming" in WebAssembly) {
    WebAssembly.instantiateStreaming(wasm, go.importObject).then(function (
        obj
    ) {
        go.run(obj.instance);
        webSocket(encryptedFunc());
    });
} else {
    wasm.then((resp) => resp.arrayBuffer()).then((bytes) =>
        WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
            go.run(obj.instance);
            webSocket(encryptedFunc());
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
function render(items) {
    main.style.display = "block";
    spinner.style.display = "none";
    tableRef.innerHTML = "";
    for (const item of items) {
        tableRef.insertRow().innerHTML = `<th scope='row'>
            ${new Date(item.unix * 1000).toLocaleString("zh-CN")}
            </th>
            <td>
                <div style="width:100%;white-space:normal;word-wrap:break-word;word-break:break-all;">
                    ${item.msg}
                </div>
            </td>
            <span class="d-grid gap-2">
                <a 
                    role="button"  
                    tabindex="0"
                    type="button" 
                    onclick="copy(this)" 
                    ${
                        extractFirstUrl(item.msg)
                            ? `class="btn btn-primary">打开`
                            : `data-bs-container="body" 
                    class="btn btn-dark" 
                    data-bs-trigger="focus" 
                    data-bs-toggle="popover" 
                    data-bs-placement="left" 
                    data-bs-content="复制成功">
                    复制`
                    }
                </a>
            </span>`;
    }
    const popoverTriggerList = document.querySelectorAll(
        '[data-bs-toggle="popover"]'
    );
    const popoverList = [...popoverTriggerList].map(
        (popoverTriggerEl) => new bootstrap.Popover(popoverTriggerEl)
    );
}

function copy(e) {
    const text = e.parentNode.parentNode.children[1].textContent;
    //如果text节点没文本.代表是图片
    if (!text.trim().length) {
        // 创建 Blob 对象
        const blob = base64ToBlob(
            e.parentNode.parentNode.children[1].children[0].children[0].src
        );
        const data = [new ClipboardItem({ ["image/png"]: blob })];
        // 将 DataTransfer 对象的数据写入剪切板
        navigator.clipboard.write(data);
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
function extractFirstUrl(text) {
    // 定义匹配 URL 的正则表达式模式
    var pattern = /https?:\/\/\S+/i;

    // 使用正则表达式进行匹配
    var match = text.match(pattern);

    // 如果找到匹配则返回第一个 URL
    if (match) {
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
function submit(value) {
    const text = value
        ? value
        : document.getElementById("imageContainer").textContent;
    // 解析 HTML 字符串
    var doc = parser.parseFromString(text, "text/html");
    if (
        //纯文本
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
    document.getElementById("imageContainer").textContent = "";
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
// 定义一个函数来获取图片类型
function getImageType(dataUrl) {
    // 匹配 data:image/ 后的内容
    var regex = /^data:image\/([a-z]+);/;
    var match = dataUrl.match(regex);
    if (match) {
        // 如果匹配成功，返回图片类型
        return match[1];
    } else {
        // 如果匹配失败，返回空字符串或者 null（根据实际情况）
        return "";
    }
}
function handleTextPaste(item) {
    item.getAsString(function (text) {
        submit(text);
    });
}

function decompressData(list) {
    list.forEach((element) => {
        // 使用 atob 将 base64 字符串解码为二进制数据
        const binaryString = atob(element.msg);
        const plain = new Uint8Array(
            [...binaryString].map((char) => char.charCodeAt(0))
        );
        try {
            // 解压gzip数据
            const uncompressedData = pako.ungzip(plain, {
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
    ws = new WebSocket(
        protocol +
            window.location.host +
            window.location.pathname +
            `ws?token=${encryptedData}`
    );
    ws.onmessage = function (e) {
        var reader = new FileReader();
        reader.onload = function () {
            var text = reader.result;
            const uncompressedData = pako.ungzip(text, { to: "string" });
            const list = JSON.parse(uncompressedData);
            decompressData(list);
            render(list);
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
        webSocket(encryptedFunc());
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
