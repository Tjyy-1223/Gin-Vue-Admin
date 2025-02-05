/** @type {import('tailwindcss').Config} */

const { iconsPlugin, getIconCollections } = require('@egoist/tailwindcss-icons');
export default {
  content: ["index.html", "./src/**/*.{html,js,ts,jsx,tsx,vue}"],
  theme: {
    extend: {
      transitionDuration: {
        '500': '500ms',
      },
    },
  },
  plugins: [
    iconsPlugin({
      collections: getIconCollections(['mdi', 'lucide']), // 或使用 "all" 来使用全部图标
    }),
  ],
}
