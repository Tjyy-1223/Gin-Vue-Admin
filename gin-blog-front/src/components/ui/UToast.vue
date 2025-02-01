<template>
    <Teleport to="body">
        <div v-show="flux.events.length" class="pointer-events-none fixed w-4/5 sm:w-[400px]" :class="{
            'left-1/2 -translate-x-1/2': align === 'center',
            'left-16': align === 'left',
            'right-16': align === 'right',
            'top-4': position === 'top',
            'bottom-6 ': position === 'bottom',
        }" :style="{ zIndex }">
            <!--  
            1. 指定渲染的根元素为 <ul> 标签
            2. enter-class -> enter-to-class: 进入时的初始状态，基于 position 的值来设置不同的动画
            3. leave-class -> leave-to-class 设置里开始的动画效果
            4. move-class 列表项移动时的过渡效果，使用平滑的过渡
        -->
            <TransitionGroup tag="ul" enter-active-class="transition ease-out duration-200"
                leave-active-class="transition ease-in duration-200 absolute w-full" :enter-class="position === 'bottom'
                    ? 'transform translate-y-3 opacity-0'
                    : 'transform -translate-y-3 opacity-0'" enter-to-class="transform translate-y-0 opacity-100"
                leave-class="transform translate-y-0 opacity-100" :leave-to-class="position === 'bottom'
                    ? 'transform translate-y-1/4 opacity-0'
                    : 'transform -translate-y-1/4 opacity-0'" move-class="ease-in-out duration-200"
                class="inline-block w-full">

                <!-- li 标签 -->
                <li v-for="event in flux.events" :key="event.id" :class="{
                    'pb-2': position === 'bottom',
                    'pt-2': position === 'top',
                }">
                    <!-- 这是一个插槽（slot）元素，用于为父组件提供一个可插入的区域 -->
                    <!-- 
                        这个模板是一个通用的通知组件（例如提示框），它：
                        根据 event.type 显示不同的图标和颜色（如成功、信息、警告、错误）。
                        显示 event.content 作为通知的文本内容。
                        提供一个关闭按钮，当点击时会调用 flux.remove(event) 来移除该通知。
                     -->
                    <slot :type="event.type" :content="event.content">
                        <div
                            class="pointer-events-auto w-full overflow-hidden rounded-lg bg-white ring-1 ring-black ring-opacity-5">
                            <div class="flex justify-between px-4 py-3">
                                <div class="flex items-center">
                                    <div class="mr-6 h-6 w-6" :class="{
                                        'i-mdi:check-circle text-green': event.type === 'success',
                                        'i-mdi:information-outline text-blue': event.type === 'info',
                                        'i-mdi:alert-outline text-yellow': event.type === 'warning',
                                        'i-mdi:alert-circle-outline text-red': event.type === 'error',
                                    }" />
                                    <div class="ml-1">
                                        <slot>
                                            <div> {{ event.content }} </div>
                                        </slot>
                                    </div>

                                </div>
                                <button v-if="closeable"
                                    class="i-mdi:close h-5 w-5 flex items-center justify-center rounded-full rounded-full p-1 font-bold text-gray-400 hover:text-gray-600"
                                    @click="flux.remove(event)" />
                            </div>
                        </div>
                    </slot>

                </li>
            </TransitionGroup>

        </div>
    </Teleport>
</template>

<script setup lang="ts">
import { reactive } from 'vue';

const props = defineProps({
    show: { type: Boolean, default: false },
    position: {
        type: String,
        default: 'top',
        validator: (value: string) => {
            const positions: string[] = ['top', 'bottom'];
            return positions.includes(value);
        },
    },
    align: {
        type: String,
        default: 'center',
        validator: (value: string) => {
            const alignments: string[] = ['left', 'center', 'right'];
            return alignments.includes(value);
        },
    },
    timeout: { type: Number, default: 2500 },
    queue: { type: Boolean, default: true },
    zIndex: { type: Number, default: 100 },
    closeable: { type: Boolean, default: false },
});


interface FluxEvent {
    id: number;
    content: string;
    type: 'success' | 'info' | 'warning' | 'error';
}

interface Flux {
    events: FluxEvent[];
    success: (content: string) => void;
    info: (content: string) => void;
    warning: (content: string) => void;
    error: (content: string) => void;
    show: (type: 'success' | 'info' | 'warning' | 'error', content: string) => void;
    add: (type: 'success' | 'info' | 'warning' | 'error', content: string) => void;
    remove: (event: FluxEvent) => void;
}

// 这段代码定义了一个响应式对象 flux，它用于管理和控制一组通知事件（例如成功、信息、警告、错误通知）。在这个对象中，定义了几个方法来添加、显示和移除这些通知。
// 假设这是一个通知管理系统，你可以通过调用 flux.success('Message') 来显示一个成功类型的通知，flux.info('Info message') 来显示一个信息类型的通知，通知会在页面上显示一段时间后自动消失。
const flux = reactive<Flux>({
    events: [],

    // success、info、warning、error：这些是 flux 对象的快捷方法，用于向 flux.events 数组中添加不同类型的通知事件。
    // flux.success('message') 会等同于 flux.add('success', 'message')
    success: (content: string) => flux.add('success', content),
    info: (content: string) => flux.add('info', content),
    warning: (content: string) => flux.add('warning', content),
    error: (content: string) => flux.add('error', content),

    /**
     * @param {'success' | 'info' | 'warning' | 'error'} type
     * @param {string} content
     */
    show: (type: 'success' | 'info' | 'warning' | 'error', content: string) => flux.add(type, content),

    // 通过 setTimeout 延迟 100 毫秒后，创建一个新的通知事件对象，并将其添加到 flux.events 中。每个通知事件有一个唯一的 id，生成方式是通过 Date.now() 来保证唯一性。
    add: (type: 'success' | 'info' | 'warning' | 'error', content: string) => {
        if (!props.queue)
            flux.events = [];  // 这里的 props.queue 假设是一个外部的响应式属性，需确保 props 已被正确传递

        setTimeout(() => {
            const event: FluxEvent = { id: Date.now(), content, type };
            flux.events.push(event);
            setTimeout(() => flux.remove(event), props.timeout); // 假设 props.timeout 是预定义的超时时间
        }, 100);
    },

    // 这个方法接受一个通知事件对象 event，然后通过过滤掉该事件的 id，从 flux.events 数组中移除该事件。
    remove: (event: FluxEvent) => {
        flux.events = flux.events.filter((e) => e.id !== event.id);
    },
});

// defineExpose 将 flux 对象的 show、success、info、warning 和 error 方法暴露给外部组件。
defineExpose({
    show: flux.show,
    success: flux.success,
    info: flux.info,
    warning: flux.warning,
    error: flux.error,
})
</script>

<style scoped></style>