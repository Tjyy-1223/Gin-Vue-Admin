const CryptoSecret = '__SecretKey__'  // 用于加密的密钥

/**
 * 存储序列化后的数据到 LocalStorage
 * @param {string} key - 存储的键
 * @param {any} value - 存储的值，可以是任何类型的对象，会被序列化为 JSON
 * @param {number} expire - 数据的过期时间，单位为秒，默认 7 天
 */
export function setLocal(key, value, expire = 60 * 60 * 24 * 7) {
  // 创建一个包含 value、存储时间和过期时间的对象
  const data = JSON.stringify({
    value,  // 存储的数据
    time: Date.now(),  // 当前时间戳，表示存储的时间
    expire: expire ? new Date().getTime() + expire * 1000 : null,  // 计算过期时间，默认不设置过期时间
  })

  // 将序列化后的数据加密后存储到 localStorage 中
  window.localStorage.setItem(key, encrypto(data))  // 使用 encrypto 函数加密数据
}

/**
 * 从 LocalStorage 中获取数据，解密后反序列化，并检查是否过期
 * @param {string} key - 存储的键
 * @returns {any} - 解密并且未过期的值，如果过期则返回 null
 */
export function getLocal(key) {
  // 获取加密后的数据
  const encryptedVal = window.localStorage.getItem(key)
  
  // 如果数据存在，则进行解密
  if (encryptedVal) {
    const val = decrypto(encryptedVal)  // 解密数据
    const { value, expire } = JSON.parse(val)  // 反序列化解密后的数据

    // 检查是否已过期，如果未过期则返回存储的值
    if (!expire || expire > new Date().getTime()) {
      return value
    }
  }

  // 如果数据已过期，则移除对应的 localStorage 项
  removeLocal(key)
  return null
}

/**
 * 从 LocalStorage 中移除指定的键值对
 * @param {string} key - 存储的键
 */
export function removeLocal(key) {
  window.localStorage.removeItem(key)  // 使用 localStorage 的 removeItem 方法移除项
}

/**
 * 清空所有存储在 LocalStorage 中的数据
 */
export function clearLocal() {
  window.localStorage.clear()  // 使用 localStorage 的 clear 方法清空所有数据
}

/**
 * 加密数据: 使用 Base64 加密
 * @param {any} data - 需要加密的数据
 * @returns {string} - 加密后的字符串
 */
function encrypto(data) {
  // 将数据转为 JSON 字符串
  const newData = JSON.stringify(data)

  // 使用 Base64 编码并在加密数据前加上一个密钥（防止直接暴力解密）
  const encryptedData = btoa(CryptoSecret + newData)  // `btoa` 是浏览器内置的 Base64 编码方法
  return encryptedData
}

/**
 * 解密数据: 使用 Base64 解密
 * @param {string} cipherText - 加密后的密文
 * @returns {any} - 解密后的原始数据，如果解密失败则返回 null
 */
function decrypto(cipherText) {
  // 使用 Base64 解码
  const decryptedData = atob(cipherText)  // `atob` 是浏览器内置的 Base64 解码方法

  // 移除加密时添加的密钥
  const originalText = decryptedData.replace(CryptoSecret, '')  // 将密钥从解密的文本中去除

  // 尝试将解密后的文本解析为 JSON 对象
  try {
    const parsedData = JSON.parse(originalText)  // 解析 JSON 数据
    return parsedData  // 返回解密后的数据
  }
  catch (error) {
    return null  // 如果解密过程中出现错误，则返回 null
  }
}
