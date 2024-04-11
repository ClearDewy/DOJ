/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import "time"

type LanguageCompileType struct {
	Command string
}
type LanguageRunType struct {
	Command string
}

type LanguageType struct {
	Language string
	Env      []string
	SrcPath  string
	ExePath  string
	Compile  *LanguageCompileType
	Run      *LanguageRunType

	MaxCpuTime  time.Duration
	MaxRealTime time.Duration
	MaxMemory   uint64
}

func FindLanguage(name string) *LanguageType {
	for _, value := range LanguageList {
		if value.Language == name {
			return value
		}
	}
	return nil
}

var (
	DefaultEnv = []string{
		"PATH=/usr/bin:/bin",
	}

	Pyhton3Env = []string{
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"LANG=en_US.UTF-8",
		"LANGUAGE=en_US:en",
		"LC_ALL=en_US.UTF-8",
		"PYTHONIOENCODING=utf-8",
	}
	GolangEnv = []string{
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"GOCACHE=/w",
		"GOPATH=/w/go",
		"LANG=en_US.UTF-8",
		"LANGUAGE=en_US:en",
		"LC_ALL=en_US.UTF-8",
		"GODEBUG=madvdontneed=1",
	}
)

// LanguageList 语言命令，不要更改顺序
var LanguageList = []*LanguageType{
	{
		Language: "C",
		Env:      DefaultEnv,
		SrcPath:  "main.c",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/gcc -DONLINE_JUDGE -w -fmax-errors=1 -std=c11 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  3 * time.Second,
		MaxRealTime: 10 * time.Second,
		MaxMemory:   256 << 20, // 256 MB
	},
	// C With O2
	{
		Language: "C With O2",
		Env:      DefaultEnv,
		SrcPath:  "main.c",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/gcc -DONLINE_JUDGE -O2 -w -fmax-errors=1 -std=c11 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  3 * time.Second,
		MaxRealTime: 10 * time.Second,
		MaxMemory:   256 << 20, // 256 MB
	},
	// C++
	{
		Language: "C++",
		Env:      DefaultEnv,
		SrcPath:  "main.cpp",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/g++ -DONLINE_JUDGE -w -fmax-errors=1 -std=c++14 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  10 * time.Second,
		MaxRealTime: 20 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// C++ With O2
	{
		Language: "C++ With O2",
		Env:      DefaultEnv,
		SrcPath:  "main.cpp",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=1 -std=c++14 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  10 * time.Second,
		MaxRealTime: 20 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// C++ 17
	{
		Language: "C++ 17",
		Env:      DefaultEnv,
		SrcPath:  "main.cpp",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/g++ -DONLINE_JUDGE -w -fmax-errors=1 -std=c++17 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  10 * time.Second,
		MaxRealTime: 20 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// C++ 17 With O2
	{
		Language: "C++ 17 With O2",
		Env:      DefaultEnv,
		SrcPath:  "main.cpp",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=1 -std=c++17 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  10 * time.Second,
		MaxRealTime: 20 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// C++ 20
	{
		Language: "C++ 20",
		Env:      DefaultEnv,
		SrcPath:  "main.cpp",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/g++ -DONLINE_JUDGE -w -fmax-errors=1 -std=c++2a {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  10 * time.Second,
		MaxRealTime: 20 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// C++ 20 With O2
	{
		Language: "C++ 20 With O2",
		Env:      DefaultEnv,
		SrcPath:  "main.cpp",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=1 -std=c++2a {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  10 * time.Second,
		MaxRealTime: 20 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// Java
	{
		Language: "Java",
		Env:      DefaultEnv,
		SrcPath:  "Main.java",
		ExePath:  "Main.jar",
		Compile: &LanguageCompileType{
			Command: "/bin/bash -c \"javac -encoding utf-8 {src_path} && jar -cvf {exe_path} *.class\"",
		},
		Run: &LanguageRunType{
			Command: "/usr/bin/java -Dfile.encoding=UTF-8 -cp {exe_path} Main",
		},
		MaxCpuTime:  10 * time.Second,
		MaxRealTime: 20 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// Python2
	{
		Language: "Python2",
		Env:      DefaultEnv,
		SrcPath:  "main.py",
		ExePath:  "main.pyc",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/python -m py_compile {src_path}",
		},
		Run: &LanguageRunType{
			Command: "/usr/bin/python {exe_path}",
		},
		MaxCpuTime:  5 * time.Second,
		MaxRealTime: 10 * time.Second,
		MaxMemory:   256 << 20, // 256 MB
	},
	// Python3
	{
		Language: "Python3",
		Env:      Pyhton3Env,
		SrcPath:  "main.py",
		ExePath:  "__pycache__/main.cpython-37.pyc",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/python3.7 -m py_compile {src_path}",
		},
		Run: &LanguageRunType{
			Command: "/usr/bin/python3.7 {exe_path}",
		},
		MaxCpuTime:  5 * time.Second,
		MaxRealTime: 10 * time.Second,
		MaxMemory:   256 << 20, // 256 MB
	},
	// PyPy2
	{
		Language: "PyPy2",
		Env:      DefaultEnv,
		SrcPath:  "main.py",
		ExePath:  "__pycache__/main.pypy-73.pyc",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/pypy -m py_compile {src_path}",
		},
		Run: &LanguageRunType{
			Command: "/usr/bin/pypy {exe_path}",
		},
		MaxCpuTime:  5 * time.Second,
		MaxRealTime: 10 * time.Second,
		MaxMemory:   256 << 20, // 256 MB
	},
	// PyPy3
	{
		Language: "PyPy3",
		Env:      Pyhton3Env,
		SrcPath:  "main.py",
		ExePath:  "__pycache__/main.pypy38.pyc",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/pypy3 -m py_compile {src_path}",
		},
		Run: &LanguageRunType{
			Command: "/usr/bin/pypy3 {exe_path}",
		},
		MaxCpuTime:  5 * time.Second,
		MaxRealTime: 10 * time.Second,
		MaxMemory:   256 << 20, // 256 MB
	},
	// Golang
	{
		Language: "Golang",
		Env:      GolangEnv,
		SrcPath:  "main.go",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/go build -o {exe_path} {src_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path}",
		},
		MaxCpuTime:  3 * time.Second,
		MaxRealTime: 5 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// C#
	{
		Language: "C#",
		Env:      DefaultEnv,
		SrcPath:  "Main.cs",
		ExePath:  "main",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/mcs -optimize+ -out:{exe_path} {src_path}",
		},
		Run: &LanguageRunType{
			Command: "/usr/bin/mono {exe_path}",
		},
		MaxCpuTime:  5 * time.Second,
		MaxRealTime: 10 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// PHP
	{
		Language: "PHP",
		Env:      DefaultEnv,
		SrcPath:  "main.php",
		ExePath:  "main.php",
		Run: &LanguageRunType{
			Command: "/usr/bin/php {exe_path}",
		},
	},
	// JavaScript Node
	{
		Language: "JavaScript Node",
		Env:      DefaultEnv,
		SrcPath:  "main.js",
		ExePath:  "main.js",
		Run: &LanguageRunType{
			Command: "/usr/bin/node {exe_path}",
		},
	},
	// JavaScript V8
	{
		Language: "JavaScript V8",
		Env:      DefaultEnv,
		SrcPath:  "main.js",
		ExePath:  "main.js",
		Run: &LanguageRunType{
			Command: "/usr/bin/jsv8/d8 {exe_path}",
		},
	},
	// SPJ-C
	{
		Language: "SPJ-C",
		Env:      DefaultEnv,
		SrcPath:  "spj.c",
		ExePath:  "spj",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/gcc -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c11 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path} {std_input} {user_output} {std_output}",
		},
		MaxCpuTime:  8 * time.Second,
		MaxRealTime: 15 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// SPJ-C++
	{
		Language: "SPJ-C++",
		Env:      DefaultEnv,
		SrcPath:  "spj.cpp",
		ExePath:  "spj",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++14 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path} {std_input} {user_output} {std_output}",
		},
		MaxCpuTime:  15 * time.Second,
		MaxRealTime: 25 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// INTERACTIVE-C
	{
		Language: "INTERACTIVE-C",
		Env:      DefaultEnv,
		SrcPath:  "interactive.c",
		ExePath:  "interactive",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/gcc -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c11 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path} {std_input} {user_output} {std_output}",
		},
		MaxCpuTime:  8 * time.Second,
		MaxRealTime: 15 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
	// INTERACTIVE-C++
	{
		Language: "INTERACTIVE-C++",
		Env:      DefaultEnv,
		SrcPath:  "interactive.cpp",
		ExePath:  "interactive",
		Compile: &LanguageCompileType{
			Command: "/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++14 {src_path} -lm -o {exe_path}",
		},
		Run: &LanguageRunType{
			Command: "{exe_path} {std_input} {user_output} {std_output}",
		},
		MaxCpuTime:  15 * time.Second,
		MaxRealTime: 25 * time.Second,
		MaxMemory:   512 << 20, // 512 MB
	},
}
