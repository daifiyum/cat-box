## sing-box 简易托盘应用

***如果你只是想使用核心，但又希望能方便一点，那还等什么，赶快下载吧！***

1. 网页订阅管理，支持添加多个订阅，切换订阅，自动更新订阅

2. 基于singbox核心运行，可自由升级核心

3. 基于模板生成配置，想怎么改就怎么改

4. 采用go开发，体积小，占用少

5. 无需安装，不乱拉屎，所有数据都放在程序当前目录，绝对绿色运行

### 快速上手

![](https://thumbsnap.com/i/hEoHi8nY.png)

### 模板

模板路径：resources\template\template.json

相比于图形化修改配置，直接修改模板配置，会更自由灵活（图形化复杂，写不出来，直接对模板进行修改不更好，我爱说实话）

### 订阅

自动更新：

支持以m（分钟）或h（小时）为单位自定义更新延时

例如：1m，32m，1h，23h，49h

订阅代理：

当开启代理时，订阅会自动使用代理，未开启则直连

### 图标

路径：resources\icons

可自定义，但图片格式必须是.ico

### 开机自启

手动修改注册表方式实现

注册表路径：HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run

进入 Run 项目，右键新建“字符串值”，右键新建的数值，点击“修改”：

数值名称：cat-box

数值数据："cat-box.exe路径" --enable-workspace

数值数据例子："E:\cat-box\cat-box.exe" --enable-workspace

参数说明：--enable-workspace 参数表示开启工作目录，否则不可以通过绝对路径运行

### 其他

windows消息通知使用了rust基于winrt编译的dll库，目前go实现的都是基于powershell脚本，不够优雅，狗都不用
