import { computed } from 'vue'
import * as NaiveUI from 'naive-ui'
import { useThemeStore } from '@/store'
import themes from '@/assets/themes'

function setupMessage(NMessage) {
    class Message {
        static instance // 单例实例
        message // 用于存储消息实例
        removeTimer // 用于存储移除定时器

        constructor() {
            // 如果实例已存在，直接返回实例
            if (Message.instance) {
                return Message.instance
            }
            Message.instance = this
            this.message = {}
            this.removeTimer = {}
        }

        // 销毁指定消息，延时一定时间后执行
        destroy(key, duration = 200) {
            setTimeout(() => {
                if (this.message[key]) {
                    this.message[key].destroy()  // 销毁消息
                    delete this.message[key]     // 删除消息记录
                }
            }, duration)
        }

        // 延时移除消息，若消息已存在则清除定时器重新计时
        removeMessage(key, duration = 5000) {
            this.removeTimer[key] && clearTimeout(this.removeTimer[key])  // 清除之前的定时器
            this.removeTimer[key] = setTimeout(() => {
                this.message[key]?.destroy()  // 超时后销毁消息
            }, duration)
        }

        // 根据类型和选项展示消息
        showMessage(type, content, option = {}) {
            if (Array.isArray(content)) {
                // 如果 content 是数组，遍历并显示每一条消息
                return content.forEach(msg => NMessage[type](msg, option))
            }

            // 如果没有指定 key，则直接显示消息
            if (!option.key) {
                return NMessage[type](content, option)
            }

            // 获取当前 key 对应的消息
            const currentMessage = this.message[option.key]
            if (currentMessage) {
                // 如果消息已存在，更新其类型和内容
                currentMessage.type = type
                currentMessage.content = content
            }
            else {
                // 如果消息不存在，创建新的消息实例
                this.message[option.key] = NMessage[type](content, {
                    ...option,
                    duration: 0, // 防止自动销毁
                    onAfterLeave: () => {
                        delete this.message[option.key]  // 销毁后删除消息实例
                    },
                })
            }
            // 设置消息移除定时器
            this.removeMessage(option.key, option.duration)
        }

        // 不同类型的消息显示方法封装
        loading(content, option = { duration: 0 }) {
            this.showMessage('loading', content, option)
        }

        success(content, option = {}) {
            this.showMessage('success', content, option)
        }

        error(content, option = {}) {
            this.showMessage('error', content, option)
        }

        info(content, option = {}) {
            this.showMessage('info', content, option)
        }

        warning(content, option = {}) {
            this.showMessage('warning', content, option)
        }
    }

    return new Message()  // 返回实例化的消息对象
}

function setupDialog(NDialog) {
    // 修改 NDialog 的 confirm 方法
    NDialog.confirm = function (option = {}) {
        const showIcon = !!(option.title)  // 如果有标题，则显示图标
        return NDialog[option.type || 'warning']({
            showIcon,  // 是否显示图标
            positiveText: '确定',  // 确认按钮文字
            negativeText: '取消',  // 取消按钮文字
            onPositiveClick: option.confirm,  // 点击确认按钮的回调
            onNegativeClick: option.cancel,  // 点击取消按钮的回调
            onMaskClick: option.cancel,  // 点击遮罩层的回调
            ...option,  // 合并额外的配置
        })
    }
    return NDialog
}

/**
 * 挂载 NaiveUI API
 */
export function setupNaiveDiscreteApi() {
    const themeStore = useThemeStore()  // 获取主题 store
    const configProviderProps = computed(() => ({
        theme: themeStore.darkMode ? NaiveUI.darkTheme : undefined,  // 根据主题状态切换暗黑模式
        themeOverrides: themes.themeOverrides,  // 主题自定义覆盖
    }))

    // 创建 Naive UI 的离散 API 实例
    const { message, dialog, notification, loadingBar } = NaiveUI.createDiscreteApi(
        ['message', 'dialog', 'notification', 'loadingBar'],
        { configProviderProps },
    )

    // 挂载到全局对象上，方便在其他地方访问
    window.$loadingBar = loadingBar
    window.$notification = notification
    window.$message = setupMessage(message)  // 初始化消息系统
    window.$dialog = setupDialog(dialog)    // 初始化对话框系统
}

/**
 * 解决 naive-ui 和 unocss 样式冲突
 */
export function setupNaiveUnocss() {
    const meta = document.createElement('meta')  // 创建一个 meta 标签
    meta.name = 'naive-ui-style'  // 设置标签的 name 属性
    document.head.appendChild(meta)  // 将标签添加到文档的头部
}