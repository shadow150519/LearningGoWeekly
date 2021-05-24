## viper笔记(艹，变成官方文档的翻译了QAQ)
### 什么是viper
viper是go用来读写配置文件的第三方库,它支持如下特性:

+ setting defaults(给属性设置默认值)
+ reading from JSON, TOML, YAML, HCL, envfile and Java properties config files(支持JSON, TOML, YAML, HCL, envfile and Java properties等配置文件格式)
+ live watching and re-reading of config files (optional)(支持在运行时监视配置文件更改，并重新读取配置，可选的)
+ reading from environment variables (从环境变量中读取)
+ reading from remote config systems (etcd or Consul), and watching changes (支持从远程服务器读取配置)
+ reading from command line flags (支持从命令行读取参数)
+ reading from buffer (从buffer读取配置)
+ setting explicit values (显式设置值)
  
可以将Viper视为满足您所有应用程序配置需求的注册表

### viper参数优先级
Viper参数遵循一下的优先级:

explicit call to Set (显示赋值)
flag (使用pflag ??? 好像是cobra里面的什么东西，不了解)
env (环境变量)
config (配置文件)
key/value store (key/value存储 ??? 不是太懂)
default (默认值)

重要：viper configuration 的 key是**大小写不敏感**的，这个特性正在下一个版本中讨论是否要修改

### 把值放到Viper
#### 默认值
```go
viper.SetDefault("ContentDir", "content")
viper.SetDefault("LayoutDir", "layouts")
viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
```

#### 读配置文件
Viper 支持JSON,TOML,YANL,INI,envfile,Java Properties.Viper可以从多个路径中搜索配置，但是目前一个viper只支持一个configuration file,
(这并不与viper中可以添加多个文件路径冲突)
下面是使用viper从配置文件读取配置的一个例子
```go
viper.SetConfigName("config") // 配置文件的名字 (可以没有拓展名)
viper.SetConfigType("yaml") // 如果SetConfigName是没有提供拓展名，则需要提供该参数
viper.AddConfigPath("$HOME/.appname")  // 调用这个函数多次来添加多个路径
viper.AddConfigPath(".")               // 在当前目录寻找配置文件
err := viper.ReadInConfig() // 寻找并且读配置文档
if err != nil { // 处理读取配置文件时的error
    panic(fmt.Errorf("Fatal error config file: %s \n", err))
}
```

你可以处理特定的错误情况比如当没有找到配置文件,向下面的例子一样
```go
if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
        // Config file not found; ignore error if desired
    } else {
        // Config file was found but another error was produced
    }
}

// Config file found and successfully parsed
```
**注意**，从1.6版本开始，可以读一个不带拓展名的配置文件，并且在你的程序中手动确定文件格式(我不太明白如何在编程中确定文件的格式,原文如下)
NOTE [since 1.6]: You can also have a file without an extension and specify the format programmaticaly. For those 
configuration files that lie in the home of the user without any extension like .bashrc

#### 写入配置文件
从配置文件都信息是有用的，但是当你想保存所有在运行时产生的修改。一些函数被提供，每一个都有特定的用处

+ WriteConfig 将当前viper的配置写入预定义的路径，如果存在。会在没有预定义路径的时候产生error。如果当前已经存在了这个配置文件，会重写它，
  相当于truncate
  
+ SafeWriteConfig 将当前viper的配置写入预定义的路径。当没有预定义的路径时会产生错误。如果已经存在了配置文件，不会重写原有的文件

+ WriteConfigAs 将当前viper的配置写入给定的路径，如果存在。会在没有预定义路径的时候产生error。如果当前已经存在了这个配置文件，会重写它，
相当于truncate

+ SafeWriteConfigAs 将当前viper的配置写入给定的路径。当没有预定义的路径时会产生错误。如果已经存在了配置文件，不会重写原有的文件

一般说来，带Safe开头的都不会重写文件，只是当文件不存在时新建这个文件，不带Safe的函数会create 或是 truncate

下面是个简单的例子
```go
viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
viper.SafeWriteConfig()
viper.WriteConfigAs("/path/to/my/.config")
viper.SafeWriteConfigAs("/path/to/my/.config") // will error since it has already been written
viper.SafeWriteConfigAs("/path/to/my/.other_config")
```

#### 监听并且从配置文件中重新读取
Viper 支持在你的应用运行时动态读取配置
只需要告诉viper的实例WatchConfig即可。你可以提供一个函数给Viper，当配置文件发生变化是，都会调用这个函数，当然这个函数是可选的，
**但是这个函数目前有个bug，配置改动时这个函数会被调用两次**

**确保你已经添加了所有的 configPath 在你调用 WatchConfig()**

```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
	fmt.Println("Config file changed:", e.Name)
})
```

#### 从io.Reader 读取配置
Viper 提供了非常多的配置源，比如files, environment variables, flags, and remote K/V store，但是你不会被此限制，你可以实习自己的配置
方法并将它提供给viper
```go
viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

// any approach to require this configuration into your program.
var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

viper.ReadConfig(bytes.NewBuffer(yamlExample))

viper.Get("name") // this would be "steve"
```

#### 重写 Setting
可以通过命令行参数或是在程序内重写
```go
viper.Set("Verbose", true)
viper.Set("LogFile", LogFile)
```

#### 注册并使用别名
别名允许一个值可以被多个键调用
```go
viper.RegisterAlias("loud", "Verbose")

viper.Set("verbose", true) // same result as next line
viper.Set("loud", true)   // same result as prior line

viper.GetBool("loud") // true
viper.GetBool("verbose") // true
```


待补充
#### 环境变量
#### Flags参数
#### 远程的Key/Value存储支持

