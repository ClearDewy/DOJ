CREATE DATABASE IF NOT EXISTS `doj`;

USE `doj`;

DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`  (
    `uid` varchar(32) NOT NULL PRIMARY KEY ,
    `username` varchar(64) NOT NULL UNIQUE COMMENT '用户名',
    `password` varchar(64) NOT NULL COMMENT '密码',
    `school` varchar(64) NOT NULL DEFAULT '' COMMENT '学校',
    `major` varchar(64) NOT NULL DEFAULT '' COMMENT '专业',
    `number` varchar(32) NOT NULL DEFAULT '' COMMENT '学号',
    `name` varchar(32) NOT NULL DEFAULT '' COMMENT '真实姓名',
    `gender` tinyint  NOT NULL DEFAULT 0 COMMENT '性别:1：男 2:女 ',
    `cf_username` varchar(64) NOT NULL DEFAULT '' COMMENT 'cf的username',
    `email` varchar(128) NOT NULL DEFAULT '' COMMENT '邮箱' UNIQUE ,
    `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像地址' UNIQUE ,
    `signature` mediumtext COMMENT '个性签名',
    `title_name` varchar(32) NOT NULL DEFAULT '' COMMENT '头衔、称号',
    `title_color` varchar(15) NOT NULL DEFAULT '' COMMENT '头衔、称号的颜色',

    `system_auth` tinyint DEFAULT 0 COMMENT '管理系统配置',
    `user_auth` tinyint DEFAULT 0 COMMENT '管理用户',
    `problem_auth` tinyint DEFAULT 0 COMMENT '0：不能创建题目；1：创建题目；2：管理所有题目',
    `context_auth` tinyint DEFAULT 0 COMMENT '0：不能创建比赛；1：创建比赛；2：管理所有比赛',
    `train_auth` tinyint DEFAULT 0 COMMENT '0：不能创建训练；1：创建训练；2：管理所有训练',
    `problem_status_auth` tinyint DEFAULT 0 COMMENT '重测所有题目，修改所有提交状态',

    `submit_auth` tinyint DEFAULT 1 COMMENT '提交题目',
    `context_attend_auth` tinyint DEFAULT 1 COMMENT '参加比赛',
    `train_attend_auth` tinyint DEFAULT 1 COMMENT '参加训练',

    `gmt_create` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',

    INDEX (`username`),
    INDEX (`email`)
);

