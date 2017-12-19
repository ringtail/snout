package num

var Linux_arm = map[int]string{
	0:      "restart_syscall",
	1:      "exit",
	2:      "fork",
	3:      "read",
	4:      "write",
	5:      "open",
	6:      "close",
	8:      "creat",
	9:      "link",
	10:     "unlink",
	11:     "execve",
	12:     "chdir",
	13:     "time",
	14:     "mknod",
	15:     "chmod",
	16:     "lchown",
	19:     "lseek",
	20:     "getpid",
	21:     "mount",
	22:     "umount",
	23:     "setuid",
	24:     "getuid",
	25:     "stime",
	26:     "ptrace",
	27:     "alarm",
	29:     "pause",
	30:     "utime",
	33:     "access",
	34:     "nice",
	36:     "sync",
	37:     "kill",
	38:     "rename",
	39:     "mkdir",
	40:     "rmdir",
	41:     "dup",
	42:     "pipe",
	43:     "times",
	45:     "brk",
	46:     "setgid",
	47:     "getgid",
	49:     "geteuid",
	50:     "getegid",
	51:     "acct",
	52:     "umount2",
	54:     "ioctl",
	55:     "fcntl",
	57:     "setpgid",
	60:     "umask",
	61:     "chroot",
	62:     "ustat",
	63:     "dup2",
	64:     "getppid",
	65:     "getpgrp",
	66:     "setsid",
	67:     "sigaction",
	70:     "setreuid",
	71:     "setregid",
	72:     "sigsuspend",
	73:     "sigpending",
	74:     "sethostname",
	75:     "setrlimit",
	77:     "getrusage",
	78:     "gettimeofday",
	79:     "settimeofday",
	80:     "getgroups",
	81:     "setgroups",
	82:     "select",
	83:     "symlink",
	85:     "readlink",
	86:     "uselib",
	87:     "swapon",
	88:     "reboot",
	89:     "readdir",
	90:     "mmap",
	91:     "munmap",
	92:     "truncate",
	93:     "ftruncate",
	94:     "fchmod",
	95:     "fchown",
	96:     "getpriority",
	97:     "setpriority",
	99:     "statfs",
	100:    "fstatfs",
	102:    "socketcall",
	103:    "syslog",
	104:    "setitimer",
	105:    "getitimer",
	106:    "stat",
	107:    "lstat",
	108:    "fstat",
	111:    "vhangup",
	114:    "wait4",
	115:    "swapoff",
	116:    "sysinfo",
	117:    "ipc",
	118:    "fsync",
	119:    "sigreturn",
	120:    "clone",
	121:    "setdomainname",
	122:    "uname",
	124:    "adjtimex",
	125:    "mprotect",
	126:    "sigprocmask",
	128:    "init_module",
	129:    "delete_module",
	131:    "quotactl",
	132:    "getpgid",
	133:    "fchdir",
	134:    "bdflush",
	135:    "sysfs",
	136:    "personality",
	138:    "setfsuid",
	139:    "setfsgid",
	140:    "_llseek",
	141:    "getdents",
	142:    "_newselect",
	143:    "flock",
	144:    "msync",
	145:    "readv",
	146:    "writev",
	147:    "getsid",
	148:    "fdatasync",
	149:    "_sysctl",
	150:    "mlock",
	151:    "munlock",
	152:    "mlockall",
	153:    "munlockall",
	154:    "sched_setparam",
	155:    "sched_getparam",
	156:    "sched_setscheduler",
	157:    "sched_getscheduler",
	158:    "sched_yield",
	159:    "sched_get_priority_max",
	160:    "sched_get_priority_min",
	161:    "sched_rr_get_interval",
	162:    "nanosleep",
	163:    "mremap",
	164:    "setresuid",
	165:    "getresuid",
	168:    "poll",
	169:    "nfsservctl",
	170:    "setresgid",
	171:    "getresgid",
	172:    "prctl",
	173:    "rt_sigreturn",
	174:    "rt_sigaction",
	175:    "rt_sigprocmask",
	176:    "rt_sigpending",
	177:    "rt_sigtimedwait",
	178:    "rt_sigqueueinfo",
	179:    "rt_sigsuspend",
	180:    "pread64",
	181:    "pwrite64",
	182:    "chown",
	183:    "getcwd",
	184:    "capget",
	185:    "capset",
	186:    "sigaltstack",
	187:    "sendfile",
	190:    "vfork",
	192:    "mmap2",
	193:    "truncate64",
	194:    "ftruncate64",
	195:    "stat64",
	196:    "lstat64",
	197:    "fstat64",
	198:    "lchown32",
	199:    "getuid32",
	200:    "getgid32",
	201:    "geteuid32",
	202:    "getegid32",
	203:    "setreuid32",
	204:    "setregid32",
	205:    "getgroups32",
	206:    "setgroups32",
	207:    "fchown32",
	208:    "setresuid32",
	209:    "getresuid32",
	210:    "setresgid32",
	211:    "getresgid32",
	212:    "chown32",
	213:    "setuid32",
	214:    "setgid32",
	215:    "setfsuid32",
	216:    "setfsgid32",
	217:    "getdents64",
	218:    "pivot_root",
	219:    "mincore",
	220:    "madvise",
	221:    "fcntl64",
	224:    "gettid",
	225:    "readahead",
	226:    "setxattr",
	227:    "lsetxattr",
	228:    "fsetxattr",
	229:    "getxattr",
	230:    "lgetxattr",
	231:    "fgetxattr",
	232:    "listxattr",
	233:    "llistxattr",
	234:    "flistxattr",
	235:    "removexattr",
	236:    "lremovexattr",
	237:    "fremovexattr",
	238:    "tkill",
	239:    "sendfile64",
	240:    "futex",
	241:    "sched_setaffinity",
	242:    "sched_getaffinity",
	243:    "io_setup",
	244:    "io_destroy",
	245:    "io_getevents",
	246:    "io_submit",
	247:    "io_cancel",
	248:    "exit_group",
	249:    "lookup_dcookie",
	250:    "epoll_create",
	251:    "epoll_ctl",
	252:    "epoll_wait",
	253:    "remap_file_pages",
	256:    "set_tid_address",
	257:    "timer_create",
	258:    "timer_settime",
	259:    "timer_gettime",
	260:    "timer_getoverrun",
	261:    "timer_delete",
	262:    "clock_settime",
	263:    "clock_gettime",
	264:    "clock_getres",
	265:    "clock_nanosleep",
	266:    "statfs64",
	267:    "fstatfs64",
	268:    "tgkill",
	269:    "utimes",
	270:    "arm_fadvise64_64",
	271:    "pciconfig_iobase",
	272:    "pciconfig_read",
	273:    "pciconfig_write",
	274:    "mq_open",
	275:    "mq_unlink",
	276:    "mq_timedsend",
	277:    "mq_timedreceive",
	278:    "mq_notify",
	279:    "mq_getsetattr",
	280:    "waitid",
	281:    "socket",
	282:    "bind",
	283:    "connect",
	284:    "listen",
	285:    "accept",
	286:    "getsockname",
	287:    "getpeername",
	288:    "socketpair",
	289:    "send",
	290:    "sendto",
	291:    "recv",
	292:    "recvfrom",
	293:    "shutdown",
	294:    "setsockopt",
	295:    "getsockopt",
	296:    "sendmsg",
	297:    "recvmsg",
	298:    "semop",
	299:    "semget",
	300:    "semctl",
	301:    "msgsnd",
	302:    "msgrcv",
	303:    "msgget",
	304:    "msgctl",
	305:    "shmat",
	306:    "shmdt",
	307:    "shmget",
	308:    "shmctl",
	309:    "add_key",
	310:    "request_key",
	311:    "keyctl",
	312:    "semtimedop",
	313:    "vserver",
	314:    "ioprio_set",
	315:    "ioprio_get",
	316:    "inotify_init",
	317:    "inotify_add_watch",
	318:    "inotify_rm_watch",
	319:    "mbind",
	320:    "get_mempolicy",
	321:    "set_mempolicy",
	322:    "openat",
	323:    "mkdirat",
	324:    "mknodat",
	325:    "fchownat",
	326:    "futimesat",
	327:    "fstatat64",
	328:    "unlinkat",
	329:    "renameat",
	330:    "linkat",
	331:    "symlinkat",
	332:    "readlinkat",
	333:    "fchmodat",
	334:    "faccessat",
	335:    "pselect6",
	336:    "ppoll",
	337:    "unshare",
	338:    "set_robust_list",
	339:    "get_robust_list",
	340:    "splice",
	341:    "arm_sync_file_range",
	342:    "tee",
	343:    "vmsplice",
	344:    "move_pages",
	345:    "getcpu",
	346:    "epoll_pwait",
	347:    "kexec_load",
	348:    "utimensat",
	349:    "signalfd",
	350:    "timerfd_create",
	351:    "eventfd",
	352:    "fallocate",
	353:    "timerfd_settime",
	354:    "timerfd_gettime",
	355:    "signalfd4",
	356:    "eventfd2",
	357:    "epoll_create1",
	358:    "dup3",
	359:    "pipe2",
	360:    "inotify_init1",
	361:    "preadv",
	362:    "pwritev",
	363:    "rt_tgsigqueueinfo",
	364:    "perf_event_open",
	365:    "recvmmsg",
	366:    "accept4",
	367:    "fanotify_init",
	368:    "fanotify_mark",
	369:    "prlimit64",
	370:    "name_to_handle_at",
	371:    "open_by_handle_at",
	372:    "clock_adjtime",
	373:    "syncfs",
	374:    "sendmmsg",
	375:    "setns",
	376:    "process_vm_readv",
	377:    "process_vm_writev",
	983041: "breakpoint",
	983042: "cacheflush",
	983043: "usr26",
	983044: "usr32",
	983045: "set_tls",
}
