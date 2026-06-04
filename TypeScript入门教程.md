# TypeScript 入门教程

## 目录

1. [TypeScript 简介](#一typescript-简介)
2. [基础类型](#二基础类型)
3. [接口 (Interface)](#三接口-interface)
4. [类型别名 (Type Alias)](#四类型别名-type-alias)
5. [泛型 (Generics)](#五泛型-generics)
6. [函数类型](#六函数类型)
7. [Class 类](#七class-类)
8. [Vue 3 中的 TypeScript](#八vue-3-中的-typescript)
9. [模块导入导出](#九模块导入导出)
10. [高级类型](#十高级类型)
11. [实战案例](#十一实战案例)
12. [常见错误与解决](#十二常见错误与解决)

---

## 一、TypeScript 简介

### 1.1 什么是 TypeScript？

TypeScript 是 JavaScript 的超集，它添加了**类型系统**和**面向对象特性**。

```
JavaScript 代码 ────────────▶ TypeScript 代码
    .js                    .ts
                            │
                            ▼
                         编译/转译
                            │
                            ▼
                      JavaScript 代码
                            .js
```

### 1.2 为什么需要 TypeScript？

| 特性 | JavaScript | TypeScript |
|------|------------|------------|
| 类型检查 | 运行时 | **编译时** |
| 错误发现 | 执行后 | **编写时** |
| 代码提示 | 一般 | **智能** |
| 重构支持 | 困难 | **容易** |
| 学习曲线 | 低 | 中等 |

### 1.3 在项目中启用 TypeScript

在 `.vue` 文件中使用：

```vue
<script lang="ts" setup>
  // 这里使用 TypeScript 语法
</script>
```

---

## 二、基础类型

### 2.1 基础类型一览

```typescript
// 字符串
let name: string = "张三";
let greeting: string = `你好，${name}`;  // 模板字符串

// 数字
let age: number = 25;
let price: number = 99.99;
let hex: number = 0xff;        // 十六进制
let binary: number = 0b1010;  // 二进制
let octal: number = 0o777;    // 八进制

// 布尔值
let isActive: boolean = true;
let isLoggedIn: boolean = false;

// 未定义
let u: undefined = undefined;

// 空值
let n: null = null;

// 任意类型（关闭类型检查）
let anything: any = "可以是任何东西";
anything = 123;
anything = true;

// 不知道什么类型（比 any 更严格）
let unknownValue: unknown = "某个人输入的值";
// unknownValue.toUpperCase();  // 错误！需要先检查类型
if (typeof unknownValue === "string") {
  console.log(unknownValue.toUpperCase());  // 正确
}

// 无返回值
function sayHello(): void {
  console.log("Hello!");
}

// 永不返回（函数永远执行不完）
function infiniteLoop(): never {
  while (true) {
    // 永远执行
  }
}

// 永远不会到达（代码永远不会被执行到）
function throwError(): never {
  throw new Error("出错了！");
}
```

### 2.2 项目中的实际用法

来自 `utils.ts`：

```typescript
// 函数参数类型注解
export function getUrlParameter(name: string) {
  //                ^^^^^^^^^^^^^
  //                参数必须是字符串
}

// 返回值类型注解
export function getDateStr(timestamp: number): string {
  //                          ^^^^^^^^           ^^^^^^
  //                          参数类型           返回值类型
}

export function deleteCookie(name: string): void {
  //                                       ^^^^^^
  //                                       返回空值
}
```

### 2.3 类型推断

TypeScript 会**自动推断**变量类型：

```typescript
// TypeScript 会推断 message 是 string 类型
let message = "你好";  // let message: string

// 如果不指定类型，TypeScript 会根据初始值推断
let count = 0;        // let count: number
let isReady = true;   // let isReady: boolean

// 数组
let numbers = [1, 2, 3];        // let numbers: number[]
let names = ["张三", "李四"];   // let names: string[]
```

### 2.4 数组类型

```typescript
// 两种写法都可以
let arr1: number[] = [1, 2, 3];
let arr2: Array<number> = [1, 2, 3];

// 字符串数组
let fruits: string[] = ["苹果", "香蕉", "橙子"];

// 任意类型数组
let mixed: any[] = [1, "hello", true];

// 元组（固定长度和类型的数组）
let tuple: [string, number] = ["张三", 25];
// tuple = [25, "张三"];  // 错误！顺序不对
```

---

## 三、接口 (Interface)

### 3.1 什么是接口？

接口用于定义**对象的结构**，就像一份"合同"或"蓝图"。

```typescript
// 定义一个用户接口
interface User {
  // 属性: 类型
  name: string;      // 必须有
  age: number;       // 必须有
  email?: string;    // 可选属性（问号表示可有可无）
  readonly id: number;  // 只读属性（创建后不能修改）
}
```

### 3.2 项目中的接口定义

来自 `api.ts`：

```typescript
/**
 * HttpParams - HTTP 请求参数接口
 *
 * 定义所有 HTTP 请求可用的配置选项。
 */
export interface HttpParams {
  // url - 请求的 URL 地址（必需）
  url: string;

  // data - 请求数据（可选）
  data?: any;

  // method - 请求方法（可选）
  method?: string;

  // headers - 自定义请求头（可选）
  headers?: any;

  // onDownloadProgress - 下载进度回调（可选）
  onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void;

  // onUploadProgress - 上传进度回调（可选）
  onUploadProgress?: (progressEvent: AxiosProgressEvent) => void;

  // signal - 请求取消信号（可选）
  signal?: GenericAbortSignal;

  // beforeRequest - 请求前回调（可选）
  beforeRequest?: () => void;

  // afterRequest - 请求后回调（可选）
  afterRequest?: () => void;
}
```

### 3.3 可选属性

用 `?` 标记可选属性：

```typescript
interface Config {
  host: string;      // 必须
  port?: number;     // 可选
  debug?: boolean;   // 可选
}

let config: Config = {
  host: "localhost"  // 只提供必需的
};
// 或者
let config2: Config = {
  host: "localhost",
  port: 8080,
  debug: true
};
```

### 3.4 只读属性

用 `readonly` 标记只读属性：

```typescript
interface Point {
  readonly x: number;  // 创建后不能修改
  readonly y: number;
}

let point: Point = { x: 10, y: 20 };
point.x = 30;  // 错误！x 是只读的
```

### 3.5 接口继承

用 `extends` 继承其他接口：

```typescript
// 基础接口
interface Animal {
  name: string;
}

// 继承 Animal
interface Dog extends Animal {
  breed: string;
}

let dog: Dog = {
  name: "旺财",    // 来自 Animal
  breed: "金毛"    // 自己定义的
};
```

### 3.6 项目中的接口继承

```typescript
// 基础接口
export interface sysParams {
  sys: string;
}

// 继承 sysParams，添加更多属性
export interface officialCallbackParams extends sysParams {
  ticket: string;  // 额外的属性
}
```

### 3.7 函数类型接口

```typescript
// 定义函数类型
interface SearchFunc {
  (source: string, subString: string): boolean;
}

let mySearch: SearchFunc;

// 实现函数
mySearch = function(source: string, subString: string) {
  let result = source.search(subString);
  return result > -1;
};
```

### 3.8 接口与类型的区别

```typescript
// 接口
interface Person {
  name: string;
  age: number;
}

// 类型别名
type Person = {
  name: string;
  age: number;
};

// 主要区别：
// 1. 接口可以声明合并（同名接口会自动合并）
// 2. 类型别名可以定义联合类型、交叉类型
// 大多数情况下可以互换使用
```

---

## 四、类型别名 (Type Alias)

### 4.1 基本语法

```typescript
// 用 type 关键字定义类型别名
type Name = string;
type NameResolver = () => string;
type MaybeNull<T> = T | null;
```

### 4.2 项目中的用法

来自 `request.ts`：

```typescript
// 定义响应类型
export interface Response<T = any> {
  // data 的类型由泛型 T 决定
  data?: T;

  // status 是 number 类型
  status: number;

  // message 是可选的 string
  message?: string;
}
```

### 4.3 联合类型

一个值可以是多种类型之一：

```typescript
// string 或 number
type StringOrNumber = string | number;

let value: StringOrNumber;
value = "hello";  // 正确
value = 123;      // 正确
// value = true;   // 错误

// 多种类型
type Result = "success" | "error" | "loading";
let status: Result;
status = "success";  // 正确
status = "pending";  // 错误！"pending" 不是有效值
```

### 4.4 交叉类型

将多个类型合并成一个：

```typescript
interface Part1 {
  name: string;
}

interface Part2 {
  age: number;
}

// 合并两个接口
type Person = Part1 & Part2;

let person: Person = {
  name: "张三",  // 来自 Part1
  age: 25        // 来自 Part2
};
```

---

## 五、泛型 (Generics)

### 5.1 什么是泛型？

泛型让我们能够创建**可复用**的组件，同时保持**类型安全**。

```typescript
// 没有泛型：类型不明确
function identity(arg: any): any {
  return arg;
}

// 使用泛型：类型被保留
function identity<T>(arg: T): T {
  return arg;
}

// 调用时指定类型
let output = identity<string>("hello");  // output 是 string 类型
let output2 = identity(123);             // TypeScript 自动推断是 number
```

### 5.2 项目中的泛型用法

来自 `api.ts`：

```typescript
/**
 * uploadFile - 上传文件
 *
 * @template T - 返回数据的类型
 * @param params - 请求参数
 * @returns Promise 包含类型为 T 的数据
 */
export function uploadFile<T = any>(params: { formData: FormData }) {
  const path = "/v1/file/upload";
  // 返回 Promise<Response<T>>
  return post<T>({ url: path, data: params.formData });
}

// 使用示例：
interface FileResult {
  url: string;
}

// 调用时指定类型
uploadFile<FileResult>({ formData }).then(res => {
  // res.data 是 FileResult 类型
  console.log(res.data?.url);
});
```

来自 `request.ts`：

```typescript
/**
 * get - GET 请求
 *
 * @template T - 响应数据的类型
 * @param params - 请求参数
 * @returns Promise<Response<T>> - 泛型 Promise
 */
export function get<T = any>({ url, data, method = "GET", ... }: HttpParams): Promise<Response<T>> {
  return http<T>({ url, method, data, ... });
}

// 使用示例：
interface UserData {
  id: number;
  name: string;
}

get<UserData>({ url: "/user" }).then(res => {
  // res.data 是 UserData 类型
  const user: UserData = res.data!;
  console.log(user.name);
});
```

### 5.3 泛型约束

限制泛型的范围：

```typescript
// 约束 T 必须有 length 属性
interface Lengthwise {
  length: number;
}

function loggingIdentity<T extends Lengthwise>(arg: T): T {
  console.log(arg.length);  // 现在可以访问 length 了
  return arg;
}

loggingIdentity("hello");     // 正确，string 有 length
loggingIdentity([1, 2, 3]);  // 正确，数组有 length
loggingIdentity(123);        // 错误！number 没有 length
```

### 5.4 泛型类

```typescript
class GenericNumber<T> {
  zeroValue: T;
  add: (x: T, y: T) => T;
}

let myGenericNumber = new GenericNumber<number>();
myGenericNumber.zeroValue = 0;
myGenericNumber.add = (x, y) => x + y;
```

### 5.5 泛型接口

```typescript
interface Container<T> {
  value: T;
  getValue(): T;
}

let container: Container<string> = {
  value: "Hello",
  getValue() {
    return this.value;
  }
};
```

### 5.6 多类型参数

```typescript
// 多个泛型参数
function pair<T, U>(first: T, second: U): [T, U] {
  return [first, second];
}

let p = pair("hello", 123);  // [string, number]
```

---

## 六、函数类型

### 6.1 函数参数和返回值

```typescript
// 参数类型注解
function greet(name: string): string {
  return "Hello, " + name;
}

// 无参数无返回值
function sayGoodbye(): void {
  console.log("Goodbye!");
}

// 可选参数（用问号）
function hello(name?: string): void {
  if (name) {
    console.log("Hello, " + name);
  } else {
    console.log("Hello!");
  }
}
```

### 6.2 箭头函数类型

```typescript
// 普通函数
function add(a: number, b: number): number {
  return a + b;
}

// 箭头函数
const add2 = (a: number, b: number): number => a + b;

// 箭头函数的类型标注
const add3: (a: number, b: number) => number = (a, b) => a + b;
```

### 6.3 项目中的函数类型

来自 `request.ts`：

```typescript
// 进度回调函数类型
onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void;
//                 ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//                 箭头函数类型：接收参数，返回 void

// 请求前后的回调
beforeRequest?: () => void;
//           ^^^^^^^^^^^^^^
//           无参数，返回 void
```

### 6.4 函数重载

同一个函数名，不同的参数类型：

```typescript
function reverse(x: number): number;      // 如果输入数字，返回数字
function reverse(x: string): string;     // 如果输入字符串，返回字符串
function reverse(x: number | string): number | string {
  // 函数实现
  if (typeof x === "number") {
    return Number(x.toString().split("").reverse().join(""));
  } else {
    return x.split("").reverse().join("");
  }
}

reverse(123);      // 调用第一个重载
reverse("hello");   // 调用第二个重载
```

---

## 七、Class 类

### 7.1 基本语法

```typescript
class Animal {
  // 属性
  name: string;
  age: number;

  // 构造函数
  constructor(name: string, age: number) {
    this.name = name;
    this.age = age;
  }

  // 方法
  speak(): void {
    console.log(`${this.name} 在叫`);
  }
}

// 创建实例
let dog = new Animal("旺财", 3);
dog.speak();  // 输出：旺财 在叫
```

### 7.2 项目中的 Class 用法

来自 `home.vue`：

```typescript
/**
 * homeRes - 首页响应数据类型
 *
 * 使用 ES6 Class 定义，匹配后端返回的数据结构
 */
class HomeRes {
  // 轮播图 URL 数组
  banners: Array<string>;

  // 第一行图片 URL 数组
  images1: Array<string>;

  // 第二行图片 URL 数组
  images2: Array<string>;

  // 构造函数
  constructor(banners: Array<string>, images1: Array<string>, images2: Array<string>) {
    this.banners = banners;
    this.images1 = images1;
    this.images2 = images2;
  }
}

// 使用
let homeData = new HomeRes(
  ["url1", "url2", "url3"],
  ["img1", "img2", "img3"],
  ["pic1", "pic2", "pic3"]
);
```

来自 `upload.vue`：

```typescript
/**
 * fileUploadRes - 文件上传响应数据类型
 */
class FileUploadRes {
  // 上传成功后的文件 URL
  url: string;

  constructor(url: string) {
    this.url = url;
  }
}

// 使用
let result = new FileUploadRes("https://example.com/file.jpg");
console.log(result.url);  // https://example.com/file.jpg
```

### 7.3 访问修饰符

| 修饰符 | 说明 | 访问范围 |
|--------|------|---------|
| `public` | 公开 | 任何地方都可以访问（默认） |
| `private` | 私有 | 只能在本类内部访问 |
| `protected` | 受保护 | 只能在本类及子类中访问 |

```typescript
class Person {
  public name: string;       // 公开
  private age: number;      // 私有
  protected id: number;      // 受保护

  constructor(name: string, age: number, id: number) {
    this.name = name;
    this.age = age;
    this.id = id;
  }
}

let person = new Person("张三", 25, 1);
console.log(person.name);    // 正确：公开属性
// console.log(person.age);  // 错误：私有属性
// console.log(person.id);   // 错误：受保护属性
```

### 7.4 readonly 属性

```typescript
class User {
  readonly id: number;  // 只读属性
  name: string;

  constructor(id: number, name: string) {
    this.id = id;
    this.name = name;
  }
}

let user = new User(1, "张三");
user.id = 2;  // 错误！id 是只读的
```

### 7.5 类的继承

```typescript
class Animal {
  constructor(public name: string) {}

  speak(): void {
    console.log(this.name + " 在叫");
  }
}

class Dog extends Animal {
  breed: string;

  constructor(name: string, breed: string) {
    super(name);  // 调用父类构造函数
    this.breed = breed;
  }

  speak(): void {
    console.log(this.name + " 汪汪叫！");
  }
}

let dog = new Dog("旺财", "金毛");
dog.speak();  // 旺财 汪汪叫！
```

---

## 八、Vue 3 中的 TypeScript

### 8.1 defineProps 定义属性

`defineProps` 用于定义组件的输入属性：

```typescript
// 方式 1：使用类型推断（推荐）
defineProps<{
  msg: string;           // 必需的属性
  count?: number;        // 可选属性
}>()

// 使用
// <MyComponent msg="Hello" />
// props.msg 是 string 类型
// props.count 是 number | undefined 类型
```

来自 `HelloWorld.vue`：

```typescript
// 定义 msg 属性，类型为 string
defineProps<{ msg: string }>()
```

### 8.2 ref 创建响应式数据

`ref` 用于创建响应式的基本类型数据：

```typescript
import { ref } from 'vue'

// ref(初始值) 创建响应式变量
let count = ref(0);           // number 类型
let message = ref("hello");   // string 类型
let isActive = ref(true);     // boolean 类型

// 访问值：使用 .value
count.value++;                // 修改值
console.log(count.value);     // 读取值

// 在模板中自动解包，不需要 .value
// {{ count }} 会显示 1
```

### 8.3 项目中的 ref 用法

```typescript
import { ref } from 'vue'

// 创建响应式数据
let data = {
  loginMethods: {
    github: ref(""),              // ref("")
    gitlab: ref(""),
    wx_qrcode: {
      expire_seconds: ref(0),      // ref(0)
      ticket: ref(""),
      qr_code_url: ref(""),
    }
  }
};

// 修改值
data.loginMethods.github.value = "https://github.com/login...";

// 读取值
console.log(data.loginMethods.github.value);
```

### 8.4 泛型 ref

```typescript
// 明确指定类型
let count = ref<number | null>(null);  // count.value 是 number | null
let name = ref<string>("张三");          // name.value 是 string

// 数组
let users = ref<Array<{ id: number; name: string }>>([]);

// 或者使用尖括号语法
let items = ref<{ id: number; title: string }[]>([]);
```

### 8.5 生命周期钩子

```typescript
import { onBeforeMount, onMounted, onBeforeUnmount } from 'vue'

// 组件挂载前
onBeforeMount(() => {
  console.log("组件即将挂载");
  // 调用 API 获取数据
});

// 组件挂载后
onMounted(() => {
  console.log("组件已挂载");
  // DOM 操作
});

// 组件卸载前
onBeforeUnmount(() => {
  console.log("组件即将卸载");
  // 清理工作，如取消定时器
});
```

### 8.6 项目中的生命周期钩子

来自 `login.vue`：

```typescript
import { onBeforeMount, ref } from 'vue'

// onBeforeMount - 组件挂载前执行
onBeforeMount(() => {
  // 从 URL 获取参数
  const sys = getUrlParameter("sys") || "ai";

  // 调用 API 获取登录方式
  loginMethods({ sys: sys }).then(function (res) {
    // 更新响应式数据
    data.loginMethods.github.value = res.data.github;
    data.loginMethods.gitlab.value = res.data.gitlab;
    // ...
  });
});
```

### 8.7 computed 计算属性

```typescript
import { ref, computed } from 'vue'

let firstName = ref("张");
let lastName = ref("三");

// 计算属性：基于响应式数据计算得出
let fullName = computed(() => {
  return firstName.value + lastName.value;
});

console.log(fullName.value);  // 张三

// 修改依赖项
firstName.value = "李";
console.log(fullName.value);  // 李三
```

### 8.8 watch 监听器

```typescript
import { ref, watch } from 'vue'

let count = ref(0);

// 监听 count 的变化
watch(count, (newValue, oldValue) => {
  console.log(`count 从 ${oldValue} 变成了 ${newValue}`);
});

count.value++;  // 输出：count 从 0 变成了 1
```

---

## 九、模块导入导出

### 9.1 导出

```typescript
// 导出函数
export function add(a: number, b: number): number {
  return a + b;
}

// 导出常量
export const PI = 3.14159;

// 导出类型/接口
export interface User {
  name: string;
  age: number;
}

// 导出类
export class Animal {
  speak(): void {
    console.log("...");
  }
}

// 默认导出（每个文件只能有一个）
export default class Animal {
  // ...
}
```

### 9.2 导入

```typescript
// 导入命名导出
import { add, PI, User } from './math';
import { add as sum } from './math';  // 重命名

// 导入默认导出
import Animal from './animal';

// 导入所有
import * as utils from './utils';

// 导入类型（编译时会被移除）
import type { User } from './types';
```

### 9.3 项目中的导入导出

```typescript
// utils.ts 导出
export function getUrlParameter(name: string) { ... }
export function getCookie(key: string) { ... }
export function getDateStr(timestamp: number) { ... }
export function setCookie(name: string, value: string, days: number) { ... }
export function deleteCookie(name: string) { ... }

// 其他文件导入
import { getUrlParameter, getCookie } from '../utils/utils.ts';
```

```typescript
// request.ts 导出
export interface HttpParams { ... }
export interface Response<T = any> { ... }
export function get<T = any>(...) { ... }
export function post<T = any>(...) { ... }
export default post;

// 其他文件导入
import { get, post } from '../request/request';
```

---

## 十、高级类型

### 10.1 可选链 (?.)

安全访问嵌套属性：

```typescript
interface Person {
  address?: {
    city: string;
  };
}

let person: Person = {};

// 传统写法：需要判断每层
if (person && person.address) {
  console.log(person.address.city);
}

// 可选链：更简洁
console.log(person.address?.city);  // undefined，不会报错

// 可选调用
let result = someFunction?.();  // 如果 someFunction 存在才调用
```

### 10.2 空值合并 (??)

```typescript
// ?? 只有在值为 null 或 undefined 时才使用默认值
let value = null;
let result = value ?? "默认值";  // "默认值"

value = "hello";
result = value ?? "默认值";  // "hello"

// vs || ：|| 将 0、""、false 也视为假值
let count = 0;
let a = count || 1;   // 1（因为 0 是假值）
let b = count ?? 1;   // 0（因为 0 不是 null/undefined）
```

### 10.3 类型守卫

```typescript
// typeof 类型守卫
function padLeft(value: string | number, padding: string | number) {
  if (typeof padding === "number") {
    // TypeScript 知道 padding 是 number
    return Array(padding + 1).join(" ") + value;
  }
  // TypeScript 知道 padding 是 string
  return padding + value;
}

// instanceof 类型守卫
class Dog {
  bark(): void {
    console.log("汪汪");
  }
}

class Cat {
  meow(): void {
    console.log("喵喵");
  }
}

function speak(animal: Dog | Cat) {
  if (animal instanceof Dog) {
    animal.bark();  // TypeScript 知道是 Dog
  } else {
    animal.meow();  // TypeScript 知道是 Cat
  }
}

// 自定义类型守卫
interface Fish {
  swim(): void;
}

interface Bird {
  fly(): void;
}

function isFish(pet: Fish | Bird): pet is Fish {
  return (pet as Fish).swim !== undefined;
}
```

### 10.4 类型断言

```typescript
// 将一种类型当作另一种类型
let someValue: any = "hello";

// 方式 1：尖括号语法
let strLength1: number = (<string>someValue).length;

// 方式 2：as 语法（推荐）
let strLength2: number = (someValue as string).length;

// 在 DOM 操作中常用
let input = document.getElementById("myInput") as HTMLInputElement;
input.value = "Hello";  // TypeScript 知道 input 是 input 元素
```

### 10.5 keyof 操作符

获取对象类型的所有键：

```typescript
interface Person {
  name: string;
  age: number;
  address: string;
}

// keyof Person 等于 "name" | "age" | "address"
type PersonKeys = keyof Person;

function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key];
}

let person: Person = { name: "张三", age: 25, address: "北京" };

let name = getProperty(person, "name");     // string
let age = getProperty(person, "age");      // number
// getProperty(person, "invalid");          // 错误！
```

### 10.6 映射类型

从现有类型创建新类型：

```typescript
// 所有属性变为可选
type Partial<T> = {
  [P in keyof T]?: T[P];
};

// 所有属性变为必需
type Required<T> = {
  [P in keyof T]-?: T[P];
};

// 所有属性变为只读
type Readonly<T> = {
  readonly [P in keyof T]: T[P];
};

// 示例
interface User {
  name: string;
  age: number;
}

type PartialUser = Partial<User>;
// { name?: string; age?: number; }

type ReadonlyUser = Readonly<User>;
// { readonly name: string; readonly age: number; }
```

---

## 十一、实战案例

### 11.1 完整 API 接口定义

来自 `request.ts`：

```typescript
// 1. 定义请求参数接口
export interface HttpParams {
  // url - 必需，字符串
  url: string;

  // data - 可选，任意类型
  data?: any;

  // method - 可选，字符串
  method?: string;

  // headers - 可选，任意类型
  headers?: any;

  // onDownloadProgress - 可选，函数类型
  onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void;

  // onUploadProgress - 可选，函数类型
  onUploadProgress?: (progressEvent: AxiosProgressEvent) => void;

  // signal - 可选，GenericAbortSignal 类型
  signal?: GenericAbortSignal;

  // beforeRequest - 可选，无参无返回值的函数
  beforeRequest?: () => void;

  // afterRequest - 可选，无参无返回值的函数
  afterRequest?: () => void;
}

// 2. 定义响应接口，包含泛型
export interface Response<T = any> {
  // data - 可选，类型由泛型 T 决定
  data?: T;

  // status - 必需，数字
  status: number;

  // message - 可选，字符串
  message?: string;
}

// 3. 定义泛型函数
export function get<T = any>(
  { url, data, method = "GET", headers, onDownloadProgress, signal, beforeRequest, afterRequest }: HttpParams
): Promise<Response<T>> {
  // 实现
  return http<T>({ url, method, data, headers, onDownloadProgress, signal, beforeRequest, afterRequest });
}

// 4. 使用泛型函数
interface User {
  id: number;
  name: string;
}

get<User>({ url: "/api/user" }).then(res => {
  // res.data 是 User | undefined
  console.log(res.data?.name);
});
```

### 11.2 完整组件类型定义

```typescript
// 1. 导入依赖
import { ref, onBeforeMount } from 'vue';
import { get } from '../api/api';

// 2. 定义响应式数据
let userInfo = ref<{
  name: string;
  avatar: string;
  user_id: number;
}>({
  name: "",
  user_id: 0,
  avatar: "",
});

// 3. 定义类
class HomeRes {
  banners: Array<string>;
  images1: Array<string>;
  images2: Array<string>;

  constructor(banners: Array<string>, images1: Array<string>, images2: Array<string>) {
    this.banners = banners;
    this.images1 = images1;
    this.images2 = images2;
  }
}

// 4. 定义函数
function getStyle(index: number): string {
  let style = "position:absolute;width:14rem;height:10rem;";
  style += " left:" + 15.25 * index + 'rem;';
  return style;
}

// 5. 使用生命周期钩子
onBeforeMount(() => {
  get<HomeRes>({ url: "/v1/home" }).then(function (res) {
    // 处理响应
    console.log(res.data);
  }).catch(function (err) {
    // 处理错误
    console.error(err);
  });
});

// 6. 定义 Props
defineProps<{
  msg: string;
  count?: number;
}>();
```

### 11.3 JWT Token 解析

```typescript
// JWT Token 格式：header.payload.signature
function parseJWT(token: string) {
  // 1. 用点号分割
  const parts = token.split('.');

  if (parts.length !== 3) {
    throw new Error('Invalid JWT token');
  }

  // 2. Base64 解码 payload
  // atob() 是浏览器内置的 Base64 解码函数
  const payload = atob(parts[1]);

  // 3. JSON 解析
  const decoded = JSON.parse(payload);

  return decoded;
}

// 使用
interface TokenPayload {
  user_id: number;
  name: string;
  avatar: string;
  exp: number;
}

let payload = parseJWT(accessToken) as TokenPayload;
console.log(payload.user_id);
console.log(payload.name);
```

---

## 十二、常见错误与解决

### 12.1 类型错误

```typescript
// 错误：不能将类型 "abc" 分配给类型 number
let num: number = "123";  // 错误！

// 解决：类型转换或正确赋值
let num: number = 123;        // 直接赋值数字
let num2: number = parseInt("123");  // 转换
let num3: number = Number("123");     // 转换
```

### 12.2 可能为空

```typescript
// 错误：对象可能为 "null"
let name: string | null = getName();
console.log(name.toUpperCase());  // 错误！

// 解决1：类型守卫
if (name !== null) {
  console.log(name.toUpperCase());  // 正确
}

// 解决2：可选链
console.log(name?.toUpperCase());  // 如果 name 是 null，返回 undefined

// 解决3：空值合并
console.log((name ?? "").toUpperCase());  // null 时使用空字符串
```

### 12.3 属性不存在

```typescript
// 错误：类型 "object" 上不存在属性 "name"
let obj: object = { name: "张三" };
console.log(obj.name);  // 错误！

// 解决：使用类型断言
console.log((obj as { name: string }).name);

// 或者定义接口
interface Person {
  name: string;
}
let person: Person = { name: "张三" };
console.log(person.name);  // 正确
```

### 12.4 函数类型不匹配

```typescript
// 错误：类型不匹配
type Callback = (data: string) => void;
let fn: Callback = (n: number) => {};  // 错误！参数类型不对

// 解决：参数类型匹配
let fn2: Callback = (n: string) => {};
```

### 12.5 泛型推断失败

```typescript
// 错误：TypeScript 无法推断类型
function createArray<T>(value: T) {
  return [value];
}

let arr = createArray(123);  // number[] - 自动推断

// 手动指定类型
let arr2 = createArray<number>(123);
```

---

## 附录：TypeScript 配置

### tsconfig.json 常用配置

```json
{
  "compilerOptions": {
    // 编译目标版本
    "target": "ES2020",

    // 模块系统
    "module": "ESNext",

    // 是否允许编译 JS 文件
    "allowJs": true,

    // 是否检查 JS 文件
    "checkJs": false,

    // 输出目录
    "outDir": "./dist",

    // 严格模式（建议开启）
    "strict": true,

    // 严格空值检查
    "strictNullChecks": true,

    // 隐式 any 检查
    "noImplicitAny": true,

    // 装饰器支持
    "experimentalDecorators": true,

    // 模块解析方式
    "moduleResolution": "node",

    // 跳过库检查
    "skipLibCheck": true
  },

  // 包含的文件
  "include": ["src/**/*"],

  // 排除的文件
  "exclude": ["node_modules", "dist"]
}
```

---

## 学习资源

| 资源 | 链接 |
|------|------|
| TypeScript 官方文档 | https://www.typescriptlang.org/docs/ |
| TypeScript 中文文档 | https://www.tslang.cn/docs/handbook/basic-types.html |
| Vue 3 + TypeScript | https://vuejs.org/guide/typescript/overview.html |
| TypeScript Playground | https://www.typescriptlang.org/play |

---

**文档版本**: 1.0
**最后更新**: 2026-06-04
**作者**: AI Assistant
