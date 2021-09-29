package art

/*
#include "picoev.h"
#include <stdlib.h>
#include <errno.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/tcp.h>
#include <fcntl.h>
#include <signal.h>

static inline void setup_sock(int fd) {
	int on = 1;
	setsockopt(fd, IPPROTO_TCP, TCP_NODELAY, &on, sizeof(int));
	fcntl(fd, F_SETFL, O_NONBLOCK);
}

typedef struct {
	int port;
	int fd;
	int reuse_addr_res;
	int reuse_port_res;
	int size;
	int bind_res;
	int listen_res;
	int tcp_nodelay_res;
	int nonblock_res;
} picoev_bind_acceptor_t;

void do_bind_acceptor(uintptr_t arg0, uintptr_t arg1) {
	picoev_bind_acceptor_t* args = (picoev_bind_acceptor_t*)(void*)arg0;
	int fd = socket(AF_INET, SOCK_STREAM, 0);
	args->fd = fd;
	if (fd == -1) {
		return;
	}
	struct sockaddr_in addr;
	addr.sin_family = AF_INET;
	addr.sin_port = htons(args->port);
	addr.sin_addr.s_addr = htonl(INADDR_ANY);
	args->bind_res = bind(fd, (const struct sockaddr*)&addr, sizeof(addr));
	if (args->bind_res != 0) {
		return;
	}
	args->listen_res = listen(fd, SOMAXCONN);
	if (args->listen_res != 0) {
		return;
	}
	int on = 1;
	args->tcp_nodelay_res = setsockopt(fd, IPPROTO_TCP, TCP_NODELAY, &on, sizeof(int));
	if (args->tcp_nodelay_res != 0) {
		return;
	}
	args->nonblock_res = fcntl(fd, F_SETFL, O_NONBLOCK);
	if (args->nonblock_res != 0) {
		return;
	}
}

typedef struct {
	int max_fd;
	int result;
} picoev_init_t;

void do_picoev_init(size_t arg0, size_t arg1) {
	picoev_init_t* args = (picoev_init_t*)(void*)arg0;
	args->result = picoev_init(args->max_fd);
}

typedef struct {
	int result;
} picoev_deinit_t;

void do_picoev_deinit(size_t arg0, size_t arg1) {
	picoev_deinit_t* args = (picoev_deinit_t*)(void*)arg0;
	args->result = picoev_deinit();
}

typedef struct {
	size_t max_timeout;
	int loop;
} picoev_create_loop_t;

void do_picoev_create_loop(size_t arg0, size_t arg1) {
	picoev_create_loop_t* args = (picoev_create_loop_t*)(void*)arg0;
	args->loop = (size_t)picoev_create_loop((int)args->max_timeout);
}

typedef struct {
	size_t loop;
	size_t result;
} picoev_destroy_loop_t;

void do_picoev_destroy_loop(size_t arg0, size_t arg1) {
	picoev_destroy_loop_t* args = (picoev_destroy_loop_t*)(void*)arg0;
	args->result = (size_t)picoev_destroy_loop((picoev_loop*)(void*)args->loop);
}

typedef struct {
	size_t loop;
	int fd;
	int secs;
} picoev_set_timeout_t;

void do_picoev_set_timeout_loop(size_t arg0, size_t arg1) {
	picoev_set_timeout_t* args = (picoev_set_timeout_t*)(void*)arg0;
	picoev_set_timeout((picoev_loop*)(void*)args->loop, args->fd, args->secs);
}

void do_picoev_accept(picoev_loop* loop, int fd, int revents, void* cb_arg) {

}

typedef struct {
	size_t loop;
	int max_wait;
	int result;
} picoev_loop_once_t;

void do_picoev_loop_once(size_t arg0, size_t arg1) {
	picoev_loop_once_t* args = (picoev_loop_once_t*)(void*)arg0;
	args->result = picoev_loop_once((picoev_loop*)(void*)args->loop, args->max_wait);
}

*/
import "C"
import (
	"errors"
	"github.com/moontrade/memory"
	"github.com/moontrade/memory/unsafecgo"
	"sync"
	"unsafe"
)

type Loop C.picoev_loop

type Conn struct {
	FD    int32
	_     int32
	Read  memory.FatPointer
	Write memory.FatPointer
}

var (
	initialized = false
	initResult  = 0
	maxFD       = 1024
	mu          sync.Mutex
)