DROP TABLE IF EXISTS `problem`;
CREATE TABLE `problem`  (
    `id` int UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT ,
    `problem_id` varchar(255) UNIQUE NOT NULL COMMENT '问题的自定义ID 例如（HOJ-1000）',
    `title` varchar(255) NOT NULL DEFAULT '',
    `author` varchar(64) COMMENT '作者',
    `type` int NOT NULL DEFAULT 0 COMMENT '0为ACM,1为OI',
    `description` longtext,
    `input` longtext,
    `output` longtext,
    `examples` longtext COMMENT '题面样例',
    `oj` varchar(255) NOT NULL DEFAULT 'Mine' COMMENT '该题目属于哪个oj，自身oj为Mine',
    `source` text,
    `difficulty` int DEFAULT 0 COMMENT '题目难度,0简单，1中等，2困难',
    `hint` longtext,
    `auth` int DEFAULT 1 COMMENT '默认为1公开，2为私有，3为比赛题目',
    `code_share` tinyint DEFAULT 1 COMMENT '该题目对应的相关提交代码，用户是否可用分享',
    `judge_mode` varchar(255) NOT NULL DEFAULT 'default' COMMENT '题目评测模式,default、spj、interactive',
    `judge_case_mode` varchar(255) NOT NULL DEFAULT 'default' COMMENT '题目样例评测模式,default(全部评测/得分和),ergodic_without_error(遇错即止),subtask_lowest,subtask_average',
    `user_extra_file` mediumtext COMMENT '题目评测时用户程序的额外文件 json key:name value:content',
    `judge_extra_file` mediumtext COMMENT '题目评测时交互或特殊程序的额外文件 json key:name value:content',
    `spj_code` longtext,
    `spj_lid` int COMMENT '特判程序或交互程序代码的语言',
    `is_remove_end_blank` tinyint DEFAULT 1 COMMENT '是否默认去除用户代码的文末空格',
    `open_case_result` tinyint DEFAULT 1 COMMENT '是否默认开启该题目的测试样例结果查看',
    `case_version` varchar(40) NOT NULL DEFAULT '0' COMMENT '题目测试数据的版本号',
    `modified_user` varchar(255) COMMENT '修改题目的管理员用户名',
    `is_file_io` tinyint DEFAULT 0 COMMENT '是否是file io自定义输入输出文件模式',
    `io_read_file_name` varchar(255) COMMENT '题目指定的file io输入文件的名称',
    `io_write_file_name` varchar(255) COMMENT '题目指定的file io输出文件的名称',
    `gmt_create` datetime NULL DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` datetime NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX (`problem_id`),

    FOREIGN KEY (`author`) REFERENCES `user_info`(`username`) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (`modified_user`) REFERENCES `user_info`(`username`) ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (`spj_lid`) REFERENCES `language`(`id`) ON DELETE SET NULL ON UPDATE CASCADE
)AUTO_INCREMENT=1000;

DROP TABLE IF EXISTS `problem_case`;
CREATE TABLE `problem_case` (
    `id` int unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键id',
    `pid` int unsigned NOT NULL COMMENT '题目id',
    `input` longtext COMMENT '测试样例的输入',
    `output` longtext COMMENT '测试样例的输出',
    `score` int COMMENT '该测试样例的IO得分',
    `status` int DEFAULT '0' COMMENT '0可用，1不可用',
    `group_num` int DEFAULT '1' COMMENT 'subtask分组的编号',
    `gmt_create` datetime DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (`pid`) REFERENCES `problem` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

DROP TABLE IF EXISTS `language`;
CREATE TABLE `language` (
    `id` int unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `content_type` varchar(255) NOT NULL DEFAULT '' COMMENT '语言类型',
    `description` varchar(255) NOT NULL DEFAULT '' COMMENT '语言描述',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '语言名字',
    `compile_command` mediumtext COMMENT '编译指令',
    `template` longtext COMMENT '模板',
    `code_template` longtext COMMENT '语言默认代码模板',
    `is_spj` tinyint DEFAULT 0 COMMENT '是否可作为特殊判题的一种语言',
    `oj` varchar(255) NOT NULL DEFAULT 'Mine' COMMENT '该题目属于哪个oj，自身oj为Mine',
    `gmt_create` datetime DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX (`oj`)
);

DROP TABLE IF EXISTS `problem_language`;
CREATE TABLE `problem_language` (
    `id` int unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `pid` int unsigned NOT NULL ,
    `lid` int unsigned NOT NULL ,
    `time_limit` bigint UNSIGNED NULL DEFAULT 1000 COMMENT '单位ms',
    `memory_limit` bigint UNSIGNED  DEFAULT 256 COMMENT '单位mb',
    `stack_limit` bigint UNSIGNED  DEFAULT 128 COMMENT '单位mb',
    `gmt_create` datetime DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (`pid`) REFERENCES `problem` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`lid`) REFERENCES `language` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

DROP TABLE IF EXISTS `tag_classification`;
CREATE TABLE `tag_classification`  (
    `id` int UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(255) NOT NULL COMMENT '标签分类名称',
    `oj` varchar(255) NOT NULL DEFAULT 'Mine' COMMENT '该题目属于哪个oj，自身oj为Mine',
    `gmt_create` datetime NULL,
    `gmt_modified` datetime NULL,
    `rank` int UNSIGNED ZEROFILL NULL COMMENT '标签分类优先级 越小越高',
    INDEX (`name`),
    INDEX (`oj`)
);

DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag`  (
    `id` int UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(255) NOT NULL COMMENT '标签名字',
    `color` varchar(10) NOT NULL COMMENT '标签颜色',
    `tcid` int UNSIGNED ,
    `gmt_create` datetime NULL DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX (`name`),
    FOREIGN KEY (`tcid`) REFERENCES `tag_classification` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
);

DROP TABLE IF EXISTS `problem_tag`;
CREATE TABLE `problem_tag`  (
    `id` int UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `pid` int UNSIGNED NOT NULL ,
    `tid` int UNSIGNED NOT NULL ,
    `gmt_create` datetime NULL DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`pid`) REFERENCES `problem` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`tid`) REFERENCES `tag` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

