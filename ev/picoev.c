#ifdef __linux__
	#include "epoll.c"
#elif __APPLE__
	#include "kqueue.c"
#elif defined(__FreeBSD__) || defined(__NetBSD__) || defined(__OpenBSD__) || defined(__DragonFly__)
	#include "kqueue.c"
#else
	#include "select.c"
#endif