func MaxFD() int {
	return maxFD
}
func SetMaxFD(maxFd int) {
	maxFD = maxFd
}

type initT struct {
	maxFd  int32
	result int32
}

func Init(maxFd int) int {
	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return initResult
	}
	args := initT{maxFd: int32(maxFd)}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.Call((*byte)(C.do_picoev_init), ptr, 0)

	if args.result == 0 {
		initialized = true
		initResult = 0
	} else {
		initResult = int(args.result)
	}
	return int(args.result)
}

type deinitT struct {
	result int32
}

func Deinit() int {
	mu.Lock()
	defer mu.Unlock()
	if !initialized {
		return 0
	}
	args := deinitT{}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.Call((*byte)(C.do_picoev_deinit), ptr, 0)
	return int(args.result)
}

type createLoop struct {
	maxTimeout uintptr
	loop       uintptr
}

func New(maxTimeout int) *Loop {
	mu.Lock()
	isInitialized := initialized
	mu.Unlock()
	if !isInitialized {
		Init(maxFD)
	}

	args := createLoop{maxTimeout: uintptr(maxTimeout)}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.Call((*byte)(C.do_picoev_create_loop), ptr, 0)
	return (*Loop)(unsafe.Pointer(args.loop))
}

/*
typedef struct {
	int port;
	int fd;
	int reuse_addr_res;
	int reuse_port_res;
	int size;
	int bind_res;
	int listen_res;
	int tcp_nodelay_res;
	int nonblock_res;
} picoev_bind_acceptor_t;
*/
type bindAcceptorT struct {
	port          int32
	fd            int32
	reuseAddrRes  int32
	reusePortRes  int32
	size          int32
	bindRes       int32
	listenRes     int32
	tcpNoDelayRes int32
	nonblockRes   int32
}

var (
	ErrFD         = errors.New("C.socket failed")
	ErrReuseAddr  = errors.New("reuse addr")
	ErrReusePort  = errors.New("reuse port")
	ErrBind       = errors.New("bind")
	ErrListen     = errors.New("listen")
	ErrTCPNoDelay = errors.New("TCP_NODELAY")
	ErrNONBLOCK   = errors.New("NONBLOCK")
)

func (l *Loop) BindAcceptor(port int) error {
	args := bindAcceptorT{port: int32(port)}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.Call((*byte)(C.do_bind_acceptor), ptr, 0)
	if args.fd < 0 {
		return ErrFD
	}
	if args.reuseAddrRes != 0 {
		return ErrReuseAddr
	}
	if args.reusePortRes != 0 {
		return ErrReusePort
	}
	if args.bindRes != 0 {
		return ErrBind
	}
	if args.listenRes != 0 {
		return ErrListen
	}
	if args.tcpNoDelayRes != 0 {
		return ErrTCPNoDelay
	}
	if args.nonblockRes != 0 {
		return ErrNONBLOCK
	}
	return nil
}

type destroyLoop struct {
	loop   uintptr
	result int32
}

func (l *Loop) Destroy() int {
	args := destroyLoop{loop: uintptr(unsafe.Pointer(l))}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.Call((*byte)(C.do_picoev_destroy_loop), ptr, 0)
	return int(args.result)
}

type setTimeout struct {
	loop uintptr
	fd   int32
	secs int32
}

func (l *Loop) SetTimeout(fd, secs int32) {
	args := setTimeout{
		loop: uintptr(unsafe.Pointer(l)),
		fd:   fd,
		secs: secs,
	}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.Call((*byte)(C.do_picoev_set_timeout_loop), ptr, 0)
}

//func (l *Loop) Add(fd, events, timeoutInSecs)

type loopOnceT struct {
	loop    uintptr
	maxWait int32
	result  int32
}

func (l *Loop) Once(maxWait int) int {
	args := loopOnceT{loop: uintptr(unsafe.Pointer(l)), maxWait: int32(maxWait)}
	ptr := uintptr(unsafe.Pointer(&args))
	unsafecgo.Call((*byte)(C.do_picoev_loop_once), ptr, 0)
	return int(args.result)
}

func (l *Loop) run() {
	var (
		args = destroyLoop{loop: uintptr(unsafe.Pointer(l))}
		ptr  = uintptr(unsafe.Pointer(&args))
	)
	for {
		unsafecgo.Call((*byte)(C.do_picoev_loop_once), ptr, 0)
	}
}