DROP TABLE IF EXISTS `judge`;
CREATE TABLE `judge` (
    `id` int UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `pid` int UNSIGNED NOT NULL ,
    `problem_id` varchar(255) NOT NULL COMMENT '问题的自定义ID',
    `uid` varchar(32) NOT NULL,
    `username` varchar(64) NOT NULL COMMENT '提交用户的用户名',
    `lid` int UNSIGNED NOT NULL ,
    `submit_time` datetime NOT NULL COMMENT '代码提交时间',
    `status` int NOT NULL  COMMENT '评测结果码，具体参考文档',
    `message` mediumtext COMMENT '评测信息提示',
    `time` int COMMENT '运行时间(ms)',
    `memory` int COMMENT '运行内存（kb）',
    `code` longtext NOT NULL COMMENT '提交的代码',
    `length` int NOT NULL default 0 COMMENT '代码长度',
    `vjudge_submit_id` varchar(127) DEFAULT '' COMMENT 'vjudge判题在其它oj的提交id',
    `is_manual` tinyint NOT NULL default 0 COMMENT '是否为人工评测',

    INDEX (`vjudge_submit_id`),
    INDEX (`problem_id`,`uid`),

    FOREIGN KEY (`pid`) REFERENCES `problem` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`uid`) REFERENCES `user_info` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`problem_id`) REFERENCES `problem` (`problem_id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`username`) REFERENCES `user_info` (`username`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`lid`) REFERENCES `language` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);


DROP TABLE IF EXISTS `judge_case`;
CREATE TABLE `judge_case`(
    `id` int UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `jid` int UNSIGNED NOT NULL COMMENT '提交判题id',
    `problem_case_id` int UNSIGNED NOT NULL COMMENT '对应测试数据的样例id',

    `status` int NOT NULL  COMMENT '该样例评测结果码，具体参考文档',
    `message` mediumtext COMMENT '评测信息提示',
    `time` int COMMENT '运行时间(ms)',
    `memory` int COMMENT '运行内存（kb）',

    UNIQUE (`jid`,`problem_case_id`),

    FOREIGN KEY (`jid`) REFERENCES `judge` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`problem_case_id`) REFERENCES `problem_case` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

# 插入数据
INSERT INTO user_info(`uid`,`username`,`password`) VALUES ('1','ClearDewy','doj123456');

INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('c', 'GCC 9.4.0', 'C', '/usr/bin/gcc -DONLINE_JUDGE -w -fmax-errors=3 -std=c11 {src_path} -lm -o {exe_path}', '#include <stdio.h>\r\nint main() {\r\n    int a,b;\r\n    scanf(\"%d %d\",&a,&b);\r\n    printf(\"%d\",a+b);\r\n    return 0;\r\n}', '//PREPEND BEGIN\r\n#include <stdio.h>\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\nint add(int a, int b) {\r\n  // Please fill this blank\r\n  return ___________;\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nint main() {\r\n  printf(\"%d\", add(1, 2));\r\n  return 0;\r\n}\r\n//APPEND END', 1, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('c', 'GCC 9.4.0', 'C With O2', '/usr/bin/gcc -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c11 {src_path} -lm -o {exe_path}', '#include <stdio.h>\r\nint main() {\r\n    int a,b;\r\n    scanf(\"%d %d\",&a,&b);\r\n    printf(\"%d\",a+b);\r\n    return 0;\r\n}', '//PREPEND BEGIN\r\n#include <stdio.h>\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\nint add(int a, int b) {\r\n  // Please fill this blank\r\n  return ___________;\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nint main() {\r\n  printf(\"%d\", add(1, 2));\r\n  return 0;\r\n}\r\n//APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'G++ 9.4.0', 'C++', '/usr/bin/g++ -DONLINE_JUDGE -w -fmax-errors=3 -std=c++14 {src_path} -lm -o {exe_path}', '#include<iostream>\r\nusing namespace std;\r\nint main()\r\n{\r\n    int a,b;\r\n    cin >> a >> b;\r\n    cout << a + b;\r\n    return 0;\r\n}', '//PREPEND BEGIN\r\n#include <iostream>\r\nusing namespace std;\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\nint add(int a, int b) {\r\n  // Please fill this blank\r\n  return ___________;\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nint main() {\r\n  cout << add(1, 2);\r\n  return 0;\r\n}\r\n//APPEND END', 1, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'G++ 9.4.0', 'C++ With O2', '/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++14 {src_path} -lm -o {exe_path}', '#include<iostream>\r\nusing namespace std;\r\nint main()\r\n{\r\n    int a,b;\r\n    cin >> a >> b;\r\n    cout << a + b;\r\n    return 0;\r\n}', '//PREPEND BEGIN\r\n#include <iostream>\r\nusing namespace std;\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\nint add(int a, int b) {\r\n  // Please fill this blank\r\n  return ___________;\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nint main() {\r\n  cout << add(1, 2);\r\n  return 0;\r\n}\r\n//APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'G++ 9.4.0', 'C++ 17', '/usr/bin/g++ -DONLINE_JUDGE -w -fmax-errors=3 -std=c++17 {src_path} -lm -o {exe_path}', '#include<iostream>\r\nusing namespace std;\r\nint main()\r\n{\r\n    int a,b;\r\n    cin >> a >> b;\r\n    cout << a + b;\r\n    return 0;\r\n}', '//PREPEND BEGIN\r\n#include <iostream>\r\nusing namespace std;\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\nint add(int a, int b) {\r\n  // Please fill this blank\r\n  return ___________;\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nint main() {\r\n  cout << add(1, 2);\r\n  return 0;\r\n}\r\n//APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'G++ 9.4.0', 'C++ 17 With O2', '/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++17 {src_path} -lm -o {exe_path}', '#include<iostream>\r\nusing namespace std;\r\nint main()\r\n{\r\n    int a,b;\r\n    cin >> a >> b;\r\n    cout << a + b;\r\n    return 0;\r\n}', '//PREPEND BEGIN\r\n#include <iostream>\r\nusing namespace std;\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\nint add(int a, int b) {\r\n  // Please fill this blank\r\n  return ___________;\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nint main() {\r\n  cout << add(1, 2);\r\n  return 0;\r\n}\r\n//APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'G++ 9.4.0', 'C++ 20', '/usr/bin/g++ -DONLINE_JUDGE -w -fmax-errors=3 -std=c++20 {src_path} -lm -o {exe_path}', '#include<iostream>\r\nusing namespace std;\r\nint main()\r\n{\r\n    int a,b;\r\n    cin >> a >> b;\r\n    cout << a + b;\r\n    return 0;\r\n}', '//PREPEND BEGIN\r\n#include <iostream>\r\nusing namespace std;\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\nint add(int a, int b) {\r\n  // Please fill this blank\r\n  return ___________;\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nint main() {\r\n  cout << add(1, 2);\r\n  return 0;\r\n}\r\n//APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'G++ 9.4.0', 'C++ 20 With O2', '/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++20 {src_path} -lm -o {exe_path}', '#include<iostream>\r\nusing namespace std;\r\nint main()\r\n{\r\n    int a,b;\r\n    cin >> a >> b;\r\n    cout << a + b;\r\n    return 0;\r\n}', '//PREPEND BEGIN\r\n#include <iostream>\r\nusing namespace std;\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\nint add(int a, int b) {\r\n  // Please fill this blank\r\n  return ___________;\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nint main() {\r\n  cout << add(1, 2);\r\n  return 0;\r\n}\r\n//APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('java', 'OpenJDK 1.8', 'Java', '/usr/bin/javac {src_path} -d {exe_dir} -encoding UTF8', 'import java.util.Scanner;\r\npublic class Main{\r\n    public static void main(String[] args){\r\n        Scanner in=new Scanner(System.in);\r\n        int a=in.nextInt();\r\n        int b=in.nextInt();\r\n        System.out.println((a+b));\r\n    }\r\n}', '//PREPEND BEGIN\r\nimport java.util.Scanner;\r\n//PREPEND END\r\n\r\npublic class Main{\r\n    //TEMPLATE BEGIN\r\n    public static Integer add(int a,int b){\r\n        return _______;\r\n    }\r\n    //TEMPLATE END\r\n\r\n    //APPEND BEGIN\r\n    public static void main(String[] args){\r\n        System.out.println(add(a,b));\r\n    }\r\n    //APPEND END\r\n}\r\n', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('python', 'Python 3.7.5', 'Python3', '/usr/bin/python3 -m py_compile {src_path}', 'a, b = map(int, input().split())\r\nprint(a + b)', '//PREPEND BEGIN\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\ndef add(a, b):\r\n    return a + b\r\n//TEMPLATE END\r\n\r\n\r\nif __name__ == \'__main__\':  \r\n    //APPEND BEGIN\r\n    a, b = 1, 1\r\n    print(add(a, b))\r\n    //APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('python', 'Python 2.7.17', 'Python2', '/usr/bin/python -m py_compile {src_path}', 'a, b = map(int, raw_input().split())\r\nprint a+b', '//PREPEND BEGIN\r\n//PREPEND END\r\n\r\n//TEMPLATE BEGIN\r\ndef add(a, b):\r\n    return a + b\r\n//TEMPLATE END\r\n\r\n\r\nif __name__ == \'__main__\':  \r\n    //APPEND BEGIN\r\n    a, b = 1, 1\r\n    print add(a, b)\r\n    //APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('go', 'Golang 1.19', 'Golang', '/usr/bin/go build -o {exe_path} {src_path}', 'package main\r\nimport \"fmt\"\r\n\r\nfunc main(){\r\n    var x int\r\n    var y int\r\n    fmt.Scanln(&x,&y)\r\n    fmt.Printf(\"%d\",x+y)  \r\n}\r\n', '\r\npackage main\r\n\r\n//PREPEND BEGIN\r\nimport \"fmt\"\r\n//PREPEND END\r\n\r\n\r\n//TEMPLATE BEGIN\r\nfunc add(a,b int)int{\r\n    return ______\r\n}\r\n//TEMPLATE END\r\n\r\n//APPEND BEGIN\r\nfunc main(){\r\n    var x int\r\n    var y int\r\n    fmt.Printf(\"%d\",add(x,y))  \r\n}\r\n//APPEND END\r\n', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('csharp', 'C# Mono 4.6.2', 'C#', '/usr/bin/mcs -optimize+ -out:{exe_path} {src_path}', 'using System;\r\nusing System.Linq;\r\n\r\nclass Program {\r\n    public static void Main(string[] args) {\r\n        Console.WriteLine(Console.ReadLine().Split().Select(int.Parse).Sum());\r\n    }\r\n}', '//PREPEND BEGIN\r\nusing System;\r\nusing System.Collections.Generic;\r\nusing System.Text;\r\n//PREPEND END\r\n\r\nclass Solution\r\n{\r\n    //TEMPLATE BEGIN\r\n    static int add(int a,int b){\r\n        return _______;\r\n    }\r\n    //TEMPLATE END\r\n\r\n    //APPEND BEGIN\r\n    static void Main(string[] args)\r\n    {\r\n        int a ;\r\n        int b ;\r\n        Console.WriteLine(add(a,b));\r\n    }\r\n    //APPEND END\r\n}', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('php', 'PHP 7.3.33', 'PHP', '/usr/bin/php {src_path}', '<?=array_sum(fscanf(STDIN, \"%d %d\"));', NULL, 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('python', 'PyPy 2.7.18 (7.3.8)', 'PyPy2', '/usr/bin/pypy -m py_compile {src_path}', 'print sum(int(x) for x in raw_input().split(\' \'))', '//PREPEND BEGIN\n//PREPEND END\n\n//TEMPLATE BEGIN\ndef add(a, b):\n    return a + b\n//TEMPLATE END\n\n\nif __name__ == \'__main__\':  \n    //APPEND BEGIN\n    a, b = 1, 1\n    print add(a, b)\n    //APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('python', 'PyPy 3.8.12 (7.3.8)', 'PyPy3', '/usr/bin/pypy3 -m py_compile {src_path}', 'print(sum(int(x) for x in input().split(\' \')))', '//PREPEND BEGIN\n//PREPEND END\n\n//TEMPLATE BEGIN\ndef add(a, b):\n    return a + b\n//TEMPLATE END\n\n\nif __name__ == \'__main__\':  \n    //APPEND BEGIN\n    a, b = 1, 1\n    print(add(a, b))\n    //APPEND END', 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('javascript', 'Node.js 14.19.0', 'JavaScript Node', '/usr/bin/node {src_path}', 'var readline = require(\'readline\');\nconst rl = readline.createInterface({\n        input: process.stdin,\n        output: process.stdout\n});\nrl.on(\'line\', function(line){\n   var tokens = line.split(\' \');\n    console.log(parseInt(tokens[0]) + parseInt(tokens[1]));\n});', NULL, 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('javascript', 'JavaScript V8 8.4.109', 'JavaScript V8', '/usr/bin/jsv8/d8 {src_path}', 'const [a, b] = readline().split(\' \').map(n => parseInt(n, 10));\nprint((a + b).toString());', NULL, 0, 'Mine');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('c', 'GCC', 'GCC', NULL, NULL, NULL, 0, 'HDU');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'G++', 'G++', NULL, NULL, NULL, 0, 'HDU');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'C++', 'C++', NULL, NULL, NULL, 0, 'HDU');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('c', 'C', 'C', NULL, NULL, NULL, 0, 'HDU');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('pascal', 'Pascal', 'Pascal', NULL, NULL, NULL, 0, 'HDU');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('java', 'Java', 'Java', NULL, NULL, NULL, 0, 'HDU');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('csharp', 'C#', 'C#', NULL, NULL, NULL, 0, 'HDU');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('c', 'GNU GCC C11 5.1.0', 'GNU GCC C11 5.1.0', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'Clang++17 Diagnostics', 'Clang++17 Diagnostics', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'GNU G++14 6.4.0', 'GNU G++14 6.4.0', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'GNU G++17 7.3.0', 'GNU G++17 7.3.0', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'GNU G++20 11.2.0 (64 bit, winlibs)', 'GNU G++20 11.2.0 (64 bit, winlibs)', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'Microsoft Visual C++ 2017', 'Microsoft Visual C++ 2017', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('csharp', 'C# Mono 6.8', 'C# Mono 6.8', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('d', 'D DMD32 v2.091.0', 'D DMD32 v2.091.0', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('go', 'Go 1.15.6', 'Go 1.15.6', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('haskell', 'Haskell GHC 8.10.1', 'Haskell GHC 8.10.1', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('java', 'Java 1.8.0_241', 'Java 1.8.0_241', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('java', 'Kotlin 1.4.0', 'Kotlin 1.4.0', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('ocaml', 'OCaml 4.02.1', 'OCaml 4.02.1', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('pascal', 'Delphi 7', 'Delphi 7', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('pascal', 'Free Pascal 3.0.2', 'Free Pascal 3.0.2', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('pascal', 'PascalABC.NET 3.4.2', 'PascalABC.NET 3.4.2', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('perl', 'Perl 5.20.1', 'Perl 5.20.1', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('php', 'PHP 7.2.13', 'PHP 7.2.13', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('python', 'Python 2.7.18', 'Python 2.7.18', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('python', 'Python 3.9.1', 'Python 3.9.1', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('python', 'PyPy 2.7 (7.3.0)', 'PyPy 2.7 (7.3.0)', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('python', 'PyPy 3.7 (7.3.0)', 'PyPy 3.7 (7.3.0)', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('ruby', 'Ruby 3.0.0', 'Ruby 3.0.0', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('rust', 'Rust 1.49.0', 'Rust 1.49.0', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('scala', 'Scala 2.12.8', 'Scala 2.12.8', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('javascript', 'JavaScript V8 4.8.0', 'JavaScript V8 4.8.0', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('javascript', 'Node.js 12.6.3', 'Node.js 12.6.3', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('csharp', 'C# 8, .NET Core 3.1', 'C# 8, .NET Core 3.1', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('csharp', 'Java 11.0.6', 'Java 11.0.6', NULL, NULL, NULL, 0, 'CF');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'G++', 'G++', NULL, NULL, NULL, 0, 'POJ');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('c', 'GCC', 'GCC', NULL, NULL, NULL, 0, 'POJ');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('java', 'Java', 'Java', NULL, NULL, NULL, 0, 'POJ');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('pascal', 'Pascal', 'Pascal', NULL, NULL, NULL, 0, 'POJ');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('cpp', 'C++', 'C++', NULL, NULL, NULL, 0, 'POJ');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('c', 'C', 'C', NULL, NULL, NULL, 0, 'POJ');
INSERT INTO `language` (`content_type`, `description`, `name`, `compile_command`, `template`, `code_template`, `is_spj`, `oj`) VALUES ('fortran', 'Fortran', 'Fortran', NULL, NULL, NULL, 0, 'POJ');