## 从viper中获取值
+ Get(key string) : interface{} 
+ GetBool(key string) : bool 
+ GetFloat64(key string) : float64
+ GetInt(key string) : int
+ GetIntSlice(key string) : []int
+ GetString(key string) : string
+ GetStringMap(key string) : map[string]interface{}
+ GetStringMapString(key string) : map[string]string
+ GetStringSlice(key string) : []string
+ GetTime(key string) : time.Time 
+ GetDuration(key string) : time.Duration
+ IsSet(key string) : bool
+ AllSettings() : map[string]interface{}

当一个GetXXX函数没有找到键对应的值是，会返回该类型的零值，可以使用IsSet()来判断这个key存不存在
```go
viper.GetString("logfile") // case-insensitive Setting & Getting
if viper.GetBool("verbose") {
    fmt.Println("verbose enabled")
}
```
### 访问嵌套的key
访问方法也接受一个规范的路径来访问嵌套的key.例如下面的JSON配置文件
```JSON
{
    "host": {
        "address": "localhost",
        "port": 5799
    },
    "datastore": {
        "metric": {
            "host": "127.0.0.1",
            "port": 3099
        },
        "warehouse": {
            "host": "198.0.0.1",
            "port": 2112
        }
    }
}
```

Viper可以访问一个嵌套的字段通过"."分隔符，就像你正常访问结构体内属性一样
```go
GetString("data.metric.host") // 返回 "127.0.0.1"
```
这个遵循上面列出的优先级原则,对该key的搜索会搜寻所有的注册配置(不仅仅是配置文件，是上面列出的所有配置可能存在的地方)，直到找到或是没找到返回零值

例如，给了一个配置文件里面有datastore.metric.host和datastore.metric.port(可能被重写了)，如果程序中通过default方法定义了
datastore.metric.protocol，那么这个值也是可以被找到的

然而，如果datastore.metric被重写了(通过flag,环境变量,set()方法等等)，那么所有的子字段变成未定义，他们被更高的配置等级的定义隐藏了

最后，如果存在一个key匹配了分隔符形式的路径，那么它会被返回
```go
{
    "datastore.metric.host": "0.0.0.0",
    "host": {
        "address": "localhost",
        "port": 5799
    },
    "datastore": {
        "metric": {
            "host": "127.0.0.1",
            "port": 3099
        },
        "warehouse": {
            "host": "198.0.0.1",
            "port": 2112
        }
    }
}

GetString("datastore.metric.host") // returns "0.0.0.0"
```
### 拉取子树
简单说就是一个Viper读一个配置文件，可以调用一个方法用这个配置文件的一部分生成一个小viper
例如，viper代表如下的配置:

```go
app:
  cache1:
    max-items: 100
    item-size: 64
  cache2:
    max-items: 200
    item-size: 80
```

在执行了下面的代码后

```go
subv := viper.Sub("app.cache1")
```

subv 代办了:

```go
max-items: 100
item-size: 64
```

假设我们有下面的函数

```go
func NewCache(cfg *Viper)*Cache{...}
```

这个函数根据一个配置返回一个缓存,有了这个特性，可以非常方便生成两个不同的缓存

```go
cfg1 := viper.Sub("app.cache1")
cache1 := NewCache(cfg1)

cfg2 := viper.Sub("app.cache2")
cache2 := NewCache(cfg2)
```

### Unmarshaling
你也可以将配置Unmarshal到一个结构体或是map等特定结构
有两个方法做这个事情:
+ Unmarshal(rawVal interface{}) : error
+ UnmarshalKey(key string, rawVal interface{}) : error

例如:
```go
type config struct {
	Port int
	Name string
	PathMap string `mapstructure:"path_map"`
}

var C config

err := viper.Unmarshal(&C)
if err != nil {
	t.Fatalf("unable to decode into struct, %v", err)
}
```

如果你想要unmarshal的配置里面某些键自己带有"."符号(默认的键的分隔符),你需要更改分隔符
```go
v := viper.NewWithOptions(viper.KeyDelimiter("::"))

v.SetDefault("chart::values", map[string]interface{}{
    "ingress": map[string]interface{}{
        "annotations": map[string]interface{}{
            "traefik.frontend.rule.type":                 "PathPrefix",
            "traefik.ingress.kubernetes.io/ssl-redirect": "true",
        },
    },
})

type config struct {
	Chart struct{
        Values map[string]interface{}
    }
}

var C config

v.Unmarshal(&C)
```

Viper也支持将配置unmarshal到内嵌的结构体

```go
/*
Example config:

module:
    enabled: true
    token: 89h3f98hbwf987h3f98wenf89ehf
*/
type config struct {
	Module struct {
		Enabled bool

		moduleConfig `mapstructure:",squash"`
	}
}

// moduleConfig could be in a module specific package
type moduleConfig struct {
	Token string
}

var C config

err := viper.Unmarshal(&C)
if err != nil {
	t.Fatalf("unable to decode into struct, %v", err)
}
```

Viper底层使用[github.com/mitchellh/mapstructure](http://github.com/mitchellh/mapstructure) 来unmarshal,默认使用mapstructure

#### Marshaling to string
你可能不想把所有的配置写入文件，而是想写道一个字符串中.你可以调用Allsettings()获取所有配置，并使用你想要的方式来marshall

```go
import (
    yaml "gopkg.in/yaml.v2"
    // ...
)

func yamlStringSettings() string {
    c := viper.AllSettings()
    bs, err := yaml.Marshal(c)
    if err != nil {
        log.Fatalf("unable to marshal config to YAML: %v", err)
    }
    return string(bs)
}
```