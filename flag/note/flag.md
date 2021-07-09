# flag包

## 基本使用

```go
// 1.定义flag
// flag有以下三种使用, "-" 和 "--" 等价
// -name 和 -name=true 等价
// -name=value
// -name value bool不能这么用


// var n = flag.Type(name,defaultValue,usage)
n = flag.String("name","tom","my name")

// flag.TypeVar(p,name,defaultValue,usage)
var n string
flag.StringVar(&n,"name","tom","my name")

// 2.解析flag，一定要将所有的flag声明完毕才能调用
flag.Parse()
```

flag包的大多函数都是对FlagSet的一层包装,flag中有一个默认的CommandLine对象，他的类型就是FlagSet.
```go
// 使用Lookup查找flagset中的某个flag
myFlagSet.Lookup("name")

// 使用Set设置flagset中某个flag的值
myFlagSet.Set("name","mike")
```


## 主要结构

### Flag
Flag结构体表示一个flag的状态
```go
type Flag struct {
	Name     string // name as it appears on command line
	Usage    string // help message
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}
```

### FlagSet 
FlagSet表示预定义的一个Flag的集合，整个flag包都是围绕它展开的
FlagSet的零值没有名字并且使用ContinueOnError模式(下面会讲到)来处理解析flag时遇到的错
一个FlagSet里的每一个flag的name必须是唯一的，当试图在一个FlagSet中定义同名的Flag将会导致Panic
```go
type FlagSet struct {
	// Usage 是当解析flags错误时调用的函数，这个字段是一个函数(不是一个方法)，
	// 我们可以自己重新编写一个函数来替换默认的，来决定如何处理解析错误.
	// 当Usage函数被调用后，之后的行为取决于ErrorHandling的值.对于命令行来说
	// 默认调用ExitOnError，即再出错调用Usage函数后，直接退出程序.
	Usage func()

	name          string // FlagSet的name
	parsed        bool  // 是否已经解析
	actual        map[string]*Flag // 定义的flags
	formal        map[string]*Flag // 实际解析到的flags
	args          []string // 解析完之后剩下的参数
	errorHandling ErrorHandling // errorHanding的模式
	output        io.Writer // nil means stderr; use Output() accessor
}
```

### Value
如果要实现一种类型的flag绑定，需要实现flag接口

如果一种类型还实现了IsBoolFlag()的方法并且在方法中返回true，此时 -flagName 与
-flagName=true 等价

Set方法在赋值时只被调用一次。String方法可能被一个零值的对象调用，例如一个nil的指针，所以在自己编写String()
方法时要注意处理这个，否则会造成运行时的错误

```go
type Value interface {
    String() string
    Set(string) error
}
```

### Getter
Getter is an interface that allows the contents of a Value to be retrieved.
Getter是在Value接口的基础上扩展的一个接口， 可以允许Value的内容被检索
```go
type Getter interface {
    Value
    Get() interface{}
}
```

### TypeValue
```go
type TypeValue Type
// 基本大部分的ValueType类型都实现了下面的四个方法，即实现了Getter接口
// funcValue 没有实现Get方法，只实现了Value接口
// boolValue 还是实现了IsBoolFlag方法

func newTypeValue(val Type, p *Type) *TypeValue

func(tv *TypeValue) String() string
func (tv *TypeValue) Set(s string)error
func (tv *TypeValue) Get()interface{}

func (tv *TypeValue)IsBoolFlag() bool
```

以BoolValue具体分析

```go
// -- bool Value
type boolValue bool

// 构造函数
func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

// Set方法，用于parse时的赋值
func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		err = errParse
	}
	*b = boolValue(v)
	return err
}

// Get方法，用于获取boolValue的值
func (b *boolValue) Get() interface{} { return bool(*b) }

// String方法，用于输出
func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }
```

### ErrorHandling

```go
// ErrorHanding 定义了当解析失败时 FlagSet.Parse如何处理
type ErrorHandling int

const (
    ContinueOnError ErrorHandling = iota // 继续解析
    ExitOnError                          // 退出程序
    PanicOnError                         // panic
)

```


## 源码分析

### 声明一个flag(以string类型为例子)
当我们使用flag.String(name,value,usage)时，发生了什么呢
函数调用过程如下
```go
func (f *FlagSet) String(name string, value string, usage string) *string {
    p := new(string)
    f.StringVar(p, name, value, usage)
    return p
}
```
可以看到String()方法只是对StringVar()的进一步包装,继续看StringVar()
newStringValue将value存放到p中，并且转换为实现了Value接口的stringValue类型
```go
func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.Var(newStringValue(value, p), name, usage)
}
```
StringVar内部调用Var，继续看Var()
```go
func (f *FlagSet) Var(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	// 初始化一个Flag
	flag := &Flag{name, usage, value, value.String()}
	// 检查FlagSet中是否存在同名flag
	_, alreadythere := f.formal[name]
	if alreadythere {
		var msg string
		if f.name == "" {
			msg = fmt.Sprintf("flag redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s flag redefined: %s", f.name, name)
		}
		fmt.Fprintln(f.Output(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if f.formal == nil {
		f.formal = make(map[string]*Flag)
	}
	// 将flag存入
	f.formal[name] = flag
}
```
说白了就是内部维护一个map存放你定义的flag

### 命令行参数解析
```go
func (f *FlagSet) Parse(arguments []string) error {
	// 将parsed标记为true
	f.parsed = true
	f.args = arguments
	// 逐个解析flag
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		switch f.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			if err == ErrHelp {
				os.Exit(0)
			}
			os.Exit(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}
```






