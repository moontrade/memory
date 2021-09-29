#ifdef __linux__
	#include "src/epoll.c"
#elif __APPLE__
	#include "src/kqueue.c"
#elif defined(__FreeBSD__) || defined(__NetBSD__) || defined(__OpenBSD__) || defined(__DragonFly__)
	#include "src/kqueue.c"
#else
	#include "src/select.c"
#endif