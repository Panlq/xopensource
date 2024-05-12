// 创建 Vue 实例
const { createApp, ref } = Vue;

createApp({
  data() {
    return {
      events: [], // 保存接收到的事件信息
      error: "", // 保存 SSE 连接错误信息
      eventSource: null, // 保存 EventSource 对象
      startButtonDisabled: false, // 控制开始按钮是否禁用
      stopButtonVisible: false, // 控制停止按钮是否可见
    };
  },
  methods: {
    start() {
      // 发起请求调用 /start 接口
      fetch("/start", { method: "GET" })
        .then((response) => {
          if (response.ok) {
            // 开始 SSE 连接
            this.eventSource = new EventSource(
              "http://127.0.0.1:8844/stream?stream=message"
            );

            // 监听 open 事件，处理连接成功
            this.eventSource.onopen = (event) => {
              console.log("连接成功");
              this.startButtonDisabled = true; // 禁用开始按钮
              this.stopButtonVisible = true; // 显示停止按钮
            };

            // 监听 error 事件，处理连接错误
            this.eventSource.addEventListener("error", (error) => {
              console.error("SSE 连接错误:", error);
              this.error = "SSE 连接错误: " + error.message;
            });

            // 监听关闭事件，处理连接关闭
            this.eventSource.addEventListener("close", (event) => {
              console.log("SSE 连接已关闭");
              this.stopButtonVisible = false; // 隐藏停止按钮
              this.startButtonDisabled = false; // 启用开始按钮
              this.events.push(event.data);
            });

            // 监听 message 事件，接收到新事件时触发
            this.eventSource.addEventListener("message", (event) => {
              // 更新事件信息到页面
              this.events.push(event.data);
            });
          } else {
            this.error = "Failed to start SSE";
          }
        })
        .catch((error) => {
          console.error("Failed to start SSE:", error);
          this.error = "Failed to start SSE";
        });
    },
    stop() {
      // 发起请求调用 /stop 接口
      fetch("/stop", { method: "GET" })
        .then((response) => {
          if (response.ok) {
            console.log("SSE 停止成功");
            // 关闭 SSE 连接
            if (this.eventSource) {
              this.eventSource.close();
              console.log("SSE 连接已关闭");
            }
          } else {
            throw new Error("Failed to stop SSE");
          }
        })
        .catch((error) => {
          console.error("Failed to stop SSE:", error);
          this.error = "Failed to stop SSE";
        });
    },
    clearMessages() {
      this.events = []; // 清空消息列表
    },
  },
}).mount("#app");
